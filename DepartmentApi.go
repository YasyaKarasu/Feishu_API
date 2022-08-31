package feishuapi

type DepartmentInfo struct {
	Name        string
	GroupId     string
	MemberCount int
}

func NewDepartmentInfo(data map[string]interface{}) *DepartmentInfo {
	dept := data["department"].(map[string]interface{})
	return &DepartmentInfo{
		Name:        dept["name"].(string),
		GroupId:     dept["chat_id"].(string),
		MemberCount: dept["member_count"].(int),
	}
}

func (c AppClient) InfoById(department_id string) *DepartmentInfo {
	data := c.Request("get", "https://open.feishu.cn/open-apis/contact/v3/departments/"+department_id, nil, nil, nil)
	if data == nil {
		return nil
	}
	return NewDepartmentInfo(data)
}
