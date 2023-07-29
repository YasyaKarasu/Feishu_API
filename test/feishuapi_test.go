package test

import (
	"fmt"
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"testing"
)

func TestDocumentGetAllRecords(t *testing.T) {
	var cli feishuapi.AppClient

	Init()
	SetAppClientConfig(&cli)

	cli.StartTokenTimer()

	bitable := cli.DocumentGetAllBitables(cli.KnowledgeSpaceGetNodeInfo("wikcnzChhgKLNBNVc8yIZAgB8fb").ObjToken)[0]
	table := cli.DocumentGetAllTables(bitable.AppToken)[0]

	allRecords := cli.DocumentGetAllRecordsWithLinks(table.AppToken, table.TableId)
	for _, record := range allRecords {
		logrus.Info(parseRecordFields(record.Fields))
	}
}

type Record struct {
	// 多行文本
	MultiLineText []interface{}
	// 维护人
	Maintainers []Maintainer
	// 一句话介绍
	OneLineIntroduction []interface{}
	// 维护的节点链接
	NodeLink []interface{}
	// 👍
	LikeCount int
}

// Maintainer 定义一个结构，用于存储维护人的信息
type Maintainer struct {
	Name string
	ID   string
}

// 从API返回的record的Fields中解析出Record信息
// 如果某个字段没写，读取map时会返回nil，所以要检查并处理
func parseRecordFields(record map[string]interface{}) Record {
	result := Record{}
	// 解析多行文本
	if record["多行文本"] != nil {
		result.MultiLineText = record["多行文本"].([]interface{})
	}
	// 解析维护人
	if record["维护人"] != nil {
		maintainers := record["维护人"].([]interface{})
		for _, maintainer := range maintainers {
			maintainerMap := maintainer.(map[string]interface{})
			result.Maintainers = append(result.Maintainers, Maintainer{
				Name: maintainerMap["name"].(string),
				ID:   maintainerMap["id"].(string),
			})
		}
	}
	// 解析一句话介绍
	if record["一句话介绍"] != nil {
		result.OneLineIntroduction = record["一句话介绍"].([]interface{})
	}
	// 解析维护的节点链接
	if record["维护节点链接"] != nil {
		result.NodeLink = record["维护节点链接"].([]interface{})
		for _, v := range result.NodeLink {
			currentMap := v.(map[string]interface{})
			fmt.Println(currentMap["type"].(string))
			fmt.Println(currentMap["text"].(string))
			if _, ok := currentMap["token"]; ok {
				fmt.Println(currentMap["token"].(string))
			}
		}
	}
	// 解析点赞数
	if record["👍"] != nil {
		result.LikeCount = int(record["👍"].(float64))
	}

	return result
}

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
