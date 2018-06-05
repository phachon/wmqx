package controllers

import (
	"github.com/valyala/fasthttp"
	"wmqx/app"
	"runtime"
	"wmqx/utils"
	"fmt"
	"os/exec"
	"bytes"
	"strings"
	"path/filepath"
	"os"
)

type LogController struct {
	BaseController
}

// return LogController
func NewLogController() *LogController {
	return &LogController{}
}

// log search
func (this *LogController) Search(ctx *fasthttp.RequestCtx) {

	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	number := this.GetCtxInt(ctx, "number")
	keyword := this.GetCtxString(ctx, "keyword")
	level := this.GetCtxString(ctx, "level")

	if number == 0 {
		number = 50
	}
	var filename string
	if level != "" {
		levelFilenameConf := app.Conf.GetStringMapString("log.file.levelFilename")
		if len(levelFilenameConf) > 0 {
			levelFilename, ok := levelFilenameConf[strings.ToLower(level)]
			if ok {
				filename = levelFilename
			}
		}
	}else {
		filename = app.Conf.GetString("log.file.filename")
	}

	logs := []string{}
	if filename == "" {
		this.jsonSuccess(ctx, "success", logs)
		return
	}

	if keyword == "" {
		var err error
		if runtime.GOOS == "windows" {
			logs, err = utils.Tail.Run(filename, number)
			if err != nil {
				this.jsonError(ctx, "read log lines error: "+err.Error(), nil)
				return
			}
		}else {
			command := fmt.Sprintf("tail -n %d %s |tac", number, filename)
			cmd := exec.Command("bash", "-c", command)
			stdOut := &bytes.Buffer{}
			stdErr := &bytes.Buffer{}
			cmd.Stdout = stdOut
			cmd.Stderr = stdErr
			err := cmd.Run()
			if err != nil {
				this.jsonError(ctx, "read log lines error: "+stdErr.String(), nil)
				return
			}
			logs = strings.Split(stdOut.String(), "\n")
		}

	}else {
		if runtime.GOOS == "windows" {
			this.jsonError(ctx, "windows not support keyword search", nil)
			return
		}
		command := fmt.Sprintf("grep \"%s\" %s |tail -n 100|tac", keyword, filename)
		cmd := exec.Command("bash", "-c", command)
		stdOut := &bytes.Buffer{}
		stdErr := &bytes.Buffer{}
		cmd.Stdout = stdOut
		cmd.Stderr = stdErr
		err := cmd.Run()
		if err != nil {
			this.jsonError(ctx, "search log error: "+stdErr.String(), nil)
			return
		}
		logs = strings.Split(stdOut.String(), "\n")
	}

	this.jsonSuccess(ctx, "success", logs)
}

// log list
func (this *LogController) List(ctx *fasthttp.RequestCtx) {

	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	filename := app.Conf.GetString("log.file.filename")
	dir, err := filepath.Abs(filename)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}
	dir = filepath.Dir(dir)

	files, err := utils.File.WalkDir(dir, ".log")
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}

	fileNames := []string{}
	for _, file := range files {
		fileInfo, _ := os.Stat(file)
		fileNames = append(fileNames, fileInfo.Name())
	}

	this.jsonSuccess(ctx, "success", fileNames)
}

// log download
func (this *LogController) Download(ctx *fasthttp.RequestCtx) {

	//r := this.AccessToken(ctx)
	//if r != true {
	//	this.jsonError(ctx, "token error", nil)
	//	return
	//}

	filename := this.GetCtxString(ctx, "filename")
	if filename == "" {
		this.jsonError(ctx, "wmqx log filename not empty!", nil)
		return
	}

	LogFilename := app.Conf.GetString("log.file.filename")
	dir, err := filepath.Abs(LogFilename)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}
	logDir := filepath.Dir(dir)

	filePath := logDir+"/"+filename

	ok, _ := utils.File.PathExists(filePath)
	if ok == false {
		this.jsonError(ctx, "filename not found!", nil)
		return
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/force-download")
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename=\""+fileInfo.Name()+"\"")
	ctx.Response.Header.Set("Content-Transfer-Encoding", "binary")
	ctx.SendFile(filePath)
}


