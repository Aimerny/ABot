package bot_task

import (
	"ABots/core/models"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
)

// 每日推荐(指定谱面类型)
func RecommendDayWithLevelType(db *gorm.DB, count int, levelTypes []models.MusicLevelType, gameType int) *[]models.LevelDO {

	res := make([]models.LevelDO, 0)
	//db.Where("level_type in ? and game_type = ?", levelTypes, gameType).Order("RAND()").Limit(count).Find(&res)
	// 圣诞节彩蛋
	db.Where("game_type = ? and music_name = 'ジングルベル'", gameType).Order("RAND()").Limit(count).Find(&res)
	log.Infof(">>>>>>>随机推荐指定类型乐谱: 数量:%d, 游戏:%s", count, models.GameType(gameType))
	for index, level := range res {
		log.Infof("%d. [%s-%s]%s", index+1, level.LevelType, level.Difficulty, level.MusicName)
	}
	return &res
}

// 每日推荐(指定难度)
func RecommendDay(db *gorm.DB, count int, rate float64, gameType int) *[]models.LevelDO {
	iRate := math.Floor(rate)
	oRate := iRate + 1
	if iRate <= 0 || iRate > 15 {
		log.Errorf("Error level of [%f]", rate)
	}
	res := make([]models.LevelDO, 0)
	db.Where("rate_point > ? and rate_point < ? and game_type = ?", iRate, oRate, gameType).Order("RAND()").Limit(count).Find(&res)
	db.Where("rate_point > ? and rate_point < ? and game_type = ?", iRate, oRate, gameType).Order("RAND()").Limit(count).Find(&res)
	log.Infof(">>>>>>>随机推荐乐谱:难度:%.0f, 数量:%d, 游戏:%s", iRate, count, models.GameType(gameType))
	for index, level := range res {
		log.Infof("%d. [%s-%s]%s", index+1, level.LevelType, level.Difficulty, level.MusicName)
	}
	return &res
}

func GenRecommendText(levels *[]models.LevelDO) string {
	var music string
	for index, level := range *levels {
		music += fmt.Sprintf("%d. [%s-%s] %s\n", index+1, level.LevelType, level.Difficulty, level.MusicName)
	}
	return fmt.Sprintf(
		"休息时间到, 接下来是出勤时间!\n" +
			"SBGA为您倾情推荐以下曲目, 请尽情游玩罢!\n" +
			music)
}

func GenRecommendChuniMarkDawn(levels *[]models.LevelDO) string {
	var music string
	for index, level := range *levels {
		music += fmt.Sprintf("%d. <font color=\"%s\">[%s-%s]</font> %s\n", index+1, GetColorByLevelType(level.LevelType), level.LevelType, level.Difficulty, level.MusicName)
	}
	return fmt.Sprintf(
		"# 出勤时间到\n\n<font color=\"orange\">**SBGA**</font>为您倾情推荐以下中二节奏曲目!**请尽情游玩罢**\n" +
			music)
}

func GenRecommendMaiMarkDawn(levels *[]models.LevelDO) string {
	var music string
	for index, level := range *levels {
		music += fmt.Sprintf("%d. <font color=\"%s\">[%s-%s]</font> %s\n", index+1, GetColorByLevelType(level.LevelType), level.LevelType, level.Difficulty, level.MusicName)
	}
	return fmt.Sprintf(
		"# 出勤时间到\n\n<font color=\"orange\">**SBGA**</font>为您倾情推荐以下舞萌DX曲目!**请尽情游玩罢**\n" +
			music)
}

func GetColorByLevelType(levelType models.MusicLevelType) string {
	var color string
	switch levelType {
	case models.Basic:
		color = "GREEN"
	case models.Advanced:
		color = "#e5792b"
	case models.Expert:
		color = "RED"
	case models.Master:
		color = "#851072"
	case models.ReMaster:
		color = "#cb9ad1"
	case models.Ultra:
		color = "#28292c"
	}
	return color
}
