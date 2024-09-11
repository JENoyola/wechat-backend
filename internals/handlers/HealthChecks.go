package handlers

import (
	"fmt"
	"net/http"
	"os"
	"wechat-back/internals/logger"
	"wechat-back/internals/tools"
)

func ServerHealthCheckEP(w http.ResponseWriter, r *http.Request) {

	logger.StartLogger().InfoLogger(fmt.Sprintf("Server v%s is running and is healthy running on PORT %s on %s ENV\n", os.Getenv("API_VERSION"), os.Getenv("PORT"), os.Getenv(("ENV"))))

	tools.WriteJSON(w, http.StatusOK, nil)

}
