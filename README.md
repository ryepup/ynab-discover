# ynab-discover

Converts CSV exported from Discover.com™ into a format suitable for importing
into YNAB™. Renames the "Trans. Date" column to "Date" and flips the sign on 
all amounts (Discover™ shows expenses as positive, YNAB™ expects negative). 

> [!IMPORTANT]
> This project is **not affiliated with, maintained by, or endorsed by** YNAB™ or Discover Financial Services™.

## Usage

### As a CLI tool

1. export transactions from the Discover™ website as a CSV
2. run `go run ./cmd/ynab-discover -src from-discover.csv -dst to-ynab.csv`
3. import the new file to YNAB™

### As a library

```go
import "github.com/ryepup/ynab-discover"

// Convert a Discover CSV to YNAB format
err := discover.ConvertCSV(ctx, reader, writer)
```

## See also

- <https://support.ynab.com/en_us/file-based-import-a-guide-Bkj4Sszyo>
- <https://support.ynab.com/en_us/formatting-a-csv-file-an-overview-BJvczkuRq>