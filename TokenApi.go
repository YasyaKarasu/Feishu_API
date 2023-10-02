package feishuapi

import "github.com/sirupsen/logrus"

type LoginSession struct {
	OpenId     string
	EmployeeId string
}

type UserAccessToken struct {
	Access_token  string
	Name          string
	Refresh_token string
	User_id       string
	Open_id       string
}

// Create a new LoginSession
func NewLoginSession(data map[string]any) *LoginSession {
	return &LoginSession{
		OpenId:     data["open_id"].(string),
		EmployeeId: data["employee_id"].(string),
	}
}

// Create a new UserAccessToken
func NewUserAccessToken(data map[string]any) *UserAccessToken {
	return &UserAccessToken{
		Access_token:  data["access_token"].(string),
		Name:          data["name"].(string),
		Refresh_token: data["refresh_token"].(string),
		User_id:       getInMap(data, "user_id", "").(string),
		Open_id:       data["open_id"].(string),
	}
}

// Get the login session by login_token
func (c AppClient) GetLoginSession(login_token string) *LoginSession {
	body := make(map[string]string)
	body["code"] = login_token

	resp := c.Request("post", "open-apis/mina/v2/tokenLoginValidate", nil, nil, body)
	if resp == nil {
		logrus.Error("nil login session return")
		return nil
	}
	return NewLoginSession(resp)
}

// Get the UserAccessToken
func (c *AppClient) GetUserAccessToken(code string) *UserAccessToken {
	u := "open-apis/authen/v1/access_token"

	body := make(map[string]string)
	body["grant_type"] = "authorization_code"
	body["code"] = code

	resp := c.Request("post", u, nil, nil, body)

	if resp == nil {
		logrus.Error("nil user access token return")
		return nil
	}

	return NewUserAccessToken(resp)
}

func (c *AppClient) GetCode(redirectURL string, appID string) string {
	u := "open-apis/authen/v1/index"

	query := make(map[string]any)
	query["redirect_uri"] = redirectURL
	query["app_id"] = appID

	resp := c.Request("post", u, query, nil, nil)

	if resp == nil {
		logrus.Error("nil code return")
		return ""
	}

	return resp["code"].(string)
}
