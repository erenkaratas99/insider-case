package pkg

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"insider/internal/apps/messengerApi/entities"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func ReadJsonFile(path string) ([]byte, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	return byteValue, err
}

func ToInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int32:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

func generateRandomTurkishGSM() string {
	rand.Seed(time.Now().UnixNano())
	phoneDigits := rand.Intn(900000000) + 100000000
	return "+905" + strconv.Itoa(phoneDigits)
}

func GenerateMockData(client *mongo.Client, dbName, collectionName string, count int) error {
	collection := client.Database(dbName).Collection(collectionName)

	if _, err := collection.DeleteMany(context.Background(), bson.M{}); err != nil {
		return fmt.Errorf("failed to delete existing data: %w", err)
	}
	fmt.Printf("Deleted all existing documents from collection %q\n", collectionName)

	baseTime := time.Now()
	var docs []interface{}

	epochTime := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < count; i++ {
		docTime := baseTime.Add(time.Duration(i) * time.Minute)

		m := entities.MessageInfo{
			Id:             uuid.NewString(),
			To:             generateRandomTurkishGSM(),
			Content:        fmt.Sprintf("This is message number %d", i),
			Status:         "pending",
			SentAt:         &epochTime,
			FailedToSentAt: &epochTime,
			CreatedAt:      &docTime,
			UpdatedAt:      &docTime,
		}
		docs = append(docs, m)
	}

	if _, err := collection.InsertMany(context.Background(), docs); err != nil {
		return fmt.Errorf("failed to insert mock data: %w", err)
	}

	fmt.Printf("Inserted %d mock documents into collection %q\n", count, collectionName)
	return nil
}
