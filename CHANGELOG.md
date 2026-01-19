# Changelog

All notable changes to aeon will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [0.2.0] - 2026-01-19

### Added

- Dynamic timezone management in Clock view
  - Add zones with `a` key (supports city names, abbreviations, IANA timezones)
  - Delete zones with `d` key
  - Navigate zones with arrow keys
  - Zones persist across sessions via config file
- Enhanced time parsing with natural language support
  - Natural language: `now`, `noon`, `midnight`
  - Relative times: `in 2 hours`, `tomorrow 3pm`, `next monday noon`
  - Date support: `2026-01-20 3pm`, `jan 20 3pm`, `1/20 3pm`
  - Multi-word timezone handling: `tomorrow 3pm New York to Los Angeles`
- Comprehensive test suite for time parsing (30+ test cases)

### Changed

- Improved Convert view with enhanced help text and examples
- Better error messages for invalid time expressions

### Technical

- New `parser.go` module with modular parsing functions
- Smart zone resolution for multi-word expressions
- Zone persistence via YAML configuration

[0.2.0]: https://github.com/othaime-en/aeon/releases/tag/v0.2.0

## [0.1.0] - 2026-01-15

### Added

- Multi-zone world clock with live updates
- Time conversion between zones with natural language input
- Basic meeting slot finder showing business hours
- Smart timezone resolution supporting 15,000+ cities
- Common city abbreviations (NYC, LA, SF, etc.)
- Keyboard-driven navigation with arrow keys and shortcuts
- Three-view interface: Clock, Convert, Meeting

### Technical

- Built with Go 1.24 and Bubble Tea TUI framework
- Integrated GeoNames city database
- Offline-first architecture using IANA tzdata

[0.1.0]: https://github.com/othaime-en/aeon/releases/tag/v0.1.0
