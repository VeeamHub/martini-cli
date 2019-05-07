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
	"strings"
	"time"

	"github.com/urfave/cli"
)

type Connection struct {
	server       string
	token        string
	renew        string
	lifetime     int64
	username     string
	password     string
	ignoressc    bool
	client       *http.Client
	serverskew   int64
	PrintOptions *PrintOptions
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

	uri := c.MakeUrl(urlEnd)
	c.PrintOptions.Verbosef("Sending POST to %s %s", uri, string(json))

	req, _ := http.NewRequest("POST", uri, bytes.NewReader(json))
	req.Header.Add("X-Authorization", fmt.Sprintf("bearer %s", c.token))

	resp, err := c.client.Do(req)
	statuscode = resp.StatusCode
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
		c.PrintOptions.Verbosef("(%d) %s", statuscode, string(body))
	} else {
		c.PrintOptions.Verbosef("Error code %v", err)
	}
	return body, statuscode, err
}
func (c *Connection) Get(urlEnd string) ([]byte, int, error) {
	var body []byte
	statuscode := -1

	uri := c.MakeUrl(urlEnd)
	c.PrintOptions.Verbosef("Sending GET to %s ", uri)

	req, _ := http.NewRequest("GET", c.MakeUrl(urlEnd), nil)
	req.Header.Add("X-Authorization", fmt.Sprintf("bearer %s", c.token))

	resp, err := c.client.Do(req)
	statuscode = resp.StatusCode
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
		c.PrintOptions.Verbosef("(%d) %s", statuscode, string(body))
	} else {
		c.PrintOptions.Verbosef("Error code %v", err)
	}
	return body, statuscode, err
}
func (c *Connection) tokenAuthentication(askpassword func() string) (bool, error) {
	authed := false
	var err error
	req, _ := http.NewRequest("GET", c.MakeUrl("login/heartbeat"), nil)

	c.PrintOptions.Verbosef("Trying to heartbeat to login/heartbeat")

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
						c.PrintOptions.Verbosef("Authenticated")
					} else if hm.Status == "expired" {
						c.PrintOptions.Verbosef("Expired token")
					}
				}
			}
		} else {
			c.PrintOptions.Verbosef("Failed heartbeat")
		}
	}
	return authed, err
}
func (c *Connection) Auth(askpassword func() string, doRenew bool) error {
	authed := false
	var rerr error

	//try to heartbeat with token
	//if ok, then we can
	//don't need to care too much about error handling because if it fails, we will just try to authenticate

	if c.lifetime != 0 {
		timenowu := time.Now().Unix()
		timenowatserver := timenowu + c.serverskew
		remaining := c.lifetime - timenowatserver
		delayRefreshUntil := int64(600)
		deadline := int64(15)

		if remaining > delayRefreshUntil && c.token != "" {
			//log.Println("Trying ", c.token)
			c.PrintOptions.Verbosef("Trying token authentication")
			authed, rerr = c.tokenAuthentication(askpassword)
		} else if remaining > deadline {
			if doRenew {
				//refresh token because we want to be kept up to date
				c.PrintOptions.Verbosef("Token is about to expire, renewing via login/renew")
				req, _ := http.NewRequest("POST", c.MakeUrl("login/renew"), strings.NewReader(fmt.Sprintf("{\"renew\":\"%s\"}", c.renew)))
				req.Header.Add("X-Authorization", fmt.Sprintf("bearer %s", c.token))
				resp, err := c.client.Do(req)
				if err == nil {
					if resp.StatusCode == 200 {
						//as close as possible next to request
						timenow := time.Now()
						body, err := ioutil.ReadAll(resp.Body)
						if err == nil {
							var token Token
							err = json.Unmarshal(body, &token)
							if err == nil {
								if token.Token != "" {
									c.token = token.Token
									c.lifetime = token.Lifetime
									c.renew = token.Renew
									c.serverskew = token.ServerTime - timenow.Unix()
									c.PrintOptions.Verbosef("Token is renewed")
									authed = true
								} else {
									rerr = errors.New("Token is empty which should not happen. Are you sure this is a martini server")
									c.PrintOptions.Verbosef("couldn't understand token %v", err)
									//log.Printf("%v", rerr)
								}
							} else {
								rerr = err
								c.PrintOptions.Verbosef("couldn't understand token %v", err)
								//log.Printf("%v", rerr)
							}
						}
					}
				}
			} else {
				c.PrintOptions.Verbosef("You should renew your token soon")
				authed, rerr = c.tokenAuthentication(askpassword)
			}
		}
	}

	//if we are not authed, our token is not fine so we should continue
	if !authed {
		c.PrintOptions.Verbosef("Was not able to authenticate via token, trying via regular admin procedure")

		var login Login
		login.Username = c.username

		if c.password != "" {
			login.Password = c.password
		} else if askpassword != nil {
			login.Password = askpassword()
		}

		lbyte, _ := json.Marshal(login)
		reader := bytes.NewReader(lbyte)
		req, _ := http.NewRequest("POST", c.MakeUrl("login/create"), reader)
		c.PrintOptions.Verbosef("Posting to login/create to make a new session")
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
							c.PrintOptions.Verbosef("Error understanding token %v", err)
						}
					} else {
						rerr = err
					}
				} else {
					rerr = fmt.Errorf("Authentication failed")
					c.PrintOptions.Verbosef("Error understanding token statuscode %d, body %s", resp.StatusCode, body)
				}
			} else {
				rerr = err
			}
		} else {
			rerr = err
		}

	}
	if !authed && rerr == nil {
		rerr = fmt.Errorf("Authentication failed for unknown reason")
	}
	return rerr
}
func NewConnectionFromCLIContext(po *PrintOptions, c *cli.Context) *Connection {
	return NewConnection(po, c.GlobalString("server"), c.GlobalString("token"), c.GlobalString("username"), c.GlobalString("password"), c.GlobalBool("ignoreSelfSignedCertificate"), c.GlobalString("renewtoken"), c.GlobalInt64("renewlifetime"), c.GlobalInt64("renewserverskew"))
}

