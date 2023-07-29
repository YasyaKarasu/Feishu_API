package feishuapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type BitableInfo struct {
	BlockId  string
	AppToken string
}

// Create a new BitableInfo
func NewBitableInfo(data map[string]any) *BitableInfo {
	bitInfo := data["bitable"].(map[string]any)
	token := bitInfo["token"].(string)
	splitToken := strings.Split(token, "_")
	return &BitableInfo{
		BlockId:  data["block_id"].(string),
		AppToken: splitToken[0],
	}
}

// Get all Bitables in a Document
func (c AppClient) DocumentGetAllBitables(DocumentId string) []BitableInfo {
	var all_bitables []BitableInfo

	l := c.GetAllPages("get", "open-apis/docx/v1/documents/"+DocumentId+"/blocks", nil, nil, nil, 100)

	if l == nil {
		logrus.WithField("DocumentID", DocumentId).Warn("nil bitable info return")
		return nil
	}

	for _, value := range l {
		info := value.(map[string]any)
		if info["block_type"].(float64) == 18 {
			all_bitables = append(all_bitables, *NewBitableInfo(info))
		}
	}

	return all_bitables
}

type TableInfo struct {
	AppToken string
	TableId  string
	Revision int
	Name     string
}

// Create a new TableInfo
func NewTableInfo(AppToken string, data map[string]any) *TableInfo {
	return &TableInfo{
		AppToken: AppToken,
		TableId:  data["table_id"].(string),
		Revision: int(data["revision"].(float64)),
		Name:     data["name"].(string),
	}
}

// Get all tables by AppToken
func (c AppClient) DocumentGetAllTables(AppToken string) []TableInfo {
	var all_tables []TableInfo

	l := c.GetAllPages("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables", nil, nil, nil, 100)

	if l == nil {
		logrus.WithField("AppToken", AppToken).Warn("nil table info return")
		return nil
	}

	for _, value := range l {
		all_tables = append(all_tables, *NewTableInfo(AppToken, value.(map[string]any)))
	}

	return all_tables
}

type RecordInfo struct {
	AppToken         string
	TableId          string
	RecordId         string
	LastModifiedTime int
	Fields           map[string]any
}

// Create a new RecordInfo
func NewRecordInfo(AppToken string, TableId string, data map[string]any) *RecordInfo {
	return &RecordInfo{
		AppToken:         AppToken,
		TableId:          TableId,
		RecordId:         data["record_id"].(string),
		LastModifiedTime: int(data["last_modified_time"].(float64)),
		Fields:           data["fields"].(map[string]any),
	}
}

func NewRecordInfoWithoutModifiedTime(AppToken string, TableId string, data map[string]any) *RecordInfo {
	return &RecordInfo{
		AppToken: AppToken,
		TableId:  TableId,
		RecordId: data["record_id"].(string),
		Fields:   data["fields"].(map[string]any),
	}
}

// Get all Records by AppToken and TableId
func (c AppClient) DocumentGetAllRecords(AppToken string, TableId string) []RecordInfo {
	var all_records []RecordInfo

	query := make(map[string]any)
	query["automatic_fields"] = "true"
	l := c.GetAllPages("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records", query, nil, nil, 100)

	if l == nil {
		logrus.WithFields(logrus.Fields{
			"AppToken": AppToken,
			"TableID":  TableId,
		}).Warn("nil record info return")
		return nil
	}

	for _, value := range l {
		all_records = append(all_records, *NewRecordInfo(AppToken, TableId, value.(map[string]any)))
	}

	return all_records
}

// DocumentGetAllRecordsWithLinks retrieves all records from a specified table in the Bitable app.
// It provides an option to include multi-line text fields as []map containing hyperlinks' URLs if available.
// The purpose of this method is to enable access to hyperlinks within multi-line text fields for specific use cases.
// Note: Adding the query parameter "text_field_as_array" will change the format of multi-line text fields
// from string to []map.
func (c AppClient) DocumentGetAllRecordsWithLinks(AppToken string, TableId string) []RecordInfo {
	var all_records []RecordInfo

	query := make(map[string]any)
	query["automatic_fields"] = "true"
	query["text_field_as_array"] = "true" // Include this query parameter to get multi-line text fields as []map with hyperlinks.

	l := c.GetAllPages("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records", query, nil, nil, 100)

	if l == nil {
		logrus.WithFields(logrus.Fields{
			"AppToken": AppToken,
			"TableID":  TableId,
		}).Warn("nil record info return")
		return nil
	}

	for _, value := range l {
		all_records = append(all_records, *NewRecordInfo(AppToken, TableId, value.(map[string]any)))
	}

	return all_records
}

