package feishuapi

import (
	"encoding/json"
	"strconv"
	"time"
)

type TimeInfo struct {
	Timestamp string `json:"timestamp"`
}

type VChat struct {
	VCType string `json:"vc_type"`
}

type EventLocation struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

type Reminder struct {
	Minutes int `json:"minutes"`
}

type EventAttendeeAbility string

const (
	AttendeeAbilityNone            EventAttendeeAbility = "none"
	AttendeeAbilityCanSeeOthers    EventAttendeeAbility = "can_see_others"
	AttendeeAbilityCanInviteOthers EventAttendeeAbility = "can_invite_others"
	AttendeeAbilityCanModifyEvent  EventAttendeeAbility = "can_modify_event"
)

type CalendarEventCreateRequest struct {
	Summary          string               `json:"summary,omitempty"`
	Description      string               `json:"description,omitempty"`
	NeedNotification bool                 `json:"need_notification"`
	StartTime        TimeInfo             `json:"start_time,omitempty"`
	EndTime          TimeInfo             `json:"end_time,omitempty"`
	VChat            VChat                `json:"vchat,omitempty"`
	AttendeeAbiliy   EventAttendeeAbility `json:"attendee_ability,omitempty"`
	Location         EventLocation        `json:"location,omitempty"`
	Reminders        []Reminder           `json:"reminders,omitempty"`
}

func DefaultCalendarEventCreateRequest() *CalendarEventCreateRequest {
	return &CalendarEventCreateRequest{
		Summary:          "Event title",
		Description:      "Event description",
		NeedNotification: true,
		StartTime:        TimeInfo{Timestamp: strconv.FormatInt(time.Now().Unix(), 10)},
		EndTime:          TimeInfo{Timestamp: strconv.FormatInt(time.Now().Add(time.Minute*30).Unix(), 10)},
		VChat: VChat{
			VCType: "vc",
		},
		AttendeeAbiliy: AttendeeAbilityCanSeeOthers,
		Location: EventLocation{
			Name:    "",
			Address: "",
		},
		Reminders: []Reminder{
			{
				Minutes: 5,
			},
		},
	}
}

func (c *CalendarEventCreateRequest) WithSummary(summary string) *CalendarEventCreateRequest {
	c.Summary = summary
	return c
}

func (c *CalendarEventCreateRequest) WithDescription(description string) *CalendarEventCreateRequest {
	c.Description = description
	return c
}

func (c *CalendarEventCreateRequest) WithNeedNotification(needNotification bool) *CalendarEventCreateRequest {
	c.NeedNotification = needNotification
	return c
}

func (c *CalendarEventCreateRequest) WithStartTime(startTime time.Time) *CalendarEventCreateRequest {
	c.StartTime.Timestamp = strconv.FormatInt(startTime.Unix(), 10)
	return c
}

func (c *CalendarEventCreateRequest) WithEndTime(endTime time.Time) *CalendarEventCreateRequest {
	c.EndTime.Timestamp = strconv.FormatInt(endTime.Unix(), 10)
	return c
}

func (c *CalendarEventCreateRequest) WithAttendeeAbility(ability EventAttendeeAbility) *CalendarEventCreateRequest {
	c.AttendeeAbiliy = ability
	return c
}

func (c *CalendarEventCreateRequest) WithLocation(name string, address string) *CalendarEventCreateRequest {
	c.Location.Name = name
	c.Location.Address = address
	return c
}

func (c *CalendarEventCreateRequest) WithReminders(minutes []int) *CalendarEventCreateRequest {
	c.Reminders = []Reminder{}
	for _, minute := range minutes {
		c.Reminders = append(c.Reminders, Reminder{Minutes: minute})
	}
	return c
}

type CalendarEvent struct {
	Id                  string
	OrganizerCalendarId string
	EventInfo           CalendarEventCreateRequest
}

func NewCalendarEvent(data map[string]interface{}) *CalendarEvent {
	eventInfo := CalendarEventCreateRequest{}
	map2struct(data, &eventInfo)
	return &CalendarEvent{
		Id:                  data["event_id"].(string),
		OrganizerCalendarId: data["organizer_calendar_id"].(string),
		EventInfo:           eventInfo,
	}
}

func (c AppClient) CalendarEventCreate(calendarId string, calendarEvent *CalendarEventCreateRequest) *CalendarEvent {
	body := make(map[string]interface{})
	struct2map(calendarEvent, &body)

	info := c.Request("post", "open-apis/calendar/v4/calendars/"+calendarId+"/events", nil, nil, body)
	return NewCalendarEvent(info["event"].(map[string]interface{}))
}

