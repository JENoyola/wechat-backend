package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(cfg ServeConfig) error {

	srv := http.Server{
		Addr:              fmt.Sprintf(":%v", cfg.PORT),
		Handler:           cfg.HANDLER,
		ReadTimeout:       cfg.READTIMEOUT * time.Minute,
		ReadHeaderTimeout: cfg.READHEADER * time.Minute,
		WriteTimeout:      cfg.WRITE * time.Minute,
		IdleTimeout:       cfg.IDLE * time.Minute,
		MaxHeaderBytes:    2048,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGSYS)

	if cfg.ENV == "PROD" || cfg.ENV == "DIST" {

		go func() {
			fmt.Printf("API v%v STARTING AT PORT %v on ENV %v\n", cfg.API_VERSION, cfg.PORT, cfg.ENV)
			err := srv.ListenAndServeTLS(cfg.TLSC, cfg.TLSK)
			if err != nil {
				return
			}
		}()
	} else {
		go func() {
			fmt.Printf("API v%v STARTING AT PORT %v on ENV %v\n", cfg.API_VERSION, cfg.PORT, cfg.ENV)
			err := srv.ListenAndServe()
			if err != nil {
				return
			}
		}()
	}

	<-stop
	// log.WarningLogger("Initializing Gracefully Server Stop Protocol")

	/*if cfg.ENV == "PROD" || cfg.ENV == "DIST" {
		// Send email to developer
	}*/

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		// log.Fatal(err.Error())
		return err
	}

	fmt.Println("server stopped gracefully")
	return nil

}
