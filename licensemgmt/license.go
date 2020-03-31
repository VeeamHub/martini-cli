package licensemgmt

import (
	"encoding/json"
	"fmt"

	"github.com/VeeamHub/martini-cli/core"
)

type MartiniLicenseUser struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	OrganizationId   string `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
	LastBackupDate   string `json:"lastBackupDate"`
}
type MartiniLicenseInfo struct {
	OrgId         string `json:"orgid"`
	OrgName       string `json:"orgname"`
	LicensedUsers int    `json:"licensedUsers"`
	NewUsers      int    `json:"newUsers"`
}

func ListUsers(conn *core.Connection, instance string) ([]MartiniLicenseUser, error) {
	var arr []MartiniLicenseUser
	var err error

	idjson, _ := json.Marshal(core.SendID{Id: instance})

	txt, sc, rerr := conn.Post("license/listusers", idjson)
	if rerr == nil {
		if sc != 200 {
			rc := core.ReturnStatus{}
			json.Unmarshal(txt, &rc)

			err = fmt.Errorf("Not valid return code %d on license user; content [%s]", sc, rc.Status)
		} else {
			err = json.Unmarshal(txt, &arr)
			if err != nil {
				err = fmt.Errorf("Unexpected output server license user %s (%v)", txt, err)
			}
		}
	} else {
		err = rerr
	}

	return arr, err
}

func ListInfo(conn *core.Connection, instance string) ([]MartiniLicenseInfo, error) {
	var arr []MartiniLicenseInfo
	var err error

	idjson, _ := json.Marshal(core.SendID{Id: instance})

	txt, sc, rerr := conn.Post("license/listinformation", idjson)
	if rerr == nil {
		if sc != 200 {
			rc := core.ReturnStatus{}
			json.Unmarshal(txt, &rc)

			err = fmt.Errorf("Not valid return code %d on license info; content [%s]", sc, rc.Status)
		} else {
			err = json.Unmarshal(txt, &arr)
			if err != nil {
				err = fmt.Errorf("Unexpected output server license info %s (%v)", txt, err)
			}
		}
	} else {
		err = rerr
	}

	return arr, err
}
