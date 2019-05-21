package proxmox

import (
	"net/http"
	"crypto/tls"
	"net/url"
	"time"
	"fmt"
	"strings"
	"errors"
	"io/ioutil"
	"encoding/json"
	"reflect"
)

const (
	HTTP_TIMEOUT = 5
	API_TARGET = "/api2/json/"
	API_AUTH_TARGET = "access/ticket"
	API_TOKEN_LIFETIME = 120
	API_TOKEN_UPDATEBEFORE = 5
)

type APICaller interface {
	fillParent(v interface{},parent interface{})
}

type APITarget string

type ProxmoxVersionInfo struct {
	Release string
	Repoid string
	Version string
}

type BasicObject struct {
	parent interface{} `json:"-"`
}

func NewBasicObject(parent interface{}) BasicObject {
	return BasicObject{parent: parent}
}

func (bo BasicObject) GetParent() interface{} {
	return bo.parent
}

func (bo BasicObject) getProxmox(bb BasicObject) *Proxmox {
	return bb.GetProxmox()
}

func (bo BasicObject) GetProxmox() *Proxmox {

	if px, ok:=bo.parent.(*Proxmox); ok && bo.parent != nil {
		return px
	} else if (bo.parent != nil) {
		r := reflect.ValueOf(bo.GetParent())
		f := reflect.Indirect(r).FieldByName("BasicObject")
		return bo.getProxmox(f.Interface().(BasicObject))
	} else {
		return nil
	}
}

func (bo BasicObject) getNode(bb BasicObject) *Node {
	return bb.GetNode()
}

func (bo BasicObject) GetNode() *Node {

	if px,ok:=bo.parent.(*Node); ok {
		return px
	} else {
		r := reflect.ValueOf(bo.GetParent())
		f := reflect.Indirect(r).FieldByName("BasicObject")
		if f.IsValid() {
			return bo.getNode(f.Interface().(BasicObject))
		}
	}

	return nil

}

type Proxmox struct {
	host string
	port string
	user string
	pass string
	realm string
	ticket string
	csrftoken string
	privs map[string]interface{}
	ticketTime time.Time
	*http.Client
}


func New(host,port,user,pass,realm string) (*Proxmox,error) {
	tr := &http.Transport{ TLSClientConfig: &tls.Config{ InsecureSkipVerify: true }, }

	client := &http.Client{
		Transport: tr,
		Timeout: HTTP_TIMEOUT * time.Second,
	}

	p := &Proxmox{
		host: host,
		port: port,
		realm: realm,
		user: user,
		pass: pass,
		Client: client,
	}

	err := p.updateTicket()

	if err != nil {
		return nil,err
	}

	return p,nil
}

func (px *Proxmox) GetAuthTicket() string {
	return px.ticket
}

func (px *Proxmox) GetHost() string {
	return px.host
}

func (px *Proxmox) APICall(method string, target APITarget, data url.Values) ([]byte,int,error){
	if time.Since(px.ticketTime) >= time.Duration(API_TOKEN_LIFETIME - API_TOKEN_UPDATEBEFORE) {
		err := px.updateTicket()
		if err != nil {
			return nil, 0, err
		}
	}

	request, err := http.NewRequest(method, string(target), strings.NewReader( data.Encode()))

	if err != nil {
		return nil, 0, err
	}

	if method == "GET" || method == "DELETE" || method == "POST" {
		request.Header.Add("CSRFPreventionToken",px.csrftoken)
	}

	cookieExpire := px.ticketTime.Add(time.Duration(API_TOKEN_LIFETIME) * time.Minute)
	cookie := &http.Cookie{
		Name: "PVEAuthCookie",
		Value: px.ticket,
		Expires: cookieExpire,
	}

	request.AddCookie(cookie)

	response, err := px.Do(request)

	if err != nil {
		if response != nil {
			return nil, response.StatusCode, err
		} else {
			return nil, 0, err
		}
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	} else {
		return responseBody, response.StatusCode, err
	}

}

func (px *Proxmox) APICall2(method string, target string, data url.Values, result interface{}, ac APICaller) (int, error) {
	apitarget,err := px.MakeAPITarget(target)
	if err != nil {
		return 0, err
	}

	responseData, httpCode, err := px.APICall(method, apitarget, data)
	if err != nil {
		return 0, err
	}
	if httpCode != 200 {
		return httpCode, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}


	jsonErr := px.DataUnmarshal(responseData, result, ac)

	if jsonErr != nil {
		return httpCode, jsonErr
	}

	return httpCode, nil
}

