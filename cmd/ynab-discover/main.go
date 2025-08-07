package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	yd "github.com/ryepup/ynab-discover"
)

var srcFlag = flag.String("src", "-", "where to read the file downloaded from discover, '-' for stdin")
var dstFlag = flag.String("dst", "", "where to write the file to import to YNAB, or empty for stdout")

func main() {
	exitCode := 0
	defer func() { os.Exit(exitCode) }()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		exitCode = -1
	}
}

func run(ctx context.Context) error {
	flag.Parse()

	reader := os.Stdin
	if *srcFlag != "-" {
		f, err := os.Open(*srcFlag)
		if err != nil {
			return fmt.Errorf("failed to read source: %w", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				slog.Warn("failed to close source file", "error", err)
			}
		}()
		reader = f
	}

	writer := os.Stdout
	if *dstFlag != "" {
		f, err := os.Create(*dstFlag)
		if err != nil {
			return fmt.Errorf("failed to open destination: %w", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				slog.Warn("failed to close destination file", "error", err)
			}
		}()
		writer = f
	}

	return yd.ConvertCSV(ctx, reader, writer)
}
