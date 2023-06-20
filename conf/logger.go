package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

func InitLogger() *zap.SugaredLogger {
	logMode := zapcore.DebugLevel
	// 定义日志级别 
	if !viper.GetBool("mode.develop") {
		logMode = zapcore.InfoLevel
	}
	core := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(getWriteSyncer(), zapcore.AddSync(os.Stdout)), logMode)
	return zap.New(core).Sugar()
}

// 定义输出格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// key值显示：time
	encoderConfig.TimeKey = "time"
	// 日志级别显示大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 时间格式化: YYYY-MM-DD HH-MM-SS; > 1.20版本
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Local().Format(time.DateTime))
	}

	return zapcore.NewJSONEncoder(encoderConfig)
}

// 定义输出
func getWriteSyncer() zapcore.WriteSyncer {
	stSeparator := string(filepath.Separator)
	stRootDir, _ := os.Getwd()
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format(time.DateOnly) + ".log"
	fmt.Println(stLogFilePath)
        // 日志切割
	luberjackSyncer := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    viper.GetInt("log.MaxSize"),    // 日志文件最大的尺寸(M), 超限后开始自动分割
		MaxBackups: viper.GetInt("log.MaxBackups"), // 保留旧文件的最大个数
		MaxAge:     viper.GetInt("log.MaxAge"),     // 保留旧文件的最大天数
		Compress:   false,
	}

	return zapcore.AddSync(luberjackSyncer)
}
