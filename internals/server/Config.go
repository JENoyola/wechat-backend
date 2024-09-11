package server

import (
	"net/http"
	"time"
)

type ServeConfig struct {
	PORT        string
	HANDLER     http.Handler
	IDLE        time.Duration
	WRITE       time.Duration
	READHEADER  time.Duration
	READTIMEOUT time.Duration
	TLSC        string
	TLSK        string
	ENV         string
	API_VERSION string
}
