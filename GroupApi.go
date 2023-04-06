package feishuapi

import "github.com/sirupsen/logrus"

type GroupInfo struct {
	ChatId    string
	Name      string
	TenantKey string
}

// Create a new GroupInfo
func NewGroupInfo(data map[string]any) *GroupInfo {
	return &GroupInfo{
		ChatId:    data["chat_id"].(string),
		Name:      data["name"].(string),
		TenantKey: data["tenant_key"].(string),
	}
}

// Get All the chat group that the feishu robot is in
func (c AppClient) GroupGetAllInfo() []GroupInfo {
	l := c.GetAllPages("get", "open-apis/im/v1/chats", nil, nil, nil, 100)
	if l == nil {
		logrus.Warn("nil group info return")
		return nil
	}

	var all_groups []GroupInfo
	for _, value := range l {
		all_groups = append(all_groups, *NewGroupInfo(value.(map[string]any)))
	}

	return all_groups
}

type GroupMember struct {
	MemberId string
	Name     string
}

// Create a new GroupMember
func NewGroupMember(data map[string]any) *GroupMember {
	return &GroupMember{
		MemberId: data["member_id"].(string),
		Name:     data["name"].(string),
	}
}

// Get all the group members in a specific group
func (c AppClient) GroupGetMembers(groupId string, userIdType UserIdType) []GroupMember {
	body := make(map[string]string)
	body["member_id_type"] = string(userIdType)

	l := c.GetAllPages("get", "open-apis/im/v1/chats/"+groupId+"/members", nil, nil, body, 100)
	if l == nil {
		logrus.WithField("GroupID", groupId).Warn("nil group member info return")
		return nil
	}

	var all_members []GroupMember
	for _, value := range l {
		all_members = append(all_members, *NewGroupMember(value.(map[string]any)))
	}

	return all_members
}

// CreateGroup Create a new group
func (c AppClient) GroupCreate(groupName string, userIdType UserIdType, ownerId string) *GroupInfo {
	query := make(map[string]any)
	query["user_id_type"] = string(userIdType)
	body := make(map[string]string)
	body["name"] = groupName
	body["owner_id"] = ownerId

	info := c.Request("post", "open-apis/im/v1/chats", query, nil, body)

	if info == nil {
		logrus.WithFields(logrus.Fields{
			"GroupName": groupName,
			"OwnerID":   ownerId,
		}).Error("create group fail")
		return nil
	}

	return NewGroupInfo(info)
}

// GetGroupInfo Get a group information
func (c AppClient) GroupGetInfo(chatId string) *GroupInfo {
	info := c.Request("get", "open-apis/im/v1/chats/"+chatId+"?user_id_type=open_id", nil, nil, nil)

	if info == nil {
		logrus.WithField("ChatID", chatId).Warn("nil group info return")
		return nil
	}

	info["chat_id"] = chatId
	info["tenant_key"] = ""
	return NewGroupInfo(info)
}

// AddMembers
// app_id to add bot
func (c AppClient) GroupAddMembers(chatId string, memberIdType UserIdType, succeedType string, idList []string) bool {
	query := make(map[string]any)
	query["member_id_type"] = string(memberIdType)
	query["succeed_type"] = succeedType

	body := make(map[string][]string)

	var result bool = true
	var idlist []string

	for len(idList) > 50 {
		idlist = idList[0:50]
		body["id_list"] = idlist
		resp := c.Request("post", "open-apis/im/v1/chats/"+chatId+"/members", query, nil, body)
		idList = append(idList[:0], idList[50:]...)
		if resp == nil {
			logrus.WithFields(logrus.Fields{
				"ChatID": chatId,
				"IdList": idList,
			}).Error("add member fail")
			result = false
		}
	}
	body["id_list"] = idList
	resp := c.Request("post", "open-apis/im/v1/chats/"+chatId+"/members", query, nil, body)
	if resp == nil {
		logrus.WithFields(logrus.Fields{
			"ChatID": chatId,
			"IdList": idList,
		}).Error("add member fail")
		result = false
	}
	return result
}

// DeleteMembers
// app_id to delete bot
func (c AppClient) GroupDeleteMembers(chatId string, memberIdType UserIdType, idList []string) bool {
	query := make(map[string]any)
	query["member_id_type"] = string(memberIdType)

	body := make(map[string][]string)

	var result bool = true
	var idlist []string

	for len(idList) > 50 {
		idlist = idList[0:50]
		body["id_list"] = idlist
		resp := c.Request("delete", "open-apis/im/v1/chats/"+chatId+"/members", query, nil, body)
		idList = append(idList[:0], idList[50:]...)
		if resp == nil {
			logrus.WithFields(logrus.Fields{
				"ChatID": chatId,
				"IdList": idList,
			}).Error("delete member fail")
			result = false
		}
	}
	body["id_list"] = idList
	resp := c.Request("delete", "open-apis/im/v1/chats/"+chatId+"/members", query, nil, body)
	if resp == nil {
		logrus.WithFields(logrus.Fields{
			"ChatID": chatId,
			"IdList": idList,
		}).Error("delete member fail")
		result = false
	}
	return result
}

func (c AppClient) GroupChangeOwner(chatId string, memberIdType UserIdType, ownerId string) {
	query := make(map[string]any)
	query["user_id_type"] = string(memberIdType)

	body := make(map[string]string)
	body["owner_id"] = ownerId

	resp := c.Request("put", "open-apis/im/v1/chats/"+chatId, query, nil, body)
	if resp == nil {
		logrus.WithFields(logrus.Fields{
			"ChatID":     chatId,
			"NewOwnerId": ownerId,
		}).Warn("change group owner fail")
	}
}
