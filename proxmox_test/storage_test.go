package proxmox_test

import (
	"reflect"
	"testing"
	. "github.com/mrgloba/proxmox-api2/proxmox"
)

func TestStorage_GetContent(t *testing.T) {
	tests := []struct {
		name    string
		param int
		wantErr bool
	}{
		{
			name:    "Storgae.GetContent() test from proxmox object",
			param: 0,
			wantErr: true,
		},
		{
			name:    "Storgae.GetContent() test from node object",
			param: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var storages []Storage
			var err error

			if (tt.param == 0) {
				storages, err = server.GetStorageList()
			} else {
				nodes, err := server.GetNodes()
				if err != nil {
					t.Log(err.Error())
					return
				}
				storages, err = nodes[0].GetStorageList()
			}

			if err != nil {
				t.Log(err.Error())
				return
			}
			storage := Storage{}
			for _, s := range storages {
				if s.Storage == "local" {
					storage = s
				}
			}

			if storage.Storage != "local" {
				t.Log("Storage \"local\" not found.")
				return
			}

			got, err := storage.GetContents()
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if reflect.DeepEqual(got, nil) {
				t.Errorf("Storage.GetContent() = %v, want %v", got, nil)
			}
		})
	}
}
