// 浣滅敤锛氬綋鍓嶅井鏈嶅姟涓撶敤鐨勫唴閮ㄥ伐鍏峰簱锛氭棩蹇楃粍浠跺皝瑁咃紙闈㈠悜渚濊禆娉ㄥ叆璁捐锛?
package logger

import (
	"os"

	"github.com/google/wire"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ProviderSet 渚?Wire 浣跨敤鐨勪緷璧栨彁渚涜€呴泦鍚?
var ProviderSet = wire.NewSet(NewLogger)

// Options 瀹氫箟鏃ュ織鐨勯厤缃弬鏁?
type Options struct {
	Filename   string // 鏃ュ織鏂囦欢璺緞 (濡?"logs/app.log")锛屽鏋滀负绌哄垯鍙緭鍑哄埌鎺у埗鍙?
	Level      string // 鏃ュ織杈撳嚭绾у埆 (濡?"debug", "info", "warn", "error")
	MaxSize    int    // 鍗曚釜鏃ュ織鏂囦欢鏈€澶у昂瀵?(MB)
	MaxBackups int    // 鏈€澶氫繚鐣欏灏戜釜鏃ф棩蹇楁枃浠?
	MaxAge     int    // 鏈€澶氫繚鐣欏灏戝ぉ鏃ф棩蹇楁枃浠?
	Compress   bool   // 鏄惁鍘嬬缉褰掓。
}

// NewLogger 鏃ュ織瀹炰緥鐨勬瀯閫犲嚱鏁帮紝浣滀负 Wire 鐨?Provider
// 閫氳繃浼犻€?Options 閫夐」鏉ョ敓鎴愬苟杩斿洖涓€涓師鐢熺殑 *zap.Logger
func NewLogger(opts *Options) (*zap.Logger, error) {
	// 瑙ｆ瀽鏃ュ織绾у埆
	var level zapcore.Level
	err := level.UnmarshalText([]byte(opts.Level))
	if err != nil {
		level = zap.InfoLevel // 榛樿鍥為€€涓?Info 绾у埆
	}

	encoder := getEncoder()
	var core zapcore.Core

	// 濡傛灉閰嶇疆浜嗘枃浠惰矾寰勶紝鍒欏悓鏃惰緭鍑哄埌鏂囦欢鍜屾帶鍒跺彴锛涘惁鍒欏彧杈撳嚭鍒版帶鍒跺彴
	if opts.Filename != "" {
		writeSyncer := getLogWriter(opts.Filename, opts.MaxSize, opts.MaxBackups, opts.MaxAge, opts.Compress)
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		)
	} else {
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	}

	// 娉ㄥ叆鏃讹紝涓嶉渶瑕?AddCallerSkip(1)锛屽洜涓鸿皟鐢ㄦ柟灏嗙洿鎺ユ寔鏈夊苟璋冪敤杩欎釜 logger 瀹炰緥
	logger := zap.New(core, zap.AddCaller())
	return logger, nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 榛樿涓?ConsoleEncoder锛屽鏋滈渶瑕?JSON 鏍煎紡鍙敼涓?zapcore.NewJSONEncoder(encoderConfig)
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackups, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
