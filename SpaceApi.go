package feishuapi

import "github.com/sirupsen/logrus"

type SpaceType string
type Visibility string

const (
	Team   SpaceType = "team"
	Person SpaceType = "person"
)

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

type SpaceInfo struct {
	Name        string
	Description string
	SpaceId     string
}

// Create a new SpaceInfo
func NewSpaceInfo(data map[string]interface{}) *SpaceInfo {
	return &SpaceInfo{
		Name:        data["name"].(string),
		Description: data["description"].(string),
		SpaceId:     data["space_id"].(string),
	}
}

// Create a Knowledge Space
func (c AppClient) KnowledgeSpaceCreate(name string, description string, user_access_token string) *SpaceInfo {
	body := make(map[string]string)
	body["name"] = name
	body["description"] = description

	headers := make(map[string]string)
	headers["Authorization"] = user_access_token

	info := c.Request("post", "open-apis/wiki/v2/spaces", nil, headers, body)

	if info == nil {
		logrus.WithFields(logrus.Fields{
			"Name":        name,
			"Description": description,
		}).Error("create knowledge space fail")
		return nil
	}

	return NewSpaceInfo(info["space"].(map[string]interface{}))
}

// Add members to a Knowledge Space
// memberType: "openchat" for chat id, "userid" for feishuapi.UserId, "unionid" for feishuapi.UnionId, "opendepartmentid" for DepartmentId
func (c AppClient) KnowledgeSpaceAddMembers(spaceId string, membersId []string, memberType string) {
	body := make(map[string]string)
	body["member_type"] = memberType
	body["member_role"] = "member"
	for _, v := range membersId {
		body["member_id"] = v
		resp := c.Request("post", "open-apis/wiki/v2/spaces/"+spaceId+"/members", nil, nil, body)

		if resp == nil {
			logrus.WithFields(logrus.Fields{
				"SpaceID":    spaceId,
				"MemberType": memberType,
				"MemberID":   v,
			}).Warn("add member fail")
		}
	}
}

// Add robots to a Knowledge Space as admin
func (c AppClient) KnowledgeSpaceAddBotsAsAdmin(spaceId string, BotsId []string, user_access_token string) {
	headers := make(map[string]string)
	headers["Authorization"] = user_access_token

	body := make(map[string]string)
	body["member_type"] = "openid"
	body["member_role"] = "admin"

	for _, v := range BotsId {
		body["member_id"] = v
		resp := c.Request("post", "open-apis/wiki/v2/spaces/"+spaceId+"/members", nil, headers, body)

		if resp == nil {
			logrus.WithFields(logrus.Fields{
				"SpaceID":  spaceId,
				"MemberID": v,
			}).Warn("add bot fail")
		}
	}
}

type Node struct {
	NodeToken       string
	ParentNodeToken string
	Title           string
}

// Create a new Node
func NewNode(data map[string]interface{}) *Node {
	return &Node{
		NodeToken:       data["node_token"].(string),
		ParentNodeToken: data["parent_node_token"].(string),
		Title:           data["title"].(string),
	}
}

// Copy a node from SpaceId/NodeToken to TargetSpaceId/TargetParentToken
func (c AppClient) KnowledgeSpaceCopyNode(SpaceId string, NodeToken string, TargetSpaceId string, TargetParentToken string, Title ...string) *Node {
	body := make(map[string]string)
	body["target_parent_token"] = TargetParentToken
	body["target_space_id"] = TargetSpaceId
	if len(Title) != 0 {
		body["title"] = Title[0]
	}

	info := c.Request("post", "open-apis/wiki/v2/spaces/"+SpaceId+"/nodes/"+NodeToken+"/copy", nil, nil, body)

	if info == nil {
		logrus.WithFields(logrus.Fields{
			"SpaceID":           SpaceId,
			"NodeToken":         NodeToken,
			"TargetSpaceID":     TargetSpaceId,
			"TargetParentToken": TargetParentToken,
		}).Error("copy node fail")
		return nil
	}

	return NewNode(info["node"].(map[string]interface{}))
}

type NodeInfo struct {
	NodeToken       string
	ObjToken        string
	ObjType         string
	ParentNodeToken string
	Title           string
	HasChild        bool
}

// Create a new NodeInfo
func NewNodeInfo(data map[string]interface{}) *NodeInfo {
	return &NodeInfo{
		NodeToken:       data["node_token"].(string),
		ObjToken:        data["obj_token"].(string),
		ObjType:         data["obj_type"].(string),
		ParentNodeToken: data["parent_node_token"].(string),
		Title:           data["title"].(string),
		HasChild:        data["has_child"].(bool),
	}
}

// Get All Nodes in target Space and under specific ParentNode(not necessary)
func (c AppClient) KnowledgeSpaceGetAllNodes(SpaceId string, ParentNodeToken ...string) []NodeInfo {
	var all_node []NodeInfo
	var l []interface{}

	if len(ParentNodeToken) != 0 {
		query := make(map[string]string)
		query["parent_node_token"] = ParentNodeToken[0]
		l = c.GetAllPages("get", "open-apis/wiki/v2/spaces/"+SpaceId+"/nodes", query, nil, nil, 10)
	} else {
		l = c.GetAllPages("get", "open-apis/wiki/v2/spaces/"+SpaceId+"/nodes", nil, nil, nil, 10)
	}

	if l == nil {
		logrus.WithField("SpaceID", SpaceId).Error("nil node info return")
		return nil
	}

	for _, value := range l {
		all_node = append(all_node, *NewNodeInfo(value.(map[string]interface{})))
	}

	return all_node
}
