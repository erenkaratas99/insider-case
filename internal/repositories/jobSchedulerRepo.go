package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"insider/configs/appConfigs"
	"insider/internal/apps/jobScheduler/entities"
	"insider/pkg"
	"time"
)

type JobSchedulerRepo struct {
	mc      *mongo.Client
	jobsCol *mongo.Collection
}

func NewJobSchedulerRepository(cfg *appConfigs.Configurations, mc *mongo.Client) (*JobSchedulerRepo, error) {
	jobsCol, err := pkg.GetMongoCollection(mc, cfg.JobScheduler.MongoDbName, cfg.JobScheduler.JobsColName)
	if err != nil {
		return nil, err
	}
	return &JobSchedulerRepo{mc: mc, jobsCol: jobsCol}, nil
}

func (r *JobSchedulerRepo) CreateJob(jobName string, params ...map[string]interface{}) (string, error) {
	if jobName == "" {
		return "", errors.New("jobName cannot be empty")
	}

	value := new(map[string]interface{})
	if len(params) > 0 {
		value = &params[0]
	}

	job := entities.Job{
		Id:        uuid.NewString(),
		JobName:   jobName,
		Params:    *value,
		StartedAt: time.Now(),
	}

	res, err := r.jobsCol.InsertOne(context.TODO(), job)
	if err != nil {
		return "", err
	}

	insertedID, ok := res.InsertedID.(string)
	if !ok {
		return "", errors.New("failed to convert inserted ID to string")
	}

	return insertedID, nil
}

func (r *JobSchedulerRepo) SetFailedAt(jobID string) error {
	if jobID == "" {
		return errors.New("jobID cannot be empty")
	}

	filter := bson.M{"_id": jobID}
	update := bson.M{"$set": bson.M{"failedAt": time.Now()}}

	_, err := r.jobsCol.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *JobSchedulerRepo) SetFinishedAt(jobID string) error {
	if jobID == "" {
		return errors.New("jobID cannot be empty")
	}

	filter := bson.M{"_id": jobID}
	update := bson.M{"$set": bson.M{"finishedAt": time.Now()}}

	_, err := r.jobsCol.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *JobSchedulerRepo) GetJob(jobId string) (*entities.Job, error) {
	if jobId == "" {
		return nil, errors.New("jobID cannot be empty")
	}
	filter := bson.M{"_id": jobId}

	var job entities.Job

	err := r.jobsCol.FindOne(context.TODO(), filter).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	return &job, nil
}

func (r *JobSchedulerRepo) UpdateJob(jobId string, updateFields map[string]interface{}) error {
	if jobId == "" {
		return errors.New("jobID cannot be empty")
	}

	if len(updateFields) == 0 {
		return errors.New("updateFields cannot be empty")
	}

	filter := bson.M{"_id": jobId}

	update := bson.M{"$set": updateFields}

	res, err := r.jobsCol.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("job not found or no changes made")
	}

	return nil
}
