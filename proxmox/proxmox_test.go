package proxmox

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	TEST_PROXMOX_VMID = 999
)

var DEBUG_TESTS bool = true
var server *Proxmox

func TestMain(m *testing.M) {

	if setup() != nil {
		println("Setup test environment failed!")
		return
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func setup() error {
	var err error
	server, err = New(TEST_PROXMOX_HOST, TEST_PROXMOX_PORT, TEST_PROXMOX_USER, TEST_PROXMOX_PASS, TEST_PROXMOX_REALM)
	return err
}

func TestProxmox_makeAPITarget(t *testing.T) {
	type fields struct {
		host       string
		port       string
		user       string
		pass       string
		realm      string
		ticket     string
		csrftoken  string
		privs      map[string]interface{}
		ticketTime time.Time
		Client     *http.Client
	}
	type args struct {
		path string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    APITarget
		wantErr bool
	}{
		{
			name: "Make API target URL",
			fields: fields{
				host: "localhost",
				port: "8006",
			},
			args: args{
				path: "test/path",
			},
			want:    APITarget("https://localhost:8006/api2/json/test/path"),
			wantErr: false,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			px := &Proxmox{
				host:       tt.fields.host,
				port:       tt.fields.port,
				user:       tt.fields.user,
				pass:       tt.fields.pass,
				realm:      tt.fields.realm,
				ticket:     tt.fields.ticket,
				csrftoken:  tt.fields.csrftoken,
				privs:      tt.fields.privs,
				ticketTime: tt.fields.ticketTime,
				Client:     tt.fields.Client,
			}
			got, err := px.makeAPITarget(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Proxmox.makeAPITarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Proxmox.makeAPITarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		host  string
		port  string
		user  string
		pass  string
		realm string
	}
	tests := []struct {
		name    string
		args    args
		want    *Proxmox
		wantErr bool
	}{
		{
			name: "Create Proxmox object and obtain API ticket",
			args: args{
				host:  TEST_PROXMOX_HOST,
				port:  TEST_PROXMOX_PORT,
				user:  TEST_PROXMOX_USER,
				pass:  TEST_PROXMOX_PASS,
				realm: TEST_PROXMOX_REALM,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.host, tt.args.port, tt.args.user, tt.args.pass, tt.args.realm)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.ticket) == 0 {
				t.Errorf("New() ticket is not set")
				return
			}
		})
	}
}

func TestProxmox_GetProxmoxVersion(t *testing.T) {
	tests := []struct {
		name    string
		want    *ProxmoxVersionInfo
		wantErr bool
	}{
		{
			name: "Get proxmox version info",
			want: &ProxmoxVersionInfo{
				Release: "18",
				Repoid:  "ef2610e8",
				Version: "4.4",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.GetProxmoxVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("Proxmox.GetProxmoxVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("Proxmox.GetProxmoxVersion() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestProxmox_GetNodes(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(got []Node) bool
		wantErr  bool
	}{
		{
			name: "Get nodes",
			testFunc: func(got []Node) bool {
				if len(got) == 0 {
					return false
				}
				if got[0].Node != "utm-other" {
					return false
				} else {
					return true
				}
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.GetNodes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Proxmox.GetNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.testFunc(got) {
				t.Errorf("Proxmox.GetNodes() error = %v", got)
				return
			}
		})
	}
}

func TestProxmox_GetStorages(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(got []Storage) bool
		wantErr  bool
	}{
		{
			name: "Get nodes",
			testFunc: func(got []Storage) bool {
				if len(got) == 0 {
					return false
				}
				if got[0].Storage == "local" || got[0].Storage == "storage" {
					return true
				} else {
					return false
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.GetStorages()
			if (err != nil) != tt.wantErr {
				t.Errorf("Proxmox.GetStorages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.testFunc(got) {
				t.Errorf("Proxmox.GetStorages() error = %v", got)
				return
			}
		})
	}
}

func TestProxmox_dataUnmarshal(t *testing.T) {
	type fields struct {
		host       string
		port       string
	}
	type args struct {
		body []byte
		v    []Node
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:"Data Unmarshal",
			fields: fields{ host:"localhost", port:"8006"},
			args: args{
				body: []byte("{\"data\":[{\"maxdisk\":49076379648,\"level\":\"\",\"cpu\":0.0366875975348971,\"node\":\"utm-other\",\"disk\":2129575936,\"id\":\"node/utm-other\",\"maxmem\":33420623872,\"type\":\"node\",\"maxcpu\":4,\"mem\":8576897024,\"uptime\":4649561}]}"),
				v: []Node{},
			},
			wantErr: false,
		},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			px := &Proxmox{
				host:       tt.fields.host,
				port:       tt.fields.port,
			}
			if err := px.dataUnmarshal(tt.args.body, &tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Proxmox.dataUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.args.v[0].px.host != "localhost" {
				t.Error("Proxmox.dataUnmarshal() error = Unmarshalled object not filled with Basics")
			}
		})
	}
}
