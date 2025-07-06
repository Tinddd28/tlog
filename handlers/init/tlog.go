package tlog

import (
	"context"
	"encoding/json"
	"fmt"
	stdLog "log"
	"os"
	"path/filepath"
	"time"

	"log/slog"
)

type LoggerOptions struct {
	SlogOpts   *slog.HandlerOptions
	LogDir     string
	Format     string
	SavingDays int
}

type HandlerOptions struct {
	opts LoggerOptions
	slog.Handler
	l               *stdLog.Logger
	attrs           []slog.Attr
	logFile         *os.File
	currentFileName string
}

func (opts LoggerOptions) NewLogger() (*HandlerOptions, error) {
	h := &HandlerOptions{
		opts: opts,
	}
	if err := h.rotateLogs(); err != nil {
		return nil, err
	}

	go h.scheduleRotation()
	go h.cleanupOldLogs()
	return h, nil
}

func (h *HandlerOptions) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[2006-01-02 15:04:05]")
	msg := r.Message

	h.l.Println(
		timeStr,
		level,
		msg,
		string(b),
	)

	return nil
}

func (h *HandlerOptions) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerOptions{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *HandlerOptions) WithGroup(name string) slog.Handler {
	// TODO: implement
	return &HandlerOptions{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
func (h *HandlerOptions) scheduleRotation() {
	ticker := time.NewTicker(time.Hour * 12)
	defer ticker.Stop()

	for range ticker.C {
		h.rotateLogs()
	}
}

func (h *HandlerOptions) rotateLogs() error {
	fileName := time.Now().Format("2006-01-02") + h.opts.Format
	if fileName == h.currentFileName {
		return nil
	}

	if h.logFile != nil {
		h.logFile.Close()
	}

	filePath := filepath.Join(h.opts.LogDir, fileName)
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	h.logFile = logFile
	h.currentFileName = fileName

	h.l = stdLog.New(logFile, "", 0)
	h.Handler = slog.NewJSONHandler(logFile, h.opts.SlogOpts)

	stdLog.Println("Change on new log-file: ", fileName)
	return nil
}

func (h *HandlerOptions) cleanupOldLogs() {
	if h == nil {
		fmt.Println("TimeHandler nil pointer")
		return
	}
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		files, err := os.ReadDir(h.opts.LogDir)
		if err != nil {
			stdLog.Println("Error read dir logs")
			continue
		}
		now := time.Now()

		for _, file := range files {
			if !file.IsDir() {
				filePath := filepath.Join(h.opts.LogDir, file.Name())
				info, err := os.Stat(filePath)
				if err != nil {
					continue
				}
				// Если существует больше указанного времени - удаляем файл
				if now.Sub(info.ModTime()) > time.Duration(h.opts.SavingDays)*time.Hour*24 {
					stdLog.Println("Deleting old log-file: ", filePath)
					os.Remove(filePath)
				}
			}
		}
	}
}
