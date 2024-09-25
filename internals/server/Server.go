package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wechat-back/internals/logger"
)

func StartServer(cfg ServeConfig) error {

	alog := logger.StartLogger()

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

	StartWebsocketService()

	if cfg.ENV == "PROD" || cfg.ENV == "DIST" {

		go func() {
			alog.InfoLogger(fmt.Sprintf("API v%v STARTING AT PORT %v on ENV %v\n", cfg.API_VERSION, cfg.PORT, cfg.ENV))
			err := srv.ListenAndServeTLS(cfg.TLSC, cfg.TLSK)
			if err != nil {
				alog.ErrorLog(err.Error())
				return
			}
		}()
	} else {
		go func() {
			alog.InfoLogger(fmt.Sprintf("API v%v STARTING AT PORT %v on ENV %v\n", cfg.API_VERSION, cfg.PORT, cfg.ENV))
			err := srv.ListenAndServe()
			if err != nil {
				alog.ErrorLog(err.Error())
				return
			}
		}()
	}

	<-stop
	alog.WarningLogger("Initializing Gracefully Server Stop Protocol")
	StopWebsocketService()

	// notify developers about server gracefully shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		alog.ErrorLog(err.Error())
		return err
	}

	alog.WarningLogger("server stopped gracefully")
	return nil

}
