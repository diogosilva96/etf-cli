# etf-cli

Simple ETF cli application that scrapes some etf data.

## TODO

- Add pretty etf/history display. Add metrics (e.g, dif last 5 days, dif last 30 days, etc.)
- Add symbol validation to ensure it exists
- Add cli commands:
  - `etf track "VWCE.DE"`
  - `etf list`
  - `etf fetch`
  - `etf untrack "VWCE.DE"`
- Run ETF scraping with go routines in parallel.