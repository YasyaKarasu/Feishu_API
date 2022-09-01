package feishuapi

type LoginSession struct {
	OpenId     string
	EmployeeId string
}

// Create a new LoginSession
func NewLoginSession(data map[string]interface{}) *LoginSession {
	return &LoginSession{
		OpenId:     data["open_id"].(string),
		EmployeeId: data["employee_id"].(string),
	}
}

// Get the login session by login_token
func (c AppClient) GetLoginSession(login_token string) *LoginSession {
	body := make(map[string]string)
	body["code"] = login_token

	resp := c.Request("post", "open-apis/mina/v2/tokenLoginValidate", nil, nil, body)
	if resp == nil {
		return nil
	}
	return NewLoginSession(resp)
}
