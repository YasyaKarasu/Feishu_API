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
	"github.com/spf13/viper"
)

type AppClient struct {
	_tenant_access_token string
	schedular            *cron.Cron
}

func _url(path string) string {
	return strings.Trim(viper.GetString("feishu.LARK_HOST"), "/") + "/" + strings.Trim(path, "/")
}

func (c *AppClient) authorizeTenantAccessToken() bool {
	u := _url("/open-apis/auth/v3/tenant_access_token/internal")

	urlValues := url.Values{}
	urlValues.Add("app_id", viper.GetString("feishu.APP_ID"))
	urlValues.Add("app_secret", viper.GetString("feishu.APP_SECRET"))

	resp, err := http.PostForm(u, urlValues)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"response": resp,
			"error":    err,
		}).Error("cannot get tenant_access_token")
		return false
	}

	defer resp.Body.Close()

	var result map[string]interface{}

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

	var result map[string]interface{}

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

func (c AppClient) Request(method string, path string, query map[string]string, headers map[string]string, body interface{}) map[string]interface{} {
	u := _url(path)

	header := make(map[string]string, len(headers))
	for k, v := range headers {
		header[k] = v
	}

	if _, ok := header["Authorization"]; !ok {
		header["Authorization"] = "Bearer " + c._tenant_access_token
	}
	header["content-type"] = "application/json; charset=utf-8"

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

		logrus.Info(body)

		req, err = http.NewRequest("POST", urlPath, bytes.NewReader(bytesData))
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

	logrus.Info(req)

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

	var result map[string]interface{}

	json.Unmarshal(respBody, &result)

	//logrus.Info(result["data"])

	return result["data"].(map[string]interface{})
}

func (c AppClient) GetAllPages(method string, path string, query map[string]string, headers map[string]string, body interface{}, page_size int) []interface{} {
	if page_size < 10 || page_size > 100 {
		return nil
	}

	var all_list []interface{}
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

		l := resp["items"].([]interface{})
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

func GetInMap(mapToSearch map[string]interface{}, key string, defaults interface{}) interface{} {
	value, ok := mapToSearch[key]
	if !ok {
		return defaults
	} else {
		return value
	}
}
