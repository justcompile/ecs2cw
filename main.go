package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justcompile/ecs2cw/lib"
)

var pollInterval = flag.Duration("interval", time.Duration(time.Second*60), "Polling interval")
var configPath = flag.String("config", "", "path to configuration file")
var namespace = flag.String("namespace", "", "namespace into which CloudWatch metrics will be published")

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	flag.Parse()

	if *configPath == "" {
		log.Fatal("You must specify --config flag")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
		done <- true
	}()

	cfg, err := lib.NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	dispatcher := lib.NewDispatcher(&lib.DispatchOptions{
		Config:    cfg,
		Ctx:       ctx,
		Interval:  *pollInterval,
		Namespace: *namespace,
	})

	if err := dispatcher.Poll(); err != nil {
		done <- true
		cancel()
		log.Fatalf("[ERROR] %s\n", err)
	}

	<-done
}
