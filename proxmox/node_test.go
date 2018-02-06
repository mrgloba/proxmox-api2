package proxmox

import (
	"testing"
)

func TestNode_GetStorageList(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Get Storages",
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := server.GetNodes()
			if err != nil {
				t.Log(err.Error())
				return
			}
			got, err := nodes[0].GetStorageList()
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.GetStorages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got[0].Storage != "local" && got[1].Storage != "local" {
				t.Errorf("Node.GetStorages() = %v, no default storage", got)
			}
		})
	}
}

func TestNode_GetLxcList(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "Get Lxcsis",
			wantErr: false,
		},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			nodes, err := server.GetNodes()
			if err != nil {
				t.Log(err.Error())
				return
			}
			got, err := nodes[0].GetLxcList()
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.GetLxcsis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) <= 0 {
				t.Errorf("Node.GetLxcsis() error = %v, wantErr %v", "not lxcs received",got)
			}
		})
	}
}
