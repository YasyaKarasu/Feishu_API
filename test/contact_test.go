package test

import (
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestUserInfoByName(t *testing.T) {
	var cli feishuapi.AppClient

	Init()
	SetAppClientConfig(&cli)

	cli.StartTokenTimer()

	userInfo := cli.UserInfoByName("张三", "")
	logrus.Info(userInfo)
}
