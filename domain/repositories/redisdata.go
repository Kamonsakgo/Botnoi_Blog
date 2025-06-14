package repositories

import (
	"bn-crud-ads/domain/datasources"
	"bn-crud-ads/domain/entities"
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

type redisConnectionRepository struct {
	Context context.Context
	Redis   *redis.Client
}

type IRedisConnectionRepository interface {
	GetAdsAlertMessage() (*entities.AdsAlertMessage, error)
	SetAdsData(AdsData []byte) bool
	// GetSpeakerData() []entities.SpeakerModel
	// SetSpeakerData(speaker []byte) bool
	// GetPrivateSpeaker() []entities.SpeakerModel
	// SetPrivateSpeaker(speaker []byte) bool
	// GetPublicSpeaker() []entities.SpeakerModel
	// SetPublicSpeaker(speaker []byte) bool

	// GetVoiceStudioCache(user_id string) *entities.VoiceStudioModel
	// SetVoiceStudioCache(user_id string, voice_studio []byte) bool
}

func NewRedisRepository(redis *datasources.RedisConnection) IRedisConnectionRepository {
	return &redisConnectionRepository{
		Context: redis.Context,
		Redis:   redis.Redis,
	}
}

func (repo redisConnectionRepository) GetAdsAlertMessage() (*entities.AdsAlertMessage, error) {
	AdsAlertMessageData, err := repo.Redis.Get(repo.Redis.Context(), "ads").Result()
	if err != nil {
		log.Println("error GetAdsFromRedis ", err.Error())
		return nil, err
	}

	var data *entities.AdsAlertMessage
	json.Unmarshal([]byte(AdsAlertMessageData), &data)
	// log.Println(data)
	// log.Printf("Type of data unmarshal: %T\n", data)
	return data, err
}

func (repo redisConnectionRepository) SetAdsData(AdsData []byte) bool {
	err := repo.Redis.Set(repo.Redis.Context(), "ads", AdsData, 0).Err()
	if err != nil {
		log.Println("error Set Ads alert message data ", err.Error())
		return false
	}
	log.Println("Set new Ads alert message data success!")
	return true
}
