package feishuapi

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type MsgReceiverType string

const (
	UserOpenId  MsgReceiverType = "open_id"
	UserUnionId MsgReceiverType = "union_id"
	UserUserId  MsgReceiverType = "user_id"
	UserEmail   MsgReceiverType = "email"
	GroupChatId MsgReceiverType = "chat_id"
)

type MsgContentType string

const Text MsgContentType = "text"

// Send a message to a person / chat group, return whether if it had been send successfully
func (c AppClient) Send(receive_id_type MsgReceiverType, receive_id string, msg_type MsgContentType, msg string) bool {
	query := make(map[string]string)
	query["receive_id_type"] = string(receive_id_type)

	content := ""

	if msg_type == Text {
		contmap := make(map[string]string)
		contmap["text"] = msg
		bytecont, err := json.Marshal(contmap)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"receive_id_type": string(receive_id_type),
				"receive_id":      receive_id,
				"msg_type":        string(msg_type),
				"msg":             msg,
			}).Error("marshal text to json fail")
			return false
		}
		content = string(bytecont)
	}
	body := make(map[string]string)
	body["receive_id"] = receive_id
	body["content"] = content
	body["msg_type"] = string(msg_type)

	resp := c.Request("post", "open-apis/im/v1/messages", query, nil, body)
	return resp != nil
}
