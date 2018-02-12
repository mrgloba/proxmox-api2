package proxmox

import (
	"reflect"
	"testing"
)

func TestTask_GetStatus(t *testing.T) {
	type fields struct {
		EndTime     int
		BaseTask    BaseTask
		BasicObject BasicObject
	}
	tests := []struct {
		name    string
		fields  fields
		want    TaskStatus
		wantErr bool
	}{
		{
			name: "Task.GetStatus() test",
			fields: fields{
				BaseTask:    BaseTask{UPid: "UPID:utm-other:0000530F:1DF56C3C:5A7DCCE8:vzdestroy:999:root@pam:"},
				BasicObject: BasicObject{px: server},
			},
			want: TaskStatus{
				ExitStatus: "OK",
				BaseTask: BaseTask{
					UPid:      "UPID:utm-other:0000530F:1DF56C3C:5A7DCCE8:vzdestroy:999:root@pam:",
					Status:    "stopped",
					Node:      "utm-other",
					Type:      "vzdestroy",
					User:      "root@pam",
					StartTime: 1518193896,
					Id:        "999",
					PID:       21263,
					PStart:    502623292,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				EndTime:     tt.fields.EndTime,
				BaseTask:    tt.fields.BaseTask,
				BasicObject: tt.fields.BasicObject,
			}
			got, err := task.GetStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("Task.GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if DEBUG_TESTS {
				t.Logf("%v\n", *got)
			}

			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Task.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_WaitForStatus(t *testing.T) {
	type fields struct {
		EndTime     int
		BaseTask    BaseTask
		BasicObject BasicObject
	}
	type args struct {
		status  string
		timeout int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		want1   TaskStatus
		wantErr bool
	}{
		{
			name: "Task.WaitForStatus() test",
			fields: fields{
				BaseTask:    BaseTask{UPid: "UPID:utm-other:0000530F:1DF56C3C:5A7DCCE8:vzdestroy:999:root@pam:"},
				BasicObject: BasicObject{px: server},
			},
			args:args{ status:"stopped", timeout:10},
			want: true,
			want1: TaskStatus{
				ExitStatus: "OK",
				BaseTask: BaseTask{
					UPid:      "UPID:utm-other:0000530F:1DF56C3C:5A7DCCE8:vzdestroy:999:root@pam:",
					Status:    "stopped",
					Node:      "utm-other",
					Type:      "vzdestroy",
					User:      "root@pam",
					StartTime: 1518193896,
					Id:        "999",
					PID:       21263,
					PStart:    502623292,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				EndTime:     tt.fields.EndTime,
				BaseTask:    tt.fields.BaseTask,
				BasicObject: tt.fields.BasicObject,
			}
			got, got1, err := task.WaitForStatus(tt.args.status, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("Task.WaitForStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if DEBUG_TESTS {
				t.Logf("%v : %v\n",got, *got1)
			}

			if got != tt.want {
				t.Errorf("Task.WaitForStatus() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(*got1, tt.want1) {
				t.Errorf("Task.WaitForStatus() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
