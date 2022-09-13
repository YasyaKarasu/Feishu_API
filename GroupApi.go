package feishuapi

type GroupInfo struct {
	ChatId    string
	Name      string
	TenantKey string
}

// Create a new GroupInfo
func NewGroupInfo(data map[string]interface{}) *GroupInfo {
	return &GroupInfo{
		ChatId:    data["chat_id"].(string),
		Name:      data["name"].(string),
		TenantKey: data["tenant_key"].(string),
	}
}

// Get All the chat group that the feishu robot is in
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

// Create a new GroupMember
func NewGroupMember(data map[string]interface{}) *GroupMember {
	return &GroupMember{
		OpenId: data["member_id"].(string),
		Name:   data["name"].(string),
	}
}

// Get all the group members in a specific group
func (c AppClient) GetGroupMembers(groupId string, userIdType UserIdType) []GroupMember {
	body := make(map[string]string)
	body["member_id_type"] = string(userIdType)

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