func (px *Proxmox) MakeAPITarget(path string) (APITarget, error){

	apiUrl := "https://" + px.host + ":" +px.port

	u, err := url.ParseRequestURI(apiUrl)

	if err != nil {
		return APITarget(""), err
	}

	u.Path = API_TARGET + path

	urlStr := fmt.Sprintf("%v", u)

	apiTarget := APITarget(urlStr)

	return apiTarget,nil
}

func (px *Proxmox) updateTicket() (error){

	var csrftoken, ticket string
	var privs map[string]interface{}

	authTarget, err := px.MakeAPITarget(API_AUTH_TARGET)

	if err != nil {
		return err
	}

	data:= url.Values{}

	data.Set("username", px.user + "@" + px.realm)
	data.Add("password", px.pass)

	request, err := http.NewRequest("POST", string(authTarget), strings.NewReader( data.Encode()))

	if err != nil {
		return err
	}

	response, err := px.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}

	body,err := ioutil.ReadAll(response.Body)


	if err != nil {
		return err
	} else {
		b := []byte( string(body) )

		var f map[string]interface{}

		err := json.Unmarshal(b,&f)

		if err != nil {
			return err
		}

		jsondata := f["data"].(map[string]interface{})
		privs = jsondata["cap"].(map[string]interface{})
		csrftoken = jsondata["CSRFPreventionToken"].(string)
		ticket = jsondata["ticket"].(string)
	}


	px.ticket = ticket
	px.csrftoken = csrftoken
	px.ticketTime = time.Now()
	px.privs = privs

	return nil
}

func (px *Proxmox) DataUnmarshal(body []byte, v interface{}, ac APICaller) error {
	var f map[string]interface{}

	err := json.Unmarshal(body, &f)
	if err != nil {
		return err
	}

	temp, err := json.Marshal(f["data"])
	if err != nil {
		return err
	}

	mErr := json.Unmarshal( temp, v )
	if mErr != nil {
		return mErr
	}

	if ac!=nil {
		ac.fillParent(v,px)
	}

	return nil
}

func (px *Proxmox) fillParent(v interface{}, parent interface{}) {
	rvt := reflect.TypeOf(v).Elem()
	rvv := reflect.ValueOf(v).Elem()

	switch rvt.Kind() {
	case reflect.Slice:
		for i:=0; i<rvv.Len(); i++ {
			px.fillParent(rvv.Index(i).Addr().Interface(),parent)
		}
	case reflect.Struct:
		value := rvv.FieldByName("BasicObject")
		if value.CanSet() {
			value.Set(reflect.ValueOf(BasicObject{parent: parent}))
		}
	default:
	}

}

func (px *Proxmox) GetProxmoxVersion() (*ProxmoxVersionInfo,error) {

	target,err := px.MakeAPITarget("version")
	if err != nil {
		return nil, err
	}

	responseData, httpCode, err := px.APICall("GET", target, nil)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}

	var versionInfo ProxmoxVersionInfo

	jsonErr := px.DataUnmarshal(responseData, &versionInfo,px)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return &versionInfo, nil
}

func (px *Proxmox) GetNodes()([]Node, error) {
	target,err := px.MakeAPITarget("nodes")
	if err != nil {
		return nil, err
	}

	responseData, httpCode, err := px.APICall("GET", target, nil)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}

	var nodes []Node

	jsonErr := px.DataUnmarshal(responseData, &nodes,px)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return nodes, nil
}

func (px *Proxmox) GetNode(nodeName string) (*Node, error){
	nodes, err := px.GetNodes()

	if err != nil { return nil, err}

	for _,v:= range nodes {
		if v.Node == nodeName {
			return &v,nil
		}
	}

	return nil, errors.New("node not found")
}

func (px *Proxmox) GetStorageList()([]Storage,error){
	target,err := px.MakeAPITarget("storage")
	if err != nil {
		return nil, err
	}

	responseData, httpCode, err := px.APICall("GET", target, nil)
	if err != nil {
		return nil, err
	}
	if httpCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Request return error: %d",httpCode))
	}

	var storages []Storage

	jsonErr := px.DataUnmarshal(responseData, &storages,px)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return storages, nil
}