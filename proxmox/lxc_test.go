package proxmox

import (
	"testing"
)

func TestLxcConfigReceiver_Parse(t *testing.T) {
	type fields struct {
		mp0           string
		mp1           string
		mp2           string
		mp3           string
		mp4           string
		mp5           string
		mp6           string
		mp7           string
		mp8           string
		mp9           string
		net0          string
		net1          string
		net2          string
		net3          string
		net4          string
		net5          string
		net6          string
		net7          string
		net8          string
		net9          string
		startup       string
		lxc           []interface{}
		BaseLxcConfig BaseLxcConfig
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "LxcConfigReceiver.Parse() test",
			fields: fields{
				mp0: "local:999/vm-999-disk-1.raw,mp=test,size=8G",
				net0: "name=eth0,bridge=vmbr0,firewall=1,hwaddr=2A:DC:21:F5:39:46,ip=dhcp,tag=12,type=veth",
				startup:"order=1,up=120,down=120",
			},
		},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lcr := &LxcConfigReceiver{
				mp0:           tt.fields.mp0,
				mp1:           tt.fields.mp1,
				mp2:           tt.fields.mp2,
				mp3:           tt.fields.mp3,
				mp4:           tt.fields.mp4,
				mp5:           tt.fields.mp5,
				mp6:           tt.fields.mp6,
				mp7:           tt.fields.mp7,
				mp8:           tt.fields.mp8,
				mp9:           tt.fields.mp9,
				net0:          tt.fields.net0,
				net1:          tt.fields.net1,
				net2:          tt.fields.net2,
				net3:          tt.fields.net3,
				net4:          tt.fields.net4,
				net5:          tt.fields.net5,
				net6:          tt.fields.net6,
				net7:          tt.fields.net7,
				net8:          tt.fields.net8,
				net9:          tt.fields.net9,
				startup:       tt.fields.startup,
				lxc:           tt.fields.lxc,
				BaseLxcConfig: tt.fields.BaseLxcConfig,
			}

			got := lcr.Parse()
			t.Logf("%v", got)
		})
	}
}
