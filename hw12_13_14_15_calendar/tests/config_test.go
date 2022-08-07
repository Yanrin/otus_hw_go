package main

import (
	"testing"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("Empty path", func(t *testing.T) {
		_, err := config.New("")
		require.Equal(t, config.ErrFilePathEmpty, err)
	})

	t.Run("Wrong json data", func(t *testing.T) {
		_, err := config.New("testdata/config/wrong.yml")
		require.Equal(t, config.ErrReadFile, err)
	})

	t.Run("Correct data", func(t *testing.T) {
		cfg, err := config.New("testdata/config/config.yml")
		require.Equal(t, nil, err)
		require.Equal(t, cfg.HTTP.Host, "localhost")
	})
}
