package feishuapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type MessageCard struct {
	Config   *MessageCardConfig   `json:"config,omitempty"`
	Header   *MessageCardHeader   `json:"header,omitempty"`
	Elements []MessageCardElement `json:"elements,omitempty"`
	CardLink *MessageCardLink     `json:"card_link,omitempty"`
}

func NewMessageCard() *MessageCard {
	return &MessageCard{}
}

func (card *MessageCard) WithConfig(config *MessageCardConfig) *MessageCard {
	card.Config = config
	return card
}

func (card *MessageCard) WithHeader(header *MessageCardHeader) *MessageCard {
	card.Header = header
	return card
}

func (card *MessageCard) WithElements(elements []MessageCardElement) *MessageCard {
	card.Elements = elements
	return card
}

func (card *MessageCard) Build() *MessageCard {
	return card
}

func (card *MessageCard) String() (string, error) {
	if len(card.Elements) == 0 {
		return "", errors.New("elements is required")
	}
	data, err := json.Marshal(card)
	return string(data), err
}

type MessageCardConfig struct {
	EnableForward *bool `json:"enable_forward,omitempty"`
	UpdateMulti   *bool `json:"update_multi,omitempty"`
}

func NewMessageCardConfig() *MessageCardConfig {
	return &MessageCardConfig{}
}

func (config *MessageCardConfig) WithEnableForward(enableForward bool) *MessageCardConfig {
	config.EnableForward = &enableForward
	return config
}

func (config *MessageCardConfig) WithUpdateMulti(updateMulti bool) *MessageCardConfig {
	config.UpdateMulti = &updateMulti
	return config
}

func (config *MessageCardConfig) Build() *MessageCardConfig {
	return config
}

type MessageCardTitleTemplate string

const (
	TemplateBlue      MessageCardTitleTemplate = "blue"
	TemplateWathet    MessageCardTitleTemplate = "wathet"
	TemplateTurquoise MessageCardTitleTemplate = "turquoise"
	TemplateGreen     MessageCardTitleTemplate = "green"
	TemplateYellow    MessageCardTitleTemplate = "yellow"
	TemplateOrange    MessageCardTitleTemplate = "orange"
	TemplateRed       MessageCardTitleTemplate = "red"
	TemplateCarmine   MessageCardTitleTemplate = "carmine"
	TemplateViolet    MessageCardTitleTemplate = "violet"
	TemplatePurple    MessageCardTitleTemplate = "purple"
	TemplateIndigo    MessageCardTitleTemplate = "indigo"
	TemplateGrey      MessageCardTitleTemplate = "grey"
	TemplateDefault   MessageCardTitleTemplate = "default"
)

type MessageCardHeader struct {
	Title    *MessageCardPlainText     `json:"title,omitempty"`
	Template *MessageCardTitleTemplate `json:"template,omitempty"`
}

func NewMessageCardHeader() *MessageCardHeader {
	return &MessageCardHeader{}
}

func (header *MessageCardHeader) WithTitle(title *MessageCardPlainText) *MessageCardHeader {
	header.Title = title
	return header
}

func (header *MessageCardHeader) WithTemplate(template MessageCardTitleTemplate) *MessageCardHeader {
	header.Template = &template
	return header
}

func (header *MessageCardHeader) Build() *MessageCardHeader {
	return header
}

type MessageCardElement interface {
	Tag() string
	MarshalJSON() ([]byte, error)
}

func messageCardElementJSON(element MessageCardElement) ([]byte, error) {
	data, err := struct2mapByReflect(element)
	if err != nil {
		return nil, err
	}
	data["tag"] = element.Tag()
	return json.Marshal(data)
}

type MessageCardText interface {
	MessageCardElement
	IsText()
}

type MessageCardPlainText struct {
	Content string `json:"content,omitempty"`
	Lines   *int   `json:"lines,omitempty"`
}

func NewMessageCardPlainText() *MessageCardPlainText {
	return &MessageCardPlainText{}
}

func (plainText *MessageCardPlainText) WithContent(content string) *MessageCardPlainText {
	plainText.Content = content
	return plainText
}

func (plainText *MessageCardPlainText) WithLines(lines int) *MessageCardPlainText {
	plainText.Lines = &lines
	return plainText
}

func (plainText *MessageCardPlainText) Build() *MessageCardPlainText {
	return plainText
}

func (plainText *MessageCardPlainText) Tag() string {
	return "plain_text"
}

func (plainText *MessageCardPlainText) MarshalJSON() ([]byte, error) {
	return messageCardElementJSON(plainText)
}

func (plainText *MessageCardPlainText) IsText() {}

func (plainText *MessageCardPlainText) IsNote() {}

type MessageCardLarkMarkdown struct {
	Content string `json:"content,omitempty"`
}

