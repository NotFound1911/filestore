package zap

import (
	"fmt"
	"github.com/NotFound1911/filestore/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	logDebugLevel string = "debug"
	logInfoLevel  string = "info"
	logWarnLevel  string = "warn"
	logErrorLevel string = "error"
	logPanicLevel string = "panic"
	logFatalLevel string = "fatal"
)

type Service struct {
	Logger *zap.Logger
}

func NewService(conf *config.Configuration, name string) *Service {
	var encoder zapcore.Encoder

	// 调整编码器默认配置 输出内容
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(name + "." + l.String())
	}

	// 设置编码器，日志的输出格式
	if conf.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 同时输出到控制台和文件
	var multiWS zapcore.WriteSyncer
	if conf.Log.EnableFile {
		log := &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s.log", conf.Log.RootDir, name),
			MaxSize:    conf.Log.MaxSize,
			MaxAge:     conf.Log.MaxSize,
			MaxBackups: conf.Log.MaxBackups,
			Compress:   conf.Log.Compress,
		}
		multiWS = zapcore.NewMultiWriteSyncer(zapcore.AddSync(log), zapcore.AddSync(os.Stdout))
	} else {
		multiWS = zapcore.AddSync(os.Stdout)
	}
	core := zapcore.NewCore(encoder, multiWS, configToZapLevel(conf))
	logger := zap.New(core, zap.AddCaller())
	return &Service{
		Logger: logger,
	}
}
func configToZapLevel(conf *config.Configuration) zapcore.Level {
	level := conf.Log.Level
	switch level {
	case logDebugLevel:
		return zap.DebugLevel
	case logInfoLevel:
		return zap.InfoLevel
	case logWarnLevel:
		return zap.WarnLevel
	case logErrorLevel:
		return zap.ErrorLevel
	case logPanicLevel:
		return zap.PanicLevel
	case logFatalLevel:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
