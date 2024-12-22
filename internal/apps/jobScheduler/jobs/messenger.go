package jobs

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"insider/configs/appConfigs"
	errorConfigs "insider/configs/errorConfigs"
	"insider/internal/apps/jobScheduler/entities"
	"insider/internal/clients"
	"insider/internal/repositories"
	"insider/pkg"
	"strconv"
	"sync"
	"time"
)

type MessengerJob struct {
	isRunning       bool
	mu              sync.Mutex
	ctx             context.Context
	cancel          context.CancelFunc
	ticker          *time.Ticker
	jsRepo          *repositories.JobSchedulerRepo
	errCodes        *errorConfigs.ErrorCode
	messengerClient *clients.MessengerClient
	webhookClient   *clients.WebHookSiteClient
	jobID           string
}

func NewMessengerJob(jsRepo *repositories.JobSchedulerRepo, errCodes *errorConfigs.ErrorCode, cfg *appConfigs.Configurations) *MessengerJob {
	return &MessengerJob{
		ticker:          time.NewTicker(time.Duration(cfg.JobScheduler.MessengerJobIntervalInSeconds) * time.Second),
		jsRepo:          jsRepo,
		errCodes:        errCodes,
		messengerClient: clients.NewMessengerClient(cfg),
		webhookClient:   clients.NewWebHookSiteClient(),
	}
}

func (mj *MessengerJob) Start() *entities.JobDTO {
	mj.mu.Lock()
	defer mj.mu.Unlock()

	if mj.isRunning {
		return &entities.JobDTO{
			Code:    mj.errCodes.MESSENGER_JOB_ALREADY_WORKING,
			Status:  "failed",
			Message: "Job is already working.",
		}
	}

	mj.isRunning = true
	mj.ctx, mj.cancel = context.WithCancel(context.Background())

	jobID, err := mj.jsRepo.CreateJob("MessengerJob", map[string]interface{}{
		"last_offset": 0,
	})
	if err != nil {
		// rollback
		mj.isRunning = false
		return &entities.JobDTO{
			Code:    mj.errCodes.JOB_COULDNT_CREATED,
			Status:  "failed",
			Message: "Job couldn't be initiated.",
		}
	}
	mj.jobID = jobID

	go mj.doJob(mj.ctx)

	return &entities.JobDTO{
		Code:    0,
		Status:  "success",
		Message: "success",
	}
}

func (mj *MessengerJob) Stop() *entities.JobDTO {
	mj.mu.Lock()
	defer mj.mu.Unlock()

	if !mj.isRunning {
		return &entities.JobDTO{
			Code:    mj.errCodes.MESSENGER_JOB_ALREADY_STOPPED,
			Status:  "failed",
			Message: "Job is already stopped.",
		}
	}

	mj.cancel()
	mj.isRunning = false

	if err := mj.jsRepo.SetFinishedAt(mj.jobID); err != nil {
		return &entities.JobDTO{
			Code:    mj.errCodes.JOB_COULDNT_UPDATED,
			Status:  "failed",
			Message: "Job couldn't be updated.",
		}
	}

	return &entities.JobDTO{
		Code:    0,
		Status:  "success",
		Message: "success",
	}
}

func (mj *MessengerJob) IsRunning() bool {
	mj.mu.Lock()
	defer mj.mu.Unlock()
	return mj.isRunning
}

func (mj *MessengerJob) doJob(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			log.Println("MessengerJob: context canceled, exiting job.")
			return

		case <-mj.ticker.C:
			log.Println("MessengerJob: Starting iteration...")

			jobInfo, err := mj.jsRepo.GetJob(mj.jobID)
			if err != nil {
				log.Errorf("Could not retrieve job info for jobID %s: %v", mj.jobID, err)
				continue
			}

			offsetVal, ok := jobInfo.Params["last_offset"]
			if !ok {
				log.Warn("last_offset not found in job params, using 0")
				offsetVal = 0
			}

			lastOffset, err := pkg.ToInt(offsetVal)
			if err != nil {
				log.Errorf("Invalid last_offset in job metadata: %v, defaulting to 0", err)
				lastOffset = 0
			}

			body, err := mj.messengerClient.GetTwoPendingMessages(strconv.Itoa(lastOffset + 2))
			if err != nil {
				log.Errorf("Failed to get pending messages: %v", err)
				continue
			}

			messengerResponse := new(entities.MessengerClientResponse)
			if err := json.Unmarshal(body, &messengerResponse); err != nil {
				log.Errorf("Failed to unmarshal messenger response: %v", err)
				continue
			}

			if len(messengerResponse.Data) == 0 {
				log.Println("No new messages to process.")
				continue
			}

			for _, msg := range messengerResponse.Data {
				message := map[string]interface{}{
					"to":      msg.To,
					"content": msg.Content,
				}

				if err := mj.webhookClient.SendMessage(message); err != nil {
					log.Errorf("Error sending message (ID=%s): %v", msg.ID, err)
					continue
				}

				if err := mj.messengerClient.CommitMessageAsSent(msg.ID); err != nil {
					log.Errorf("Error committing message as sent (ID=%s): %v", msg.ID, err)
					continue
				}
			}

			newOffset := lastOffset + len(messengerResponse.Data)

			err = mj.jsRepo.UpdateJob(mj.jobID, map[string]interface{}{
				"last_offset": newOffset,
			})
			if err != nil {
				log.Errorf("Failed to update last_offset in DB: %v", err)
			}

			log.Printf("MessengerJob: iteration finished. Updated last_offset to %d\n", newOffset)
		}
	}
}
