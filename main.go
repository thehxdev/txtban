package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thehxdev/txtban/txtban"
)

func main() {
	configPath := flag.String("c", "./config.toml", "Path to config file")
	showVersion := flag.Bool("v", false, "Show version info")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Txtban v%s\nhttps://github.com/thehxdev/txtban\n", txtban.VERSION)
		return
	}

	tb := txtban.Init(*configPath)
	defer tb.CloseDB()

	serverCtx, serverCtxStop := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		<-sigChan
		tb.InfLogger.Println("Shutting down the server...")
		shutdownCtx, shutdownCancelFunc := context.WithTimeout(serverCtx, 10*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				tb.ErrLogger.Fatal("server shutdown timeout... force exit")
			}
			serverCtxStop()
		}()

		err := tb.Server.Shutdown(shutdownCtx)
		if err != nil {
			tb.ErrLogger.Fatal(err)
		}

		shutdownCancelFunc()
	}()

	err := tb.Run()
	if err != nil && err != http.ErrServerClosed {
		tb.ErrLogger.Fatal(err)
	}

	<-serverCtx.Done()
}
