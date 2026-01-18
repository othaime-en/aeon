# aeon

A minimalistic TUI time zone converter for the terminal.

[![Commits since latest](https://img.shields.io/github/commits-since/othaime-en/aeon/latest)](https://github.com/othaime-en/aeon/commits/latest)

## Features

- **World Clock**: Live updating clocks for multiple time zones
- **Quick Conversion**: Convert times between zones (`3pm NYC to Berlin`)
- **Meeting Finder**: Find overlapping business hours across zones
- **Smart Resolution**: Recognizes 15,000+ cities and common abbreviations

## Installation

```bash
git clone https://github.com/othaime-en/aeon.git
cd aeon
go build -o aeon
./aeon
```

## Usage

```bash
./aeon
```

### Navigation

- `←/→` or `Tab` - Switch views
- `1/2/3` - Jump to Clock/Convert/Meeting
- `c` - Convert view
- `m` - Meeting view
- `Enter` - Start input
- `Esc` or `q` - Cancel/Quit

### Examples

**Convert times:**

```
3pm NYC to Berlin
9:30am Tokyo to London
15:00 UTC to PST
tomorrow 3pm NYC to Berlin
in 2 hours Tokyo to London
next monday noon San Francisco to Hong Kong
2026-01-20 3pm LA to NYC
noon NYC to Berlin
```

**Find meeting slots:**

```
NYC, London, Tokyo
San Francisco, Berlin, Singapore
```

## Requirements

- Go 1.24+
- Terminal with color support

## License

MIT License - see [LICENSE](LICENSE)

## Credits

Timezone data from [GeoNames](https://www.geonames.org/) (CC BY 4.0)
