package feishuapi

type GroupInfo struct {
	ChatId    string
	Name      string
	TenantKey string
}

func NewGroupInfo(data map[string]interface{}) *GroupInfo {
	return &GroupInfo{
		ChatId:    data["chat_id"].(string),
		Name:      data["name"].(string),
		TenantKey: data["tenant_key"].(string),
	}
}

func (c AppClient) GetGroupList() []GroupInfo {
	l := c.GetAllPages("get", "open-apis/im/v1/chats", nil, nil, nil, 100)
	if l == nil {
		return nil
	}

	var all_groups []GroupInfo
	for _, value := range l {
		all_groups = append(all_groups, *NewGroupInfo(value.(map[string]interface{})))
	}
	return all_groups
}

type GroupMember struct {
	OpenId string
	Name   string
}

func NewGroupMember(data map[string]interface{}) *GroupMember {
	return &GroupMember{
		OpenId: data["member_id"].(string),
		Name:   data["name"].(string),
	}
}

func (c AppClient) GetGroupMembers(groupId string) []GroupMember {
	body := make(map[string]string)
	body["member_id_type"] = "open_id"

	l := c.GetAllPages("get", "open-apis/im/v1/chats/"+groupId+"/members", nil, nil, body, 100)
	if l == nil {
		return nil
	}

	var all_members []GroupMember
	for _, value := range l {
		all_members = append(all_members, *NewGroupMember(value.(map[string]interface{})))
	}
	return all_members
}
