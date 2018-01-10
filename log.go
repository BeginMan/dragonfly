package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	LDebug = iota
	LInfo
	LWarn
	LError
	LFatal
)

type Logger struct {
	FileName      string
	File          *os.File
	Level         uint
	Flag          string
	PrintToStdout bool
}

func newLogger(flag, path string, printToStd bool) *Logger {
	logger := Logger{FileName: path, PrintToStdout: printToStd, Flag: flag}
	var err error

	if path != "" {
		logger.File, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		checkErr(err)
	}
	return &logger
}

func (l Logger) Debug(args ...interface{}) {
	if l.Level <= LDebug {
		l.Log("DEBUG", args...)
	}
}

func (l Logger) Info(args ...interface{}) {
	if l.Level <= LInfo {
		l.Log("Info", args...)
	}
}
func (l Logger) Warn(args ...interface{}) {
	if l.Level <= LWarn {
		l.Log("Warn", args...)
	}
}
func (l Logger) Error(args ...interface{}) {
	if l.Level <= LError {
		l.Log("Error", args...)
	}
}
func (l Logger) Fatal(args ...interface{}) {
	if l.Level <= LFatal {
		l.Log("Fatal", args...)
	}
}

func (l Logger) Log(level string, args ...interface{}) {
	tfmt := "2006-01-02 15:04:05"
	prefix := time.Now().Format(tfmt)
	msg := l.Flag + " | " + level + " | " + prefix + " | " + fmtLogMsg(args...)

	if l.PrintToStdout {
		fmt.Println(msg)
	}

	if l.File != nil {
		l.File.WriteString(msg + "\n")
		l.File.Sync()
	}
}

func fmtLogMsg(args ...interface{}) string {
	var argsStr = []string{}
	for _, v := range args {
		if _, ok := v.(error); ok {
			argsStr = append(argsStr, fmt.Sprintf("%v", v))
			continue
		}
		argsStr = append(argsStr, fmt.Sprintf("%#v", v))
	}
	return strings.Join(argsStr, ", ")
}
