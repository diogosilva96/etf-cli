# etf-cli

Simple ETF cli that displays ETF data & allows management of tracked etfs.
The ETF data is scraped from [Yahoo Finance](https://finance.yahoo.com/).

Available commands:

- `report` - generates a report for the ETFs in the configuration
- `add <etf-symbol>` - adds the specified ETF symbol to the configuration. The symbol must exist in Yahoo finance, such as [VWCE.DE](https://finance.yahoo.com/quote/VWCE.DE/)
- `remove <index>` - removes the ETF with the specified index from the configuration
- `list` - lists the ETF symbols currently in the configuration
- `help` - help command

## Installation

Steps:

1. Download the source code.
2. Navigate to the root project folder.
3. Run `go build .`
4. Run `go install .`
5. You can then run any command using `etf-cli [command]`

## TODO
- Add tests
- Add html output option for the reports with graphs
- Investigate how to publish cli by using for example homebrew, goreleaser