// Get A Record by AppToken, TableId and RecordId
func (c AppClient) DocumentGetRecord(AppToken string, TableId string, RecordId string) *RecordInfo {
	query := make(map[string]any)
	query["automatic_fields"] = "true"
	record := c.Request("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records/"+RecordId, query, nil, nil)

	if record == nil {
		logrus.WithFields(logrus.Fields{
			"AppToken": AppToken,
			"TableID":  TableId,
			"RecordID": RecordId,
		}).Warn("nil record info return")
		return nil
	}

	return NewRecordInfo(AppToken, TableId, record["record"].(map[string]any))
}

func (c AppClient) DocumentGetRecordWithoutModifiedTime(AppToken string, TableId string, RecordId string) *RecordInfo {
	query := make(map[string]any)
	query["automatic_fields"] = "true"
	record := c.Request("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records/"+RecordId, query, nil, nil)

	if record == nil {
		logrus.WithFields(logrus.Fields{
			"AppToken": AppToken,
			"TableID":  TableId,
			"RecordID": RecordId,
		}).Warn("nil record info return")
		return nil
	}

	return NewRecordInfoWithoutModifiedTime(AppToken, TableId, record["record"].(map[string]any))
}

// Get a []Byte form Record by AppToken, TableId and RecordId
func (c AppClient) DocumentGetRecordInByte(AppToken string, TableId string, RecordId string) []byte {
	client := &http.Client{}
	url := "https://open.feishu.cn/open-apis/bitable/v1/apps/" + AppToken + "/tables/" + TableId + "/records/" + RecordId
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		logrus.Error(err)
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+c._tenant_access_token)

	res, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return body
}

