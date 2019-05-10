package instance

import (
	"encoding/json"
	"fmt"

	"github.com/tdewin/martini-cli/core"
)

type MartiniInstance struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	TenantId string `json:"tenant_id"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Location string `json:"location"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type MartiniInstanceIdentification struct {
	Id         string `json:"id,omitempty"`
	InstanceId string `json:"instanceid,omitempty"`
	Data       string `json:"data,omitempty"`
}

func (i *MartiniInstance) Create(conn *core.Connection) error {
	var err error

	b, err := json.Marshal(i)
	returnstatus := core.ReturnStatus{}
	//json.Unmarshal(txt, returnstatus)
	//err = fmt.Errorf("Not valid return code %d on tenant create %s", sc, returnstatus.Status)
	if err == nil {
		txt, sc, rerr := conn.Post("instance/create", b)

		if rerr == nil {
			json.Unmarshal(txt, &returnstatus)

			if sc != 200 {
				err = fmt.Errorf("Not valid return code %d on tenant create [%s]", sc, returnstatus.Status)
			} else {
				i.Id = returnstatus.Id
			}
		} else {
			err = rerr
		}
	}

	return err
}

func Mappings(conn *core.Connection, tenantid string) (map[string]string, map[string]string, error) {
	instances, err := List(conn, tenantid)

	var idtoname map[string]string
	var nametoid map[string]string

	idtoname = make(map[string]string)
	nametoid = make(map[string]string)

	for _, t := range instances {
		idtoname[t.Id] = t.Name
		nametoid[t.Name] = t.Id
	}
	return idtoname, nametoid, err
}

func List(conn *core.Connection, tenantid string) ([]MartiniInstance, error) {
	var arr []MartiniInstance
	var err error

	idjson, _ := json.Marshal(MartiniInstanceIdentification{Id: tenantid})

	txt, sc, rerr := conn.Post("instance/list", idjson)
	if rerr == nil {
		if sc != 200 {
			rc := core.ReturnStatus{}
			json.Unmarshal(txt, &rc)

			err = fmt.Errorf("Not valid return code %d on instance list; content [%s]", sc, rc.Status)
		} else {
			err = json.Unmarshal(txt, &arr)
			if err != nil {
				err = fmt.Errorf("Unexpected output server on instance list %s (%v)", txt, err)
			}
		}
	} else {
		err = rerr
	}

	return arr, err
}

type BrokerType struct {
	Id       string `json:"id"`
	Clientip string `json:"clientip,omitempty"`
}
type MartiniBrokerEndpoint struct {
	Id             string `json:"id"`
	Status         string `json:"status"`
	Port           string `json:"port,omitempty"`
	ExpectedClient string `json:"expectedclient,omitempty"`
}

func Broker(conn *core.Connection, id string, clientip string) (MartiniBrokerEndpoint, error) {
	var err error
	brokerendpoint := MartiniBrokerEndpoint{}
	b, _ := json.Marshal(BrokerType{id, clientip})

	txt, sc, rerr := conn.Post("instance/broker", b)
	if rerr == nil {
		je := json.Unmarshal(txt, &brokerendpoint)
		if je != nil {
			err = fmt.Errorf("Could not understand result; %v", txt)
		}
		if sc != 200 {
			err = fmt.Errorf("Not valid return code %d on tenant broker; content %s", sc, brokerendpoint.Status)
		}
	} else {
		err = rerr
	}

	return brokerendpoint, err
}
