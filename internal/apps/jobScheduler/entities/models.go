package entities

import "time"

type Job struct {
	Id         string                 `json:"id" bson:"_id"`
	JobName    string                 `json:"jobName" bson:"jobName"`
	Params     map[string]interface{} `json:"params" bson:"params"`
	StartedAt  time.Time              `json:"startedAt" bson:"startedAt"`
	FinishedAt time.Time              `json:"finishedAt" bson:"finishedAt"`
	FailedAt   time.Time              `json:"failedAt" bson:"failedAt"`
}