func NewMessageCardLarkMarkdown() *MessageCardLarkMarkdown {
	return &MessageCardLarkMarkdown{}
}

func (larkMarkdown *MessageCardLarkMarkdown) WithContent(content string) *MessageCardLarkMarkdown {
	larkMarkdown.Content = content
	return larkMarkdown
}

func (larkMarkdown *MessageCardLarkMarkdown) Build() *MessageCardLarkMarkdown {
	return larkMarkdown
}

func (larkMarkdown *MessageCardLarkMarkdown) Tag() string {
	return "lark_md"
}

func (larkMarkdown *MessageCardLarkMarkdown) MarshalJSON() ([]byte, error) {
	return messageCardElementJSON(larkMarkdown)
}

func (larkMarkdown *MessageCardLarkMarkdown) IsText() {}

func (larkMarkdown *MessageCardLarkMarkdown) IsNote() {}

type MessageCardColumnSetFlexMode string

const (
	FlexModeNone    MessageCardColumnSetFlexMode = "none"
	FlexModeStretch MessageCardColumnSetFlexMode = "stretch"
	FlexModeFlow    MessageCardColumnSetFlexMode = "flow"
	FlexModeBisect  MessageCardColumnSetFlexMode = "bisect"
	FlexModeTrisect MessageCardColumnSetFlexMode = "trisect"
)

type MessageCardColumnSetBackgroundStyle string

const (
	BackgroundStyleDefault MessageCardColumnSetBackgroundStyle = "default"
	BackgroundStyleGrey    MessageCardColumnSetBackgroundStyle = "grey"
)

type MessageCardColumnSetHorizontalSpacing string

const (
	HorizontalSpacingDefault MessageCardColumnSetHorizontalSpacing = "default"
	HorizontalSpacingSmall   MessageCardColumnSetHorizontalSpacing = "small"
)

type MessageCardColumnSet struct {
	FlexMode          *MessageCardColumnSetFlexMode          `json:"flex_mode,omitempty"`
	BackgroundStyle   *MessageCardColumnSetBackgroundStyle   `json:"background_style,omitempty"`
	HorizontalSpacing *MessageCardColumnSetHorizontalSpacing `json:"horizontal_spacing,omitempty"`
	Columns           []MessageCardColumn                    `json:"columns,omitempty"`
}

func NewMessageCardColumnSet() *MessageCardColumnSet {
	return &MessageCardColumnSet{}
}

func (columnSet *MessageCardColumnSet) WithFlexMode(flexMode MessageCardColumnSetFlexMode) *MessageCardColumnSet {
	columnSet.FlexMode = &flexMode
	return columnSet
}

func (columnSet *MessageCardColumnSet) WithBackgroundStyle(backgroundStyle MessageCardColumnSetBackgroundStyle) *MessageCardColumnSet {
	columnSet.BackgroundStyle = &backgroundStyle
	return columnSet
}

func (columnSet *MessageCardColumnSet) WithHorizontalSpacing(horizontalSpacing MessageCardColumnSetHorizontalSpacing) *MessageCardColumnSet {
	columnSet.HorizontalSpacing = &horizontalSpacing
	return columnSet
}

func (columnSet *MessageCardColumnSet) WithColumns(columns []MessageCardColumn) *MessageCardColumnSet {
	columnSet.Columns = columns
	return columnSet
}

func (columnSet *MessageCardColumnSet) Build() *MessageCardColumnSet {
	return columnSet
}

func (columnSet *MessageCardColumnSet) Tag() string {
	return "column_set"
}

func (columnSet *MessageCardColumnSet) MarshalJSON() ([]byte, error) {
	if columnSet.FlexMode == nil {
		return nil, errors.New("flex_mode is required")
	}
	return messageCardElementJSON(columnSet)
}

type MessageCardColumnWidth string

const (
	WidthAuto     MessageCardColumnWidth = "auto"
	WidthWeighted MessageCardColumnWidth = "weighted"
)

type MessageCardVerticalAlign string

const (
	VerticalAlignTop    MessageCardVerticalAlign = "top"
	VerticalAlignCenter MessageCardVerticalAlign = "center"
	VerticalAlignBottom MessageCardVerticalAlign = "bottom"
)

type MessageCardColumn struct {
	Width         *MessageCardColumnWidth   `json:"width,omitempty"`
	Weight        *int                      `json:"weight,omitempty"`
	VerticalAlign *MessageCardVerticalAlign `json:"vertical_align,omitempty"`
	Elements      []MessageCardElement      `json:"elements,omitempty"`
}

func NewMessageCardColumn() *MessageCardColumn {
	return &MessageCardColumn{}
}

