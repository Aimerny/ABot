package entity

import (
	bot_task "ABots/core/bot-task"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

type ArcherBot struct {
	baseBot
	archerDS      *gorm.DB
	mainScheduler *gocron.Scheduler
}

func NewArcherBot(db *gorm.DB, templateBot baseBot) *ArcherBot {
	// 初始化定时调度器
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	mainScheduler := gocron.NewScheduler(timezone)
	log.Infof("ArcherBot main cron scheduler loaded!")
	return &ArcherBot{baseBot: templateBot, archerDS: db, mainScheduler: mainScheduler}
}

func (bot *ArcherBot) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	bot.mainScheduler.Name("Archer Scheduler")

	recommendRest := func() {
		rests := bot_task.ArcherRecommendRest(bot.archerDS, 2)
		log.Infof("推荐的餐厅为:%v", rests)
		msg := bot_task.GenArcherRecommendRestMarkdown(rests)
		bot.SendMarkDownMsg(msg)
	}
	// 工作日的12,18点发送
	_, err := bot.mainScheduler.CronWithSeconds("0 0 12,18 ? * 1,2,3,4,5").Do(recommendRest)
	//_, err = bot.mainScheduler.Every(5).Second().Do(recommendRest)

	if err != nil {
		log.Errorf("Add job failed: %e", err)
		return
	}
	bot.mainScheduler.StartBlocking()
}
