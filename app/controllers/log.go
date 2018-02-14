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

// log index
func (this *LogController) Index(ctx *fasthttp.RequestCtx) {

	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	number := this.GetCtxInt(ctx, "number")
	if number == 0 {
		number = 50
	}
	filename := app.Conf.GetString("log.filename")

	var err error
	logs := []string{}
	if runtime.GOOS == "windows" {
		logs, err = utils.NewTail().Run(filename, number)
		if err != nil {
			this.jsonError(ctx, "read log lines error: "+err.Error(), nil)
			return
		}
	}else {
		command := fmt.Sprintf("tail -n %d %s", number, filename)
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

	this.jsonSuccess(ctx, "success", logs)
}

// log search
func (this *LogController) Search(ctx *fasthttp.RequestCtx) {

	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	number := this.GetCtxInt(ctx, "number")
	if number == 0 {
		number = 50
	}
	keyword := this.GetCtxString(ctx, "keyword")
	if keyword == "" {
		this.jsonError(ctx, "keyword not empty", nil)
		return
	}

	filename := app.Conf.GetString("log.filename")
	logs := []string{}
	if runtime.GOOS == "windows" {
		this.jsonError(ctx, "windows not support", nil)
		return
	}
	command := fmt.Sprintf("grep \"%s\" %s |tail -n 100", keyword, filename)
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

	this.jsonSuccess(ctx, "success", logs)
}

// log list
func (this *LogController) List(ctx *fasthttp.RequestCtx) {

	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	filename := app.Conf.GetString("log.filename")
	dir, err := filepath.Abs(filename)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}
	dir = filepath.Dir(dir)

	files, err := utils.NewFile().WalkDir(dir, ".log")
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

	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	filename := this.GetCtxString(ctx, "filename")
	if filename == "" {
		this.jsonError(ctx, "filename not empty!", nil)
		return
	}

	LogFilename := app.Conf.GetString("log.filename")
	dir, err := filepath.Abs(LogFilename)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}
	logDir := filepath.Dir(dir)

	if runtime.GOOS == "windowds" {}
	filePath := logDir+"/"+filename

	ok, _ := utils.NewFile().PathExists(filePath)
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


