package app

import (
	"github.com/spf13/viper"
	"os"
	"flag"
	"github.com/phachon/go-logger"
	"log"
	"strings"
)

// app bootstrap init

var (
	flagConf = flag.String("conf", "config.toml", "please input conf path")
)

var (
	RootPath = ""
	AppPath = ""
	Conf = viper.New()
	Log = go_logger.NewLogger()
)

// 启动初始化
func init()  {

	initFlag()
	initPath()
	initConfig()
	initLog()
}

// init flag
func initFlag() {
	flag.Parse()
}

// init dir and path
func initPath() {
	AppPath, _ = os.Getwd()
	RootPath = strings.Replace(AppPath, "app", "", 1)
}

// init config
func initConfig()  {

	if *flagConf == "" {
		log.Println("config file not found!")
		os.Exit(1)
	}

	Conf.SetConfigType("toml")
	Conf.SetConfigFile(*flagConf)
	err := Conf.ReadInConfig()
	if err != nil {
		log.Println("Fatal error config file: "+err.Error())
		os.Exit(1)
	}

	file := Conf.ConfigFileUsed()
	if(file != "") {
		log.Println("Use config file: " + file)
	}
}

// init log
func initLog() {

	Log.Detach("console")
	consoleConfig := &go_logger.ConsoleConfig{
		Color: true, // 显示颜色
	}
	Log.Attach("console", go_logger.NewConfigConsole(consoleConfig))

	filename := Conf.GetString("log.filename")
	maxSize := Conf.GetInt64("log.maxSize")
	maxLine := Conf.GetInt64("log.maxLine")
	dateSlice := Conf.GetString("log.dateSlice")
	jsonFormat := Conf.GetBool("log.jsonFormat")

	fileConfig := &go_logger.FileConfig{
		Filename: filename,
		MaxSize: maxSize,
		MaxLine: maxLine,
		DateSlice: dateSlice,
		JsonFormat: jsonFormat,
	}

	Log.Attach("file", go_logger.NewConfigFile(fileConfig))

	// 设置日志级别
	Log.SetLevel(go_logger.LOGGER_LEVEL_DEBUG)
}