package pkg

import (
	"os"
	"time"

	"github.com/serialt/sync/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func LevelToZapLevel(level string) zapcore.Level {
	// 转换日志级别
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel

	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}

}

func NewLogger() *zap.Logger {

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(LevelToZapLevel(config.LogLevel))

	// 输出的消息
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // zapcore.CapitalLevelEncoder //按级别显示不同颜色
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05")) //指定时间格式
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 日志输出类型
	var encoder zapcore.Encoder
	switch config.LogType {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)

	}

	var core zapcore.Core
	if config.LogFile == "" {
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), atomicLevel)
	} else {
		file := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.LogFile,       // 日志文件
			MaxSize:    config.LogMaxSize,    // 单个日志文件大小，单位M
			MaxBackups: config.LogMaxBackups, // 轮转保留个数
			MaxAge:     config.LogMaxAge,     // 最长保留时间，单位天
			Compress:   config.LogCompress,   // 默认不压缩
		})
		core = zapcore.NewCore(encoder, zapcore.AddSync(file), atomicLevel)
	}

	// 开启开发模式，堆栈跟踪: [zap.AddCaller()]
	Logger := zap.New(core, zap.AddCaller())
	return Logger
}

func NewSugarLogger() *zap.SugaredLogger {
	sugraLog := NewLogger()
	return sugraLog.Sugar()

}

var Logger *zap.Logger
var Sugar *zap.SugaredLogger
