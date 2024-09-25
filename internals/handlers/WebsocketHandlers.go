package handlers

import (
	"net/http"
	"wechat-back/internals/database"
	"wechat-back/internals/logger"
	"wechat-back/internals/server"
	"wechat-back/internals/tools"

	"github.com/gorilla/websocket"
)

/*
HandleP2PConnectionEP
Handles the p2p connection and adds the connection to the pool of users p2p
*/
func HandleP2PConnectionEP(w http.ResponseWriter, r *http.Request, db database.DBHUB) {

	alog := logger.StartLogger()

	var upgrader = websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteWebsocketJSON(conn, tools.FormatErrResponse(server.SERVER_ERROR, err))
		return
	}

	// get user and user

	u, exist, err := db.FindUserDB(r.URL.Query().Get("email"))
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteWebsocketJSON(conn, tools.FormatErrResponse(server.BAD_REQUEST, err))
		conn.Close()
		return
	} else if !exist {
		alog.ErrorLog("forbidden")
		tools.WriteWebsocketJSON(conn, tools.FormatCustomErrResponse("forbidden", server.NOT_ALLOWED))
		return
	}

	payload := server.P2PConnectionCredentials{
		Conn:            conn,
		AuthorID:        u.ID.Hex(),
		TargetID:        r.URL.Query().Get("tar"),
		AuthorData:      &u,
		TargetPushToken: r.URL.Query().Get("tar_tkn"),
	}

	server.WebsocketHUB.P2PConnections[payload.AuthorID] = payload
	server.WebsocketHUB.DBConn = db

	go server.ListenForP2PActivity(payload)

	tools.WriteJSON(w, http.StatusOK, nil)

}

/*
HandleGroupConnectionsEP
Handles the group connection and adds the connection to the group pool
*/
func HandleGroupConnectionsEP(w http.ResponseWriter, r *http.Request, db database.DBHUB) {

	alog := logger.StartLogger()

	var upgrader = websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteWebsocketJSON(conn, tools.FormatErrResponse(server.SERVER_ERROR, err))
	}

	groupID := r.URL.Query().Get("gi")
	uid := r.URL.Query().Get("ui")

	author, exist, err := db.FindUserDB(uid)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteWebsocketJSON(conn, tools.FormatErrResponse(server.BAD_REQUEST, err))
	} else if !exist {
		alog.ErrorLog("forbidden")
		tools.WriteWebsocketJSON(conn, tools.FormatCustomErrResponse("forbidden", server.NOT_ALLOWED))
	}

	group, err := db.GetGroupDB(groupID)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteWebsocketJSON(conn, tools.FormatErrResponse(server.DB_ERROR, err))
	}

	payload := server.GroupConnectionCredentials{
		Conn:       conn,
		AuthorID:   author.ID.Hex(),
		TargetID:   group.ID.Hex(),
		AuthorData: &author,
		TargetData: group,
	}
	server.WebsocketHUB.GroupConnections[author.ID.Hex()] = payload
	server.WebsocketHUB.DBConn = db

	go server.ListenForGroupActivity(payload)
	tools.WriteJSON(w, http.StatusOK, nil)

}