func (column *MessageCardColumn) WithWidth(width MessageCardColumnWidth) *MessageCardColumn {
	column.Width = &width
	return column
}

func (column *MessageCardColumn) WithWeight(weight int) *MessageCardColumn {
	column.WithWidth(WidthWeighted)
	column.Weight = &weight
	return column
}

func (column *MessageCardColumn) WithVerticalAlign(verticalAlign MessageCardVerticalAlign) *MessageCardColumn {
	column.VerticalAlign = &verticalAlign
	return column
}

func (column *MessageCardColumn) WithElements(elements []MessageCardElement) *MessageCardColumn {
	column.Elements = elements
	return column
}

func (column *MessageCardColumn) Build() *MessageCardColumn {
	return column
}

func (column *MessageCardColumn) Tag() string {
	return "column"
}

func (column *MessageCardColumn) MarshalJSON() ([]byte, error) {
	return messageCardElementJSON(column)
}

type MessageCardDiv struct {
	Text   MessageCardText     `json:"text,omitempty"`
	Fields []*MessageCardField `json:"fields,omitempty"`
	Extra  MessageCardExtra    `json:"extra,omitempty"`
}

func NewMessageCardDiv() *MessageCardDiv {
	return &MessageCardDiv{}
}

func (div *MessageCardDiv) WithText(text MessageCardText) *MessageCardDiv {
	div.Text = text
	return div
}

func (div *MessageCardDiv) WithFields(fields []*MessageCardField) *MessageCardDiv {
	div.Fields = fields
	return div
}

func (div *MessageCardDiv) WithExtra(extra MessageCardExtra) *MessageCardDiv {
	div.Extra = extra
	return div
}

func (div *MessageCardDiv) Build() *MessageCardDiv {
	return div
}

func (div *MessageCardDiv) Tag() string {
	return "div"
}

func (div *MessageCardDiv) MarshalJSON() ([]byte, error) {
	if div.Text == nil && len(div.Fields) == 0 {
		return nil, errors.New("text or fields is required")
	}
	return messageCardElementJSON(div)
}

type MessageCardField struct {
	IsShort *bool           `json:"is_short,omitempty"`
	Text    MessageCardText `json:"text,omitempty"`
}

func NewMessageCardField() *MessageCardField {
	return &MessageCardField{}
}

func (field *MessageCardField) WithIsShort(isShort bool) *MessageCardField {
	field.IsShort = &isShort
	return field
}

func (field *MessageCardField) WithText(text MessageCardText) *MessageCardField {
	field.Text = text
	return field
}

func (field *MessageCardField) Build() *MessageCardField {
	return field
}

func (field *MessageCardField) MarshalJSON() ([]byte, error) {
	if field.Text == nil {
		return nil, errors.New("text is required")
	}
	return json.Marshal(field)
}

type MessageCardExtra interface {
	MessageCardElement
	IsExtra()
}

type MessageCardMarkdownTextAlign string

const (
	TextAlignLeft   MessageCardMarkdownTextAlign = "left"
	TextAlignCenter MessageCardMarkdownTextAlign = "center"
	TextAlignRight  MessageCardMarkdownTextAlign = "right"
)

type MessageCardMarkdown struct {
	Content   string                        `json:"content,omitempty"`
	TextAlign *MessageCardMarkdownTextAlign `json:"text_align,omitempty"`
	Href      map[string]*MessageCardURL    `json:"href,omitempty"`
}

func NewMessageCardMarkdown() *MessageCardMarkdown {
	return &MessageCardMarkdown{}
}

func (markdown *MessageCardMarkdown) WithContent(content string) *MessageCardMarkdown {
	markdown.Content = content
	return markdown
}

func (markdown *MessageCardMarkdown) WithTextAlign(textAlign MessageCardMarkdownTextAlign) *MessageCardMarkdown {
	markdown.TextAlign = &textAlign
	return markdown
}

func (markdown *MessageCardMarkdown) WithHref(href map[string]*MessageCardURL) *MessageCardMarkdown {
	markdown.Href = href
	return markdown
}

func (markdown *MessageCardMarkdown) Build() *MessageCardMarkdown {
	return markdown
}

func (markdown *MessageCardMarkdown) Tag() string {
	return "markdown"
}

func (markdown *MessageCardMarkdown) MarshalJSON() ([]byte, error) {
	return messageCardElementJSON(markdown)
}

type MessageCardURL struct {
	URL        *string `json:"url,omitempty"`
	AndroidURL *string `json:"android_url,omitempty"`
	IOSURL     *string `json:"ios_url,omitempty"`
	PCURL      *string `json:"pc_url,omitempty"`
}

func NewMessageCardURL() *MessageCardURL {
	return &MessageCardURL{}
}

func (url *MessageCardURL) WithURL(urlStr string) *MessageCardURL {
	url.URL = &urlStr
	return url
}

