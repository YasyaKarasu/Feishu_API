package feishuapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

type RobotInfo struct {
	Name   string
	OpenId string
}

func NewRobotInfo(data map[string]interface{}) *RobotInfo {
	return &RobotInfo{
		Name:   data["app_name"].(string),
		OpenId: data["open_id"].(string),
	}
}

func (c AppClient) RobotGetInfo() *RobotInfo {
	u := "https://open.feishu.cn/open-apis/bot/v3/info"

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c._tenant_access_token

	var req *http.Request

	Url, err := url.Parse(u)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": u,
			"err": err,
		}).Error("url parse error")
		return nil
	}

	urlPath := Url.String()

	req, err = http.NewRequest("GET", urlPath, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": u,
			"err": err,
		}).Error("Request create error")
		return nil
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
			"url": u,
			"err": err,
		}).Error("client do error")
		return nil
	}
	defer resp.Body.Close()

	if !responseOk(resp) {
		respBody, _ := ioutil.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"url":      u,
			"headers":  header,
			"error":    "response incorrect",
			"response": string(respBody),
		}).Error("response status error")
		return nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": u,
			"err": err,
		}).Error("read response body error")
		return nil
	}

	var result map[string]interface{}

	json.Unmarshal(respBody, &result)

	return NewRobotInfo(result["bot"].(map[string]interface{}))
}
