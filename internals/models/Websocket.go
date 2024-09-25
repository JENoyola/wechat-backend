package models

const (
	// MESSAGE_TYPE_TEXT Type of message that is text based
	MESSAGE_TYPE_TEXT = 8
	// MESSAGE_TYPE_FILE Type of message is a file
	MESSAGE_TYPE_FILE = 22
	// MESSAGE_TYPE_MEDIA_VIDEOS type of message that is an array of video content
	MESSAGE_TYPE_MEDIA_VIDEOS = 59
	// MESSAGE_TYPE_MEDIA_IMAGES Type of message that is an array of image content
	MESSAGE_TYPE_MEDIA_IMAGES = 65
)

// InboundP2PMessage base structure to websocket message model peer 2 peer
type InboundP2PMessage struct {
	MessageID       string `json:"message_id"`
	AuthorID        string `json:"author_id"`
	TargetID        string `json:"target_id"`
	TargetPushToken string `json:"target_pushtoken"`
	Body            string `json:"body"`
}

// InboundGroupTextMessage base structure to websicket message model peer to group
type InboundGroupTextMessage struct {
	MessageID  string            `json:"message_id"`
	AuthorID   string            `json:"author_id"`
	GroupID    string            `json:"group"`
	PushTokens map[string]string `json:"push_tokens"`
	Body       string            `json:"body"`
}

// OutboundP2PTextMessage base structure to websocket outbound message model for peer to peer
type OutboundP2PTextMessage struct {
	Error       bool   `json:"error"`
	Message     string `json:"message"`
	Code        int    `json:"code"`
	MessageType int    `json:"message_type"`
	MessageID   string `json:"message_id"`
	AuthorID    string `json:"author_id"`
	Body        string `json:"body"`
}

func (o *OutboundP2PTextMessage) FormatOutboundMessage(author, body, msgID string) {
	o.MessageType = MESSAGE_TYPE_TEXT
	o.MessageID = msgID
	o.AuthorID = author
	o.Body = body
}

func (o *OutboundP2PTextMessage) FormatErrorOutboundMessage(code int, err error) *OutboundP2PTextMessage {
	o.Error = true
	o.Message = err.Error()
	o.Code = code
	o.MessageType = MESSAGE_TYPE_TEXT
	o.MessageID = ""
	o.AuthorID = ""
	o.Body = ""
	return o
}

// OutboundGroupTextMessage base structure to websocket outbound message model for peer to peer
type OutboundGroupTextMessage struct {
	Error       bool   `json:"error"`
	Message     string `json:"message"`
	Code        int    `json:"code"`
	MessageType int    `json:"message_type"`
	MessageID   string `json:"message_id"`
	AuthorID    string `json:"author_id"`
	Body        string `json:"body"`
}

// FormatOutboundTextMessage formats outboundtext message for client delivery
func (o *OutboundGroupTextMessage) FormatOutboundTextMessage(author, body, msgID string, code int) {
	o.Error = false
	o.Message = ""
	o.Code = code
	o.MessageType = MESSAGE_STANDARD_TYPE
	o.MessageID = msgID
	o.AuthorID = author
	o.Body = body
}

func (o *OutboundGroupTextMessage) FormatErrorOutboundTextMessage(code int, err error) *OutboundGroupTextMessage {
	o.Error = false
	o.Message = err.Error()
	o.Code = code
	o.MessageType = MESSAGE_STANDARD_TYPE
	o.MessageID = ""
	o.AuthorID = ""
	o.Body = ""
	return o
}
