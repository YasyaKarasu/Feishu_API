package test

import (
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestGetCode(t *testing.T) {
	var cli feishuapi.AppClient

	Init()
	SetAppClientConfig(&cli)

	cli.StartTokenTimer()

	code := cli.GetCode("https://open.feishu.cn/uqwQjLasdq04CN%2fucDOz4yN4MjL3gzM", "cli_a4da07bdbe38d00e")
	logrus.Info(code)
}
