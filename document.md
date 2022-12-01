# FeishuAPI

## 代码架构

在项目文件夹下有平行的 $9$ 个 $go$ 文件，分别为 $ApiBase.go$ ，$DepartmentApi.go$ , $DocumentApi.go$ , $EmployeeApi.go$ , $GroupApi.go$ , $MessageApi.go$ , $RobotApi.go$ , $SpaceApi.go$ 。其中 $ApiBase.go$ 为包含了请求的发送、接受与解析等其他文件的基础，另外 8 个文件为平行关系，分别包含各自相关的飞书 API 。

## 如何使用

1. 在项目中导入此包
2. 在项目配置文件中配置好 FeishuAPI 相关配置
3. 在项目中创建一个 feishuapi AppClient 实例并设置配置字段
4. 调用 AppClient.StartTokenTimer 函数来获取 tenant access token

示例：

**main.go** :

```go
package main

import (
	"test/conf"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func main() {
	var cli feishuapi.AppClient

	conf.Init()
	conf.SetAppClientConfig(&cli)

	cli.StartTokenTimer()

	employee := cli.GetAllEmployees(feishuapi.UserId)
	logrus.Info(employee)
}
```

**conf.go** :

```go
package conf

import (
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic(err)
	}

	logrus.Info("Configuration file loaded")

	var confItems = map[string][]string{
		"feishu": {"APP_ID", "APP_SECRET", "VERIFICATION_TOKEN", "ENCRYPT_KEY"},
	}

	for k, v := range confItems {
		checkConfIsSet(k, v)
	}

	logrus.Info("All required values in configuration file are set")
}

func checkConfIsSet(name string, keys []string) {
	for i := range keys {
		wholeKey := name + "." + keys[i]
		if !viper.IsSet(wholeKey) {
			logrus.WithField(wholeKey, nil).
				Fatal("The following item of your configuration file hasn't been set properly: ")
		}
	}
}

func SetAppClientConfig(c *feishuapi.AppClient) {
	c.Conf.AppId = viper.GetString("feishu.APP_ID")
	c.Conf.AppSecret = viper.GetString("feishu.APP_SECRET")
	c.Conf.VerificationToken = viper.GetString("feishu.VERIFICATION_TOKEN")
	c.Conf.EncryptKey = viper.GetString("feishu.ENCRYPT_KEY")
}
```

## API 列表

### ApiBase.go

#### Types

##### type Config

```go
type Config struct {
	AppId             string
	AppSecret         string
	VerificationToken string
	EncryptKey        string
}
```

包含飞书应用信息的配置结构，具体信息可见 https://open.feishu.cn/app?lang=zh-CN。

##### type AppClient

```go
type AppClient struct {
	_tenant_access_token string
	schedular            *cron.Cron
	Conf                 Config
}
```

飞书 API 客户端，所有的 API 调用都要通过此类型进行，其中的 Conf 为相关配置信息，_tenant_access_token 为应用身份信息，schedular 为定时更新 token 的计时器。

#### Function

#### func (AppClient) StartTokenTimer

```go
func (c *AppClient) StartTokenTimer()
```

立即获取一次 tenant_access_token 并开启计时器，每 105 分钟更新一次。

#### func (AppClient) Request

```go
func (c AppClient) Request(method string, path string, query map[string]string, headers map[string]string, body interface{}) map[string]interface{}
```

向指定路径发送请求。请求字符串可以为 get, post 或 delete （不区分大小写），请求路径只需要 "https://open.feishu.cn/" 后面的部分，例如 "open-apis/im/v1/messages" 。header 中 Authorization 和 Content-Type 不需要传递（除非需要用 user_access_token 进行鉴权）。

返回值为 api 响应体的 data 部分。

#### func (AppClient) GetAllPages

```go
func (c AppClient) GetAllPages(method string, path string, query map[string]string, headers map[string]string, body interface{}, page_size int) []interface{}
```

GetAllPages 是对 Request 的封装，针对的是拉取整个列表，需要分页的飞书 API 。相比 Request ，GetAllPages 多了一个 page_size 参数，表示分页大小。

返回值为 api 响应体的 data/items 部分。

