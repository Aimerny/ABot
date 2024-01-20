package entity

import (
	bot_task "ABots/core/bot-task"
	"ABots/core/models"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

type SegaBot struct {
	baseBot
	segaDS        *gorm.DB
	mainScheduler *gocron.Scheduler
}

func NewSegaBot(db *gorm.DB, templateBot baseBot) *SegaBot {
	// 初始化定时调度器
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	mainScheduler := gocron.NewScheduler(timezone)
	log.Infof("SegaBot main cron scheduler loaded!")
	return &SegaBot{baseBot: templateBot, segaDS: db, mainScheduler: mainScheduler}
}

func (b SegaBot) Start(group *sync.WaitGroup) {
	defer group.Done()

	b.mainScheduler.Name("Sega Scheduler")
	recommendTask := func() {
		levels := bot_task.RecommendDayWithLevelType(
			b.segaDS, 10, []models.MusicLevelType{models.Master, models.Expert, models.ReMaster}, models.TYPE_CHUNI)
		// 获得谱面列表,拼接消息
		chuniMsg := bot_task.GenRecommendChuniMarkDawn(levels)
		maiLevels := bot_task.RecommendDayWithLevelType(
			b.segaDS, 10, []models.MusicLevelType{models.Master, models.ReMaster}, models.TYPE_MAIMAI)
		maiMsg := bot_task.GenRecommendMaiMarkDawn(maiLevels)
		b.SendMarkDownMsg(chuniMsg)
		b.SendMarkDownMsg(maiMsg)
	}

	_, err := b.mainScheduler.CronWithSeconds("0 30 12 ? * 1,2,3,4,5").Do(recommendTask)
	_, err = b.mainScheduler.CronWithSeconds("0 0 19 ? * 1,2,3,4,5").Do(recommendTask)
	//_, err = b.mainScheduler.CronWithSeconds("0 20 20 ? * 1,2,3,4,5").Do(recommendTask)
	//_, err = b.mainScheduler.Every(5).Seconds().Do(recommendTask)
	if err != nil {
		log.Errorf("Add job failed: %e", err)
		return
	}
	b.mainScheduler.StartBlocking()

}
