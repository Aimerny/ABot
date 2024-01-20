package models

import (
	"bufio"
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
)

func FlushSegaLevels(db *gorm.DB) {
	db.AutoMigrate(&LevelDO{})
	// 先刷中二
	chuniFile, err := os.Open("resources/中二节奏.csv")
	if err != nil {
		log.Errorf("open file error")
	}
	// 次数对应的枚举
	levelTypeTimeMap := map[int]MusicLevelType{
		0: Basic,
		1: Advanced,
		2: Expert,
		3: Master,
		4: Ultra,
	}

	// 由于没有等级,初始化一个等级列表
	levelNameTimeMap := make(map[string]int)

	getType := func(musicName string, ratePoint float64) MusicLevelType {
		nowTime := levelNameTimeMap[musicName]
		if ratePoint <= 0 {
			return WorldEnd
		}
		return levelTypeTimeMap[nowTime]
	}
	scanner := bufio.NewScanner(chuniFile)
	for scanner.Scan() {
		linestr := scanner.Text()
		if err != nil {
			break
		}
		line := strings.Split(linestr, ",")
		ratePoint, err := strconv.ParseFloat(line[3], 32)
		if err != nil {
			log.Errorf("Error parse rate [%s]", line)
		}

		db.Create(&LevelDO{
			MusicName:  line[1],
			LevelType:  getType(line[1], ratePoint),
			Difficulty: line[2],
			RatePoint:  float32(ratePoint),
			GameType:   TYPE_CHUNI,
		})
		levelNameTimeMap[line[1]]++
		//break
		log.Infof("写入中二歌曲: %s", line)
	}
	log.Infof("结束写入中二")
}

func FlushSegaMaiMaiLevels(db *gorm.DB) {
	// maimai
	maimaiFile, err := os.Open("resources/乐谱.csv")
	if err != nil {
		log.Errorf("open file error")
	}

	reader := csv.NewReader(maimaiFile)
	reader.LazyQuotes = true
	for {
		// 曲名,类别,难度,等级,定数
		line, err := reader.Read()
		if err != nil {
			break
		}
		ratePoint, err := strconv.ParseFloat(line[4], 32)
		if err != nil {
			log.Errorf("Error parse rate [%s]", line)
		}
		db.Create(&LevelDO{
			MusicName:  line[0],
			LevelType:  MusicLevelType(strings.ToUpper(line[2])),
			Difficulty: line[3],
			RatePoint:  float32(ratePoint),
			GameType:   TYPE_MAIMAI,
			DXLevel:    line[1] == "DX",
		})
		//break
		log.Infof("写入MaiMai歌曲: %s", line)
	}

}
