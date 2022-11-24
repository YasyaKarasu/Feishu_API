package feishuapi

import (
	"io"
	"net/http"
	"strings"
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
func (c AppClient) GetAllBitables(DocumentId string) []BitableInfo {
	var all_bitables []BitableInfo

	l := c.GetAllPages("get", "open-apis/docx/v1/documents/"+DocumentId+"/blocks", nil, nil, nil, 100)

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
func (c AppClient) GetAllTables(AppToken string) []TableInfo {
	var all_tables []TableInfo

	l := c.GetAllPages("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables", nil, nil, nil, 100)

	for _, value := range l {
		all_tables = append(all_tables, *NewTableInfo(AppToken, value.(map[string]interface{})))
	}

	return all_tables
}

type RecordInfo struct {
	AppToken string
	TableId  string
	RecordId string
	Fields   map[string]interface{}
}

// Create a new RecordInfo
func NewRecordInfo(AppToken string, TableId string, data map[string]interface{}) *RecordInfo {
	return &RecordInfo{
		AppToken: AppToken,
		TableId:  TableId,
		RecordId: data["record_id"].(string),
		Fields:   data["fields"].(map[string]interface{}),
	}
}

// Get all Records by AppToken and TableId
func (c AppClient) GetAllRecords(AppToken string, TableId string) []RecordInfo {
	var all_records []RecordInfo

	l := c.GetAllPages("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records", nil, nil, nil, 100)

	for _, value := range l {
		all_records = append(all_records, *NewRecordInfo(AppToken, TableId, value.(map[string]interface{})))
	}

	return all_records
}

// Get A Record by AppToken, TableId and RecordId
func (c AppClient) GetRecord(AppToken string, TableId string, RecordId string) *RecordInfo {
	record := c.Request("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records/"+RecordId, nil, nil, nil)

	return NewRecordInfo(AppToken, TableId, record)
}

// Get a []Byte form Record by AppToken, TableId and RecordId
func (c AppClient) GetRecordInByte(AppToken string, TableId string, RecordId string) []byte {
	client := &http.Client{}
	url := "https://open.feishu.cn/open-apis/bitable/v1/apps/" + AppToken + "/tables/" + TableId + "/records/" + RecordId
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+c._tenant_access_token)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func (c AppClient) GetRawContent(DocumentId string) string {
	resp := c.Request("get", "open-apis/docx/v1/documents/"+DocumentId+"/raw_content", nil, nil, nil)
	content := resp["content"].(string)

	return content
}