func NewConnection(po *PrintOptions, server string, token string, login string, password string, ignoressc bool, renew string, lifetime int64, serverskew int64) *Connection {

	if server == "" {
		server = "https://localhost/api"
	}
	if login == "" {
		login = "admin"
	}
	tr := &http.Transport{}
	if ignoressc {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		po.Verbose("Ignoring Self Sign Certificate; Consider a real certificate")
	}
	c := Connection{server, token, renew, lifetime, login, password, ignoressc, &http.Client{Transport: tr}, serverskew, po}
	return &c
}

type PrintOptions struct {
	Json    bool
	verbose bool
}

func (p *PrintOptions) Verbose(txt string) {
	if p.verbose {
		log.Print(txt)
	}
}
func (p *PrintOptions) Verbosef(txt string, v ...interface{}) {
	if p.verbose {
		log.Printf(txt, v...)
	}
}

func (p *PrintOptions) Println(a ...interface{}) {
	if !p.Json {
		fmt.Println(a...)
	}
}
func (p *PrintOptions) Print(txt string) {
	if !p.Json {
		fmt.Print(txt)
	}
}
func (p *PrintOptions) Printf(txt string, v ...interface{}) {
	if !p.Json {
		fmt.Printf(txt, v...)
	}
}

func (p *PrintOptions) PrintJSON(txt string) {
	if p.Json {
		fmt.Println(txt)
	}
}
func (p *PrintOptions) MarshalPrintJSON(m interface{}) {
	if p.Json {
		txt, _ := json.Marshal(m)
		fmt.Println(string(txt))
	}
}

func NewPrintOptionsFromCLIContext(c *cli.Context) PrintOptions {
	return NewPrintOptions(c.GlobalBool("json"), c.GlobalBool("verbose"))
}
func NewPrintOptions(json bool, verbose bool) PrintOptions {
	return PrintOptions{json, verbose}
}
