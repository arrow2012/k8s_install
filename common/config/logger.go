package config

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func initLogger(logpath string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logpath, // 日志文件路径
		MaxSize:    1,     // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,      // 日志文件最多保存多少个备份
		MaxAge:     7,       // 文件最多保存多少天
		Compress:   true,    // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,      // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)

	//core := zapcore.NewCore(
	//	zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
	//	zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
	//	atomicLevel,                                                                     // 日志级别
	//)

	Encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewTee(
		// 打印在kafka topic中（伪造的case）
		//zapcore.NewCore(kafkaEncoder, topicErrors, atomicLevel),
		// 打印在控制台
		zapcore.NewCore(Encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), atomicLevel),
		// 打印在文件中
		zapcore.NewCore(Encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), atomicLevel),
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "k8s-install"))
	// 构造日志
	logger := zap.New(core, caller, development, filed)
	return logger
}

var (
	Logger = initLogger("/tmp/zap.log")
)

