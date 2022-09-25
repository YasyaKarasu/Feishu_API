package feishuapi

import "strings"

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
	TableId string
	Name    string
}

// Create a new TableInfo
func NewTableInfo(data map[string]interface{}) *TableInfo {
	return &TableInfo{
		TableId: data["table_id"].(string),
		Name:    data["name"].(string),
	}
}

// Get all tables by AppToken
func (c AppClient) GetAllTables(AppToken string) []TableInfo {
	var all_tables []TableInfo

	l := c.GetAllPages("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables", nil, nil, nil, 100)

	for _, value := range l {
		all_tables = append(all_tables, *NewTableInfo(value.(map[string]interface{})))
	}

	return all_tables
}

type RecordInfo struct {
	RecordId string
	Fields   map[string]interface{}
}

// Create a new RecordInfo
func NewRecordInfo(data map[string]interface{}) *RecordInfo {
	return &RecordInfo{
		RecordId: data["record_id"].(string),
		Fields:   data["fields"].(map[string]interface{}),
	}
}

// Get all Records by AppToken and TableId
func (c AppClient) GetAllRecords(AppToken string, TableId string) []RecordInfo {
	var all_records []RecordInfo

	l := c.GetAllPages("get", "open-apis/bitable/v1/apps/"+AppToken+"/tables/"+TableId+"/records", nil, nil, nil, 100)

	for _, value := range l {
		all_records = append(all_records, *NewRecordInfo(value.(map[string]interface{})))
	}

	return all_records
}
