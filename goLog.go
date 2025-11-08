package goLog

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	Reset	= "\033[0m"
	Red		= "\033[31m"
	Green	= "\033[32m"
	Cyan	= "\033[36m"
	Yellow	= "\033[33m"
	Blue	= "\033[94m"
)

type Logger struct {
	fileLogger		*log.Logger
	screenLogger	*log.Logger
	startTime		time.Time
}

// New creates a logger that writes to both stdout and a file (when present)
func New(logFilePath string) (*Logger, error) {
	var fileLogger *log.Logger

	if logFilePath != "" {
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		fileLogger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Always create a screen logger
	screenLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	return &Logger{
		fileLogger:   fileLogger,
		screenLogger: screenLogger,
	}, nil
}


func (l *Logger) Fatal(logFn func(level, msg string), format string, v ...any) {
    msg := fmt.Sprintf(format, v...)
    l.ToBoth("FATAL", msg)
    os.Exit(1)
}

// ToBoth logs to both file and screen
func (l *Logger) ToBoth(level string, msg string) {
	l.ToFile(level, msg)
	l.ToScreen(level, msg)
}

// ToFile logs a message to the file with level-based formatting
func (l *Logger) ToFile(level string, msg string) {
	if l.fileLogger != nil {
		formatted := l.format(level, msg, false)
		l.fileLogger.Println(formatted)
	}
}

// ToScreen logs a message to stdout with color-coded level formatting
func (l *Logger) ToScreen(level string, msg string) {
	formatted := l.format(level, msg, true)
	l.screenLogger.Println(strings.TrimSpace(formatted))
}

func (l *Logger) ResetTimer() {
	l.startTime = time.Time{} // zero value
}

func (l *Logger) Debug(v ...any) {
	l.ToScreen("DEBUG", spew.Sdump(v...))
}

// format applies level prefix and optional color
func (l *Logger) format(level string, msg string, colorize bool) string {
	level = strings.ToUpper(level)
	var prefix string

	switch level {
	case "INFO": // GREEN
		prefix = "  [INFO]"
		if colorize {
			prefix = Cyan + prefix + Reset
		}
	case "WARN": // YELLOW
		prefix = "  [WARN]"
		if colorize {
			prefix = Yellow + prefix + Reset
		}
	case "ERROR": // RED
		prefix = " [ERROR]"
		if colorize {
			prefix = Red + prefix + Reset
		}
	case "FATAL": // RED
		prefix = " [FATAL]"
		if colorize {
			prefix = Red + prefix + Reset
		}
	case "DEBUG": // CYAN
		prefix = " [DEBUG]"
		if colorize {
			prefix = Blue + prefix + Reset
		}
	case "TIME": // CYAN
		prefix = "[TIMING]"
		if colorize {
			prefix = Green + prefix + Reset
		}		// handle timing
		if l.startTime.IsZero() {
			// first TIME call sets the baseline
			l.startTime = time.Now()
			msg = fmt.Sprintf("%s (start)", msg)
		} else {
			elapsed := time.Since(l.startTime)
			msg = fmt.Sprintf("%s (+%s)", msg, elapsed.Truncate(time.Millisecond))
		}
	default:
		prefix = "   [LOG]"
	}

	return fmt.Sprintf("%s %s", prefix, msg)
}