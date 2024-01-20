package models

import "gorm.io/gorm"

type MusicLevelType string

const (
	Basic    MusicLevelType = "BASIC"
	Advanced MusicLevelType = "ADVANCED"
	Expert   MusicLevelType = "EXPERT"
	Master   MusicLevelType = "MASTER"
	ReMaster MusicLevelType = "RE:MASTER"
	Ultra    MusicLevelType = "ULTRA"
	WorldEnd MusicLevelType = "WORLD END"

	// 游戏类型
	TYPE_MAIMAI int = 0
	TYPE_CHUNI  int = 1
)

func GameType(gameType int) string {
	if gameType == 0 {
		return "舞萌DX"
	} else if gameType == 1 {
		return "中二节奏2024"
	}
	return "不存在的游戏"
}

type LevelDO struct {
	gorm.Model
	// 曲名
	MusicName string
	// 谱面类型
	LevelType MusicLevelType
	// 难度
	Difficulty string
	// 定数
	RatePoint float32
	// 游戏 ( 0-舞萌, 1-中二 )
	GameType int
	// 是否DX
	DXLevel bool
}

func (level *LevelDO) TableName() string {
	return "tb_level"
}