func (url *MessageCardURL) WithAndroidURL(androidURL string) *MessageCardURL {
	url.AndroidURL = &androidURL
	return url
}

func (url *MessageCardURL) WithIOSURL(iosURL string) *MessageCardURL {
	url.IOSURL = &iosURL
	return url
}

func (url *MessageCardURL) WithPCURL(pcURL string) *MessageCardURL {
	url.PCURL = &pcURL
	return url
}

func (url *MessageCardURL) Build() *MessageCardURL {
	return url
}

func (url *MessageCardURL) MarshalJSON() ([]byte, error) {
	if url.URL == nil {
		return nil, errors.New("url is required")
	}
	if url.AndroidURL == nil {
		url.AndroidURL = url.URL
	}
	if url.IOSURL == nil {
		url.IOSURL = url.URL
	}
	if url.PCURL == nil {
		url.PCURL = url.URL
	}
	return json.Marshal(url)
}

type MessageCardHr struct {
}

func NewMessageCardHr() *MessageCardHr {
	return &MessageCardHr{}
}

func (hr *MessageCardHr) Tag() string {
	return "hr"
}

func (hr *MessageCardHr) Build() *MessageCardHr {
	return hr
}

func (hr *MessageCardHr) MarshalJSON() ([]byte, error) {
	return messageCardElementJSON(hr)
}

type MessageCardImageMode string

const (
	ModeCropCenter    MessageCardImageMode = "crop_center"
	ModeFitHorizontal MessageCardImageMode = "fit_horizontal"
)

type MessageCardImage struct {
	ImageKey     *string               `json:"img_key,omitempty"`
	Alt          *MessageCardPlainText `json:"alt,omitempty"`
	Title        MessageCardText       `json:"title,omitempty"`
	CustomWidth  *int                  `json:"custom_width,omitempty"`
	CompactWidth *bool                 `json:"compact_width,omitempty"`
	Mode         *MessageCardImageMode `json:"mode,omitempty"`
	Preview      *bool                 `json:"preview,omitempty"`
}

func NewMessageCardImage() *MessageCardImage {
	return &MessageCardImage{}
}

func (image *MessageCardImage) WithImageKey(imageKey string) *MessageCardImage {
	image.ImageKey = &imageKey
	return image
}

func (image *MessageCardImage) WithAlt(alt *MessageCardPlainText) *MessageCardImage {
	image.Alt = alt
	return image
}

func (image *MessageCardImage) WithTitle(title MessageCardText) *MessageCardImage {
	image.Title = title
	return image
}

func (image *MessageCardImage) WithCustomWidth(customWidth int) *MessageCardImage {
	image.CustomWidth = &customWidth
	return image
}

func (image *MessageCardImage) WithCompactWidth(compactWidth bool) *MessageCardImage {
	image.CompactWidth = &compactWidth
	return image
}

func (image *MessageCardImage) WithMode(mode MessageCardImageMode) *MessageCardImage {
	image.Mode = &mode
	return image
}

func (image *MessageCardImage) WithPreview(preview bool) *MessageCardImage {
	image.Preview = &preview
	return image
}

func (image *MessageCardImage) Build() *MessageCardImage {
	return image
}

func (image *MessageCardImage) Tag() string {
	return "img"
}

func (image *MessageCardImage) MarshalJSON() ([]byte, error) {
	if image.ImageKey == nil {
		return nil, errors.New("image key is required")
	}
	if image.Alt == nil {
		return nil, errors.New("alt is required")
	}
	return messageCardElementJSON(image)
}

func (image *MessageCardImage) IsNote() {}

func (image *MessageCardImage) IsExtra() {}

type MessageCardNote struct {
	Elements []MessageCardNoteElement `json:"elements,omitempty"`
}

func NewMessageCardNote() *MessageCardNote {
	return &MessageCardNote{}
}

func (note *MessageCardNote) WithElements(elements []MessageCardNoteElement) *MessageCardNote {
	note.Elements = elements
	return note
}

func (note *MessageCardNote) Build() *MessageCardNote {
	return note
}

func (note *MessageCardNote) Tag() string {
	return "note"
}

func (note *MessageCardNote) MarshalJSON() ([]byte, error) {
	return messageCardElementJSON(note)
}

type MessageCardNoteElement interface {
	MessageCardElement
	IsNote()
}

type MessageCardActionLayout string

const (
	LayoutBisected   MessageCardActionLayout = "bisected"
	LayoutTrisection MessageCardActionLayout = "trisection"
	LayoutFlow       MessageCardActionLayout = "flow"
)

type MessageCardAction struct {
	Actions []MessageCardActionElement `json:"actions,omitempty"`
	Layout  *MessageCardActionLayout   `json:"layout,omitempty"`
}

