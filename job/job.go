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
type MartiniJobIdentification struct {
	Id    string `json:"id,omitempty"`
	JobId string `json:"jobid,omitempty"`
	Data  string `json:"data,omitempty"`
}

func List(conn *core.Connection, instance string) ([]MartiniJob, error) {
	var arr []MartiniJob
	var err error

	idjson, _ := json.Marshal(MartiniJobIdentification{Id: instance})

	txt, sc, rerr := conn.Post("job/list", idjson)
	if rerr == nil {
		if sc != 200 {
			rc := core.ReturnStatus{}
			json.Unmarshal(txt, &rc)

			err = fmt.Errorf("Not valid return code %d on job start; content [%s]", sc, rc.Status)
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
func Resolve(conn *core.Connection, instance string, jobname string) (string, error) {
	//job.Resolve(conn, id, jobname)
	var err error
	var id = "-1"

	jobs, rerr := List(conn, instance)
	if rerr == nil {
		for _, j := range jobs {
			if j.Name == jobname {
				id = j.Id
			}
		}
	} else {
		err = rerr
	}
	if id == "-1" {
		err = fmt.Errorf("Was not able to resolve id for %s", jobname)
	}

	return id, err
}

func Start(conn *core.Connection, instance string, jobid string) error {
	var err error

	idjson, _ := json.Marshal(MartiniJobIdentification{Id: instance, JobId: jobid})

	txt, sc, rerr := conn.Post("job/start", idjson)
	if rerr == nil {
		if sc != 200 {
			rc := core.ReturnStatus{}
			json.Unmarshal(txt, &rc)

			err = fmt.Errorf("Not valid return code %d on job start; content [%s]", sc, rc.Status)
		}
	} else {
		err = rerr
	}

	return err
}
