package logger

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/kick-project/kick/internal/resources/exit"
)

// TODO: Increase unit test coverage from 40% to 80%

type Level int

const (
	InvalidLevel Level = iota - 1
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// LogIface
type LogIface interface {
	LoggerIface
	LogLevelIface
}

// OutputIface router interface
type OutputIface interface {
	Debug(string)
	Debugf(string, ...interface{})
	Error(string)
	Errorf(string, ...interface{})
	Output(int, string) error
	Printf(string, ...interface{})
}

// LoggerIface A interface that matches log.Logger
type LoggerIface interface {
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
	Flags() int
	Output(int, string) error
	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
	Prefix() string
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
	SetFlags(int)
	SetOutput(io.Writer)
	SetPrefix(string)
	Writer() io.Writer
}

// LogLevelIface functions run when a level is set
type LogLevelIface interface {
	Error(string)
	Errorf(string, ...interface{})
	Debug(string)
	Debugf(string, ...interface{})
}

// Log logging interface
type Log struct {
	LogIface
	lvl       Level
	stdLogger *log.Logger
	eh        exit.HandlerIface
	mu        *sync.Mutex
}

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line, or
// after the log header if the Lmsgprefix flag is provided.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int, level Level, eh exit.HandlerIface) *Log {
	return &Log{
		lvl:       level,
		stdLogger: log.New(out, prefix, flag),
		eh:        eh,
		mu:        &sync.Mutex{},
	}
}

//
// Standard logging functions
//

// Fatal is equivalent to l.Print() followed by a call to exit.Handler.Exit(1).
func (l *Log) Fatal(v ...interface{}) {
	_ = l.Output(2, fmt.Sprint(v...))
	l.eh.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to exit.Handler.Exit(1).
func (l *Log) Fatalf(format string, v ...interface{}) {
	_ = l.Output(2, fmt.Sprintf(format, v...))
	l.eh.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to exit.Handler.Exit(1).
func (l *Log) Fatalln(v ...interface{}) {
	_ = l.Output(2, fmt.Sprintln(v...))
	l.eh.Exit(1)
}

// Flags returns the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *Log) Flags() int {
	return l.stdLogger.Flags()
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Log) Output(calldepth int, s string) error {
	return l.stdLogger.Output(calldepth+1, s)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Log) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	_ = l.Output(2, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Log) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	_ = l.Output(2, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Log) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	_ = l.Output(2, s)
	panic(s)
}

// Prefix returns the output prefix for the logger.
func (l *Log) Prefix() string { return l.stdLogger.Prefix() }

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Log) Print(v ...interface{}) { _ = l.Output(2, fmt.Sprint(v...)) }

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Log) Printf(format string, v ...interface{}) { _ = l.Output(2, fmt.Sprintf(format, v...)) }

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Log) Println(v ...interface{}) { _ = l.Output(2, fmt.Sprintln(v...)) }

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *Log) SetFlags(flag int) {
	l.stdLogger.SetFlags(flag)
}

// SetOutput sets the output destination for the logger.
func (l *Log) SetOutput(w io.Writer) {
	l.stdLogger.SetOutput(w)
}

// SetPrefix sets the output prefix for the logger.
func (l *Log) SetPrefix(prefix string) {
	l.stdLogger.SetPrefix(prefix)
}

// Writer returns the output destination for the logger.
func (l *Log) Writer() io.Writer {
	return l.stdLogger.Writer()
}

//
// Custom Level functions
//

func (l *Log) levelOutput(level Level, s string) {
	oldPrefix := l.Prefix()
	l.mu.Lock()
	defer func() {
		l.SetPrefix(oldPrefix)
		l.mu.Unlock()
	}()
	if l.lvl > level {
		return
	}
	prefix := l.levelText(level)
	l.SetPrefix(prefix)

	_ = l.Output(3, s)
}

func (l *Log) levelText(level Level) (newPrefix string) {
	switch level {
	case DebugLevel:
		newPrefix = "DEBUG "
	case InfoLevel:
		newPrefix = "INFO "
	case WarnLevel:
		newPrefix = "WARN "
	case ErrorLevel:
		newPrefix = "ERROR "
	case FatalLevel:
		newPrefix = "FATAL "
	default:
		panic("Unknown debug level")
	}
	return
}

// Error error level message.
func (l *Log) Error(s string) {
	l.levelOutput(ErrorLevel, s)
}

// Errorf error level message.
func (l *Log) Errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.levelOutput(ErrorLevel, s)
}

// Debug debug level message.
func (l *Log) Debug(s string) {
	l.levelOutput(DebugLevel, s)
}

// Debugf debug level message.
func (l *Log) Debugf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.levelOutput(DebugLevel, s)
}
