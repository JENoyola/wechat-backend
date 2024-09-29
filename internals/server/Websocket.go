package server

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"wechat-back/internals/database"
	"wechat-back/internals/logger"
	"wechat-back/internals/models"
	"wechat-back/internals/tools"
	"wechat-back/providers/media"

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

	// MediaProvider access for media provider
	MediaProvider media.MediaHUB

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
			var payload models.InboundP2PTextMessage

			err = json.Unmarshal(data, &payload)
			if err != nil {
				alog.ErrorLog(err.Error())
				tools.WriteWebsocketJSON(c.Conn, models.FormatWebsocketErrResponse(err, BAD_REQUEST))
			}

			c.HandleP2PTextContent(payload)

		case websocket.BinaryMessage:
			var payload models.InboundP2PContentMessage
			cont, err := tools.ReadBinaryWebsocketMessage(data, &payload)
			if err != nil {
				alog.ErrorLog(err.Error())
				tools.WriteWebsocketJSON(c.Conn, models.FormatWebsocketErrResponse(err, BAD_REQUEST))
			}

			c.HandleP2PMediaContent(payload, cont)

		}

	}

}

func ListenForGroupActivity(c GroupConnectionCredentials) {

	alog := logger.StartLogger()
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
			var payload models.InboundGroupTextMessage

			err = json.Unmarshal(data, &payload)
			if err != nil {
				alog.ErrorLog(err.Error())
				// tools.WriteWebsocketJSON(c.Conn, msg.FormatErrorOutboundTextMessage(BAD_FIELD, err))
				continue
			}

			// broadcase
			c.HandleGroupTextContent(payload)

		case websocket.BinaryMessage:

			var payload models.InboundGroupContentMessage
			cont, err := tools.ReadBinaryWebsocketMessage(data, &payload)
			if err != nil {
				alog.ErrorLog(err.Error())
				msg := models.OutboundP2PTextMessage{}
				tools.WriteWebsocketJSON(c.Conn, msg.FormatErrorOutboundMessage(BAD_FIELD, err))
				continue
			}

			alog.InfoLogger(fmt.Sprint("body ------>", payload))

			c.HandleGroupMediaContent(payload, cont)

		}

	}
}

func (p *P2PConnectionCredentials) HandleP2PTextContent(msg models.InboundP2PTextMessage) {
	alog := logger.StartLogger()

	var payload models.P2PTextChatLog

	tarID, err := primitive.ObjectIDFromHex(p.TargetID)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteWebsocketJSON(p.Conn, models.FormatWebsocketErrResponse(err, BAD_FIELD))
	}

	payload.FormatTextLog(tarID, p.AuthorData.ID, p.AuthorData.Name, msg.Body)

	_, err = WebsocketHUB.DBConn.InsertP2PMessageDB(payload)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteWebsocketJSON(p.Conn, models.FormatWebsocketErrResponse(err, DB_ERROR))
	}

	p.BroadcastToP2P(payload)
}

func (p *P2PConnectionCredentials) HandleP2PMediaContent(msg models.InboundP2PContentMessage, binaryContent [][]byte) {
	alog := logger.StartLogger()

	var payload models.P2PContentChatLog

	tarID, err := primitive.ObjectIDFromHex(p.TargetID)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteWebsocketJSON(p.Conn, models.FormatWebsocketErrResponse(err, BAD_FIELD))
	}

	switch msg.ContentType {

	case models.MESSAGE_TYPE_MEDIA_VIDEOS:

		videoPlay, err := WebsocketHUB.MediaProvider.StoreVideo(fmt.Sprintf("%s*%d", p.TargetID, time.Now().Unix()), msg.Filename[0], binaryContent[0])
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(p.Conn, models.FormatWebsocketErrResponse(err, PROVIDER_ERROR))
		}

		payload.FormatContentChatLog(tarID, p.AuthorData.ID, p.AuthorData.Name, msg.Body, videoPlay.GUID, []string{videoPlay.Src}, []string{videoPlay.Thumbnail}, models.MESSAGE_TYPE_MEDIA_VIDEOS)

		_, err = WebsocketHUB.DBConn.InsertP2PMessageDB(payload)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(p.Conn, models.FormatWebsocketErrResponse(err, DB_ERROR))
		}

	case models.MESSAGE_TYPE_MEDIA_IMAGES:
		ImageInfo, err := WebsocketHUB.MediaProvider.InsetImages(binaryContent, msg.Filename)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(p.Conn, models.FormatWebsocketErrResponse(err, PROVIDER_ERROR))
		}

		payload.FormatContentChatLog(tarID, p.AuthorData.ID, p.AuthorData.Name, msg.Body, ImageInfo.ContentID, ImageInfo.MediaSource, ImageInfo.Thumbnails, models.MESSAGE_TYPE_MEDIA_IMAGES)

		_, err = WebsocketHUB.DBConn.InsertP2PMessageDB(payload)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(p.Conn, models.FormatWebsocketErrResponse(err, DB_ERROR))
		}

	case models.MESSAGE_TYPE_FILE:

		fileURL, err := WebsocketHUB.MediaProvider.InsertFile(binaryContent[0], msg.Filename[0])
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(p.Conn, models.FormatWebsocketErrResponse(err, PROVIDER_ERROR))
		}
		payload.FormatContentChatLog(tarID, p.AuthorData.ID, p.AuthorData.Name, msg.Body, "N/A", []string{fileURL}, []string{""}, models.MESSAGE_TYPE_FILE)

		_, err = WebsocketHUB.DBConn.InsertP2PMessageDB(payload)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(p.Conn, models.FormatWebsocketErrResponse(err, DB_ERROR))
		}
	}

	p.BroadcastToP2P(payload)

}

