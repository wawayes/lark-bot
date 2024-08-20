package domain

import (
	"encoding/json"
)

type MessageType string

const (
	MsgTypeText        MessageType = "text"
	MsgTypePost        MessageType = "post"
	MsgTypeImage       MessageType = "image"
	MsgTypeFile        MessageType = "file"
	MsgTypeAudio       MessageType = "audio"
	MsgTypeMedia       MessageType = "media"
	MsgTypeSticker     MessageType = "sticker"
	MsgTypeInteractive MessageType = "interactive"
	MsgTypeShareChat   MessageType = "share_chat"
	MsgTypeShareUser   MessageType = "share_user"
	MsgTypeLocation    MessageType = "location"
)

type ChatType string

const (
	P2PChat   ChatType = "p2p"
	GroupChat ChatType = "group"
)

type Message struct {
	MessageID  string
	RootID     string
	ParentID   string
	MsgType    MessageType
	ChatID     string
	ChatType   ChatType
	Content    string
	Sender     Sender
	CreateTime string
	UpdateTime string
	ThreadID   string
	Mentions   []Mention
	UserAgent  string
}

type Mention struct {
	Key       string
	ID        IDObject
	Name      string
	TenantKey string
}

type Sender struct {
	SenderID   IDObject
	SenderType string
	TenantKey  string
}

type IDObject struct {
	UnionID string
	UserID  string
	OpenID  string
}

type Reply struct {
	ReceiveIDType string
	ReceiveID     string
	Content       string
	MsgType       MessageType
}

func (m *Message) UnmarshalContent() (map[string]interface{}, error) {
	var content map[string]interface{}
	err := json.Unmarshal([]byte(m.Content), &content)
	return content, err
}

func (m *Message) MarshalContent(content map[string]interface{}) error {
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}
	m.Content = string(jsonContent)
	return nil
}

func (m *Message) IsMentioned() (bool, error) {
	return len(m.Mentions) > 0, nil
}
