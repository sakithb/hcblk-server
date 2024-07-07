package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer(h http.Handler) {
	s := http.Server{
		Addr: ":8080",
		Handler: h,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go WaitForInterrupt(c, &s)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func WaitForInterrupt(c chan os.Signal, s *http.Server) {
	<-c

	ctx, cancelCtx := context.WithTimeout(context.Background(), 30 * time.Second)
	err := s.Shutdown(ctx)

	if err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}

	cancelCtx()
}
