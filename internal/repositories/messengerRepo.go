package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"insider/configs/appConfigs"
	"insider/internal/apps/messengerApi/entities"
	"insider/pkg"
	"time"
)

type MessengerRepo struct {
	mc        *mongo.Client
	rc        *redis.Client
	msgCol    *mongo.Collection
	validator entities.Validator
}

func NewMessengerRepository(cfg *appConfigs.Configurations, mc *mongo.Client, rc *redis.Client) (*MessengerRepo, error) {
	msgCol, err := pkg.GetMongoCollection(mc, cfg.MessengerApi.MongoDbName, cfg.MessengerApi.MessagesColName)
	if err != nil {
		return nil, err
	}
	return &MessengerRepo{
		mc:        mc,
		rc:        rc,
		msgCol:    msgCol,
		validator: entities.NewValidator(),
	}, nil
}

func (r *MessengerRepo) GetAllMessagesMongo(l, o int64) ([]*entities.MessageInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(l)
	findOptions.SetSkip(o)

	cursor, err := r.msgCol.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Info(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*entities.MessageInfo
	for cursor.Next(ctx) {
		var message entities.MessageInfo
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	if err := cursor.Err(); err != nil {
		log.Info(err.Error())
		return nil, err
	}

	return messages, nil
}

func (r *MessengerRepo) GetTwoPendingMessagesMongo(offset int64) ([]*entities.MessageInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"status": "pending"}

	findOptions := options.Find().
		SetLimit(3).
		SetSkip(offset).
		SetProjection(bson.M{"to": 1, "content": 1})

	cursor, err := r.msgCol.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error fetching messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*entities.MessageInfo
	validCount := 0
	maxValid := 2

	for cursor.Next(ctx) {
		var msg entities.MessageInfo
		if err := cursor.Decode(&msg); err != nil {
			log.Printf("error decoding message: %v", err)
			continue
		}

		if err := r.validator.Validate(&msg); err != nil {
			err := r.ChangeMsgStatusMongo(msg.Id, "invalid")
			if err != nil {
				log.Printf("error on updating invalid message: %v", err)
			}
			continue
		}

		messages = append(messages, &msg)
		validCount++

		if validCount >= maxValid {
			break
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return messages, nil
}

func (r *MessengerRepo) ChangeMsgStatusMongo(msgId, status string) error {
	if msgId == "" || status == "" {
		return errors.New("msgId or status cannot be empty")
	}

	filter := bson.M{"_id": msgId}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := r.msgCol.UpdateOne(context.TODO(), filter, update)
	return err
}

func (m *MessengerRepo) SetRedisKey(key string, value string, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.rc.Set(ctx, key, value, expiry).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to set key: %s", err))
	}
	return nil
}

func (m *MessengerRepo) GetRedisKey(key string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := m.rc.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		return "", false, errors.New(fmt.Sprintf("failed to get key: %s", err))
	}
	return val, true, nil
}
