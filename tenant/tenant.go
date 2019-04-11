package tenant

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tdewin/martini-cli/core"
)

type MartiniTenant struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Instancefqdn     string `json:"instancefqdn"`
	Instanceusername string `json:"instanceusername"`
	Instancepassword string `json:"instancepassword"`
	Id               string `json:"id"`
}

func (m *MartiniTenant) Create(conn *core.Connection) error {
	b, err := json.Marshal(m)

	if err == nil {
		txt, sc, rerr := conn.Post("tenant/create", b)
		if rerr == nil {
			if sc != 200 {
				log.Println("Not valid return code %d", sc)
			} else {
				txtstr := strings.TrimSpace(string(txt))
				if txtstr != "" {
					fmt.Println(txtstr)
				}
			}
		} else {
			err = rerr
		}
	}

	return err
}
func (m *MartiniTenant) Deploy(conn *core.Connection) error {
	b, err := json.Marshal(m)

	if err == nil {
		txt, sc, rerr := conn.Post("tenant/deploy", b)
		if rerr == nil {
			if sc != 200 {
				log.Println("Not valid return code %d", sc)
			} else {
				txtstr := strings.TrimSpace(string(txt))
				if txtstr != "" {
					fmt.Println(txtstr)
				}
			}
		} else {
			err = rerr
		}
	}

	return err
}

func List(conn *core.Connection) ([]MartiniTenant, error) {
	var arr []MartiniTenant
	var err error

	txt, sc, rerr := conn.Get("tenant/list")
	if rerr == nil {
		if sc != 200 {
			log.Println("Not valid return code %d", sc)
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

	txt, sc, rerr := conn.Post("tenant/delete", []byte(fmt.Sprintf("{\"id\":\"%s\"}", id)))
	if rerr == nil {
		if sc != 200 {
			log.Println("Not valid return code %d", sc)
		} else {
			txtstr := strings.TrimSpace(string(txt))
			if txtstr != "" {
				fmt.Println(txtstr)
			}
		}
	} else {
		err = rerr
	}

	return err
}
