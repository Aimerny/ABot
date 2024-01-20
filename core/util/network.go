package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// HttpPostWithJsonBody 发送Post请求到指定url,将body进行json序列化作为请求体
func HttpPostWithJsonBody(url string, body any) {

	bodyData, err := json.Marshal(body)
	if err != nil {
		log.Panicf("marshal body to json failed! body:%v err:%e", body, err)
		return
	}
	buffer := bytes.NewBuffer(bodyData)
	// 请求
	log.Println()
	resp, err := http.Post(url, "application/json", buffer)

	if err != nil {
		log.Panicf("HTTP POST error: %e", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Panicf("HTTP POST status code not success")
	}
}

func CalcMd5(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	hashBytes := hasher.Sum(nil)
	md5Str := hex.EncodeToString(hashBytes)
	return md5Str
}
