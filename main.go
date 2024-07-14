package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thehxdev/txtban/txtban"
)

func main() {
	configPath := flag.String("c", "./config.toml", "Path to config file")
	flag.Parse()

	tb := txtban.Init(*configPath)
	defer tb.Conn.DB.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := tb.Run()
		if err != nil {
			tb.ErrLogger.Println(err)
		}
	}()

	_ = <-sigChan
	tb.InfLogger.Println("Shutting down the server...")
	err := tb.App.ShutdownWithTimeout(time.Second * 10)
	if err != nil {
		tb.ErrLogger.Println(err)
	}
}
