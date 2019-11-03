package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/quickfixgo/quickfix"
	"os"
	"os/signal"
	"path"
	"quickFix/app"
	app2 "quickFix/internal"
	"quickFix/internal/adapter"
)

func main() {
	flag.Parse()

	cfgFileName := path.Join("config", "quickFix.cfg")
	if flag.NArg() > 0 {
		cfgFileName = flag.Arg(0)
	}

	cfg, err := os.Open(cfgFileName)
	if err != nil {
		fmt.Printf("Error opening %v, %v\n", cfgFileName, err)
		return
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		fmt.Println("Error reading cfg,", err)
		return
	}

	logFactory := quickfix.NewScreenLogFactory()

	lmax := adapter.NewLmaxAdapter()
	application := app.NewApplication(lmax)

	acceptor, err := quickfix.NewAcceptor(application, quickfix.NewMemoryStoreFactory(), appSettings, logFactory)
	if err != nil {
		fmt.Printf("Unable to create Acceptor: %s\n", err)
		return
	}

	err = acceptor.Start()
	if err != nil {
		fmt.Printf("Unable to start Acceptor: %s\n", err)
		return
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	go func() {
		<-interrupt
		acceptor.Stop()
		os.Exit(0)
	}()
}
