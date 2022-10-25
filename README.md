# feishuapi

feishuapi is a go module packing up the basic feishu apis.

## External Packages

| Scheduler      | Log Manager     |
| -------------- | --------------- |
| robfig/cron/v3 | sirupsen/logrus |

## Guide

### Installation

```
go get github.com/YasyaKarasu/feishuapi
```

### Quick Start

1. Import this module in your project
2. Set configuration file properly
3. Create a feishuapi AppClient in your project and set the configuration parameters
4. Use AppClient.StartTokenTimer function to get tenant access token

### Example

**main.go**:

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

**conf.go**

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

**config.yaml**:

```yaml
feishu:
  APP_ID: cli_a385f***********
  APP_SECRET: RxvCp***************************
  VERIFICATION_TOKEN: CX5YW***************************
  ENCRYPT_KEY: ''
```
