package appConfigs

import (
	"encoding/json"
	"fmt"
	"insider/pkg"
	"os"
	"time"
)

type Configurations struct {
	Mongo        MongoConfigs `json:"mongo,omitempty"`
	RedisConfigs RedisConfigs `json:"redis,omitempty"`
	MessengerApi struct {
		Url             string `json:"url"`
		ApiPrefix       string `json:"api_prefix,omitempty"`
		RoutePrefix     string `json:"route_prefix,omitempty"`
		MongoDbName     string `json:"mongo_db_name,omitempty"`
		MessagesColName string `json:"messages_col_name,omitempty"`
	} `json:"messenger_api,omitempty"`
	JobScheduler struct {
		Url                           string `json:"url"`
		RoutePrefix                   string `json:"route_prefix,omitempty"`
		MongoDbName                   string `json:"mongo_db_name,omitempty"`
		JobsColName                   string `json:"jobs_col_name,omitempty"`
		MessengerJobIntervalInSeconds int    `json:"messenger_job_interval_in_seconds,omitempty"`
	} `json:"job_scheduler,omitempty"`
}

type MongoConfigs struct {
	ConnectionURI            string        `json:"connection_uri,omitempty"`
	TimeOutDurationInSeconds time.Duration `json:"time_out_duration_in_seconds,omitempty"`
}

type RedisConfigs struct {
	ConnectionURI string `json:"connection_uri,omitempty"`
}

func GetConfigs() Configurations {
	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "test"
		fmt.Println("ENV is not set, defaulting to test")
	}
	conf, err := pkg.ReadJsonFile(fmt.Sprintf("configs/appConfigs/%sconfigs.json", env))
	if err != nil {
		panic(err)
	}
	values := new(Configurations)
	if err = json.Unmarshal(conf, &values); err != nil {
		panic(err)
	}
	return *values
}