func (c AppClient) CalendarEventList(calendarId string) []CalendarEvent {
	query := make(map[string]string)
	query["anchor_time"] = strconv.FormatInt(time.Now().Unix(), 10)

	events := c.GetAllPages("get", "open-apis/calendar/v4/calendars/"+calendarId+"/events", query, nil, nil, 100)
	var calendarEvents []CalendarEvent
	for _, event := range events {
		calendarEvents = append(calendarEvents, *NewCalendarEvent(event.(map[string]interface{})))
	}
	return calendarEvents
}

type CalendarEventAttendeeType string

const (
	AttendeeTypeUser       CalendarEventAttendeeType = "user"
	AttendeeTypeChat       CalendarEventAttendeeType = "chat"
	AttendeeTypeResource   CalendarEventAttendeeType = "resource"
	AttendeeTypeThirdParty CalendarEventAttendeeType = "third_party"
)

type CalendarEventAttendeeRSVPStatus string

const (
	NeedsAction CalendarEventAttendeeRSVPStatus = "needs_action"
	Accept      CalendarEventAttendeeRSVPStatus = "accept"
	Tentative   CalendarEventAttendeeRSVPStatus = "tentative"
	Decline     CalendarEventAttendeeRSVPStatus = "decline"
	Removed     CalendarEventAttendeeRSVPStatus = "removed"
)

type CalendarEventAttendee struct {
	Type            CalendarEventAttendeeType       `json:"type,omitempty"`
	AttendeeId      string                          `json:"attendee_id,omitempty"`
	UserId          string                          `json:"user_id,omitempty"`
	ChatId          string                          `json:"chat_id,omitempty"`
	RoomId          string                          `json:"room_id,omitempty"`
	ThirdPartyEmail string                          `json:"third_party_email,omitempty"`
	OperateId       string                          `json:"operate_id,omitempty"`
	RSVPStatus      CalendarEventAttendeeRSVPStatus `json:"rsvp_status,omitempty"`
}

type CalendarEventAttendeeCreateRequest struct {
	Attendees        []CalendarEventAttendee `json:"attendees,omitempty"`
	NeedNotification bool                    `json:"need_notification"`
}

func DefaultCalendarEventAttendeeCreateRequest() *CalendarEventAttendeeCreateRequest {
	return &CalendarEventAttendeeCreateRequest{
		Attendees:        []CalendarEventAttendee{},
		NeedNotification: true,
	}
}

func (c *CalendarEventAttendeeCreateRequest) WithAttendee(attendee CalendarEventAttendee) *CalendarEventAttendeeCreateRequest {
	c.Attendees = append(c.Attendees, attendee)
	return c
}

func (c *CalendarEventAttendeeCreateRequest) WithAttendees(attendees []CalendarEventAttendee) *CalendarEventAttendeeCreateRequest {
	c.Attendees = append(c.Attendees, attendees...)
	return c
}

func (c *CalendarEventAttendeeCreateRequest) WithNeedNotification(needNotification bool) *CalendarEventAttendeeCreateRequest {
	c.NeedNotification = needNotification
	return c
}

func (c AppClient) CalendarEventAttendeeCreate(calendarId string, eventId string, userIdType UserIdType, attendee *CalendarEventAttendeeCreateRequest) []CalendarEventAttendee {
	query := make(map[string]string)
	query["user_id_type"] = string(userIdType)

	body := make(map[string]interface{})
	struct2map(attendee, &body)

	info := c.Request("post", "open-apis/calendar/v4/calendars/"+calendarId+"/events/"+eventId+"/attendees", query, nil, body)
	attendees_ := info["attendees"].([]interface{})

	attendees := []CalendarEventAttendee{}
	for _, attendee_ := range attendees_ {
		attendee := CalendarEventAttendee{}
		map2struct(attendee_.(map[string]interface{}), &attendee)
		attendees = append(attendees, attendee)
	}
	return attendees
}

func (c AppClient) CalendarEventAttendeeQuery(calendarId string, eventId string, userIdType UserIdType) []CalendarEventAttendee {
	query := make(map[string]string)
	query["user_id_type"] = string(userIdType)

	info := c.GetAllPages("get", "open-apis/calendar/v4/calendars/"+calendarId+"/events/"+eventId+"/attendees", query, nil, nil, 100)

	attendees := []CalendarEventAttendee{}
	for _, attendee_ := range info {
		attendee := CalendarEventAttendee{}
		map2struct(attendee_.(map[string]interface{}), &attendee)
		attendees = append(attendees, attendee)
	}
	return attendees
}

func struct2map(s interface{}, m interface{}) {
	bytes, _ := json.Marshal(s)
	json.Unmarshal(bytes, m)
}

func map2struct(m map[string]interface{}, s interface{}) {
	bytes, _ := json.Marshal(m)
	json.Unmarshal(bytes, s)
}
