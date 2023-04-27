package feishuapi

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
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

type Participant struct {
	ParticipantName string `json:"participant_name"`
	Department      string `json:"department"`
	UserId          string `json:"user_id"`
	EmployeeId      string `json:"employee_id"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Device          string `json:"device"`
	AppVersion      string `json:"app_version"`
	PublicIp        string `json:"public_ip"`
	InternalIp      string `json:"internal_ip"`
	UseRtcProxy     bool   `json:"use_rtc_proxy"`
	Location        string `json:"location"`
	NetworkType     string `json:"network_type"`
	Protocol        string `json:"protocol"`
	Microphone      string `json:"microphone"`
	Speaker         string `json:"speaker"`
	Camera          string `json:"camera"`
	Audio           bool   `json:"audio"`
	Video           bool   `json:"video"`
	Sharing         bool   `json:"sharing"`
	JoinTime        string `json:"join_time"`
	LeaveTime       string `json:"leave_time"`
	TimeInMeeting   string `json:"time_in_meeting"`
	LeaveReason     string `json:"leave_reason"`
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

func NewVCReserve(data map[string]any) *VCReserve {
	return &VCReserve{
		Id:        data["id"].(string),
		MeetingNo: data["meeting_no"].(string),
		URL:       data["url"].(string),
	}
}

func (c AppClient) VCReserve(reserveRequest *VCReserveRequest) *VCReserve {
	body := make(map[string]any)
	struct2map(reserveRequest, &body)

	info := c.Request("post", "/open-apis/vc/v1/reserves/apply", nil, nil, body)
	return NewVCReserve(info)
}

func (c AppClient) VCReserveWithTopic(topic string, endTime time.Time) *VCReserve {
	return c.VCReserve(DefaultVCReserveRequest().WithTopic(topic).WithEndTime(endTime))
}

func (c AppClient) VCQueryParticipantList(meetingStartTime int64, meetingEndTime int64, meetingNo string) []Participant {
	query := make(map[string]any)
	query["meeting_start_time"] = strconv.FormatInt(meetingStartTime, 10)
	query["meeting_end_time"] = strconv.FormatInt(meetingEndTime, 10)
	query["meeting_no"] = meetingNo

	info := c.getAllParticipants("get", "open-apis/vc/v1/participant_list", query, nil, nil, 100)

	var participants []Participant
	for _, participant_ := range info {
		participant := Participant{}
		data, err := json.Marshal(participant_)
		if err != nil {
			return nil
		}
		err = json.Unmarshal(data, &participant)
		if err != nil {
			return nil
		}
		participants = append(participants, participant)
	}
	return participants
}

// [Rewrite], because GetAllPages can't be used to VCQueryParticipantList
func (c AppClient) getAllParticipants(method string, path string, query map[string]any, headers map[string]string, body any, page_size int) []any {
	if page_size < 10 || page_size > 100 {
		logrus.Info("page_size should be between 10 and 100")
		return nil
	}

	var all_list []any
	page_token := ""
	has_more := true
	queries := make(map[string]any)
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
		if resp["participants"] == nil {
			return []any{}
		}

		l := resp["participants"].([]any)
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
