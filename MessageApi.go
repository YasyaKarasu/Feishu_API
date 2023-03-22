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

const (
	Text        MsgContentType = "text"
	Interactive MsgContentType = "interactive"
)

// Send a message to a person / chat group, return whether if it had been send successfully
func (c AppClient) MessageSend(receiveIdType MsgReceiverType, receiveId string, msgType MsgContentType, msg string) (string, bool) {
	query := make(map[string]string)
	query["receive_id_type"] = string(receiveIdType)

	content := ""

	switch msgType {
	case Text:
		contmap := make(map[string]string)
		contmap["text"] = msg
		bytecont, err := json.Marshal(contmap)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"receive_id_type": string(receiveIdType),
				"receive_id":      receiveId,
				"msg_type":        string(msgType),
				"msg":             msg,
			}).Error("marshal text to json fail")
			return "", false
		}
		content = string(bytecont)
	case Interactive:
		content = msg
	default:
		logrus.WithField("msgType", msgType).Error("message type unsupported")
		return "", false
	}
	body := make(map[string]string)
	body["receive_id"] = receiveId
	body["content"] = content
	body["msg_type"] = string(msgType)

	resp := c.Request("post", "open-apis/im/v1/messages", query, nil, body)

	if resp == nil {
		logrus.WithFields(logrus.Fields{
			"ReceiveID": receiveId,
			"Msg":       msg,
		}).Error("message send error")
		return "", false
	}

	return resp["message_id"].(string), true
}

func (c AppClient) UpdateMessage(mid string, content string) {
	query := make(map[string]string, 0)
	query["content"] = content

	c.Request("patch", "open-apis/im/v1/messages/"+mid, query, nil, nil)
}
