package server

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

type HttpSrv struct {
	addr string
}

func NewSrv(addr string) *HttpSrv {
	return &HttpSrv{
		addr: addr,
	}
}

func (s *HttpSrv) Start(ctx context.Context) context.Context {
	ownctx, shtdwn := context.WithCancel(context.Background())
	rtr := mux.NewRouter()

	rtr.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		dmp, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println("dump request failed:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(dmp)

		if err != nil {
			log.Println("write failed:", err)
		}
	})

	httpSrv := &http.Server{
		Addr:    s.addr,
		Handler: rtr,
	}

	go func() {
		<-ctx.Done()
		if err := httpSrv.Shutdown(context.Background()); err != nil {
			log.Println("HTTP service shutdown failure:", err)
		}
		shtdwn()
	}()

	go func() {
		log.Println("starting HTTP server on:", s.addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("HTTP server serve failure:", err)
		}
		log.Println("HTTP server terminated")
		shtdwn()
	}()

	return ownctx
}
