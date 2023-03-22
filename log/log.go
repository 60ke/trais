package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger
var StdoutCore zapcore.Core

func LogInit(level, path string) {
	writeSyncer := getLogWriter(path)
	// writeSyncer := zapcore.NewMultiWriteSyncer(
	// 	zapcore.AddSync(os.Stdout),
	// 	getLogWriter(path),
	// )
	encoder := getEncoder()
	colorEncoder := getColorEncoder()
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}
	fileCore := zapcore.NewCore(encoder, writeSyncer, l)
	StdoutCore = zapcore.NewCore(colorEncoder, zapcore.Lock(os.Stdout), l)
	teeCore := zapcore.NewTee(
		fileCore,
		StdoutCore,
	)
	Logger = zap.New(teeCore, zap.WithCaller(true)).Sugar()
}

func LogReload(level, path string) {
	LogInit(level, path)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getColorEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logPath string) zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     24,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
