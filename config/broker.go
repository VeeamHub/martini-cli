package config

import (
	"encoding/json"
	"fmt"

	"github.com/tdewin/martini-cli/core"
)

type Port struct {
	Port string `json:"port"`
}

func BrokerAddPort(conn *core.Connection, port string) error {
	var err error

	json, _ := json.Marshal(Port{port})

	txt, sc, rerr := conn.Post("brokerendpoint/add", json)
	if rerr == nil {
		if sc != 200 {
			err = fmt.Errorf("Not valid return code %d on broker add; content %s", sc, txt)
		}
	} else {
		err = rerr
	}

	return err
}

func BrokerDeletePort(conn *core.Connection, port string) error {
	var err error

	json, _ := json.Marshal(Port{port})

	txt, sc, rerr := conn.Post("brokerendpoint/delete", json)
	if rerr == nil {
		if sc != 200 {
			err = fmt.Errorf("Not valid return code %d on broker add; content %s", sc, txt)
		}
	} else {
		err = rerr
	}

	return err
}
