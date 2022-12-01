package feishuapi

import "github.com/sirupsen/logrus"

type DepartmentInfo struct {
	Name        string
	GroupId     string
	MemberCount int
}

// Create a new DepartmentInfo
func NewDepartmentInfo(data map[string]interface{}) *DepartmentInfo {
	dept := data["department"].(map[string]interface{})
	return &DepartmentInfo{
		Name:        dept["name"].(string),
		GroupId:     dept["chat_id"].(string),
		MemberCount: dept["member_count"].(int),
	}
}

// Send a request to get the information of a department by department_id
func (c AppClient) DepartmentGetInfoById(DepartmentId string) *DepartmentInfo {
	data := c.Request("get", "open-apis/contact/v3/departments/"+DepartmentId, nil, nil, nil)
	if data == nil {
		logrus.WithField("DepartmentID", DepartmentId).Warn("nil department info return")
		return nil
	}
	return NewDepartmentInfo(data)
}
