package entity

type BotResp struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
}

// ==== 默认构造函数 ====
func NewTextRequest(text Text) BotTextRequest {
	return BotTextRequest{MsgType: "text", Text: text}
}

func NewMarkDownRequest(markdown Markdown) BotMarkDownRequest {
	return BotMarkDownRequest{MsgType: "markdown", MarkDown: markdown}
}

func NewImageRequest(image Image) BotImageRequest {
	return BotImageRequest{MsgType: "image", Image: image}
}

func NewBotNewsRequest(news News) BotNewsRequest {
	return BotNewsRequest{MsgType: "news", News: news}
}

// Text请求
type BotTextRequest struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
}

type Text struct {
	Content             string `json:"content"`
	MentionedList       string `json:"mentioned_list"`
	MentionedMobileList string `json:"mentioned_mobile_list"`
}

// Markdown请求
type BotMarkDownRequest struct {
	MsgType  string   `json:"msgtype"`
	MarkDown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

// 图片类型
type BotImageRequest struct {
	MsgType string `json:"msgtype"`
	Image   Image  `json:"image"`
}

type Image struct {
	Base64 string `json:"base64"`
	Md5    string `json:"md5"`
}

// 图文类型
type BotNewsRequest struct {
	MsgType string `json:"msgtype"`
	News    News   `json:"news"`
}

type News struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Picurl      string `json:"picurl"`
}

// todo ing
// 文件类型
// 上传文件
