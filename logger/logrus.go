package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	ReportCaller bool         `json:",optional"`
	PrettyPrint  bool         `json:",optional"`
	Level        logrus.Level `json:",default=4,options=[0,1,2,3,4,5,6]"`
	ObjectName   string       `json:"object_name"`
	WriteFile    bool         `json:"write_file"`
}

// InitLogrus 初始化配置logrus
func InitLogrus(c Config) {
	// 设置输出文件
	if c.WriteFile {
		setLogOutput(fmt.Sprintf("./logs/%s.log", c.ObjectName))
	}
	logrus.SetLevel(c.Level)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		DisableTimestamp:  false,
		DisableHTMLEscape: true,
		DataKey:           "",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "msg",
			logrus.FieldKeyFile:  "func",
		},
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", strings.Replace(frame.Function, filepath.Ext(frame.Function),
				fmt.Sprintf("/%s:(%d)", filepath.Base(frame.File), frame.Line), -1)
		},
		PrettyPrint: false,
	})
}

// mkdir 检查目录是否存在，不存在则创建
func mkdir(filePath string) {
	dir := filepath.Dir(filePath)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0664)
			if err != nil {
				logrus.Fatalf("[Mkdir Error]: make dir %s error : %s", dir, err)
			}
		} else {
			logrus.Fatalf("[Stat dir Error]: stat dir %s error : %s", dir, err)
		}
	}
}

// setLogOutput 设置输出文件
func setLogOutput(filePath string) {
	mkdir(filePath)
	// 创建今天的日志文件
	createNewLogFile(filePath)
	//开启定时任务，每天0点替换日志
	now := time.Now()
	// 计算下一个0点
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	t := time.NewTimer(next.Sub(now))
	go func() {
		<-t.C
		createNewLogFile(filePath)
		tc := time.NewTicker(24 * time.Hour)
		for range tc.C {
			createNewLogFile(filePath)
		}
	}()
}

// createNewLogFile 创建日志文件&设置
func createNewLogFile(filePath string) {
	dateStr := time.Now().Format("2006-01-02")
	filePath = fmt.Sprintf("%s-%s", filePath, dateStr)
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		logrus.Errorf("open %s file failed [error: %s]", filePath, err.Error())
		return
	}
	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)
}
