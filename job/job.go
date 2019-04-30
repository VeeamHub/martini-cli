package job

import (
	"encoding/json"
	"fmt"

	"github.com/tdewin/martini-cli/core"
)

type MartiniJob struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	LastRun    string `json:"lastRun"`
	LastStatus string `json:"lastStatus"`
}
type IDWrapper struct {
	Id    string `json:"id,omitempty"`
	JobId string `json:"jobid,omitempty"`
	Data  string `json:"data,omitempty"`
}

func List(conn *core.Connection, instance string) ([]MartiniJob, error) {
	var arr []MartiniJob
	var err error

	idjson, _ := json.Marshal(IDWrapper{Id: instance})

	txt, sc, rerr := conn.Post("job/list", idjson)
	if rerr == nil {
		if sc != 200 {
			err = fmt.Errorf("Not valid return code %d on job list; content %s", sc, txt)
		} else {
			err = json.Unmarshal(txt, &arr)
			if err != nil {
				err = fmt.Errorf("Unexpected output server job list %s (%v)", txt, err)
			}
		}
	} else {
		err = rerr
	}

	return arr, err
}

func Start(conn *core.Connection, instance string, jobid string) error {
	var err error

	idjson, _ := json.Marshal(IDWrapper{Id: instance, JobId: jobid})

	txt, sc, rerr := conn.Post("job/start", idjson)
	if rerr == nil {
		if sc != 200 {
			err = fmt.Errorf("Not valid return code %d on job start; content %s", sc, txt)
		}
	} else {
		err = rerr
	}

	return err
}