### DepartmentApi.go

#### Type

##### type DepartmentInfo

```go
type DepartmentInfo struct {
	Name        string
	GroupId     string
	MemberCount int
}
```

#### Function

##### func NewDepartmentInfo

```go
func NewDepartmentInfo(data map[string]interface{}) *DepartmentInfo
```

用于从 api 响应体 map 中提取信息，返回一个 DepartmentInfo 结构体指针。

##### func (AppClient) DepartmentGetInfoById

```go
func (c AppClient) DepartmentGetInfoById(DepartmentId string) *DepartmentInfo
```

通过 DepartmentId 获取 department 信息。

### DocumentApi.go

#### Type

##### type BitableInfo

```go
type BitableInfo struct {
	BlockId  string
	AppToken string
}
```

多维表格信息，具体可参考 https://open.feishu.cn/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/bitable-overview。

##### type TableInfo

```go
type TableInfo struct {
	AppToken string
	TableId  string
	Revision int
	Name     string
}
```

多维表格内的数据表信息。

##### type RecordInfo

```go
type RecordInfo struct {
	AppToken         string
	TableId          string
	RecordId         string
	LastModifiedTime int
	Fields           map[string]interface{}
}
```

数据表中的记录信息。

#### Function

##### func NewBitableInfo

```go
func NewBitableInfo(data map[string]interface{}) *BitableInfo
```

从 api 响应体 map 中提取信息，返回一个 BitableInfo 指针。

##### func (AppClient) DocumentGetAllBitables

```go
func (c AppClient) DocumentGetAllBitables(DocumentId string) []BitableInfo
```

根据 DocumentId 获取文档中的所有多维表格，返回值为 BitableInfo 切片。

##### func NewTableInfo

```go
func NewTableInfo(AppToken string, data map[string]interface{}) *TableInfo
```

从 api 响应体 map 中提取信息，返回一个 TableInfo 指针。

##### func (AppClient) DocumentGetAllTables

```go
func (c AppClient) DocumentGetAllTables(AppToken string) []TableInfo
```

根据多维表格 AppToken 获取其中的所有数据表，返回值为 TableInfo 切片。

##### func NewRecord

```go
func NewRecordInfo(AppToken string, TableId string, data map[string]interface{}) *RecordInfo
```

从 api 响应体 map 中提取信息，返回一个 RecordInfo 指针。

##### func (AppClient) DocumentGetAllRecords

```go
func (c AppClient) DocumentGetAllRecords(AppToken string, TableId string) []RecordInfo
```

根据多维表格 AppToken 和 数据表 TableId 获取表格内所有记录，返回值为 RecordInfo 切片。

##### func (AppClient) DocumentGetRecord

```go
func (c AppClient) DocumentGetRecord(AppToken string, TableId string, RecordId string) *RecordInfo
```

根据多维表格 AppToken, 数据表 TableId 和 记录 RecordId 获取记录信息，返回值为 RecordInfo 指针。

##### func (AppClient) DocumentGetRecordInByte

```go
func (c AppClient) DocumentGetRecordInByte(AppToken string, TableId string, RecordId string) []byte
```

获取一条 Record ，但是返回值为完整的返回体的 byte 切片

##### func (AppClient) DocumentGetRawContent

```go
func (c AppClient) DocumentGetRawContent(DocumentId string) string
```

根据 DocumentId 获取文档纯文本内容，返回值为字符串。

### EmployeeApi.go

#### Constants

```go
const (
	FullTime    EmployeeType = 1
	Internship  EmployeeType = 2
	Consultant  EmployeeType = 3
	OutSourcing EmployeeType = 4
	Laboring    EmployeeType = 5
)
```

在职职位。

```go
const (
	BeforeJob   EmployeeStatus = 1
	AtJob       EmployeeStatus = 2
	DenyJob     EmployeeStatus = 3
	BeforeLeave EmployeeStatus = 4
	Left        EmployeeStatus = 5
)
```

在职情况。

```go
const (
	OpenId  UserIdType = "open_id"
	UnionId UserIdType = "union_id"
	UserId  UserIdType = "user_id"
)
```

用户 ID 类型。

#### Type

##### type EmployeeInfo

