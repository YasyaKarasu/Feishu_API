package feishuapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppId             string
	AppSecret         string
	VerificationToken string
	EncryptKey        string
}

type AppClient struct {
	_tenant_access_token string
	schedular            *cron.Cron
	Conf                 Config
}

func (c AppClient) url(path string) string {
	return "https://open.feishu.cn" + "/" + strings.Trim(path, "/")
}

func (c *AppClient) authorizeTenantAccessToken() bool {
	u := c.url("/open-apis/auth/v3/tenant_access_token/internal")

	urlValues := url.Values{}
	urlValues.Add("app_id", c.Conf.AppId)
	urlValues.Add("app_secret", c.Conf.AppSecret)

	resp, err := http.PostForm(u, urlValues)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"response": resp,
			"error":    err,
		}).Error("cannot get tenant_access_token")
		return false
	}

	defer resp.Body.Close()

	var result map[string]any

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"response": resp,
			"error":    err,
		}).Error("cannot get tenant_access_token")
		return false
	}
	json.Unmarshal(body, &result)

	c._tenant_access_token = result["tenant_access_token"].(string)
	logrus.Info("got tenant_access_token: " + c._tenant_access_token)
	return true
}

// Start a schedular to get tenant_access_token every 105 minutes
func (c *AppClient) StartTokenTimer() {
	if !c.authorizeTenantAccessToken() {
		logrus.Error("cannot get feishu token")
	}

	c.schedular = cron.New()
	c.schedular.AddFunc("@every 105m", func() {
		if !c.authorizeTenantAccessToken() {
			logrus.Error("cannot get feishu token")
		}
	})

	c.schedular.Start()
}

func responseOk(resp *http.Response) bool {
	if resp.StatusCode != 200 {
		return false
	}

	var result map[string]any

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"response": resp,
			"error":    err,
		}).Error("read response body error")
		return false
	}
	json.Unmarshal(body, &result)

	interface_code, ok := result["code"]
	if !ok {
		logrus.WithField("response", resp).Error("read response code error")
		return false
	}

	code := interface_code.(float64)

	return code == 0
}

// Send a GET / POST / DELETE string to a specific path, with header of authorization and content-type
// (in other words, authorization and content-type should not to be passed)
// On the other hand, if the api needs user_access_token, then you can pass it by headers param
func (c AppClient) Request(method string, path string, query map[string]string, headers map[string]string, body any) map[string]any {
	u := c.url(path)

	header := make(map[string]string, len(headers))
	for k, v := range headers {
		header[k] = v
	}

	if _, ok := header["Authorization"]; !ok {
		header["Authorization"] = "Bearer " + c._tenant_access_token
	}
	header["Content-Type"] = "application/json; charset=utf-8"

	var req *http.Request

	urlValues := url.Values{}
	Url, err := url.Parse(u)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":   u,
			"query": query,
			"err":   err,
		}).Error("url parse error")
		return nil
	}

	for k, v := range query {
		urlValues.Set(k, v)
	}

	Url.RawQuery = urlValues.Encode()
	urlPath := Url.String()

	if strings.EqualFold(method, "get") {
		req, err = http.NewRequest("GET", urlPath, nil)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"url":   u,
				"query": query,
				"err":   err,
			}).Error("Request create error")
			return nil
		}

	} else {
		bytesData, err := json.Marshal(body)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"url":   u,
				"query": query,
				"err":   err,
			}).Error("marshal query error")
			return nil
		}

		req, err = http.NewRequest(strings.ToUpper(method), urlPath, bytes.NewReader(bytesData))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"url":   u,
				"query": query,
				"err":   err,
			}).Error("Request create error")
			return nil
		}
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	cli := &http.Client{
		Timeout: time.Second * 15,
	}

	resp, err := cli.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":   u,
			"query": query,
			"err":   err,
		}).Error("client do error")
		return nil
	}
	defer resp.Body.Close()

	if !responseOk(resp) {
		jsonBody, _ := json.Marshal(body)
		respBody, _ := ioutil.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"url":      u,
			"query":    query,
			"headers":  headers,
			"data":     jsonBody,
			"error":    "response incorrect",
			"response": string(respBody),
		}).Error("response status error")
		return nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":   u,
			"query": query,
			"err":   err,
		}).Error("read response body error")
		return nil
	}

	var result map[string]any

	json.Unmarshal(respBody, &result)

	return result["data"].(map[string]any)
}

// Send Request several times until all the pages of information are got
func (c AppClient) GetAllPages(method string, path string, query map[string]string, headers map[string]string, body any, page_size int) []any {
	if page_size < 10 || page_size > 100 {
		logrus.Info("page_size should be between 10 and 100")
		return nil
	}

	var all_list []any
	page_token := ""
	has_more := true
	queries := make(map[string]string)
	if len(query) != 0 {
		for k, v := range query {
			queries[k] = v
		}
	}
	queries["page_size"] = strconv.Itoa(page_size)

	for {
		if !has_more {
			break
		}
		if page_token != "" {
			queries["page_token"] = page_token
		}

		resp := c.Request(method, path, queries, headers, body)
		if len(resp) == 0 {
			return nil
		}

		l := resp["items"].([]any)
		has_more = resp["has_more"].(bool)
		interface_page_token, ok := resp["page_token"]
		if ok {
			page_token = interface_page_token.(string)
		} else {
			page_token = ""
		}

		all_list = append(all_list, l...)
	}

	return all_list
}

// Get the value of provided key in a map, if there's no such key than return provided defaults
func getInMap(mapToSearch map[string]any, key string, defaults any) any {
	value, ok := mapToSearch[key]
	if !ok {
		return defaults
	} else {
		return value
	}
}
