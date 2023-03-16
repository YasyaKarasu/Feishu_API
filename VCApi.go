package feishuapi

import (
	"strconv"
	"time"
)

type ReserveMeetingSetting struct {
	Topic      string `json:"topic,omitempty"`
	AutoRecord bool   `json:"auto_record"`
}

type VCReserveRequest struct {
	EndTime         string                `json:"end_time,omitempty"`
	OwnerId         string                `json:"owner_id,omitempty"`
	MeetingSettings ReserveMeetingSetting `json:"meeting_settings,omitempty"`
}

func DefaultVCReserveRequest() *VCReserveRequest {
	return &VCReserveRequest{
		EndTime: strconv.FormatInt(time.Now().Add(time.Minute*30).Unix(), 10),
		OwnerId: "",
		MeetingSettings: ReserveMeetingSetting{
			Topic:      "Meeting title",
			AutoRecord: false,
		},
	}
}

func (vc *VCReserveRequest) WithEndTime(endTime time.Time) *VCReserveRequest {
	vc.EndTime = strconv.FormatInt(endTime.Unix(), 10)
	return vc
}

func (vc *VCReserveRequest) WithOwnerId(ownerId string) *VCReserveRequest {
	vc.OwnerId = ownerId
	return vc
}

func (vc *VCReserveRequest) WithTopic(topic string) *VCReserveRequest {
	vc.MeetingSettings.Topic = topic
	return vc
}

func (vc *VCReserveRequest) WithAutoRecord(autoRecord bool) *VCReserveRequest {
	vc.MeetingSettings.AutoRecord = autoRecord
	return vc
}

type VCReserve struct {
	Id        string
	MeetingNo string
	URL       string
}

func NewVCReserve(data map[string]interface{}) *VCReserve {
	return &VCReserve{
		Id:        data["id"].(string),
		MeetingNo: data["meeting_no"].(string),
		URL:       data["url"].(string),
	}
}

func (c AppClient) VCReserve(reserveRequest *VCReserveRequest) *VCReserve {
	body := make(map[string]interface{})
	struct2map(reserveRequest, &body)

	info := c.Request("post", "/open-apis/vc/v1/reserves/apply", nil, nil, body)
	return NewVCReserve(info)
}

func (c AppClient) VCReserveWithTopic(topic string, endTime time.Time) *VCReserve {
	return c.VCReserve(DefaultVCReserveRequest().WithTopic(topic).WithEndTime(endTime))
}
