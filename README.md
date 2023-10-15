# etf-cli

Simple ETF cli application that scrapes some etf data.

## TODO

- Add pretty etf/history display. Add metrics (e.g, dif last 5 days, dif last 30 days, etc.)
- Add symbol validation to ensure it exists
- Tidy up error handling
- Add cli commands:
  - `etf track "VWCE.DE"` - done
  - `etf list`
  - `etf get` - done
  - `etf untrack "VWCE.DE"`