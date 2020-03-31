package config

import (
	"encoding/json"
	"fmt"

	"github.com/VeeamHub/martini-cli/core"
)

type PortList struct {
	PortList []Port `json:"portlist"`
}
type Port struct {
	Port string `json:"port"`
}

func BrokerAddPort(conn *core.Connection, port string) error {
	var err error

	jsonp, _ := json.Marshal(Port{port})

	txt, sc, rerr := conn.Post("brokerendpoint/add", jsonp)
	if rerr == nil {
		if sc != 200 {
			rc := core.ReturnStatus{}
			json.Unmarshal(txt, &rc)
			err = fmt.Errorf("Not valid return code %d on broker endpoint add; content [%s]", sc, rc.Status)
		}
	} else {
		err = rerr
	}

	return err
}

func BrokerDeletePort(conn *core.Connection, port string) error {
	var err error

	jsonp, _ := json.Marshal(Port{port})

	txt, sc, rerr := conn.Post("brokerendpoint/delete", jsonp)
	if rerr == nil {
		if sc != 200 {
			rc := core.ReturnStatus{}
			json.Unmarshal(txt, &rc)
			err = fmt.Errorf("Not valid return code %d on broker endpoint delete; content [%s]", sc, rc.Status)
		}
	} else {
		err = rerr
	}

	return err
}

func BrokerList(conn *core.Connection) (PortList, error) {
	var err error

	ports := PortList{}

	txt, sc, rerr := conn.Get("brokerendpoint/list")
	if rerr == nil {
		if sc != 200 {
			rc := core.ReturnStatus{}
			json.Unmarshal(txt, &rc)
			err = fmt.Errorf("Not valid return code %d on broker list; content [%s]", sc, rc.Status)
		} else {
			err = json.Unmarshal(txt, &ports)
		}

	} else {
		err = rerr
	}

	return ports, err
}
