package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer c.Stop()
	DBInit()
	SyncTasks()

	chSig := make(chan os.Signal,1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	_ = <-chSig
}
