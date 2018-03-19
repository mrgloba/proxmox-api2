package proxmox

import (
	"strconv"
	"strings"
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

			if DEBUG_TESTS {
				t.Logf("StorageList: %v\n", got)
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
			name:    "GetLxcList",
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
				t.Errorf("Node.GetLxcList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if DEBUG_TESTS {
				t.Logf("LxcList: %v\n", got)
			}

			if len(got) <= 0 {
				t.Errorf("Node.GetLxcList() error = %v, wantErr %v", "not lxcs received", got)
			}
		})
	}
}

func TestNode_GetLxc(t *testing.T) {
	type args struct {
		vmid int64
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Get Lxc container",
			args:    args{vmid: 101},
			want:    101,
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

			got, err := nodes[0].GetLxc(tt.args.vmid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.GetLxc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if DEBUG_TESTS {
				t.Logf("Lxc: %v\n", got)
			}

			if int(got.VmId) != tt.want {
				t.Errorf("Node.GetLxc() = %v, want %v", got.VmId, tt.want)
			}
		})
	}
}

func TestNode_CreateLxc(t *testing.T) {
	type args struct {
		lxcParams LxcConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Node.CreateLxc() test",
			args: args{
				lxcParams: LxcConfig{
					BaseLxcConfig: BaseLxcConfig{
						VmId:         TEST_PROXMOX_VMID,
						Hostname:     "test1",
						Password:     "111111",
						OSTemplate:   "local:vztmpl/debian-9.0-standard_9.0-2_amd64.tar.gz",
						RootFS:       "storage:10",
						Cores:        2,
						Memory:       2048,
						Swap:         1024,
						SearchDomain: "test.loc",
						NameServer:   "8.8.8.8",
					},
					Networks: []NetworkConfig{
						{
							Name:      "eth0",
							Bridge:    "vmbr0",
							IPAddress: "10.10.10.0/24",
							Gateway:   "10.10.10.1",
							Tag:       99,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := server.GetNodes()
			if err != nil {
				t.Log(err.Error())
				return
			}
			got, err := nodes[0].CreateLxc(tt.args.lxcParams)

			if (err != nil) != tt.wantErr {
				t.Errorf("Node.CreateLxc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if DEBUG_TESTS {
				t.Logf("%v\n", *got)
			}

			parts := strings.Split(string(*got), ":")
			idx, _ := strconv.Atoi(parts[6])

			if idx != TEST_PROXMOX_VMID {
				t.Errorf("Node.CreateLxc() test failed")
			}

		})
	}
}

func TestNode_RemoveLxc(t *testing.T) {
	type args struct {
		vmid int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Remove test lxc",
			args:    args{TEST_PROXMOX_VMID},
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

			got, err := nodes[0].RemoveLxc(tt.args.vmid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.RemoveLxc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if DEBUG_TESTS {
				t.Logf("Lxc: %v\n", *got)
			}

			if strings.Index(string(*got), nodes[0].Node) <= 0 {
				t.Errorf("Node.RemoveLxc() = %v", got)
			}
		})
	}
}

func TestNode_GetTasks(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Node.GetTasks() test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := server.GetNodes()
			if err != nil {
				t.Log(err.Error())
				return
			}
			got, err := nodes[0].GetTasks()
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.GetTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if DEBUG_TESTS {
				t.Logf("%v\n", got)
			}

			if len(got) == 0 {
				t.Errorf("Node.GetTasks() = %v", got)
			}
		})
	}
}

func TestNode_ScanUsb(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{ name: "Node.ScanUSB test", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := server.GetNodes()
			if err != nil {
				t.Log(err.Error())
				return
			}
			got, err := nodes[0].ScanUsb()
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.ScanUsb() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if DEBUG_TESTS {
				t.Logf("%v\n", got)
			}

			if len(got) == 0 {
				t.Errorf("Node.ScanUsb() = %v", got)
			}
		})
	}
}
