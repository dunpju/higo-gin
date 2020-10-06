package higo

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"higo.yumi.com/src/higo/utils"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var Logrus *logrus.Logger

func Log(root string)  {
	// 日志文件
	logs := root + "runtime/logs"

	// 目录不存在，并创建
	if _, err := os.Stat(logs); os.IsNotExist(err) {
		if os.Mkdir(logs, os.ModePerm) != nil {}
	}

	fileName := logs + "/higo" // + fmt.Sprintf(".%s.log", time.Now().Format("20060102"))

	// 实例化
	Logrus = logrus.New()

	/**
	// 写入文件
	src, _ := os.Create(fileName)
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	// 设置输出
	Logrus.Out = src
	*/

	// 设置日志级别
	Logrus.SetLevel(logrus.DebugLevel)

	// 设置时间格式
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true // INFO[2020-09-28 13:20:14]
	Logrus.SetFormatter(customFormatter)

	// 设置 rotatelogs
	logWriter, _ := rotatelogs.New(
		// 分割后的文件名称
		fileName + ".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp: true,
	})

	// 新增 Hook
	Logrus.AddHook(lfHook)
}

// 输出换行debug调用栈
func PrintlnStack()  {
	ds := fmt.Sprintf("%s", debug.Stack())
	dss := strings.Split(ds,"\n")
	Logrus.Info(fmt.Sprintf("=== DEBUG STACK Bigin goroutine %d ===", utils.GoroutineID()))
	for _, b := range dss {
		Logrus.Info(fmt.Sprintf("%s", b))
	}
	Logrus.Info(fmt.Sprintf("=== DEBUG STACK End goroutine %d ===", utils.GoroutineID()))
}