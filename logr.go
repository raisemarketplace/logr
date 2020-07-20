package logr

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

var (
	messages chan *Message

	logr    Logger = &Logr{}
	pool           = &sync.Pool{}
	writers        = make(map[*WriterConfig]io.Writer)
)

// SetMeta sets the global meta data attached to every log message
func SetMeta(data map[string]interface{}) {
	logr = logr.With(data)
}

// SetBufferSize updates the message queue buffer size. Default size is 10,000 (10K)
func SetBufferSize(size int) {
	if messages != nil {
		close(messages)
		Wait()
	}
	msgs := make(chan *Message, size)
	go listen(msgs)
	messages = msgs
}

// Wait for log messages to be processed
func Wait() {
	for {
		time.Sleep(time.Millisecond)
		if messages == nil || len(messages) == 0 {
			break
		}
	}
}

// Logger defines the methods available both by the logr package and Logr containing additional meta data.
type Logger interface {
	Panic(v ...interface{})
	Panicf(msg string, v ...interface{})
	Error(v ...interface{}) string
	Errorf(msg string, v ...interface{}) string
	Warn(v ...interface{}) string
	Warnf(msg string, v ...interface{}) string
	Info(v ...interface{}) string
	Infof(msg string, v ...interface{}) string
	Debug(v ...interface{}) string
	Debugf(msg string, v ...interface{}) string
	Success(v ...interface{}) string
	Successf(msg string, v ...interface{}) string
	With(data map[string]interface{}) Logger
}

// Logr implements the Logger interface
type Logr struct {
	meta map[string]interface{}
}

// Panic logs inputs as panics and panics
func (l *Logr) Panic(v ...interface{}) {
	code := log(P, true, Interfaces(v).SSV(), l.meta)
	panic(code)
}

// Panicf logs a formatted message as a panic and panics
func (l *Logr) Panicf(msg string, v ...interface{}) {
	code := logf(P, true, msg, v, l.meta)
	panic(code)
}

// Error logs inputs as errors
func (l *Logr) Error(v ...interface{}) string {
	return log(E, false, Interfaces(v).SSV(), l.meta)
}

// Errorf logs a formatted message as an error
func (l *Logr) Errorf(msg string, v ...interface{}) string {
	return logf(E, false, msg, v, l.meta)
}

// Warn logs inputs as warnings
func (l *Logr) Warn(v ...interface{}) string {
	return log(W, false, Interfaces(v).SSV(), l.meta)
}

// Warnf logs a formatted message as a warning
func (l *Logr) Warnf(msg string, v ...interface{}) string {
	return logf(W, false, msg, v, l.meta)
}

// Info logs inputs as info messages
func (l *Logr) Info(v ...interface{}) string {
	return log(I, false, Interfaces(v).SSV(), l.meta)
}

// Infof logs a formatted message as an info message
func (l *Logr) Infof(msg string, v ...interface{}) string {
	return logf(I, false, msg, v, l.meta)
}

// Debug logs inputs as debug messages
func (l *Logr) Debug(v ...interface{}) string {
	return log(D, false, Interfaces(v).SSV(), l.meta)
}

// Debugf logs a formatted message as a debug message
func (l *Logr) Debugf(msg string, v ...interface{}) string {
	return logf(D, false, msg, v, l.meta)
}

// Success logs inputs as success messages
func (l *Logr) Success(v ...interface{}) string {
	return log(S, false, Interfaces(v).SSV(), l.meta)
}

// Successf logs a formatted message as a success message
func (l *Logr) Successf(msg string, v ...interface{}) string {
	return logf(S, false, msg, v, l.meta)
}

// With metadata in the log messages
func (l *Logr) With(data map[string]interface{}) Logger {
	meta := make(map[string]interface{})
	for k, v := range l.meta {
		meta[k] = v
	}
	for k, v := range data {
		meta[k] = v
	}
	return &Logr{
		meta: meta,
	}
}

// format a msg and log as given type
func logf(t Type, wait bool, msg string, args []interface{}, meta map[string]interface{}) string {
	return log(t, wait, fmt.Sprintf(msg, args...), meta)
}

// log inputs to given type
func log(t Type, wait bool, msg string, meta map[string]interface{}) string {
	now := time.Now()
	m := pool.Get().(*Message)
	m.Type = t
	m.Time = now.Format("Jan 02 2006 15:04:05.9999")
	m.Code = strconv.FormatInt(now.UnixNano(), 36)
	m.Desc = msg
	m.Meta = meta
	m.done = make(chan struct{})
	messages <- m

	if wait {
		<-m.done
	}

	return m.Code
}

// Panic logs inputs as panics and panics
func Panic(v ...interface{}) {
	logr.Panic(v...)
}

// Panicf logs a formatted message as a panic and panics
func Panicf(msg string, v ...interface{}) {
	logr.Panicf(msg, v...)
}

// Error logs inputs as errors
func Error(v ...interface{}) string {
	return logr.Error(v...)
}

// Errorf logs a formatted message as an error
func Errorf(msg string, v ...interface{}) string {
	return logr.Errorf(msg, v...)
}

// Warn logs inputs as warnings
func Warn(v ...interface{}) string {
	return logr.Warn(v...)
}

// Warnf logs a formatted message as a warning
func Warnf(msg string, v ...interface{}) string {
	return logr.Warnf(msg, v...)
}

// Info logs inputs as info messages
func Info(v ...interface{}) string {
	return logr.Info(v...)
}

// Infof logs a formatted message as an info message
func Infof(msg string, v ...interface{}) string {
	return logr.Infof(msg, v...)
}

// Debug logs inputs as debug messages
func Debug(v ...interface{}) string {
	return logr.Debug(v...)
}

// Debugf logs a formatted message as a debug message
func Debugf(msg string, v ...interface{}) string {
	return logr.Debugf(msg, v...)
}

// Success logs inputs as success messages
func Success(v ...interface{}) string {
	return logr.Success(v...)
}

// Successf logs a formatted message as a success message
func Successf(msg string, v ...interface{}) string {
	return logr.Successf(msg, v...)
}

// With metadata in the log messages
func With(data map[string]interface{}) Logger {
	return logr.With(data)
}
