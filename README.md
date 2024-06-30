# ynab-discover

Converts CSV exported from Discover.com into a format suitable for importing
into YNAB. This project is not maintained or endorsed by YNAB.

## Usage

1. export transactions from the Discover website as a CSV
2. run `go run . -src from-discover.csv -dst to-ynab.csv`
    - renames the "Trans. Date" column to "Date"
    - flips the sign on the "Amount" column
3. import the new file to YNAB

## See also

- <https://support.ynab.com/en_us/file-based-import-a-guide-Bkj4Sszyo>
- <https://support.ynab.com/en_us/formatting-a-csv-file-an-overview-BJvczkuRq>