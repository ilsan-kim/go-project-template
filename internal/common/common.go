package common

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

var GracefulShutdownJob chan bool

func RegisterSignal(stopFunc func()) (done chan bool) {
	done = make(chan bool)
	go func() {

		c := make(chan os.Signal, 1)
		signal.Notify(c,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		s := <-c
	B:
		for {
			switch s {
			case syscall.SIGHUP:
				fmt.Println("I got a SIGHUP signal")
			case syscall.SIGINT:
				fmt.Println("I got a SIGINT signal")
				break B
			case syscall.SIGTERM:
				fmt.Println("I got a SIGTERM signal")
				break B
			case syscall.SIGQUIT:
				fmt.Println("I got a SIGQUIT signal")
			default:
				fmt.Println("I got a signal", s)
			}
			time.Sleep(time.Second)
		}
		stopFunc()
		done <- true
	}()
	return done
}

func RunFuncSafely(f func()) {
	GracefulShutdownJob <- true
	defer func() {
		r := recover()
		if r != nil {
			log.Println(string(debug.Stack()))
		}
		<-GracefulShutdownJob
	}()

	f()
}
