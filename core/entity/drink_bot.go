package entity

import (
	bot_task "ABots/core/bot-task"
	log "github.com/sirupsen/logrus"
	"sync"
)

type DrinkBot struct {
	ArcherBot
}

func NewDrinkBot(archerBot ArcherBot) *DrinkBot {
	log.Infof("DrinkBot main cron scheduler loaded!")
	return &DrinkBot{ArcherBot: archerBot}
}

func (bot *DrinkBot) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	bot.mainScheduler.Name("Archer Scheduler")

	recommendDrink := func() {
		picPath := "resources/drink.png"
		err := bot.SendImageMsgWithLocalFilePath(picPath)
		if err != nil {
			log.Errorf("Send drink req failed: %e", err)
			return
		}
	}
	recommendRest := func() {
		rests := bot_task.ArcherRecommendRest(bot.archerDS, 2)
		log.Infof("推荐的餐厅为:%v", rests)
		msg := bot_task.GenArcherRecommendRestMarkdown(rests)
		bot.SendMarkDownMsg(msg)
	}
	// 每天上班时间整点发送
	_, err := bot.mainScheduler.CronWithSeconds("0 0 11,12,14,15,16,17,18,19 ? * *").Do(recommendDrink)
	_, err = bot.mainScheduler.CronWithSeconds("0 0 12,18 ? * 1,2,3,4,5").Do(recommendRest)
	// 测试
	//_, err = bot.mainScheduler.Every(5).Second().Do(recommendRest)

	if err != nil {
		log.Errorf("Add job failed: %e", err)
		return
	}
	bot.mainScheduler.StartBlocking()
}
