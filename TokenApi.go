package feishuapi

type LoginSession struct {
	OpenId     string
	EmployeeId string
}

func NewLoginSession(data map[string]interface{}) *LoginSession {
	return &LoginSession{
		OpenId:     data["open_id"].(string),
		EmployeeId: data["employee_id"].(string),
	}
}

func (c AppClient) GetLoginSession(login_token string) *LoginSession {
	body := make(map[string]string)
	body["code"] = login_token

	resp := c.Request("post", "open-apis/mina/v2/tokenLoginValidate", nil, nil, body)
	if resp == nil {
		return nil
	}
	return NewLoginSession(resp)
}