type FieldStaff struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Create a record in bitable
func (c AppClient) DocumentCreateRecord(AppToken string, TableId string, Fields map[string]any) *RecordInfo {
	body := make(map[string]any)
	body["fields"] = Fields

	resp := c.Request("post", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records", nil, nil, body)
	logrus.Debug("Created a record: ", resp)

	return NewRecordInfoWithoutModifiedTime(AppToken, TableId, resp["record"].(map[string]any))
}

// Update a record in bitable
func (c AppClient) DocumentUpdateRecord(AppToken string, TableId string, RecordId string, Fields map[string]any) bool {
	body := make(map[string]any)
	body["fields"] = Fields

	resp := c.Request("put", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records/"+RecordId, nil, nil, body)
	logrus.Debug("Updated record "+RecordId+": ", resp)

	return true
}

// Deleta records in bitable
func (c AppClient) DocumentDeleteRecords(AppToken string, TableId string, RecordIds []string) bool {
	body := make(map[string]any)
	body["records"] = RecordIds

	resp := c.Request("put", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records/batch_delete", nil, nil, body)
	logrus.Debug("Deleted records : ", RecordIds, resp)

	return true
}

type TextStyle struct {
	Align    int  `json:"align"`
	Done     bool `json:"done"`
	Folded   bool `json:"folded"`
	Language int  `json:"language"`
	Wrap     bool `json:"wrap"`
}

type TextElement struct {
	TextRun *struct {
		Content          *string `json:"content,omitempty"`
		TextElementStyle *struct {
			Bold            *bool `json:"bold,omitempty"`
			Italic          *bool `json:"italic,omitempty"`
			Strikethrough   *bool `json:"strikethrough,omitempty"`
			Underline       *bool `json:"underline,omitempty"`
			InlineCode      *bool `json:"inline_code,omitempty"`
			BackgroundColor *int  `json:"background_color,omitempty"`
			TextColor       *int  `json:"text_color,omitempty"`
			Link            *struct {
				URL *string `json:"url,omitempty"`
			} `json:"link,omitempty"`
		} `json:"text_element_style,omitempty"`
	} `json:"text_run,omitempty"`
	MentionUser *struct {
		UserID *string `json:"user_id,omitempty"`
	} `json:"mention_user,omitempty"`
	MentionDoc *struct {
		Token   *string `json:"token,omitempty"`
		ObjType *int    `json:"obj_type,omitempty"`
		URL     *string `json:"url,omitempty"`
		Title   *string `json:"title,omitempty"`
	} `json:"mention_doc,omitempty"`
}

type BlockText struct {
	Style    *TextStyle    `json:"style,omitempty"`
	Elements []TextElement `json:"elements,omitempty"`
}

type BlockTextElementsUpdate struct {
	Elements []TextElement `json:"elements,omitempty"`
}

type BlockISV struct {
	ComponentID     string `json:"component_id"`
	ComponentTypeID string `json:"component_type_id"`
}

type BlockCreate struct {
	BlockType int        `json:"block_type"`
	BlockText *BlockText `json:"text,omitempty"`
	BlockISV  *BlockISV  `json:"isv,omitempty"`
}

type BlockUpdate struct {
	UpdateTextElements *BlockTextElementsUpdate `json:"update_text_elements,omitempty"`
}

type BlockInfo struct {
	BlockId   string     `json:"block_id"`
	BlockType int        `json:"block_type"`
	Text      *BlockText `json:"text"`
}

func (c AppClient) DocumentGetAllBlocks(DocumentId string, userIdType UserIdType) []BlockInfo {
	query := make(map[string]any)
	query["user_id_type"] = string(userIdType)

	l := c.GetAllPages("get", "open-apis/docx/v1/documents/"+DocumentId+"/blocks", query, nil, nil, 100)
	b, _ := json.Marshal(l)
	blocks := make([]BlockInfo, 0)
	json.Unmarshal(b, &blocks)
	return blocks
}

func (c AppClient) DocumentCreateBlock(DocumentId string, BlockId string, userIdType UserIdType, blocks []BlockCreate, index int) {
	query := make(map[string]any)
	query["user_id_type"] = string(userIdType)

	body := make(map[string]any)
	body["children"] = blocks
	body["index"] = index

	c.Request("post", "open-apis/docx/v1/documents/"+DocumentId+"/blocks/"+BlockId+"/children", query, nil, body)
}

func (c AppClient) DocumentUpdateBlock(DocumentId string, BlockId string, userIdType UserIdType, update *BlockUpdate) {
	query := make(map[string]any)
	query["user_id_type"] = string(userIdType)

	body := make(map[string]any)
	struct2map(update, &body)

	c.Request("patch", "open-apis/docx/v1/documents/"+DocumentId+"/blocks/"+BlockId, query, nil, body)
}

func (c AppClient) DocumentGetRawContent(DocumentId string) string {
	resp := c.Request("get", "open-apis/docx/v1/documents/"+DocumentId+"/raw_content", nil, nil, nil)
	content := resp["content"].(string)

	return content
}

// Append data to a sheet in a spreadsheet, return the actual range of the appended data
func (c AppClient) SheetAppendData(SpreadSheetToken string, SheetId string, Range string, Data [][]interface{}) string {
	body := make(map[string]interface{})
	body["range"] = SheetId + "!" + Range
	body["values"] = Data
	abody := make(map[string]interface{})
	abody["valueRange"] = body

	resp := c.Request("post", "open-apis/sheet/v2/spreadsheets/"+SpreadSheetToken+"/values_append", nil, nil, abody)
	return resp["tableRange"].(string)
}

func (c AppClient) SheetGetData(SpreadSheetToken string, SheetId string, Range string) []interface{} {
	resp := c.Request("get", "open-apis/sheet/v2/spreadsheets/"+SpreadSheetToken+"/values/"+SheetId+"!"+Range, nil, nil, nil)
	return resp["valueRange"].(map[string]interface{})["values"].([]interface{})
}

func (c AppClient) SheetWriteData(SpreadSheetToken string, SheetId string, Range string, Data [][]interface{}) {
	body := make(map[string]interface{})
	body["range"] = SheetId + "!" + Range
	body["values"] = Data
	abody := make(map[string]interface{})
	abody["valueRange"] = body

	c.Request("put", "open-apis/sheet/v2/spreadsheets/"+SpreadSheetToken+"/values", nil, nil, abody)
}
