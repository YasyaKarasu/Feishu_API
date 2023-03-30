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
func NewBitableInfo(data map[string]interface{}) *BitableInfo {
	bitInfo := data["bitable"].(map[string]interface{})
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
		info := value.(map[string]interface{})
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
func NewTableInfo(AppToken string, data map[string]interface{}) *TableInfo {
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
		all_tables = append(all_tables, *NewTableInfo(AppToken, value.(map[string]interface{})))
	}

	return all_tables
}

type RecordInfo struct {
	AppToken         string
	TableId          string
	RecordId         string
	LastModifiedTime int
	Fields           map[string]interface{}
}

// Create a new RecordInfo
func NewRecordInfo(AppToken string, TableId string, data map[string]interface{}) *RecordInfo {
	return &RecordInfo{
		AppToken:         AppToken,
		TableId:          TableId,
		RecordId:         data["record_id"].(string),
		LastModifiedTime: int(data["last_modified_time"].(float64)),
		Fields:           data["fields"].(map[string]interface{}),
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
		all_records = append(all_records, *NewRecordInfo(AppToken, TableId, value.(map[string]interface{})))
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
