# aeon

A minimalistic terminal time zone converter with natural language support.

[![Commits since latest](https://img.shields.io/github/commits-since/othaime-en/aeon/latest)](https://github.com/othaime-en/aeon/commits/latest)

## Features

- **World Clock**: Live multi-timezone display with persistent custom zones
- **Smart Conversion**: Natural language time parsing with relative and absolute dates
- **Meeting Finder**: Identify overlapping business hours across zones
- **Intelligent Resolution**: Recognizes 15,000+ cities, common abbreviations, and IANA timezones

## Installation

```bash
git clone https://github.com/othaime-en/aeon.git
cd aeon
go build -o aeon
./aeon
```

## Usage

Launch the application:

```bash
./aeon
```

### Navigation

- `Tab` / `←/→` - Switch between views
- `1/2/3` - Jump to Clock/Convert/Meeting view
- `a` - Add timezone (Clock view)
- `d` - Delete selected zone (Clock view)
- `↑/↓` - Navigate zones (Clock view)
- `Enter` - Start input (Convert/Meeting views)
- `Esc` - Cancel input
- `q` - Quit

### Time Conversion Examples

The Convert view supports flexible time expressions:

**Relative times:**
```
tomorrow 3pm NYC to Berlin
in 2 hours Tokyo to London
next monday noon San Francisco to Hong Kong
```

**Absolute dates:**
```
2026-01-20 3pm LA to NYC
Jan 20 3pm NYC to Berlin
1/20 3pm NYC to Berlin
```

**Natural language:**
```
noon NYC to Berlin
midnight Tokyo to Sydney
now UTC to PST
```

**Traditional formats:**
```
3pm NYC to Berlin
9:30am Tokyo to London
15:00 UTC to PST
```

### Meeting Slots

Find overlapping business hours across multiple timezones:

```
NYC, London, Tokyo
San Francisco, Berlin, Singapore
```

## Timezone Resolution

Supports multiple input formats:

- **Cities**: New York, Tokyo, London
- **Abbreviations**: NYC, LA, SF, HK
- **IANA timezones**: America/New_York, Asia/Tokyo
- **Common aliases**: EST, PST, UTC, GMT

## Configuration

Zones added in the Clock view persist automatically in `~/.aeon.yaml`.

## Requirements

- Go 1.24+
- Terminal with color support

## License

MIT License - see [LICENSE](LICENSE)

## Credits

Timezone data from [GeoNames](https://www.geonames.org/) (CC BY 4.0)