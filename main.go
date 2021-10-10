package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ddpub/request-info/internal/server"
)

func main() {

	addr := flag.String("srv.addr", "127.0.0.1:7777", "server address")
	flag.Parse()

	tctx, shtdwn := sigterm()

	httpSrv := server.NewSrv(*addr)
	hctx := httpSrv.Start(tctx)

	<-hctx.Done()
	shtdwn()
	<-hctx.Done()

	log.Println("Terminated")
}

func sigterm() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("interrupt signal has been caught")
		cancel()
	}()
	return ctx, cancel
}
