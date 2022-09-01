package feishuapi

type EmployeeType int

const (
	FullTime    EmployeeType = 1
	Internship  EmployeeType = 2
	Consultant  EmployeeType = 3
	OutSourcing EmployeeType = 4
	Laboring    EmployeeType = 5
)

type EmployeeStatus int

const (
	BeforeJob   EmployeeStatus = 1
	AtJob       EmployeeStatus = 2
	DenyJob     EmployeeStatus = 3
	BeforeLeave EmployeeStatus = 4
	Left        EmployeeStatus = 5
)

type EmployeeInfo struct {
	Id           string
	Name         string
	DepartmentId string
	EmployeeType EmployeeType
	Status       EmployeeStatus
}

// Create a new EmployeeInfo
func NewEmployeeInfo(data map[string]interface{}) *EmployeeInfo {
	sf := data["system_fields"].(map[string]interface{})
	return &EmployeeInfo{
		Id:           data["user_id"].(string),
		Name:         sf["name"].(string),
		DepartmentId: GetInMap(sf, "department_id", "").(string),
		EmployeeType: EmployeeType(int(sf["employee_type"].(float64))),
		Status:       EmployeeStatus(int(sf["status"].(float64))),
	}
}

type UserIdType string

const (
	OpenId  UserIdType = "open_id"
	UnionId UserIdType = "union_id"
	UserId  UserIdType = "user_id"
)

// Get all employees' information by specific user id type
func (c AppClient) GetAllEmployees(id_type UserIdType) []EmployeeInfo {
	query := make(map[string]string)
	query["user_id_type"] = string(id_type)
	l := c.GetAllPages("get", "open-apis/ehr/v1/employees", nil, nil, nil, 100)
	if l == nil {
		return nil
	}

	var all_employees []EmployeeInfo
	for _, value := range l {
		all_employees = append(all_employees, *NewEmployeeInfo(value.(map[string]interface{})))
	}
	return all_employees
}