```go
type EmployeeInfo struct {
	Id           string
	Name         string
	DepartmentId string
	EmployeeType EmployeeType
	Status       EmployeeStatus
}
```

表示一个职员的信息。

#### Function

##### func NewEmployeeInfo

```go
func NewEmployeeInfo(data map[string]interface{}) *EmployeeInfo
```

从 api 响应体map 中提取信息，返回值为 EmployeeInfo 指针。

##### func (AppClient) EmployeeGetAllInfo

```go
func (c AppClient) EmployeeGetAllInfo(id_type UserIdType) []EmployeeInfo
```

获取所有职员信息。返回的信息中的用户 ID 类型会根据输入参数决定。

### GroupApi.go

#### Type

##### type GroupInfo

```go
type GroupInfo struct {
	ChatId    string
	Name      string
	TenantKey string
}
```

群组信息，具体可参考 https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/im-v1/chat/chat-info/intro。

##### type GroupMember

```go
type GroupMember struct {
	MemberId string
	Name     string
}
```

群组内成员信息结构体。

#### Function

##### func NewGroupInfo

```go
func NewGroupInfo(data map[string]interface{}) *GroupInfo
```

从 api 返回体 map 中提取信息，返回值为 GroupInfo 指针。

##### func (AppClient) GroupGetAllInfo

```go
func (c AppClient) GroupGetAllInfo() []GroupInfo
```

获取所有群组信息，返回值为 GroupInfo 切片。

##### func NewGroupMember

```go
func NewGroupMember(data map[string]interface{}) *GroupMember
```

从 api 返回体 map 中提取信息，返回值为 GroupMember 指针。

##### func (AppClient) GroupGetMembers

```go
func (c AppClient) GroupGetMembers(groupId string, userIdType UserIdType) []GroupMember
```

根据 groupId 获取群组内成员信息，返回值为 GroupMember 切片。

##### func (AppClient) GroupCreate

```go
func (c AppClient) GroupCreate(groupName string, userIdType UserIdType, ownerId string) *GroupInfo
```

根据参数信息创建群组，返回值为 GroupInfo 指针。

##### func (AppClient) GroupGetInfo

```go
func (c AppClient) GroupGetInfo(chatId string) *GroupInfo
```

根据 chatId 查询群组信息，返回值为 GroupInfo 指针。

##### func (AppClient) GroupAddMembers

```go
func (c AppClient) GroupAddMembers(chatId string, memberIdType UserIdType, succeedType string, idList []string) bool
```

将给定用户加入群聊，具体参数可参考 https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat-members/create。返回值为 bool 变量，只要有用户添加出错就会返回 false ，具体出错可参考 log 。

##### func (AppClient) GroupDeleteMembers

```go
func (c AppClient) GroupDeleteMembers(chatId string, memberIdType UserIdType, idList []string) bool
```

将指定用户移除群聊，返回值为 bool 。

### MessageApi.go

#### Constant

```go
const (
	UserOpenId  MsgReceiverType = "open_id"
	UserUnionId MsgReceiverType = "union_id"
	UserUserId  MsgReceiverType = "user_id"
	UserEmail   MsgReceiverType = "email"
	GroupChatId MsgReceiverType = "chat_id"
)
```

消息接受者类型。

```go
const Text MsgContentType = "text"
```

消息类型为文本。

#### Function

##### func (AppClient) MessageSend

```go
func (c AppClient) MessageSend(receiveIdType MsgReceiverType, receiveId string, msgType MsgContentType, msg string) bool
```

发送消息到指定对象，返回值为 bool ，表示是否发送成功。

### RobotApi.go

#### Type

##### type RobotInfo

```go
type RobotInfo struct {
	Name   string
	OpenId string
}
```

表示机器人信息，具体可参考 https://open.feishu.cn/document/ukTMukTMukTM/uAjMxEjLwITMx4CMyETM。

#### Function

##### func NewRobotInfo

```go
func NewRobotInfo(data map[string]interface{}) *RobotInfo
```

从 api 返回体 map 中提取信息，返回值为 RobotInfo 指针。

##### func (AppClient) RobotGetInfo

