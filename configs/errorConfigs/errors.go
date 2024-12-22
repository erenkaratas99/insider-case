package appConfigs

import (
	"encoding/json"
	"fmt"
	"insider/pkg"
)

// In case of there is a BFF structure between the UI and backend to manage error responses
type ErrorCode struct {
	MESSENGER_JOB_ALREADY_WORKING int `json:"messenger_job_already_working"`
	MESSENGER_JOB_ALREADY_STOPPED int `json:"messenger_job_already_stopped"`
	JOB_COULDNT_CREATED           int `json:"job_couldnt_created"`
	JOB_COULDNT_UPDATED           int `json:"job_couldnt_updated"`
}

func GetErrorCodes() *ErrorCode {
	conf, err := pkg.ReadJsonFile(fmt.Sprintf("configs/errorConfigs/error-codes.json"))
	if err != nil {
		panic(err)
	}
	values := new(ErrorCode)
	if err = json.Unmarshal(conf, &values); err != nil {
		panic(err)
	}
	return values
}
