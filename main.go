package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"slices"
	"strings"
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

	var reader io.Reader = os.Stdin
	if *srcFlag != "-" {
		f, err := os.Open(*srcFlag)
		if err != nil {
			return fmt.Errorf("failed to read source: %w", err)
		}
		defer f.Close()
		reader = f
	}

	var writer io.Writer = os.Stdout
	if *dstFlag != "" {
		f, err := os.Create(*dstFlag)
		if err != nil {
			return fmt.Errorf("failed to open destination: %w", err)
		}
		defer f.Close()
		writer = f
	}

	src := csv.NewReader(reader)
	hdr, err := src.Read()
	if err != nil {
		return fmt.Errorf("could not read header: %w", err)
	}
	idx := slices.Index(hdr, "Trans. Date")
	if idx == -1 {
		return fmt.Errorf("could not find transaction date column")
	}
	hdr[idx] = "Date" // rename column

	idx = slices.Index(hdr, "Amount")
	if idx == -1 {
		return fmt.Errorf("could not find amount column")
	}

	dst := csv.NewWriter(writer)
	dst.Write(hdr)

	for ctx.Err() == nil {
		row, err := src.Read()
		if row == nil && errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read next row: %w", err)
		}

		// flip the sign in string land so we don't have to worry about floating
		// point
		amt, wasNegative := strings.CutPrefix(row[idx], "-")
		if wasNegative {
			row[idx] = amt
		} else {
			row[idx] = "-" + amt
		}
		if err := dst.Write(row); err != nil {
			return fmt.Errorf("could not write next row: %w", err)
		}
	}

	dst.Flush()
	return dst.Error()
}
