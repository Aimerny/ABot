package entity

import (
	"ABots/core/util"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"sync"
)

const (
	NONE   BotTriggerType = ""
	SINGLE BotTriggerType = "SINGLE"
	CRON   BotTriggerType = "CRON"
	ALWAYS BotTriggerType = "ALWAYS"

	BotKey_SEGA string = "sega"
	//BotKey_SEGA string = "测试"

	//BotKey_ARCHER string = "uat测试"
	BotKey_ARCHER string = "调度小管家"
	//BotKey_DRINK  string = "测试"
	BotKey_DRINK string = "喝水"
)

var DefaultTriggerType BotTriggerType = SINGLE

type BotTriggerType string

// BOT接口
type IBot interface {
	Start(*sync.WaitGroup)
	SendTextMsg(string)
	SendImageMsgWithLocalFilePath(string) error
	SendMarkDownMsg(string)
	String() string
}

func InitBotByKey(key string, hookUrl string, triggerType BotTriggerType, dsMap *map[string]*gorm.DB) IBot {
	// 获得默认triggerType
	if triggerType == "" {
		triggerType = SINGLE
	}
	templateBot := newBot(key, hookUrl, triggerType)
	// 根据key获得对应的bot
	switch key {
	case BotKey_SEGA:
		return NewSegaBot((*dsMap)[BotKey_SEGA], *templateBot)
	case BotKey_ARCHER:
		return NewArcherBot((*dsMap)[BotKey_ARCHER], *templateBot)
	case BotKey_DRINK:
		return NewDrinkBot(*NewArcherBot((*dsMap)[BotKey_ARCHER], *templateBot))
	}
	log.Errorf("Init bot by key [%s] failed !", key)
	return nil
}

type baseBot struct {
	Name        string
	HookUrl     string
	TriggerType BotTriggerType
}

func newBot(name, hookUrl string, triggerType BotTriggerType) *baseBot {

	return &baseBot{
		Name: name, HookUrl: hookUrl, TriggerType: triggerType,
	}
}

func (b baseBot) String() string {
	return fmt.Sprintf("[Name: %s, HookUrl: %s, TriggerType: %s]", b.Name, b.HookUrl, b.TriggerType)
}

func (b baseBot) SendTextMsg(text string) {
	request := NewTextRequest(Text{
		Content: text,
	})
	util.HttpPostWithJsonBody(b.HookUrl, &request)
}

func (b baseBot) SendImageMsgWithLocalFilePath(path string) error {
	imageData, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("Read image path failed: %s", path)
		return err
	}

	// 将图片转成base64
	base64Str := base64.StdEncoding.EncodeToString(imageData)
	md5Str := util.CalcMd5(imageData)
	request := NewImageRequest(Image{
		Base64: base64Str,
		Md5:    md5Str,
	})
	util.HttpPostWithJsonBody(b.HookUrl, &request)
	return nil
}

func (b baseBot) Start(wg *sync.WaitGroup) {
	log.Errorf("[%s] do nothing, Please overwrite it!", b.Name)
}

func (b baseBot) SendMarkDownMsg(text string) {
	request := NewMarkDownRequest(Markdown{
		Content: text,
	})
	util.HttpPostWithJsonBody(b.HookUrl, &request)
}
