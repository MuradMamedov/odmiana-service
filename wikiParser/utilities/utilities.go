package utilities

import (
	"log"
	"reflect"
	"runtime"
	"time"
)

func StopWatch[Results interface{}](l *log.Logger, f func() (Results, error), message string) (Results, error) {
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	if message != "" {
		funcName = message
	}
	start := time.Now()
	value, err := f()
	end := time.Now()
	l.Printf("%s - Elapsed Time: %v\n", funcName, end.Sub(start))
	return value, err
}

func StopWatchParametrized[Results interface{}, Param interface{}](l *log.Logger, f func(p Param) (Results, error), param Param, message string) (Results, error) {
	wrappedFunc := func() (Results, error) {
		return f(param)
	}

	return StopWatch(l, wrappedFunc, "")
}
