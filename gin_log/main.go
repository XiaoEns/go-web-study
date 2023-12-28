package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"strings"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync() // 保存缓存中的日志
	simpleHttpGet("www.5lmh.com")
	simpleHttpGet("http://www.google.com")
}

func InitLogger() {
	// 日志级别
	logLevel := strings.ToUpper("info")

	atomicLevel := zap.NewAtomicLevel()

	switch logLevel {
	case "DEBUG":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "WARN":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case "DPANIC":
		atomicLevel.SetLevel(zapcore.DPanicLevel)
	case "PANIC":
		atomicLevel.SetLevel(zapcore.PanicLevel)
	case "FATAL":
		atomicLevel.SetLevel(zapcore.FatalLevel)
	}

	writeSyncer := getLogWriter()
	encoder := getEncoder()
	// 输出到日志
	//core := zapcore.NewCore(encoder, writeSyncer, atomicLevel)
	// 输出到日志和控制台
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writeSyncer, atomicLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), atomicLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()

	logger.Sugar().Debug("test1280's debug")
	logger.Sugar().Infof("test1280's %s", "infof")
	logger.Sugar().Warnf("test1280's %s", "warnf")

}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	//如果想要追加写入可以查看我的博客文件操作那一章
	file, _ := os.Create("./gin_log/test.log")
	return zapcore.AddSync(file)
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

//func main() {
//	gin.DisableConsoleColor()
//
//	// Logging to a file.
//	f, _ := os.Create("./gin_log/gin.log")
//	//gin.DefaultWriter = io.MultiWriter(f)
//
//	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
//	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
//	r := gin.Default()
//	r.GET("/ping", func(c *gin.Context) {
//		c.String(200, "pong")
//	})
//	r.Run()
//}
