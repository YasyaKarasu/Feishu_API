package feishuapi

import "github.com/sirupsen/logrus"

type UserInfo struct {
	UnionId       string
	OpenId        string
	UserId        string
	Name          string
	DepartmentIds []interface{}
}

func (c AppClient) UserInfoById(UserId string, IdType UserIdType) *UserInfo {
	query := make(map[string]string)
	query["user_id_type"] = string(IdType)
	data := c.Request("get", "open-apis/contact/v3/users/"+UserId, query, nil, nil)
	if data == nil {
		logrus.WithField("UserId", UserId).Warn("nil user info return")
		return nil
	}
	user := data["user"].(map[string]interface{})
	return &UserInfo{
		UnionId:       user["union_id"].(string),
		OpenId:        user["open_id"].(string),
		UserId:        user["user_id"].(string),
		Name:          user["name"].(string),
		DepartmentIds: user["department_ids"].([]interface{}),
	}
}