func NewMessageCardAction() *MessageCardAction {
	return &MessageCardAction{}
}

func (action *MessageCardAction) WithActions(actions []MessageCardActionElement) *MessageCardAction {
	action.Actions = actions
	return action
}

func (action *MessageCardAction) WithLayout(layout MessageCardActionLayout) *MessageCardAction {
	action.Layout = &layout
	return action
}

func (action *MessageCardAction) Build() *MessageCardAction {
	return action
}

func (action *MessageCardAction) Tag() string {
	return "action"
}

func (action *MessageCardAction) MarshalJSON() ([]byte, error) {
	if len(action.Actions) == 0 {
		return nil, errors.New("actions is required")
	}
	return messageCardElementJSON(action)
}

type MessageCardActionElement interface {
	MessageCardElement
	IsAction()
}

type MessageCardDatePickerBase struct {
	// format "yyyy-MM-dd"
	InitialDate *string `json:"initial_date,omitempty"`
	// format "HH:mm"
	InitialTime *string `json:"initial_time,omitempty"`
	// format "yyyy-MM-dd HH:mm"
	InitialDateTime *string                `json:"initial_datetime,omitempty"`
	PlaceHolder     *MessageCardPlainText  `json:"placeholder,omitempty"`
	Value           map[string]interface{} `json:"value,omitempty"`
	Confirm         *MessageCardConfirm    `json:"confirm,omitempty"`
}

func NewMessageCardDatePickerBase() *MessageCardDatePickerBase {
	return &MessageCardDatePickerBase{}
}

func (datePickerBase *MessageCardDatePickerBase) WithInitialDate(initialDate string) *MessageCardDatePickerBase {
	datePickerBase.InitialDate = &initialDate
	return datePickerBase
}

func (datePickerBase *MessageCardDatePickerBase) WithInitialTime(initialTime string) *MessageCardDatePickerBase {
	datePickerBase.InitialTime = &initialTime
	return datePickerBase
}

func (datePickerBase *MessageCardDatePickerBase) WithInitialDateTime(initialDateTime string) *MessageCardDatePickerBase {
	datePickerBase.InitialDateTime = &initialDateTime
	return datePickerBase
}

func (datePickerBase *MessageCardDatePickerBase) WithPlaceHolder(placeHolder *MessageCardPlainText) *MessageCardDatePickerBase {
	datePickerBase.PlaceHolder = placeHolder
	return datePickerBase
}

func (datePickerBase *MessageCardDatePickerBase) WithValue(value map[string]interface{}) *MessageCardDatePickerBase {
	datePickerBase.Value = value
	return datePickerBase
}

func (datePickerBase *MessageCardDatePickerBase) WithConfirm(confirm *MessageCardConfirm) *MessageCardDatePickerBase {
	datePickerBase.Confirm = confirm
	return datePickerBase
}

func (datePickerBase *MessageCardDatePickerBase) Build() *MessageCardDatePickerBase {
	return datePickerBase
}

func (datePickerBase *MessageCardDatePickerBase) IsAction() {}

func (datePickerBase *MessageCardDatePickerBase) IsExtra() {}

type MessageCardDatePicker struct {
	*MessageCardDatePickerBase
}

func NewMessageCardDatePicker() *MessageCardDatePicker {
	return &MessageCardDatePicker{}
}

func (datePicker *MessageCardDatePicker) WithMessageCardDatePickerBase(datePickerBase *MessageCardDatePickerBase) *MessageCardDatePicker {
	datePicker.MessageCardDatePickerBase = datePickerBase
	return datePicker
}

func (datePicker *MessageCardDatePicker) Tag() string {
	return "date_picker"
}

func (datePicker *MessageCardDatePicker) MarshalJSON() ([]byte, error) {
	if datePicker.InitialDate == nil && datePicker.PlaceHolder == nil {
		return nil, errors.New("initial_date or placeholder is required")
	}
	return messageCardElementJSON(datePicker)
}

type MessageCardPickerTime struct {
	*MessageCardDatePickerBase
}

func NewMessageCardPickerTime() *MessageCardPickerTime {
	return &MessageCardPickerTime{}
}

func (pickerTime *MessageCardPickerTime) WithMessageCardDatePickerBase(datePickerBase *MessageCardDatePickerBase) *MessageCardPickerTime {
	pickerTime.MessageCardDatePickerBase = datePickerBase
	return pickerTime
}

func (pickerTime *MessageCardPickerTime) Tag() string {
	return "picker_time"
}

func (pickerTime *MessageCardPickerTime) MarshalJSON() ([]byte, error) {
	if pickerTime.InitialTime == nil && pickerTime.PlaceHolder == nil {
		return nil, errors.New("initial_time or placeholder is required")
	}
	return messageCardElementJSON(pickerTime)
}