```go
func (c AppClient) RobotGetInfo() *RobotInfo
```

获取此机器人的信息，返回值为 RobotInfo 指针。

### SpaceApi.go

#### Type

##### type SpaceInfo

```go
type SpaceInfo struct {
	Name        string
	Description string
	SpaceId     string
}
```

知识空间信息结构体，具体可参考 https://open.feishu.cn/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/wiki-overview。

##### type NodeInfo

```go
type NodeInfo struct {
	NodeToken       string
	ObjToken        string
	ObjType         string
	ParentNodeToken string
	Title           string
	HasChild        bool
}
```

知识空间节点信息。

#### Function

##### func NewSpaceInfo

```go
func NewSpaceInfo(data map[string]interface{}) *SpaceInfo
```

从 api 返回体 map 中拉取信息，返回值为 SpaceInfo 指针。

##### func (AppClient) KnowledgeSpaceCreate

```go
func (c AppClient) KnowledgeSpaceCreate(name string, description string, user_access_token string) *SpaceInfo
```

根据所给参数信息创建知识空间，返回值为 SpaceInfo 指针。注意知识空间的创建需要用户授权，获取 user_access_token 作为参数输入。

##### func (AppClient) KnowledgeSpaceAddMembers

```go
func (c AppClient) KnowledgeSpaceAddMembers(spaceId string, membersId []string, memberType string)
```

向知识空间添加成员， memberType 可选值可参考 https://open.feishu.cn/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/wiki-v2/space-member/create。

##### func (AppClient) KnowledgeSpaceAddBotsAsAdmin

```go
func (c AppClient) KnowledgeSpaceAddBotsAsAdmin(spaceId string, BotsId []string, user_access_token string)
```

向知识空间添加机器人作为管理员。注意此操作需要用户授权，获取 user_access_token 作为参数输入。

##### func NewNodeInfo

```go
func NewNodeInfo(data map[string]interface{}) *NodeInfo
```

从 api 返回体 map 中拉取信息，返回值为 NodeInfo 指针。

##### func (AppClient) KnowledgeSpaceCopyNode

```go
func (c AppClient) KnowledgeSpaceCopyNode(SpaceId string, NodeToken string, TargetSpaceId string, TargetParentToken string, Title ...string) *NodeInfo
```

将节点从 SpaceId/NodeToken 复制到 TargetSpaceId/TargetParentToken ，可选参数 Title 表示修改标题。返回值为副本节点 NodeInfo 指针。

##### func (AppClient) KnowledgeSpaceGetAllNodes

```go
func (c AppClient) KnowledgeSpaceGetAllNodes(SpaceId string, ParentNodeToken ...string) []NodeInfo
```

获取知识空间所有节点信息，返回值为 NodeInfo 切片。

### TokenApi.go

#### Type

##### type LoginSession

```go
type LoginSession struct {
	OpenId     string
	EmployeeId string
}
```

用户登录 Session 信息。具体可参考 https://open.feishu.cn/document/uYjL24iN/ukjM04SOyQjL5IDN。

##### type UserAccessToken

```go
type UserAccessToken struct {
	Access_token  string
	Name          string
	Refresh_token string
	User_id       string
}
```

user_access_token 信息。具体可参考 https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/authen/access_token。

#### Function

##### func NewLoginSession

```go
func NewLoginSession(data map[string]interface{}) *LoginSession
```

从 api 返回体 map 中拉取信息，返回值为 LoginSession 指针。

##### func NewUserAccessToken

```go
func NewUserAccessToken(data map[string]interface{}) *UserAccessToken
```

从 api 返回体 map 中拉取信息，返回值为 UserAccessToken 指针。

##### func (AppClient) GetLoginSession

```go
func (c AppClient) GetLoginSession(login_token string) *LoginSession
```

根据 login code 获取 LoginSession 信息，返回值为 LoginSession 指针。

##### func (AppClient) GetUserAccessToken

```go
func (c *AppClient) GetUserAccessToken(code string) *UserAccessToken
```

根据登录预授权码 code 获取 user_access_token 。获取 code 方法可参考 https://open.feishu.cn/document/ukTMukTMukTM/ukzN4UjL5cDO14SO3gTN。