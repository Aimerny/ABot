package util

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

const (
	WeatherUrl         string = "https://restapi.amap.com/v3/weather/weatherInfo"
	WeatherAPIKey      string = "ad0105162363f6fb94e168c0506d296f"
	CityCodeShanghaiYp string = "310110" // 杨浦代码
	Extensions         string = "base"   // all 返回预报
)

// ==== 天气api ====
type WeatherResp struct {
	Status string         `json:"status"`
	Infos  []*LiveWeather `json:"lives"`
}

type LiveWeather struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Temp     string `json:"temperature"`
	Weather  string `json:"weather"`
}

func GetWeatherInfo() (string, error) {

	url := fmt.Sprintf("%s?key=%s&city=%s&extensions=%s", WeatherUrl, WeatherAPIKey, CityCodeShanghaiYp, Extensions)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Get Weather Info failed: URL:[%s], err:[%e]", url, err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("Get Weather Info failed, Resp:[%v]", resp)
		return "", fmt.Errorf("get Weather Info failed, Resp:[%v]", resp)
	}
	//result := make(map[string]any)
	result := WeatherResp{}
	data, err := io.ReadAll(resp.Body)
	log.Debugf("resp:%s", string(data))
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("read weather body failed, Resp:[%v]", resp)
		return "", fmt.Errorf("read weather body failed, Resp:[%v]", resp)
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Errorf("unmarshal data failed: %s", err.Error())
		return "", fmt.Errorf("unmarshal data failed")
	}
	liveWeather := result.Infos[0]
	msg := fmt.Sprintf("当前室外温度为%s°C, 天气状况为: %s\n", liveWeather.Temp, liveWeather.Weather)
	if strings.Contains(liveWeather.Weather, "雨") {
		msg += "> 今日有雨,请不要忘记带伞>_<"
	}
	return msg, nil
}
