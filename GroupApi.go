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
	MemberId string
	Name     string
}

// Create a new GroupMember
func NewGroupMember(data map[string]interface{}) *GroupMember {
	return &GroupMember{
		MemberId: data["member_id"].(string),
		Name:     data["name"].(string),
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

// CreateGroup Create a new group
func (c AppClient) CreateGroup(groupName string, user_id_type string, owner_id string) *GroupInfo {
	query := make(map[string]string)
	query["user_id_type"] = string(user_id_type)
	body := make(map[string]string)
	body["name"] = groupName
	body["owner_id"] = owner_id

	info := c.Request("post", "open-apis/im/v1/chats", query, nil, body)

	return NewGroupInfo(info)
}

// GetGroupInfo Get a group information
func (c AppClient) GetGroupInfo(chat_id string) *GroupInfo {
	info := c.Request("get", "open-apis/im/v1/chats/"+chat_id+"?user_id_type=open_id", nil, nil, nil)
	info["chat_id"] = chat_id
	info["tenant_key"] = ""
	return NewGroupInfo(info)
}

// AddMembers
// app_id to add bot
func (c AppClient) AddMembers(chat_id string, member_id_type string, succeed_type string, id_list []string) bool {
	query := make(map[string]string)
	query["member_id_type"] = string(member_id_type)
	query["succeed_type"] = string(succeed_type)

	body := make(map[string][]string)

	var result bool = true
	var idlist []string

	for len(id_list) > 50 {

		idlist = id_list[0:50]
		body["id_list"] = idlist
		resp := c.Request("post", "open-apis/im/v1/chats/"+chat_id+"/members", query, nil, body)
		id_list = append(id_list[:0], id_list[50:]...)
		if resp == nil {
			result = false
		}
	}
	body["id_list"] = id_list
	resp := c.Request("post", "open-apis/im/v1/chats/"+chat_id+"/members", query, nil, body)
	if resp == nil {
		result = false
	}
	return result
}

// DeleteMembers
// app_id to delete bot
func (c AppClient) DeleteMembers(chat_id string, member_id_type string, id_list []string) bool {
	query := make(map[string]string)
	query["member_id_type"] = string(member_id_type)

	body := make(map[string][]string)

	var result bool = true
	var idlist []string

	for len(id_list) > 50 {

		idlist = id_list[0:50]
		body["id_list"] = idlist
		resp := c.Request("delete", "open-apis/im/v1/chats/"+chat_id+"/members", query, nil, body)
		id_list = append(id_list[:0], id_list[50:]...)
		if resp == nil {
			result = false
		}
	}
	body["id_list"] = id_list
	resp := c.Request("delete", "open-apis/im/v1/chats/"+chat_id+"/members", query, nil, body)
	if resp == nil {
		result = false
	}
	return result
}
