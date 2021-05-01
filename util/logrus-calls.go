package util

import (
	log "github.com/sirupsen/logrus"
)

type LogrusCall struct {
	Entry *log.Entry
	Args  []interface{}
}

func MakeLogrusCall(entry *log.Entry, args ...interface{}) LogrusCall {
	return LogrusCall{
		Entry: entry,
		Args:  args,
	}
}

type LogrusCalls struct {
	Info, Warn, Trace, Panic, Error, Fatal []LogrusCall
}

func NewLogrusCalls() *LogrusCalls {
	return &LogrusCalls{
		Info:  make([]LogrusCall, 0),
		Warn:  make([]LogrusCall, 0),
		Trace: make([]LogrusCall, 0),
		Panic: make([]LogrusCall, 0),
		Error: make([]LogrusCall, 0),
		Fatal: make([]LogrusCall, 0),
	}
}

func (l *LogrusCalls) Call() {
	for _, c := range l.Info {
		c.Entry.Info(c.Args)
	}
	for _, c := range l.Warn {
		c.Entry.Warn(c.Args)
	}
	for _, c := range l.Trace {
		c.Entry.Trace(c.Args)
	}
	for _, c := range l.Panic {
		c.Entry.Panic(c.Args)
	}
	for _, c := range l.Error {
		c.Entry.Error(c.Args)
	}
	for _, c := range l.Fatal {
		c.Entry.Fatal(c.Args)
	}
}
