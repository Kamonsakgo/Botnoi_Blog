package entities

import "time"

type Model struct{}

func (m Model) MessageURLProfile(userID, audioID, url string) MessageURLProfile {
	return MessageURLProfile{
		UserID:   userID,
		AudioID:  audioID,
		URL:      url,
		LastDate: time.Now().UTC().Add(time.Hour * 7),
	}
}

func (m Model) FormatMessage(userID, message, channel string, count int, speaker, audioID string, url *string, page string) FormatMessage {
	return FormatMessage{
		UserID:   userID,
		AudioID:  audioID,
		Message:  message,
		Channel:  channel,
		Count:    count,
		Speaker:  speaker,
		Datetime: DateTimeBangkok(),
		Page:     page,
		URL:      url,
	}
}

func (m Model) NewGenVoiceUnit(url, speaker, text, audioID string, duration float64, typeMedia string, speed float64, volume float64) NewGenVoiceUnit {
	return NewGenVoiceUnit{
		URL:       url,
		Speaker:   speaker,
		Text:      text,
		TypeMedia: typeMedia,
		Datetime:  DateTimeBangkok(),
		AudioID:   audioID,
		Speed:     speed,
		Volume:    volume,
		Duration:  duration,
	}
}

func (m Model) StorageProfileModel(userID string) StorageProfile {
	return StorageProfile{
		UserID:        userID,
		Download:      []string{},
		Length:        0,
		MaxStorageURL: 100,
	}
}

func DateTimeBangkok() time.Time {
	return time.Now().UTC().Add(time.Hour * 7)
}

type MessageURLProfile struct {
	UserID   string    `json:"user_id"`
	AudioID  string    `json:"audio_id"`
	URL      string    `json:"url"`
	LastDate time.Time `json:"last_date"`
}

type FormatMessage struct {
	UserID   string    `json:"user_id"`
	AudioID  string    `json:"audio_id"`
	Message  string    `json:"message"`
	Channel  string    `json:"channel"`
	Count    int       `json:"count"`
	Speaker  string    `json:"speaker"`
	Datetime time.Time `json:"datetime"`
	Page     string    `json:"page"`
	URL      *string   `json:"url,omitempty"`
}

type NewGenVoiceUnit struct {
	URL       string    `json:"url"`
	Speaker   string    `json:"speaker"`
	Text      string    `json:"text"`
	TypeMedia string    `json:"type_media"`
	Datetime  time.Time `json:"datetime"`
	AudioID   string    `json:"audio_id"`
	Speed     float64   `json:"speed"`
	Volume    float64   `json:"volume"`
	Duration  float64   `json:"duration"`
}

type StorageProfile struct {
	UserID        string   `json:"user_id"`
	Download      []string `json:"download"`
	Length        int      `json:"length"`
	MaxStorageURL int      `json:"maxStorageURL"`
}
