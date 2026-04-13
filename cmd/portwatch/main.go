package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/output"
	"github.com/user/portwatch/internal/scanner"
)

func main() {
	cfg, err := config.LoadOrDefault(configPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: failed to load config: %v\n", err)
		os.Exit(1)
	}

	sc, err := scanner.NewScanner(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: failed to create scanner: %v\n", err)
		os.Exit(1)
	}

	f := filter.New(cfg)
	fmt := output.NewFormatter(cfg)
	writer := output.NewConsoleNotifierWithWriter(os.Stdout, fmt)

	logNotifier := alert.NewLogNotifier()
	dispatcher := alert.NewDispatcher()
	dispatcher.AddNotifier(writer)
	dispatcher.AddNotifier(logNotifier)

	d := daemon.New(cfg, sc, f, dispatcher)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt2 := "portwatch: starting — interval %s, scanning all ports\n"
	fmt.Fprintf(os.Stderr, fmt2, cfg.Interval)

	if err := d.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "portwatch: stopped\n")
}

func configPath() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return ""
}
