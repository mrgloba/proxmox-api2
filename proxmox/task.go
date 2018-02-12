package proxmox

import (
	"github.com/kataras/iris/core/errors"
	"fmt"
	"strings"
	"time"
)

type TaskID string

type BaseTask struct {
	Id string 		`json:"id"`
	Node string 	`json:"node"`
	PID int 		`json:"pid"`
	PStart int 		`json:"pstart"`
	StartTime int 	`json:"starttime"`
	Status string 	`json:"status"` 	// Task("OK"),TaskStatus("running","stopped")
	Type string 	`json:"type"`
	UPid TaskID		`json:"upid"`
	User string 	`json:"user"`
}

type TaskStatus struct {
	ExitStatus string `json:"exitstatus,omitempty"`
	BaseTask
}

type Task struct {
	EndTime int		`json:"endtime"`
	BaseTask
	BasicObject
}

func (t *Task) GetStatus() (*TaskStatus,error){
	if len(string(t.UPid)) == 0 { return nil, errors.New("Can't get status of nil")}

	upparts := strings.Split(string(t.UPid),":")

	target := "nodes/" + upparts[1] + "/tasks/" + string(t.UPid) + "/status"

	var taskStatus TaskStatus

	httpCode, err := t.px.APICall2("GET", target, nil, &taskStatus)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	return &taskStatus, nil
}

func (t *Task) WaitForStatus(status string, timeout int) (bool,*TaskStatus,error){
	var taskStatus *TaskStatus
	var err error
	for i:=0; i<timeout; i++ {
		taskStatus,err = t.GetStatus()
		if err !=nil {
			return false,nil,err
		}
		if taskStatus.Status == status {
			return true,taskStatus,nil
		}
		time.Sleep(1 * time.Second)
	}

	return false, taskStatus,errors.New("Timeout reached, status not get")
}