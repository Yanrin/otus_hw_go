package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	ConfigPath = "testdata/config/config.yml"
)

func TestLog(t *testing.T) {
	cfg, cfgerr := config.New(ConfigPath)

	t.Run("cfg err", func(t *testing.T) {
		require.Nil(t, cfgerr)
	})
	t.Run("write to file", func(t *testing.T) {
		testfile, err := setTestLogFile(cfg)
		require.Nil(t, err)
		defer os.Remove(testfile)

		log, err := logger.New(cfg)
		require.Nil(t, err)

		msg := "first info"
		log.Info(msg)
		log.Close()

		writed, err := readLog(cfg.Logger.Path)
		require.Nil(t, err)

		require.Contains(t, writed, msg)
	})
	t.Run("write to file with fields", func(t *testing.T) {
		testfile, err := setTestLogFile(cfg)
		require.Nil(t, err)
		defer os.Remove(testfile)

		log, err := logger.New(cfg)
		require.Nil(t, err)

		log.SetEntry(logrus.Fields{
			"k": "key",
			"v": "value",
		})
		msg := "second info"
		log.Info(msg)
		log.Close()

		writed, err := readLog(cfg.Logger.Path)
		require.Nil(t, err)

		require.Contains(t, writed, msg)
		require.Contains(t, writed, "k=key")
		require.Contains(t, writed, "v=value")
	})
	t.Run("set level", func(t *testing.T) {
		testfile, err := setTestLogFile(cfg)
		require.Nil(t, err)
		defer os.Remove(testfile)

		cfg.Logger.Level = "warning"
		log, err := logger.New(cfg)
		require.Nil(t, err)

		msgInfo := "unimportant information"
		log.Info(msgInfo)
		msgWarn := "alert"
		log.Warn(msgWarn)
		log.Close()

		writed, err := readLog(cfg.Logger.Path)
		require.Nil(t, err)

		require.NotContains(t, writed, msgInfo)
		require.Contains(t, writed, msgWarn)
	})
}

func readLog(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("can't read logfile to log %s: %w", path, err)
	}

	return string(content), nil
}

func setTestLogFile(cfg *config.Config) (string, error) {
	fn, err := os.CreateTemp("/tmp", "testlog")
	if err != nil {
		return "", err
	}
	cfg.Logger.Path = fn.Name()
	return fn.Name(), nil
}