func (g *GroupConnectionCredentials) HandleGroupTextContent(msg models.InboundGroupTextMessage) {

	alog := logger.StartLogger()

	var payload models.GroupChatTextLog

	payload.FormatTextChatLog(g.TargetData.ID, g.AuthorData.ID, g.AuthorData.Name, msg.Body)

	go func() {
		_, err := WebsocketHUB.DBConn.InsertGroupMessageDB(payload)
		if err != nil {
			alog.ErrorLog(err.Error())
			// tools.WriteWebsocketJSON(g.Conn, res.FormatErrorOutboundTextMessage(BAD_FIELD, err))
		}
	}()

	g.BroadcastToParticipants(g.TargetData.Participants, msg.PushTokens, payload)

}

func (g *GroupConnectionCredentials) HandleGroupMediaContent(msg models.InboundGroupContentMessage, binaryContent [][]byte) {

	alog := logger.StartLogger()

	var payload models.GroupChatContentLog

	switch msg.ContentType {

	case models.MESSAGE_TYPE_MEDIA_VIDEOS:

		videoPlay, err := WebsocketHUB.MediaProvider.StoreVideo(fmt.Sprintf("%s*%d", g.TargetData.GroupID, time.Now().Unix()), msg.Filename[0], binaryContent[0])
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(g.Conn, models.FormatWebsocketErrResponse(err, PROVIDER_ERROR))
		}

		payload.FormatContentChatLog(g.TargetData.ID, g.AuthorData.ID, g.AuthorData.Name, msg.Body, videoPlay.GUID, []string{videoPlay.Src}, []string{videoPlay.Thumbnail}, models.MESSAGE_TYPE_MEDIA_VIDEOS)

		_, err = WebsocketHUB.DBConn.InsertGroupMessageDB(payload)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(g.Conn, models.FormatWebsocketErrResponse(err, DB_ERROR))
		}

	case models.MESSAGE_TYPE_MEDIA_IMAGES:
		ImageInfo, err := WebsocketHUB.MediaProvider.InsetImages(binaryContent, msg.Filename)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(g.Conn, models.FormatWebsocketErrResponse(err, PROVIDER_ERROR))
		}

		payload.FormatContentChatLog(g.TargetData.ID, g.AuthorData.ID, g.AuthorData.Name, msg.Body, ImageInfo.ContentID, ImageInfo.MediaSource, ImageInfo.Thumbnails, models.MESSAGE_TYPE_MEDIA_IMAGES)

		_, err = WebsocketHUB.DBConn.InsertGroupMessageDB(payload)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(g.Conn, models.FormatWebsocketErrResponse(err, DB_ERROR))

		}

	case models.MESSAGE_TYPE_FILE:

		fileURL, err := WebsocketHUB.MediaProvider.InsertFile(binaryContent[0], msg.Filename[0])
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(g.Conn, models.FormatWebsocketErrResponse(err, PROVIDER_ERROR))

		}

		payload.FormatContentChatLog(g.TargetData.ID, g.AuthorData.ID, g.AuthorData.Name, msg.Body, "N/A", []string{fileURL}, []string{}, models.MESSAGE_TYPE_FILE)

		_, err = WebsocketHUB.DBConn.InsertGroupMessageDB(payload)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteWebsocketJSON(g.Conn, models.FormatWebsocketErrResponse(err, DB_ERROR))

		}

	}

	g.BroadcastToParticipants(g.TargetData.Participants, msg.PushTokens, payload)

}

// BROADCASTING FUNCTIONS
func (g *GroupConnectionCredentials) BroadcastToParticipants(participants []primitive.ObjectID, tokens map[string]string, payload any) {

	alog := logger.StartLogger()

	for _, usr := range g.TargetData.Participants {

		user, ok := WebsocketHUB.GroupConnections[usr.Hex()]
		if ok {
			user.Conn.WriteJSON(payload)
		} else {
			u, ok := tokens[usr.Hex()]
			if !ok {
				alog.ErrorLog(fmt.Sprintf("User %v is not a participant of group %v", usr.Hex(), g.TargetID))
			}
			alog.InfoLogger(fmt.Sprintf("about to send push notification to user %s", u))
		}

	}
}

func (p *P2PConnectionCredentials) BroadcastToP2P(payload any) {

	WebsocketHUB.mux.Lock()
	defer WebsocketHUB.mux.Unlock()

	p.Conn.WriteJSON(payload)

	user, ok := WebsocketHUB.P2PConnections[p.TargetID]
	if ok {

		user.Conn.WriteJSON(payload)
	}
	logger.StartLogger().InfoLogger(fmt.Sprintf("about to send push notification to user %s", p.TargetPushToken))
}
