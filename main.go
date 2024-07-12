package main

import (
	"flag"
	"github.com/thehxdev/txtban/txtban"
)

func main() {
	configPath := flag.String("c", "", "Path to config file")
	flag.Parse()

	tb := txtban.Init(*configPath)
	defer tb.Conn.DB.Close()

	tb.ErrLogger.Fatal(tb.Run())
}
