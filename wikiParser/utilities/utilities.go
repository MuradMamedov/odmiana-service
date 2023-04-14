package utilities

import (
	"log"
	"reflect"
	"runtime"
	"time"
)

func StopWatch(l *log.Logger, f func() (interface{}, error), message string) (interface{}, error) {
	if message == "" {
		message = "Elapsed time:"
	}
	start := time.Now()
	value, err := f()
	end := time.Now()
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	l.Printf("%s - %s %v\n", funcName, message, end.Sub(start))
	return value, err
}
