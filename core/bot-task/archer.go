package bot_task

import (
	"ABots/core/models"
	"ABots/core/util"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func ArcherRecommendRest(db *gorm.DB, count int) *[]models.Restaurant {

	res := make([]models.Restaurant, 0)
	// 非高优的先统一提高一次优先级
	db.Model(&models.Restaurant{}).Where("priority > ?", models.P_HIGH).Update("priority", gorm.Expr("priority - 1"))
	// 在高优的记录中找两个
	db.Where("priority = ?", models.P_HIGH).Order("RAND()").Limit(count).Find(&res)
	// 将选出的设为低优
	for i := range res {
		db.Model(&res[i]).Update("priority", models.P_LOW)
	}
	return &res
}

func GenArcherRecommendRestMarkdown(rests *[]models.Restaurant) string {
	// 判断是中午还是下午
	nowHour := time.Now().Hour()
	msg := fmt.Sprintf("# <font color=\"%s\">吃饭时间到</font>\n", "#016dca")
	if nowHour < 13 {
		msg = msg + "中午好,"
	} else {
		msg = msg + "下午好,"
	}
	msg = msg + "吃饭时间到了,还在纠结吃什么?为您推荐以下餐厅:\n"
	for i, rest := range *rests {
		msg += fmt.Sprintf("%d. **%s**\n", i+1, rest.Name)
	}
	addMsg, err := util.GetWeatherInfo()
	if err != nil {
		log.Warn("get weather failed:[%e]", err)
		return msg
	}

	msg += "<font color=\"red\">温馨提示</font>\n"
	return msg + addMsg
}
