package entities

type AdsDataFormat struct {
	Id          string  `json:"id" bson:"id,omitempty"`
	Url         string  `json:"url" bson:"url,omitempty"`
	Play        int     `json:"play" bson:"play"`
	Description string  `json:"description" bson:"description,omitempty"`
	Level       float32 `json:"level" bson:"level,omitempty"`
	Status      bool    `json:"status" bson:"status,omitempty"`
	Language    string  `json:"language" bson:"language,omitempty"`
	Image_url   string  `json:"image_url" bson:"image_url,omitempty"`
	Nameads     string  `json:"nameads" bson:"nameads,omitempty"`
}

type AdsModel struct {
	AdsPlay []*UpdateData `json:"ads_play" bson:"ads_play"`
}

type UpdateData struct {
	ID   string `json:"id" bson:"id"`
	Play int    `json:"play" bson:"play"`
}
type AdsAlertMessage struct {
	Ads Ads `json:"ads"`
}

type Ads struct {
	Studio     string    `json:"studio" bson:"studio,omitempty"`
	Max        int       `json:"max" bson:"max,omitempty"`
	LandingMax int       `json:"landing_max" bson:"landing_max,omitempty"`
	Level      []int     `json:"level" bson:"level,omitempty"`
	Percent    []float64 `json:"percent" bson:"percent,omitempty"`
}

// "id" : "9",
// "url" : "https://botnoi-tts.s3.us-west-2.amazonaws.com/audio/speaker_id_1/sell",
// "play" : NumberInt(1),
// "description" : "หม้อไฟฟ้าตัวนี้น่ารักปุ้กปิ้กมาก เหมาะกับคนอยู่หอหรือคอนโดสุดสุดทำเป็นชาบู ต้มมาม่า-ผัด-ทอด-กินได้หมดเลย กินกับเพื่อนสองสามคนกำลังโอเค น้ำหนักเบา ทำความสะอาดง่าย",
// "level" : 1.0,
// "status" : true,
// "language" : "th"
