package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"wechat-back/internals/database"
	"wechat-back/internals/logger"
	"wechat-back/internals/models"
	"wechat-back/internals/tools"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WebsocketHUB central entrypoint for websocket server
var WebsocketHUB *WebsocketPanel

// WebsocketPanel central pannel that holds information about the websocket server
type WebsocketPanel struct {
	// P2PConnections holds all the p2p connections
	P2PConnections map[string]P2PConnectionCredentials

	// GroupConnections holds all the group connections
	GroupConnections map[string]GroupConnectionCredentials

	// DBConn access for database
	DBConn database.DBHUB

	// mux mutext
	mux sync.Mutex
}

// P2PConnectionCredentials holds necessary information about the user connection to a peer
type P2PConnectionCredentials struct {
	// Conn websocket connection
	Conn            *websocket.Conn
	AuthorID        string
	TargetID        string
	AuthorData      *models.User
	TargetPushToken string
}

// GroupConnectionCredentials holds necessary information about the user connection to a group
type GroupConnectionCredentials struct {
	Conn       *websocket.Conn
	AuthorID   string
	TargetID   string
	AuthorData *models.User
	TargetData *models.Group
}

// StartWebsocketService starts websocket service
func StartWebsocketService() {
	WebsocketHUB = &WebsocketPanel{
		mux:              sync.Mutex{},
		P2PConnections:   make(map[string]P2PConnectionCredentials),
		GroupConnections: make(map[string]GroupConnectionCredentials),
	}
}

// StopWebsocketService stops the WebSocket server and deletes all connections
func StopWebsocketService() {
	alog := logger.StartLogger()
	if WebsocketHUB == nil {
		alog.WarningLogger("WebsocketHUB is nil, no connections to stop.")
		return
	}

	WebsocketHUB.mux.Lock()
	defer WebsocketHUB.mux.Unlock()

	for key, conn := range WebsocketHUB.P2PConnections {
		err := conn.Conn.Close()
		if err != nil {
			alog.ErrorLog(fmt.Sprintf("Error closing P2P connection %s: %v", key, err))
		}
		delete(WebsocketHUB.P2PConnections, key)
	}

	for key, conn := range WebsocketHUB.GroupConnections {
		err := conn.Conn.Close()
		if err != nil {
			alog.WarningLogger(fmt.Sprintf("Error closing Group connection %s: %v", key, err))
		}
		delete(WebsocketHUB.GroupConnections, key)
	}

	alog.InfoLogger("All WebSocket connections closed and cleaned up.")
}

func ListenForP2PActivity(c P2PConnectionCredentials) {

	alog := logger.StartLogger()
	var payload models.InboundP2PMessage
	defer c.Conn.Close()

	for {

		msgType, data, err := c.Conn.ReadMessage()
		if err != nil {
			c.Conn.Close()
			delete(WebsocketHUB.P2PConnections, c.AuthorID)
			break
		}

		switch msgType {
		case websocket.TextMessage:

			err = json.Unmarshal(data, &payload)
			if err != nil {
				alog.ErrorLog(err.Error())
				msg := models.OutboundP2PTextMessage{}
				tools.WriteWebsocketJSON(c.Conn, msg.FormatErrorOutboundMessage(BAD_FIELD, err))
			}

			c.BroadcastText(payload, &c)

		case websocket.BinaryMessage:
			alog.InfoLogger("got a binary message")

		}

	}

}

func ListenForGroupActivity(c GroupConnectionCredentials) {

	alog := logger.StartLogger()
	var payload models.InboundGroupTextMessage
	defer c.Conn.Close()

	for {

		msgType, data, err := c.Conn.ReadMessage()
		if err != nil {
			c.Conn.Close()
			delete(WebsocketHUB.GroupConnections, c.AuthorID)
			break
		}

		switch msgType {
		case websocket.TextMessage:

			err = json.Unmarshal(data, &payload)
			if err != nil {
				alog.ErrorLog(err.Error())
				msg := models.OutboundGroupTextMessage{}
				tools.WriteWebsocketJSON(c.Conn, msg.FormatErrorOutboundTextMessage(BAD_FIELD, err))
			}

			// broadcase
			c.Broadcast(payload)

		case websocket.BinaryMessage:
			alog.InfoLogger("group got a binary message")

		}

	}

}

/*
	BROADCAST METHODS
*/

func (p *P2PConnectionCredentials) BroadcastText(msg models.InboundP2PMessage, c *P2PConnectionCredentials) {
	alog := logger.StartLogger()

	var payload models.P2PTextChatLog
	var outbound models.OutboundP2PTextMessage

	tarID, err := primitive.ObjectIDFromHex(p.TargetID)
	if err != nil {
		alog.ErrorLog(err.Error())

		tools.WriteWebsocketJSON(c.Conn, outbound.FormatErrorOutboundMessage(BAD_FIELD, err))
	}

	payload.FormatTextLog(tarID, c.AuthorData.ID, c.AuthorData.Name, msg.Body)

	go func() {
		_, err = WebsocketHUB.DBConn.InsertP2PMessageDB(payload)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(c.Conn, outbound.FormatErrorOutboundMessage(DB_ERROR, err))

		}
	}()
	outbound.FormatOutboundMessage(payload.AuthorID.Hex(), payload.Body, payload.ID.Hex())

	user, ok := WebsocketHUB.P2PConnections[p.TargetID]
	p.Conn.WriteJSON(outbound)
	if ok {
		// send an update in-app state for target and author
		WebsocketHUB.mux.Lock()
		defer WebsocketHUB.mux.Unlock()

		user.Conn.WriteJSON(outbound)
	}
	// send push notification to target user
	// send push notification to user that is not connected to the group
	alog.InfoLogger(fmt.Sprintf("about to send push notification to user %s", p.TargetPushToken))
}

func (g *GroupConnectionCredentials) Broadcast(msg models.InboundGroupTextMessage) {

	alog := logger.StartLogger()

	var payload models.GroupChatTextLog
	var res models.OutboundGroupTextMessage

	payload.FormatTextChatLog(g.TargetData.ID, g.AuthorData.ID, g.AuthorData.Name, msg.Body)

	go func() {
		_, err := WebsocketHUB.DBConn.InsertGroupMessageDB(payload)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(g.Conn, res.FormatErrorOutboundTextMessage(BAD_FIELD, err))
		}
	}()

	res.FormatOutboundTextMessage(payload.AuthorID.Hex(), payload.Body, payload.ID.Hex(), OK)

	for _, usr := range g.TargetData.Participants {

		user, ok := WebsocketHUB.GroupConnections[usr.Hex()]
		if ok {
			defer WebsocketHUB.mux.Unlock()
			WebsocketHUB.mux.Lock()
			user.Conn.WriteJSON(res)

		} else {
			// send push notification to user if not connected
			u, ok := msg.PushTokens[usr.Hex()]
			if !ok {
				alog.ErrorLog(fmt.Sprintf("user %v not exist as participant", u))
				tools.WriteWebsocketJSON(g.Conn, res.FormatErrorOutboundTextMessage(BAD_FIELD, errors.New("user is not a participant of this group")))
			}
			alog.InfoLogger(fmt.Sprintf("about to send push notification to user %s", u))
		}

	}

}