type MessageCardPickerDateTime struct {
	*MessageCardDatePickerBase
}

func NewMessageCardPickerDateTime() *MessageCardPickerDateTime {
	return &MessageCardPickerDateTime{}
}

func (pickerDateTime *MessageCardPickerDateTime) WithMessageCardDatePickerBase(datePickerBase *MessageCardDatePickerBase) *MessageCardPickerDateTime {
	pickerDateTime.MessageCardDatePickerBase = datePickerBase
	return pickerDateTime
}

func (pickerDateTime *MessageCardPickerDateTime) Tag() string {
	return "picker_datetime"
}

func (pickerDateTime *MessageCardPickerDateTime) MarshalJSON() ([]byte, error) {
	if pickerDateTime.InitialDateTime == nil && pickerDateTime.PlaceHolder == nil {
		return nil, errors.New("initial_datetime or placeholder is required")
	}
	return messageCardElementJSON(pickerDateTime)
}

type MessageCardConfirm struct {
	Title *MessageCardPlainText `json:"title,omitempty"`
	Text  *MessageCardPlainText `json:"text,omitempty"`
}

func NewMessageCardConfirm() *MessageCardConfirm {
	return &MessageCardConfirm{}
}

func (confirm *MessageCardConfirm) WithTitle(title *MessageCardPlainText) *MessageCardConfirm {
	confirm.Title = title
	return confirm
}

func (confirm *MessageCardConfirm) WithText(text *MessageCardPlainText) *MessageCardConfirm {
	confirm.Text = text
	return confirm
}

func (confirm *MessageCardConfirm) Build() *MessageCardConfirm {
	return confirm
}

func (confirm *MessageCardConfirm) MarshalJSON() ([]byte, error) {
	if confirm.Title == nil {
		return nil, errors.New("title is required")
	}
	if confirm.Text == nil {
		return nil, errors.New("text is required")
	}
	return json.Marshal(confirm)
}

type MessageCardOverflow struct {
	Options []MessageCardOption    `json:"options,omitempty"`
	Value   map[string]interface{} `json:"value,omitempty"`
	Confirm *MessageCardConfirm    `json:"confirm,omitempty"`
}

func NewMessageCardOverflow() *MessageCardOverflow {
	return &MessageCardOverflow{}
}

func (overflow *MessageCardOverflow) WithOptions(options []MessageCardOption) *MessageCardOverflow {
	overflow.Options = options
	return overflow
}

func (overflow *MessageCardOverflow) WithValue(value map[string]interface{}) *MessageCardOverflow {
	overflow.Value = value
	return overflow
}

func (overflow *MessageCardOverflow) WithConfirm(confirm *MessageCardConfirm) *MessageCardOverflow {
	overflow.Confirm = confirm
	return overflow
}

func (overflow *MessageCardOverflow) Build() *MessageCardOverflow {
	return overflow
}

func (overflow *MessageCardOverflow) Tag() string {
	return "overflow"
}

func (overflow *MessageCardOverflow) MarshalJSON() ([]byte, error) {
	if len(overflow.Options) == 0 {
		return nil, errors.New("options is required")
	}
	for _, option := range overflow.Options {
		if option.Text == nil {
			return nil, errors.New("text is required")
		}
	}
	return messageCardElementJSON(overflow)
}

func (overflow *MessageCardOverflow) IsAction() {}

func (overflow *MessageCardOverflow) IsExtra() {}

type MessageCardOption struct {
	Text     *MessageCardPlainText `json:"text,omitempty"`
	Value    *string               `json:"value,omitempty"`
	URL      *string               `json:"url,omitempty"`
	MultiURL *MessageCardURL       `json:"multi_url,omitempty"`
}

func NewMessageCardOption() *MessageCardOption {
	return &MessageCardOption{}
}

func (option *MessageCardOption) WithText(text *MessageCardPlainText) *MessageCardOption {
	option.Text = text
	return option
}

func (option *MessageCardOption) WithValue(value string) *MessageCardOption {
	option.Value = &value
	return option
}

func (option *MessageCardOption) WithURL(url string) *MessageCardOption {
	option.URL = &url
	return option
}

func (option *MessageCardOption) WithMultiURL(multiURL *MessageCardURL) *MessageCardOption {
	option.MultiURL = multiURL
	return option
}

func (option *MessageCardOption) Build() *MessageCardOption {
	return option
}

func (option *MessageCardOption) MarshalJSON() ([]byte, error) {
	if option.URL != nil && option.MultiURL != nil {
		return nil, errors.New("url and multi_url can not be set at the same time")
	}
	return json.Marshal(option)
}

