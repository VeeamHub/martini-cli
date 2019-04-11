package core

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/urfave/cli"
)

type Connection struct {
	server     string
	token      string
	renew      string
	lifetime   int64
	username   string
	password   string
	ignoressc  bool
	client     *http.Client
	serverskew int64
}

func (c *Connection) GetToken() string {
	return c.token
}
func (c *Connection) GetServer() string {
	return c.server
}
func (c *Connection) GetUsername() string {
	return c.username
}
func (c *Connection) GetLifetime() int64 {
	return c.lifetime
}
func (c *Connection) GetRenew() string {
	return c.renew
}
func (c *Connection) GetServerSkew() int64 {
	return c.serverskew
}
func (c *Connection) MakeUrl(urlEnd string) string {
	return fmt.Sprintf("%s/%s", c.server, urlEnd)
}
func (c *Connection) Post(urlEnd string, json []byte) ([]byte, int, error) {
	var body []byte
	statuscode := -1
	req, _ := http.NewRequest("POST", c.MakeUrl(urlEnd), bytes.NewReader(json))
	req.Header.Add("X-Authorization", fmt.Sprintf("bearer %s", c.token))
	resp, err := c.client.Do(req)
	statuscode = resp.StatusCode
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
	}
	return body, statuscode, err
}
func (c *Connection) Get(urlEnd string) ([]byte, int, error) {
	var body []byte
	statuscode := -1
	req, _ := http.NewRequest("GET", c.MakeUrl(urlEnd), nil)
	req.Header.Add("X-Authorization", fmt.Sprintf("bearer %s", c.token))
	resp, err := c.client.Do(req)
	statuscode = resp.StatusCode
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
	}
	return body, statuscode, err
}

func (c *Connection) Auth() error {
	authed := false
	var rerr error

	//first try to heartbeat with token
	//if ok, then we can
	//don't need to care too much about error handling because if it fails, we will just try to authenticate
	if c.token != "" {

		
		//log.Println("Trying ", c.token)
		req, _ := http.NewRequest("GET", c.MakeUrl("login/heartbeat"), nil)

		req.Header.Add("X-Authorization", fmt.Sprintf("bearer %s", c.token))
		resp, err := c.client.Do(req)
		if err == nil {
			if resp.StatusCode == 200 {
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					var hm HeartbeatMessage
					err = json.Unmarshal(body, &hm)
					if err == nil {
						if hm.Status == "ok" {
							authed = true
						} else if hm.Status == "expired" {
							log.Printf("expired token")
						}
					}
				}

			}
		}
	}

	//if we are not authed, our token is not fine so we should continue
	if !authed {
		var login Login
		login.Username = c.username
		login.Password = c.password

		lbyte, _ := json.Marshal(login)
		reader := bytes.NewReader(lbyte)
		req, _ := http.NewRequest("POST", c.MakeUrl("login/create"), reader)
		resp, err := c.client.Do(req)

		if err == nil {
			timenow := time.Now()
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				if resp.StatusCode == 200 {
					var token Token
					err = json.Unmarshal(body, &token)
					if err == nil {
						if token.Token != "" {
							c.token = token.Token
							c.lifetime = token.Lifetime
							c.renew = token.Renew
							c.serverskew = token.ServerTime - timenow.Unix()

							authed = true
						} else {
							rerr = errors.New("Token is empty which should not happen. Are you sure this is a martini server")
						}
					} else {
						rerr = err
					}
				} else {
					rerr = errors.New(fmt.Sprintf("Authentication failed, statuscode %d, body %s", resp.StatusCode, body))
				}
			}
		} else {
			rerr = err
		}

	}

	return rerr
}
func NewConnectionFromCLIContext(c *cli.Context) *Connection {
	return NewConnection(c.GlobalString("server"), c.GlobalString("token"), c.GlobalString("username"), c.GlobalString("password"), c.GlobalBool("ignoreSelfSignedCertificate"))
}

func NewConnection(server string, token string, login string, password string, ignoressc bool) *Connection {
	if server == "" {
		server = "https://localhost/api"
	}
	if login == "" {
		login = "admin"
	}
	tr := &http.Transport{}
	if ignoressc {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	c := Connection{server, token, "", 0, login, password, ignoressc, &http.Client{Transport: tr}, 0}
	return &c
}