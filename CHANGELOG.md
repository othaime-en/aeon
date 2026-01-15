# Changelog

All notable changes to aeon will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

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
