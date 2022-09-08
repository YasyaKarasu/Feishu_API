package feishuapi

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
	SpaceType   SpaceType
	Visibility  Visibility
}

func NewSpaceInfo(data map[string]interface{}) *SpaceInfo {
	return &SpaceInfo{
		Name:        data["name"].(string),
		Description: data["description"].(string),
		SpaceId:     data["space_id"].(string),
		SpaceType:   data["space_type"].(SpaceType),
		Visibility:  data["visibility"].(Visibility),
	}
}

func (c AppClient) CreateKnowledgeSpace(name string, description string, user_access_token string) *SpaceInfo {
	body := make(map[string]string)
	body["name"] = name
	body["description"] = description

	headers := make(map[string]string)
	headers["Authorization"] = user_access_token

	info := c.Request("post", "open-apis/wiki/v2/spaces", nil, headers, body)

	return NewSpaceInfo(info)
}

type Node struct {
	NodeToken       string
	ParentNodeToken string
	Title           string
}

func NewNode(data map[string]interface{}) *Node {
	return &Node{
		NodeToken:       data["node_token"].(string),
		ParentNodeToken: data["parent_node_token"].(string),
		Title:           data["title"].(string),
	}
}

func (c AppClient) CopyNode(SpaceId string, NodeToken string, TargetSpaceId string, TargetParentToken string, Title ...string) *Node {
	query := make(map[string]string)
	query["target_parent_token"] = TargetParentToken
	query["target_space_id"] = TargetSpaceId
	if len(Title) != 0 {
		query["title"] = Title[0]
	}

	info := c.Request("post", "open-apis/wiki/v2/spaces/"+SpaceId+"/nodes/"+NodeToken+"/copy", query, nil, nil)

	return NewNode(info)
}
