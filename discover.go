package discover

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

// ConvertCSV converts a Discover Bank CSV file to YNAB conventions.
// It renames the "Trans. Date" column to "Date" and flips the sign of amounts.
func ConvertCSV(ctx context.Context, reader io.Reader, writer io.Writer) error {
	src := csv.NewReader(reader)
	hdr, err := src.Read()
	if err != nil {
		return fmt.Errorf("could not read header: %w", err)
	}

	// Rename "Trans. Date" column to "Date"
	idx := slices.Index(hdr, "Trans. Date")
	if idx == -1 {
		return fmt.Errorf("could not find transaction date column")
	}
	hdr[idx] = "Date"

	// Find the Amount column
	amountIdx := slices.Index(hdr, "Amount")
	if amountIdx == -1 {
		return fmt.Errorf("could not find amount column")
	}

	dst := csv.NewWriter(writer)
	if err := dst.Write(hdr); err != nil {
		return fmt.Errorf("could not write header: %w", err)
	}

	for ctx.Err() == nil {
		row, err := src.Read()
		if row == nil && errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read next row: %w", err)
		}

		// Flip the sign in string land so we don't have to worry about floating point
		amt, wasNegative := strings.CutPrefix(row[amountIdx], "-")
		if wasNegative {
			row[amountIdx] = amt
		} else {
			row[amountIdx] = "-" + amt
		}

		if err := dst.Write(row); err != nil {
			return fmt.Errorf("could not write next row: %w", err)
		}
	}

	dst.Flush()
	return dst.Error()
}
