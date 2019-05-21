package proxmox_test

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
	. "github.com/mrgloba/proxmox-api2/proxmox"
)

const (
	TEST_PROXMOX_VMID = 999
	TEST_PROXMOX_HOST  = "localhost"
	TEST_PROXMOX_PORT  = "8006"
	TEST_PROXMOX_USER  = "testuser"
	TEST_PROXMOX_PASS  = "testuser"
	TEST_PROXMOX_REALM = "pve"

	TEST_PROXMOX_NODE = "pve"

	TEST_PROXMOX_TEMPLATE = "local:vztmpl/debian-9.0-standard_9.0-2_amd64.tar.gz"
	TEST_PROXMOX_STORAGE = "local-lvm"
	TEST_PROXMOX_TEMPLATE_STORAGE = "local"

	TEST_PROXMOX_RELEASE = "5"
	TEST_PROXMOX_REPOID = "c6fdb264"
	TEST_PROXMOX_VERSION = "5.4"

)

var DEBUG_TESTS = true
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

func TestProxmox_MakeAPITarget(t *testing.T) {
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
				host: TEST_PROXMOX_HOST,
				port: TEST_PROXMOX_PORT,
			},
			args: args{
				path: "test/path",
			},
			want:    APITarget("https://" + TEST_PROXMOX_HOST + ":" + TEST_PROXMOX_PORT + "/api2/json/test/path"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := server.MakeAPITarget(tt.args.path)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.host, tt.args.port, tt.args.user, tt.args.pass, tt.args.realm)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.GetAuthTicket()) == 0 {
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
				Release: TEST_PROXMOX_RELEASE,
				Repoid:  TEST_PROXMOX_REPOID,
				Version: TEST_PROXMOX_VERSION,
			},
			wantErr: false,
		},
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
				if got[0].Node != TEST_PROXMOX_NODE {
					return false
				} else {
					return true
				}
			},
			wantErr: false,
		},
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
			got, err := server.GetStorageList()
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

func TestProxmox_DataUnmarshal(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := server.DataUnmarshal(tt.args.body, &tt.args.v, server); (err != nil) != tt.wantErr {
				t.Errorf("Proxmox.DataUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.args.v[0].GetProxmox().GetHost() != TEST_PROXMOX_HOST {
				t.Error("Proxmox.DataUnmarshal() error = Unmarshalled object not filled with Basics")
			}
		})
	}
}
