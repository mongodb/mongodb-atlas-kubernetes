package log

import (
	"flag"
	"fmt"
	"io"
	"strings"

	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug  = "debug"
	LevelInfo   = "info"
	LevelWarn   = "warn"
	LevelError  = "error"
	LevelDPanic = "dpanic"
	LevelPanic  = "panic"
	LevelFatal  = "fatal"
)

const (
	FormatJSON    = "json"
	FormatConsole = "console"
)

type Config struct {
	Level   string
	Format  string
	EnvTest bool
}

var AvailableLogLevels = []string{
	LevelDebug,
	LevelInfo,
	LevelWarn,
	LevelError,
	LevelDPanic,
	LevelPanic,
	LevelFatal,
}

var AvailableLogFormats = []string{
	FormatJSON,
	FormatConsole,
}

func DefaultConfig() Config {
	return Config{
		Level:   LevelInfo,
		Format:  FormatJSON,
		EnvTest: false,
	}
}

func RegisterFlags(conf *Config, fs *flag.FlagSet) {
	fs.StringVar(&conf.Level, "log-level", LevelInfo, fmt.Sprintf("Log level to use. Available values: %s", strings.Join(AvailableLogLevels, ", ")))
	fs.StringVar(&conf.Format, "log-format", FormatJSON, fmt.Sprintf("Log format to use. Available values: %s", strings.Join(AvailableLogFormats, ", ")))
}

func NewLogger(conf Config, out io.Writer) (*zap.Logger, error) {
	level := zap.AtomicLevel{}
	err := level.UnmarshalText([]byte(strings.ToLower(conf.Level)))
	if err != nil {
		return nil, fmt.Errorf("log level unknown, supported values are: %s\n%w", strings.Join(AvailableLogLevels, ","), err)
	}

	format := strings.ToLower(conf.Format)
	if format != "json" && format != "console" {
		return nil, fmt.Errorf("log format unknown, supported values are: %s", strings.Join(AvailableLogFormats, ","))
	}

	ctrzap.NewRaw(ctrzap.UseDevMode(conf.EnvTest), ctrzap.WriteTo(out), ctrzap.StacktraceLevel(level))

	cfg := zap.Config{
		Level:             level,
		OutputPaths:       []string{"stdout"},
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}

	return cfg.Build()
}
