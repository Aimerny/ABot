package core

import (
	"ABots/common"
	"ABots/core/entity"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

type Abot struct {
	conf        *common.Config
	BotMap      map[string]entity.IBot
	DataSources map[string]*gorm.DB
}

func InitAbot(conf *common.Config) *Abot {
	// 创建Abot
	abot := &Abot{conf: conf}
	// 初始化datasource
	abot.InitDataSource()
	// 初始化bots
	abot.InitBots()
	return abot
}

func (abot *Abot) InitBots() {
	abot.BotMap = make(map[string]entity.IBot)
	for _, botConf := range abot.conf.Bots {
		bot := entity.InitBotByKey(botConf.Name, botConf.HookUrl, entity.BotTriggerType(botConf.TriggerType), &abot.DataSources)
		abot.BotMap[botConf.Name] = bot
		log.Infof("Abot load bot: %s, %v", botConf.Name, bot)
	}
}

func (abot *Abot) InitDataSource() {
	// 初始化
	abot.DataSources = make(map[string]*gorm.DB)
	// 找到db
	dbConfig := abot.conf.DSConfigs
	for _, dbConf := range dbConfig {
		config := dbConf.DBConfig
		db := config.ConnectDB()
		log.Infof("Abot load data source: %s, %s", dbConf.SourceName, dbConf.DBType)
		abot.DataSources[dbConf.SourceName] = db
	}
}

func (abot *Abot) Run() {
	// uat
	wg := sync.WaitGroup{}

	// 启动所有机器人
	for botName, bot := range abot.BotMap {
		if bot == nil {
			log.Warnf("Abot skip start [%s], cause of nil", botName)
			continue
		}
		wg.Add(1)
		log.Infof("Abot start [%s]", botName)
		go bot.Start(&wg)
	}

	wg.Wait()
}