type MessageCardSelectMenuBase struct {
	PlaceHolder   *MessageCardPlainText  `json:"placeholder,omitempty"`
	InitialOption *string                `json:"initial_option,omitempty"`
	Options       []MessageCardOption    `json:"options,omitempty"`
	Value         map[string]interface{} `json:"value,omitempty"`
	Confirm       *MessageCardConfirm    `json:"confirm,omitempty"`
}

func NewMessageCardSelectMenuBase() *MessageCardSelectMenuBase {
	return &MessageCardSelectMenuBase{}
}

func (selectMenuBase *MessageCardSelectMenuBase) WithPlaceHolder(placeHolder *MessageCardPlainText) *MessageCardSelectMenuBase {
	selectMenuBase.PlaceHolder = placeHolder
	return selectMenuBase
}

func (selectMenuBase *MessageCardSelectMenuBase) WithInitialOption(initialOption string) *MessageCardSelectMenuBase {
	selectMenuBase.InitialOption = &initialOption
	return selectMenuBase
}

func (selectMenuBase *MessageCardSelectMenuBase) WithOptions(options []MessageCardOption) *MessageCardSelectMenuBase {
	selectMenuBase.Options = options
	return selectMenuBase
}

func (selectMenuBase *MessageCardSelectMenuBase) WithValue(value map[string]interface{}) *MessageCardSelectMenuBase {
	selectMenuBase.Value = value
	return selectMenuBase
}

func (selectMenuBase *MessageCardSelectMenuBase) WithConfirm(confirm *MessageCardConfirm) *MessageCardSelectMenuBase {
	selectMenuBase.Confirm = confirm
	return selectMenuBase
}

func (selectMenuBase *MessageCardSelectMenuBase) Build() *MessageCardSelectMenuBase {
	return selectMenuBase
}

func (selectMenuBase *MessageCardSelectMenuBase) IsAction() {}

func (selectMenuBase *MessageCardSelectMenuBase) IsExtra() {}

type MessageCardSelectStatic struct {
	*MessageCardSelectMenuBase
}

func NewMessageCardSelectStatic() *MessageCardSelectStatic {
	return &MessageCardSelectStatic{}
}

func (selectStatic *MessageCardSelectStatic) WithMessageCardSelectMenuBase(selectMenuBase *MessageCardSelectMenuBase) *MessageCardSelectStatic {
	selectStatic.MessageCardSelectMenuBase = selectMenuBase
	return selectStatic
}

func (selectStatic *MessageCardSelectStatic) Tag() string {
	return "select_static"
}

func (selectStatic *MessageCardSelectStatic) MarshalJSON() ([]byte, error) {
	if selectStatic.InitialOption == nil && selectStatic.PlaceHolder == nil {
		return nil, errors.New("placeholder is required")
	}
	if len(selectStatic.Options) == 0 {
		return nil, errors.New("options is required")
	}
	for _, option := range selectStatic.Options {
		if option.Text == nil {
			return nil, errors.New("text is required")
		}
	}
	return messageCardElementJSON(selectStatic)
}

type MessageCardSelectPerson struct {
	*MessageCardSelectMenuBase
}

func NewMessageCardSelectPerson() *MessageCardSelectPerson {
	return &MessageCardSelectPerson{}
}

func (selectPerson *MessageCardSelectPerson) WithMessageCardSelectMenuBase(selectMenuBase *MessageCardSelectMenuBase) *MessageCardSelectPerson {
	selectPerson.MessageCardSelectMenuBase = selectMenuBase
	return selectPerson
}

func (selectPerson *MessageCardSelectPerson) Tag() string {
	return "select_person"
}

func (selectPerson *MessageCardSelectPerson) MarshalJSON() ([]byte, error) {
	if selectPerson.InitialOption == nil && selectPerson.PlaceHolder == nil {
		return nil, errors.New("placeholder is required")
	}
	if len(selectPerson.Options) == 0 {
		return nil, errors.New("options is required")
	}
	return messageCardElementJSON(selectPerson)
}

type MessageCardButtonType string

const (
	TypeDefault MessageCardButtonType = "default"
	TypePrimary MessageCardButtonType = "primary"
	TypeDanger  MessageCardButtonType = "danger"
)

type MessageCardButton struct {
	Text     MessageCardText        `json:"text,omitempty"`
	URL      *string                `json:"url,omitempty"`
	MultiURL *MessageCardURL        `json:"multi_url,omitempty"`
	Type     *MessageCardButtonType `json:"type,omitempty"`
	Value    map[string]interface{} `json:"value,omitempty"`
	Confirm  *MessageCardConfirm    `json:"confirm,omitempty"`
}

func NewMessageCardButton() *MessageCardButton {
	return &MessageCardButton{}
}

