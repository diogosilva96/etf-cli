# etf-cli

Simple ETF cli that displays ETF data & allows management of tracked etfs.
The ETF data is scraped from [Yahoo Finance](https://finance.yahoo.com/).

Available commands:

- `etf-cli report` - generates a report for the ETFs in the configuration.
- `etf-cli add <etf-symbol>` - adds the specified ETF symbol to the configuration. The symbol must exist in Yahoo finance, such as [VWCE.DE](https://finance.yahoo.com/quote/VWCE.DE/)
- `etf-cli remove <etf-symbol>` - removes the specified ETF symbol from the configuration.
- `etf-cli list` - lists the ETF symbols currently in the configuration.
- `etf-cli help` - help command.

## Installation

Steps:

1. Download the source code.
2. In the root folder run `go build .`
3. Then run `go install .`
4. You can then run any command using `etf-cli [command]`

## TODO

- Tidy up error handling
- Improve command error messages
- Add functional options pattern for client.go ?
- Remove ETF command - allow removal with an index flag or by symbol?
- Publish as pkg for easy distribution? (To investigate)
