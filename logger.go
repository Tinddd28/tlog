package logger

import (
	"log/slog"
	"os"

	"github.com/Tinddd28/tlog/handlers/tlog"
)

type LogOpts struct {
	level      string `yaml:"level"`
	dir        string `yaml:"dir"`
	format     string `yaml:"format"`
	savingDays int    `yaml:"saving_days"`
}

func NewLogOpts(level string, dir string, format string, savingDays int) *LogOpts {
	return &LogOpts{
		level:      level,
		dir:        dir,
		format:     format,
		savingDays: savingDays,
	}
}

func SetupLogger(opts LogOpts) (*slog.Logger, error) {
	var log_ *slog.Logger
	log_, err := SetupTlogLogger(opts)
	if err != nil {
		return nil, err
	}

	return log_, nil
}

func SetupTlogLogger(op LogOpts) (*slog.Logger, error) {
	if err := os.MkdirAll(op.dir, os.ModePerm); err != nil {
		return nil, err
	}
	var lvl slog.Level
	switch op.level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "error":
		lvl = slog.LevelError
	case "warn":
		lvl = slog.LevelWarn
	}
	opts := tlog.LoggerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: lvl,
		},
		LogDir:     op.dir,
		Format:     op.format,
		SavingDays: op.savingDays,
	}

	handler, err := opts.NewLogger()
	if err != nil {
		return nil, err
	}

	return slog.New(handler), nil
}
