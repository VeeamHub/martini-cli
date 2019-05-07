package tenant

import (
	"encoding/json"
	"fmt"

	"github.com/tdewin/martini-cli/core"
)

type MartiniTenant struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Instancefqdn     string `json:"instancefqdn"`
	Instanceport     string `json:"instanceport"`
	Instanceusername string `json:"instanceusername"`
	Instancepassword string `json:"instancepassword"`
	Id               string `json:"id"`
}

func (m *MartiniTenant) Create(conn *core.Connection) error {
	b, err := json.Marshal(m)
	returnstatus := core.ReturnStatus{}
	//json.Unmarshal(txt, returnstatus)
	//err = fmt.Errorf("Not valid return code %d on tenant create %s", sc, returnstatus.Status)
	if err == nil {
		txt, sc, rerr := conn.Post("tenant/create", b)

		if rerr == nil {
			json.Unmarshal(txt, &returnstatus)

			if sc != 200 {
				err = fmt.Errorf("Not valid return code %d on tenant create [%s]", sc, returnstatus.Status)
			}
		} else {
			err = rerr
		}
	}

	return err
}
func Mappings(conn *core.Connection) (map[string]string, map[string]string, error) {
	tenants, err := List(conn)

	var idtoname map[string]string
	var nametoid map[string]string

	idtoname = make(map[string]string)
	nametoid = make(map[string]string)

	for _, t := range tenants {
		idtoname[t.Id] = t.Name
		nametoid[t.Name] = t.Id
	}
	return idtoname, nametoid, err
}
func List(conn *core.Connection) ([]MartiniTenant, error) {
	var arr []MartiniTenant
	var err error

	returnstatus := core.ReturnStatus{}
	//json.Unmarshal(txt, returnstatus)
	//err = fmt.Errorf("Not valid return code %d on tenant create %s", sc, returnstatus.Status)

	txt, sc, rerr := conn.Get("tenant/list")
	if rerr == nil {
		if sc != 200 {
			json.Unmarshal(txt, &returnstatus)
			err = fmt.Errorf("Not valid return code %d on tenant list [%s]", sc, returnstatus.Status)
		} else {
			err = json.Unmarshal(txt, &arr)
		}
	} else {
		err = rerr
	}

	return arr, err
}

func Delete(conn *core.Connection, id string) error {
	var err error
	returnstatus := core.ReturnStatus{}
	//json.Unmarshal(txt, returnstatus)
	//err = fmt.Errorf("Not valid return code %d on tenant create %s", sc, returnstatus.Status)

	txt, sc, rerr := conn.Post("tenant/delete", []byte(fmt.Sprintf("{\"id\":\"%s\"}", id)))
	if rerr == nil {
		json.Unmarshal(txt, &returnstatus)
		if sc != 200 {
			err = fmt.Errorf("Not valid return code %d on tenant delete [%s]", sc, returnstatus.Status)
		}
	} else {
		err = rerr
	}

	return err
}

func Resolve(conn *core.Connection, tenantname string) (string, error) {
	tenants, err := List(conn)
	tenantid := "-1"
	if err == nil {
		for _, t := range tenants {
			if t.Name == tenantname {
				tenantid = t.Id
			}
		}
	}
	if tenantid == "-1" {
		err = fmt.Errorf("Could not find tenant with id %s", tenantname)
	}
	return tenantid, err
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

	txt, sc, rerr := conn.Post("tenant/broker", b)
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

type MartiniDeploy struct {
	Type   string      `json:"type"`
	Name   string      `json:"name"`
	Email  string      `json:"email"`
	Config interface{} `json:"config"`
}
type MartiniAmazon struct {
	Region string `json:"region"`
}

func NewAWSConfig(name string, email string, region string) *MartiniDeploy {
	return &MartiniDeploy{"aws", name, email, MartiniAmazon{region}}
}

func (m *MartiniDeploy) Deploy(conn *core.Connection) error {
	b, err := json.Marshal(m)
	returnstatus := core.ReturnStatus{}
	//json.Unmarshal(txt, returnstatus)
	//err = fmt.Errorf("Not valid return code %d on tenant create %s", sc, returnstatus.Status)
	if err == nil {
		txt, sc, rerr := conn.Post("tenant/deploy", b)
		if rerr == nil {
			json.Unmarshal(txt, &returnstatus)
			if sc != 200 {
				err = fmt.Errorf("Not valid return code %d on tenant deploy [%s]", sc, returnstatus.Status)
			}
		} else {
			err = rerr
		}
	}

	return err
}
