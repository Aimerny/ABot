package common

import (
	"ABots/core/entity"
	"bufio"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

const configPath string = "conf.json"

type Config struct {
	Bots        []*BotConfig        `json:"bots"`
	DSConfigs   []*DatasourceConfig `json:"datasource"`
	ProjectName string              `json:"projectName"`
}

type BotConfig struct {
	Name        string `json:"name"`
	HookUrl     string `json:"hookUrl"`
	TriggerType string `json:"triggerType"`
}

func LoadDefaultConf() *Config {

	path := configPath
	// 判断文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 不存在,新建默认文件
		defaultConf := Config{
			ProjectName: "Example",
			Bots: []*BotConfig{
				{Name: "bot-example", HookUrl: "http://xxx", TriggerType: string(entity.CRON)},
			},
		}
		defaultConf.WriteConfig()
		log.Fatalf("Config file [%s] not found! Gerneated new file already! Check it", configPath)
	}

	confFile, err := os.Open(path)
	if err != nil {
		log.Fatal("Open config failed!")
	}
	defer confFile.Close()

	// 读取文件
	content, err := io.ReadAll(confFile)
	if err != nil {
		log.Fatal("Read config failed!")
	}

	var conf Config
	// 反序列化该json
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatal("Unmarshal file content to json failed:", err)
	}
	log.Infof("Load resources from path:[%s] succeed!", path)
	return &conf
}

// 写config文件
func (c Config) WriteConfig() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 不存在,新建默认文件
		os.Create(configPath)
	}
	confFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		log.Errorf("Can't open configFile: %s, err: %e", configPath, err)
		return
	}
	defer confFile.Close()

	configData, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		log.Errorf("Marshal config to json failed: %v, err: %e", c, err)
	}

	writer := bufio.NewWriter(confFile)
	_, err = writer.Write(configData)
	if err != nil {
		log.Errorf("Write config file failed: %s, err: %e", confFile, err)
	}
	writer.Flush()

}
