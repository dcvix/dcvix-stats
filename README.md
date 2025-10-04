dcvix Stats
===========

Dcvix Stats is a Go application that provides a graphical user interface (GUI) to display statistics from a NICE DCV server log file. It uses the Fyne toolkit for the GUI and the go-charts library to render line charts of various metrics.

The application works by parsing a DCV server log file, extracting statistical data using regular expressions, and then displaying this data in a series of line charts. The user can show or hide them from the "Show" menu.

![screenshot](assets/screenshot.png)

### Command-line Flags

The Dcvix Stats accepts the following command-line flags:

*   `--version`: Show version information.
*   `--verbose`: Enable verbose logging.
*   `--entries`: How many entries/minutes to evaluate (default 120).
*   `--logfile`: Path to the DCV server log file.


## Building

Preferred building environment Linux

### Requirements

*   Go 1.23 or later
*   Fyne library dependencies installed 

### To be able to cross compile

go install github.com/fyne-io/fyne-cross@latest

### Build the application

The project uses a `Makefile` to simplify common tasks:
```bash
make build
```
This will create an executable file in the `bin` directory.

