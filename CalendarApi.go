package feishuapi

type CalendarPermission string

const (
	CalendarPrivate          CalendarPermission = "private"
	CalendarShowOnlyFreeBusy CalendarPermission = "show_only_free_busy"
	CalendarPublic           CalendarPermission = "public"
)

type CalendarCreateRequest struct {
	Summary     string             `json:"summary,omitempty"`
	Description string             `json:"description,omitempty"`
	Permissions CalendarPermission `json:"permissions,omitempty"`
}

func DefaultCalendarCreateRequest() *CalendarCreateRequest {
	return &CalendarCreateRequest{
		Summary:     "Calendar title",
		Description: "Calendar description",
		Permissions: CalendarPrivate,
	}
}

func (c *CalendarCreateRequest) WithSummary(summary string) *CalendarCreateRequest {
	c.Summary = summary
	return c
}

func (c *CalendarCreateRequest) WithDescription(description string) *CalendarCreateRequest {
	c.Description = description
	return c
}

func (c *CalendarCreateRequest) WithPermissions(permissions CalendarPermission) *CalendarCreateRequest {
	c.Permissions = permissions
	return c
}

type Calendar struct {
	Id           string
	CalendarInfo CalendarCreateRequest
}

func NewCalendar(data map[string]any) *Calendar {
	return &Calendar{
		Id: data["calendar_id"].(string),
		CalendarInfo: CalendarCreateRequest{
			Summary:     data["summary"].(string),
			Description: data["description"].(string),
			Permissions: CalendarPermission(data["permissions"].(string)),
		},
	}
}

func (c AppClient) CalendarCreate(calendar *CalendarCreateRequest) *Calendar {
	body := make(map[string]any)
	struct2map(calendar, &body)

	info := c.Request("post", "open-apis/calendar/v4/calendars", nil, nil, body)
	return NewCalendar(info["calendar"].(map[string]any))
}

func (c AppClient) CalendarSubscribe(calendarId string, user_access_token string) {
	headers := make(map[string]string)
	headers["Authorization"] = user_access_token

	c.Request("post", "open-apis/calendar/v4/calendars/"+calendarId+"/subscribe", nil, headers, nil)
}
