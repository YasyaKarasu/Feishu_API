package feishuapi

import (
	"encoding/json"
	"strconv"
	"time"
)

type TimelineNode struct {
	Type   string
	OpenId string
}

type ApprovalInstanceInfo struct {
	Status       string
	StartTime    time.Time
	EndTime      time.Time
	DepartmentId string
	Timeline     []TimelineNode
	Form         []map[string]interface{}
}

func (c AppClient) ApprovalInstanceById(InstanceCode string) *ApprovalInstanceInfo {
	resp := c.Request("get", "open-apis/approval/v4/instances/"+InstanceCode, nil, nil, nil)
	if resp == nil {
		return nil
	}
	start_time, _ := strconv.ParseInt(resp["start_time"].(string), 10, 64)
	end_time, _ := strconv.ParseInt(resp["end_time"].(string), 10, 64)
	var timeline []TimelineNode
	for _, v := range resp["timeline"].([]interface{}) {
		timeline = append(timeline, TimelineNode{
			Type:   v.(map[string]interface{})["type"].(string),
			OpenId: v.(map[string]interface{})["open_id"].(string),
		})
	}
	var form []map[string]interface{}
	json.Unmarshal([]byte(resp["form"].(string)), &form)
	return &ApprovalInstanceInfo{
		Status:       resp["status"].(string),
		StartTime:    time.Unix(start_time/1000, 0),
		EndTime:      time.Unix(end_time/1000, 0),
		DepartmentId: resp["department_id"].(string),
		Timeline:     timeline,
		Form:         form,
	}
}
