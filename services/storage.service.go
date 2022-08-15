package services

import (
	"context"
	"errors"
	models "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

func InitMongoDB() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, Config.MongodbDatabase, options.Client().ApplyURI(Config.MongodbUri))
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB!")
}

var redisDefaultClient *redis.Client
var redisDefaultOnce sync.Once

var redisCache *cache.Cache
var redisCacheOnce sync.Once

func GetRedisDefaultClient() *redis.Client {
	redisDefaultOnce.Do(func() {
		redisDefaultClient = redis.NewClient(&redis.Options{
			Addr: Config.RedisDefaultAddr,
		})
	})

	return redisDefaultClient
}

func GetRedisCache() *cache.Cache {
	redisCacheOnce.Do(func() {
		redisCache = cache.New(&cache.Options{
			Redis:      GetRedisDefaultClient(),
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	})

	return redisCache
}

func CheckRedisConnection() {
	redisClient := GetRedisDefaultClient()
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to Redis!")
}

func getNoteCacheKey(userId primitive.ObjectID, noteId primitive.ObjectID) string {
	return "req:cache:note:" + userId.Hex() + ":" + noteId.Hex()
}

func CacheOneNote(userId primitive.ObjectID, note *models.Note) {
	if !Config.UseRedis {
		return
	}

	noteCacheKey := getNoteCacheKey(userId, note.ID)

	_ = GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   noteCacheKey,
		Value: note,
		TTL:   time.Minute,
	})
}

func GetNoteFromCache(userId primitive.ObjectID, noteId primitive.ObjectID) (*models.Note, error) {
	if !Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	note := &models.Note{}
	noteCacheKey := getNoteCacheKey(userId, noteId)
	err := GetRedisCache().Get(context.TODO(), noteCacheKey, note)
	return note, err
}
