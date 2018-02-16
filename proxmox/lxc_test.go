package proxmox

import (
	"reflect"
	"strings"
	"testing"
)

func TestLxcConfigReceiver_Parse(t *testing.T) {
	type fields struct {
		Mp0           string
		Mp1           string
		Mp2           string
		Mp3           string
		Mp4           string
		Mp5           string
		Mp6           string
		Mp7           string
		Mp8           string
		Mp9           string
		Net0          string
		Net1          string
		Net2          string
		Net3          string
		Net4          string
		Net5          string
		Net6          string
		Net7          string
		Net8          string
		Net9          string
		Startup       string
		Lxc           []interface{}
		BaseLxcConfig BaseLxcConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   LxcConfig
	}{
		{
			name: "LxcConfigReceiver.Parse() test",
			fields: fields{
				Mp1:     "local:999/vm-999-disk-1.raw,mp=test,size=8G",
				Net0:    "name=eth0,bridge=vmbr0,firewall=1,hwaddr=2A:DC:21:F5:39:46,ip=dhcp,tag=12,type=veth",
				Startup: "order=1,up=120,down=120",
			},
			want: LxcConfig{
				MountPoints: []MountPoint{{Index: 1, Volume: "local:999/vm-999-disk-1.raw", Path: "test", BaseStorageItem: BaseStorageItem{Size: 8}}},
				Networks:    []NetworkConfig{{Index: 0, Name: "eth0", Bridge: "vmbr0", Firewall: true, HWAddr: "2A:DC:21:F5:39:46", IPAddress: "dhcp", Tag: 12, Type: "veth"}},
				Startup:     StartupConfig{Order: 1, UpDelay: 120, DownDelay: 120},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lcr := &LxcConfigReceiver{
				Mp0:           tt.fields.Mp0,
				Mp1:           tt.fields.Mp1,
				Mp2:           tt.fields.Mp2,
				Mp3:           tt.fields.Mp3,
				Mp4:           tt.fields.Mp4,
				Mp5:           tt.fields.Mp5,
				Mp6:           tt.fields.Mp6,
				Mp7:           tt.fields.Mp7,
				Mp8:           tt.fields.Mp8,
				Mp9:           tt.fields.Mp9,
				Net0:          tt.fields.Net0,
				Net1:          tt.fields.Net1,
				Net2:          tt.fields.Net2,
				Net3:          tt.fields.Net3,
				Net4:          tt.fields.Net4,
				Net5:          tt.fields.Net5,
				Net6:          tt.fields.Net6,
				Net7:          tt.fields.Net7,
				Net8:          tt.fields.Net8,
				Net9:          tt.fields.Net9,
				Startup:       tt.fields.Startup,
				Lxc:           tt.fields.Lxc,
				BaseLxcConfig: tt.fields.BaseLxcConfig,
			}

			got := lcr.Parse()
			if DEBUG_TESTS {
				t.Logf("%v", got)
			}

			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("LxcConfigReceiver.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartupConfig_SetFromString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want StartupConfig
	}{
		{
			name: "StartupConfig.SetFromString() test",
			args: args{str: "order=1,up=120,down=120"},
			want: StartupConfig{Order: 1, UpDelay: 120, DownDelay: 120},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := StartupConfig{}
			sc.SetFromString(tt.args.str)

			if DEBUG_TESTS {
				t.Logf("%v\n", sc)
			}
			if !reflect.DeepEqual(sc, tt.want) {
				t.Errorf("StartupConfig.SetFromString() = %v, want %v", sc, tt.want)
			}
		})
	}
}

func TestNetworkConfig_SetFromString(t *testing.T) {
	type args struct {
		idx int
		str string
	}
	tests := []struct {
		name string
		args args
		want NetworkConfig
	}{
		{
			name: "NetworkConfig.SetFromString() test",
			args: args{
				idx: 0,
				str: "name=eth0,bridge=vmbr0,firewall=1,gw=10.10.10.1,hwaddr=2A:DC:21:F5:39:46,ip=10.10.10.2/24,ip6=auto,tag=12,type=veth",
			},
			want: NetworkConfig{
				Name:       "eth0",
				Bridge:     "vmbr0",
				Firewall:   true,
				Gateway:    "10.10.10.1",
				HWAddr:     "2A:DC:21:F5:39:46",
				IPAddress:  "10.10.10.2/24",
				IPAddresV6: "auto",
				Tag:        12,
				Type:       "veth",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := NetworkConfig{}
			nc.SetFromString(tt.args.idx, tt.args.str)

			if DEBUG_TESTS {
				t.Logf("%v\n", nc)
			}

			if !reflect.DeepEqual(nc, tt.want) {
				t.Errorf("NetworkConfig.SetFromString() = %v, want %v", nc, tt.want)
			}
		})
	}
}

func TestMountPoint_SetFromString(t *testing.T) {
	type args struct {
		idx int
		str string
	}
	tests := []struct {
		name string
		args args
		want MountPoint
	}{
		{
			name: "MountPoint.SetFromString() test",
			args: args{idx: 0, str: "local:999/vm-999-disk-2.raw,mp=/var/lib/vz/1,size=8G"},
			want: MountPoint{
				Index:           0,
				Volume:          "local:999/vm-999-disk-2.raw",
				Path:            "/var/lib/vz/1",
				BaseStorageItem: BaseStorageItem{Size: 8},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := MountPoint{}
			mp.SetFromString(tt.args.idx, tt.args.str)

			if DEBUG_TESTS {
				t.Logf("%v\n", mp)
			}

			if !reflect.DeepEqual(mp, tt.want) {
				t.Errorf("MountPoint.SetFromString() = %v, want %v", mp, tt.want)
			}
		})
	}
}

func TestStartupConfig_String(t *testing.T) {
	type fields struct {
		Order     int
		UpDelay   int
		DownDelay int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "StartupConfig.String() test1",
			fields: fields{Order: 1, UpDelay: 120, DownDelay: 120},
			want:   "order=1,up=120,down=120",
		},
		{
			name:   "StartupConfig.String() test2",
			fields: fields{Order: 1, DownDelay: 120},
			want:   "order=1,down=120",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &StartupConfig{
				Order:     tt.fields.Order,
				UpDelay:   tt.fields.UpDelay,
				DownDelay: tt.fields.DownDelay,
			}

			got := sc.String()

			if DEBUG_TESTS {
				t.Logf("%v\n", got)
			}

			if got != tt.want {
				t.Errorf("StartupConfig.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkConfig_String(t *testing.T) {
	type fields struct {
		Index      int
		Name       string
		Bridge     string
		Firewall   bool
		Gateway    string
		GatewayV6  string
		HWAddr     string
		IPAddress  string
		IPAddresV6 string
		MTU        int
		Rate       int
		Tag        int
		Trunks     string
		Type       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "NetworkConfig.String() test1",
			fields: fields{
				Name:       "eth0",
				Bridge:     "vmbr0",
				Firewall:   true,
				Gateway:    "10.10.10.1",
				HWAddr:     "2A:DC:21:F5:39:46",
				IPAddress:  "10.10.10.2/24",
				IPAddresV6: "auto",
				Tag:        12,
				Type:       "veth",
			},
			want: "name=eth0,bridge=vmbr0,firewall=1,gw=10.10.10.1,hwaddr=2A:DC:21:F5:39:46,ip=10.10.10.2/24,ip6=auto,tag=12,type=veth",
		},
		{
			name: "NetworkConfig.String() test2",
			fields: fields{
				Name:       "eth0",
				Bridge:     "vmbr0",
				Firewall:   false,
				Gateway:    "10.10.10.1",
				HWAddr:     "2A:DC:21:F5:39:46",
				IPAddress:  "10.10.10.2/24",
				IPAddresV6: "auto",
				Tag:        0,
				Type:       "veth",
			},
			want: "name=eth0,bridge=vmbr0,gw=10.10.10.1,hwaddr=2A:DC:21:F5:39:46,ip=10.10.10.2/24,ip6=auto,type=veth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := &NetworkConfig{
				Index:      tt.fields.Index,
				Name:       tt.fields.Name,
				Bridge:     tt.fields.Bridge,
				Firewall:   tt.fields.Firewall,
				Gateway:    tt.fields.Gateway,
				GatewayV6:  tt.fields.GatewayV6,
				HWAddr:     tt.fields.HWAddr,
				IPAddress:  tt.fields.IPAddress,
				IPAddresV6: tt.fields.IPAddresV6,
				MTU:        tt.fields.MTU,
				Rate:       tt.fields.Rate,
				Tag:        tt.fields.Tag,
				Trunks:     tt.fields.Trunks,
				Type:       tt.fields.Type,
			}
			got := nc.String()
			if DEBUG_TESTS {
				t.Logf("%v\n", got)
			}

			if got != tt.want {
				t.Errorf("NetworkConfig.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMountPoint_String(t *testing.T) {
	type fields struct {
		Index           int
		Volume          string
		Path            string
		BaseStorageItem BaseStorageItem
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "MountPoint.String() test1",
			fields: fields{Index: 0, Volume: "local:999/vm-999-disk-2.raw", Path: "/var/lib/vz/1", BaseStorageItem: BaseStorageItem{Size: 8}},
			want:   "local:999/vm-999-disk-2.raw,mp=/var/lib/vz/1,size=8G",
		},
		{
			name:   "MountPoint.String() test2",
			fields: fields{Index: 0, Volume: "local:8"},
			want:   "local:8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &MountPoint{
				Index:           tt.fields.Index,
				Volume:          tt.fields.Volume,
				Path:            tt.fields.Path,
				BaseStorageItem: tt.fields.BaseStorageItem,
			}

			got := mp.String()

			if DEBUG_TESTS {
				t.Logf("%v\n", got)
			}

			if got != tt.want {
				t.Errorf("MountPoint.String() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestLxc_Start(t *testing.T) {
	type fields struct {
		Pid         int
		LxcBase     LxcBase
		BasicObject BasicObject
	}
	tests := []struct {
		name    string
		fields  fields
		want    *TaskID
		wantErr bool
	}{
		{
			name:    "Lxc.Start() test",
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

			lxc, err := nodes[0].GetLxc(TEST_PROXMOX_VMID)

			if err != nil {
				t.Log(err.Error())
				return
			}

			got, err := lxc.Start(false)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lxc.Start() error = %v, wantErr %v", err, tt.wantErr)
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

func TestLxc_Stop(t *testing.T) {
	type fields struct {
		Pid         int
		LxcBase     LxcBase
		BasicObject BasicObject
	}
	tests := []struct {
		name    string
		fields  fields
		want    *TaskID
		wantErr bool
	}{
		{
			name:    "Lxc.Stop() test",
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

			lxc, err := nodes[0].GetLxc(TEST_PROXMOX_VMID)

			if err != nil {
				t.Log(err.Error())
				return
			}

			got, err := lxc.Stop(false)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lxc.Stop() error = %v, wantErr %v", err, tt.wantErr)
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

func TestLxc_Shutdown(t *testing.T) {
	type fields struct {
		Pid         int
		LxcBase     LxcBase
		BasicObject BasicObject
	}
	tests := []struct {
		name    string
		fields  fields
		want    *TaskID
		wantErr bool
	}{
		{
			name:    "Lxc.Shutdown() test",
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

			lxc, err := nodes[0].GetLxc(TEST_PROXMOX_VMID)

			if err != nil {
				t.Log(err.Error())
				return
			}

			got, err := lxc.Shutdown(false,0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lxc.Shutdown() error = %v, wantErr %v", err, tt.wantErr)
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
