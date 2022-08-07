package main

import (
	"context"
	"flag"
	"fmt"
	commonLog "log"
	"os"
	"os/signal"
	"time"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/app"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/logger"
	httpserver "github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/sqls"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Config
	cfg, err := config.New(configFile)
	if err != nil {
		commonLog.Fatal(fmt.Errorf("initialization config error: %w", err))
		return
	}

	// Logger
	log, err := logger.New(cfg)
	if err != nil {
		commonLog.Fatal(fmt.Errorf("initialization log error: %w", err))
		return
	}

	// Storage
	var storage model.EventsM
	switch cfg.StorageMode {
	case "in_memory":
		storage, err = memory.NewConnection()
		if err != nil {
			log.Error(err.Error())
			return
		}
	case "sql":
		storage, err = sqls.NewConnection(cfg)
		if err != nil {
			log.Error(err.Error())
			return
		}
	default:
		log.Error("Incorrect storage mode")
		return
	}

	calendarApp := app.NewCalendar(storage)

	chErrors := make(chan error)

	// HTTP Server
	httpServer := httpserver.NewServer(cfg, log, calendarApp)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := httpServer.Start(ctx); err != nil {
			chErrors <- err
		}
	}()
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			log.Error(fmt.Sprintf("stopping of http server error: %s", err))

			return
		}
	}()
}
