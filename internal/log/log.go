package log

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	initOnce    sync.Once
	initialized = atomic.Bool{}
)

func Initialize(logFile string, debug bool) {
	logRotator := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	logger := slog.NewJSONHandler(logRotator, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})
	slog.SetDefault(slog.New(logger))
	initialized.Store(true)
}

func Setup(logFile string, debug bool) {
	initOnce.Do(func() {
		Initialize(logFile, debug)
	})
}

func Initialized() bool {
	return initialized.Load()
}

func RecoverPanic(functionName string, onError func(err error)) {
	r := recover()
	if r == nil {
		return
	}

	timestamp := time.Now().Format("20060102-150405")

	filename := fmt.Sprintf("%s_%s.log", timestamp, functionName)

	file, err := os.Create(filename)
	if err != nil {
		slog.Error("Error creating log file", "error", err)
		return
	}

	defer file.Close()

	fmt.Fprintf(file, "Panic in %s: %v\n\n", functionName, r)
	fmt.Fprintf(file, "Timestamp: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "Stack trace:\n%s\n", string(debug.Stack()))

	if onError != nil {
		onError(r.(error))
		return
	}

	slog.Error(functionName, "error", r)
}