func (button *MessageCardButton) WithText(text MessageCardText) *MessageCardButton {
	button.Text = text
	return button
}

func (button *MessageCardButton) WithURL(url string) *MessageCardButton {
	button.URL = &url
	return button
}

func (button *MessageCardButton) WithMultiURL(multiURL *MessageCardURL) *MessageCardButton {
	button.MultiURL = multiURL
	return button
}

func (button *MessageCardButton) WithType(buttonType MessageCardButtonType) *MessageCardButton {
	button.Type = &buttonType
	return button
}

func (button *MessageCardButton) WithValue(value map[string]interface{}) *MessageCardButton {
	button.Value = value
	return button
}

func (button *MessageCardButton) WithConfirm(confirm *MessageCardConfirm) *MessageCardButton {
	button.Confirm = confirm
	return button
}

func (button *MessageCardButton) Build() *MessageCardButton {
	return button
}

func (button *MessageCardButton) Tag() string {
	return "button"
}

func (button *MessageCardButton) MarshalJSON() ([]byte, error) {
	if button.Text == nil {
		return nil, errors.New("text is required")
	}
	if button.URL != nil && button.MultiURL != nil {
		return nil, errors.New("url and multi_url can not be set at the same time")
	}
	return messageCardElementJSON(button)
}

func (button *MessageCardButton) IsAction() {}

func (button *MessageCardButton) IsExtra() {}

type MessageCardLink struct {
	URL        *string `json:"url,omitempty"`
	AndroidURL *string `json:"android_url,omitempty"`
	IOSURL     *string `json:"ios_url,omitempty"`
	PCURL      *string `json:"pc_url,omitempty"`
}

func NewCardLink() *MessageCardLink {
	return &MessageCardLink{}
}

func (cardLink *MessageCardLink) WithURL(url string) *MessageCardLink {
	cardLink.URL = &url
	return cardLink
}

func (cardLink *MessageCardLink) WithAndroidURL(androidURL string) *MessageCardLink {
	cardLink.AndroidURL = &androidURL
	return cardLink
}

func (cardLink *MessageCardLink) WithIOSURL(iosURL string) *MessageCardLink {
	cardLink.IOSURL = &iosURL
	return cardLink
}

func (cardLink *MessageCardLink) WithPCURL(pcURL string) *MessageCardLink {
	cardLink.PCURL = &pcURL
	return cardLink
}

func (cardLink *MessageCardLink) Build() *MessageCardLink {
	return cardLink
}

func (cardLink *MessageCardLink) MarshalJSON() ([]byte, error) {
	if cardLink.URL == nil {
		return nil, errors.New("url is required")
	}
	return json.Marshal(cardLink)
}

func struct2mapByReflect(val interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	s := reflect.Indirect(reflect.ValueOf(val))
	st := s.Type()
	for i := 0; i < s.NumField(); i++ {
		fieldDesc := st.Field(i)
		fieldVal := s.Field(i)
		if fieldDesc.Anonymous {
			embeddedMap, err := struct2mapByReflect(fieldVal.Interface())
			if err != nil {
				return nil, err
			}
			for k, v := range embeddedMap {
				m[k] = v
			}
			continue
		}
		jsonTag := fieldDesc.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		tag, err := parseJSONTag(jsonTag)
		if err != nil {
			return nil, err
		}
		if tag.ignore {
			continue
		}
		if fieldDesc.Type.Kind() == reflect.Ptr && fieldVal.IsNil() {
			continue
		}
		// nil maps are treated as empty maps.
		if fieldDesc.Type.Kind() == reflect.Map && fieldVal.IsNil() {
			continue
		}
		if fieldDesc.Type.Kind() == reflect.Slice && fieldVal.IsNil() {
			continue
		}
		if tag.stringFormat {
			m[tag.name] = formatAsString(fieldVal, fieldDesc.Type.Kind())
		} else {
			m[tag.name] = fieldVal.Interface()
		}
	}
	return m, nil
}

func formatAsString(v reflect.Value, kind reflect.Kind) string {
	if kind == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return fmt.Sprintf("%v", v.Interface())
}

type jsonTag struct {
	name         string
	stringFormat bool
	ignore       bool
}

func parseJSONTag(val string) (jsonTag, error) {
	if val == "-" {
		return jsonTag{ignore: true}, nil
	}
	var tag jsonTag
	i := strings.Index(val, ",")
	if i == -1 || val[:i] == "" {
		return tag, fmt.Errorf("malformed json tag: %s", val)
	}
	tag = jsonTag{
		name: val[:i],
	}
	switch val[i+1:] {
	case "omitempty":
	case "omitempty,string":
		tag.stringFormat = true
	default:
		return tag, fmt.Errorf("malformed json tag: %s", val)
	}
	return tag, nil
}
