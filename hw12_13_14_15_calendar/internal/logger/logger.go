package logger

import (
	"errors"
	"fmt"
	commonLog "log"
	"os"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/sirupsen/logrus"
)

var ErrOpenFailed = errors.New("can't open logfile")

var level = map[string]logrus.Level{
	"info":    logrus.InfoLevel,
	"error":   logrus.ErrorLevel,
	"warning": logrus.WarnLevel,
	"debug":   logrus.DebugLevel,
}

type Logger struct {
	Path     string
	log      *logrus.Logger
	entry    *logrus.Entry
	logClose func()
}

func New(cfg *config.Config) (*Logger, error) {
	if _, ok := level[cfg.Logger.Level]; !ok {
		return nil, fmt.Errorf("unexpected logger level %s", cfg.Logger.Level)
	}
	l := &Logger{}
	l.log = logrus.New()
	l.log.SetLevel(level[cfg.Logger.Level])
	l.logClose = func() {}

	if len(cfg.Logger.Path) > 0 {
		l.Path = cfg.Logger.Path
		fh, err := os.OpenFile(l.Path, os.O_WRONLY|os.O_CREATE, 0o755)
		if err != nil {
			return l, ErrOpenFailed
		}
		l.log.SetOutput(fh)

		l.logClose = func() {
			if err := fh.Close(); err != nil {
				commonLog.Println(fmt.Errorf("can't close logfile: %w", err))
			}
		}
	}

	l.SetEntry(logrus.Fields{}) // default

	return l, nil
}

func (l *Logger) SetEntry(fields logrus.Fields) {
	l.entry = l.log.WithFields(fields)
}

func (l *Logger) Info(msg string) {
	l.entry.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.entry.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.entry.Error(msg)
}

func (l *Logger) Log(msg string) {
	l.entry.Error(msg)
}

func (l *Logger) Close() {
	l.logClose()
}
