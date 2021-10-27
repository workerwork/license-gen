package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"license-gen/conf"
	"time"
)

// Setup 日志初始化设置
func Setup() {
	switch strings.ToLower(conf.LoggerConf.Level) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	//zerolog.TimeFieldFormat = time.RFC3339
	if conf.LoggerConf.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    !conf.LoggerConf.Color,
		})
	}
	log.Info().Msg("logger on work now!")
}
