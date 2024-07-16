package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thehxdev/txtban/txtban"
)

const VERSION string = "1.0.0"

func main() {
	configPath := flag.String("c", "./config.toml", "Path to config file")
	showVersion := flag.Bool("v", false, "Show version info")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Txtban v%s\nhttps://github.com/thehxdev/txtban\n", VERSION)
		return
	}

	tb := txtban.Init(*configPath)
	defer tb.Conn.DB.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer close(sigChan)
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
