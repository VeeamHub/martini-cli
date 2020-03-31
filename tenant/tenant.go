package tenant

import (
	"encoding/json"
	"fmt"

	"github.com/VeeamHub/martini-cli/core"
)

//password only when new
type MartiniTenant struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Registered string `json:"registered"`
	Password   string `json:"password"`
	Id         string `json:"id"`
}

func (m *MartiniTenant) Create(conn *core.Connection) error {
	b, err := json.Marshal(m)
	returnstatus := core.ReturnStatus{}
	//json.Unmarshal(txt, returnstatus)
	//err = fmt.Errorf("Not valid return code %d on tenant create %s", sc, returnstatus.Status)
	if err == nil {
		txt, sc, rerr := conn.Post("tenant/create", b)

		if rerr == nil {

			if sc != 200 {
				err = fmt.Errorf("Not valid return code %d on tenant create [%s]", sc, returnstatus.Status)
			} else {
				rerr := json.Unmarshal(txt, m)
				if rerr != nil {
					err = fmt.Errorf("Could not understand response from server")
				}

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

	b, _ := json.Marshal(core.SendID{Id: id})

	txt, sc, rerr := conn.Post("tenant/delete", b)
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

func ReverseResolve(conn *core.Connection, tenantid string) (string, error) {
	tenants, err := List(conn)
	tenantname := "<unresolved>"
	if err == nil {
		for _, t := range tenants {
			if t.Id == tenantid {
				tenantname = t.Name
			}
		}
	}
	if tenantname == "<unresolved>" {
		err = fmt.Errorf("Could not find tenant with id %s", tenantid)
	}
	return tenantname, err
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
		err = fmt.Errorf("Could not find tenant with name %s", tenantname)
	}
	return tenantid, err
}
