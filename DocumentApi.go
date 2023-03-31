package feishuapi

import (
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

// Get all Records by AppToken and TableId
func (c AppClient) DocumentGetAllRecords(AppToken string, TableId string) []RecordInfo {
	var all_records []RecordInfo

	query := make(map[string]string)
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

// Get A Record by AppToken, TableId and RecordId
func (c AppClient) DocumentGetRecord(AppToken string, TableId string, RecordId string) *RecordInfo {
	query := make(map[string]string)
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

	return NewRecordInfo(AppToken, TableId, record)
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

func NewCreatedRecordInfo(AppToken string, TableId string, data map[string]any) *RecordInfo {
	return &RecordInfo{
		AppToken: AppToken,
		TableId:  TableId,
		RecordId: data["record_id"].(string),
		Fields:   data["fields"].(map[string]any),
	}
}

// Create a record in bitable
func (c AppClient) DocumentCreateRecord(AppToken string, TableId string, Fields map[string]any) *RecordInfo {
	body := make(map[string]any)
	body["fields"] = Fields

	resp := c.Request("post", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records", nil, nil, body)
	logrus.Debug("Created a record: ", resp)

	return NewCreatedRecordInfo(AppToken, TableId, resp["record"].(map[string]any))
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

func (c AppClient) DocumentGetRawContent(DocumentId string) string {
	resp := c.Request("get", "open-apis/docx/v1/documents/"+DocumentId+"/raw_content", nil, nil, nil)
	content := resp["content"].(string)

	return content
}
