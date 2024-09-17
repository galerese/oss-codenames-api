package logging

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func New(level string, format string) Logger {
	const MAX_CALLER_SIZE = 25
	const MAX_CALLER_SIZE_FORMAT = "%25.25s"

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "date"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	config := zap.NewProductionConfig()
	config.Sampling = nil

	if level == "debug" {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	if format == "console" {
		config.Encoding = "console"
		encoderConfig.ConsoleSeparator = " "

		encoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
			pae.AppendString(color.HiCyanString(t.Format("2006-01-02 15:04:05")))
		}

		encoderConfig.EncodeLevel = func(l zapcore.Level, pae zapcore.PrimitiveArrayEncoder) {
			zapcore.CapitalColorLevelEncoder(l, pae)
			if len(l.String()) < 5 {
				pae.AppendString("")
			}
		}
		encoderConfig.EncodeCaller = func(ec zapcore.EntryCaller, pae zapcore.PrimitiveArrayEncoder) {
			text := ec.TrimmedPath()
			if len(text) > MAX_CALLER_SIZE {
				text = text[len(text)-MAX_CALLER_SIZE:]
			}
			pae.AppendString(color.HiBlackString(fmt.Sprintf(MAX_CALLER_SIZE_FORMAT, text)))

		}
	}

	config.EncoderConfig = encoderConfig

	var err error
	log, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	// Making Zap the default logger :)
	zap.RedirectStdLog(log)
	return Logger{log.Sugar()}
}
