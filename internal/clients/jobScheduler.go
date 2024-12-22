package clients

import (
	"errors"
	"github.com/valyala/fasthttp"
	"insider/configs/appConfigs"
)

type JobSchedulerClient struct {
	client *BaseClient
}

func NewJobSchedulerClient(cfg *appConfigs.Configurations) *JobSchedulerClient {
	return &JobSchedulerClient{client: NewBaseClient(cfg.JobScheduler.Url)}
}

func (jc *JobSchedulerClient) ToggleJob(jobName, command string) (*fasthttp.Response, error) {
	opts := map[string]string{
		"Content-Type": "application/json",
	}

	if command != "start" && command != "stop" {
		return nil, errors.New("invalid command")
	}

	res, err := jc.client.GET("/job-scheduler/"+jobName+"/"+command, opts)
	if res != nil {
		defer fasthttp.ReleaseResponse(res)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
