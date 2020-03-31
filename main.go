package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"syscall"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/VeeamHub/martini-cli/commands"
	"github.com/VeeamHub/martini-cli/core"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/VeeamHub/martini-cli/setup"

	"github.com/urfave/cli"
)

func main() {
	//cause altsrc sucks monkeyballs (no i don't want any yaml dependencies)
	hdir, _ := homedir.Dir()
	cfile := path.Join(hdir, ".martiniconfig")
	var cc core.ClientConfig
	tokenDefault := ""
	serverDefault := "https://localhost/api"
	usernameDefault := "admin"
	renewtokenDefault := ""
	renewlifetimeDefault := int64(0)
	renewserverskewDefault := int64(0)

	if _, err := os.Stat(cfile); err == nil {
		body, err := ioutil.ReadFile(cfile)
		if err == nil {
			err = json.Unmarshal(body, &cc)
			if err == nil {
				if cc.Token != "" {
					tokenDefault = cc.Token
				}
				if cc.Server != "" {
					serverDefault = cc.Server
				}
				if cc.Username != "" {
					usernameDefault = cc.Username
				}
				if cc.Lifetime != 0 {
					renewlifetimeDefault = cc.Lifetime
				}
				if cc.ServerSkew != 0 {
					renewserverskewDefault = cc.ServerSkew
				}
				if cc.Renew != "" {
					renewtokenDefault = cc.Renew
				}

			} else {
				log.Printf("Was not able to read config %v", err)
			}
		} else {
			log.Printf("Was not able to read config %v", err)
		}
	}

	app := cli.NewApp()
	app.Name = "martini"
	app.Usage = "For remote management and initial setup of martini vbo manager"
	app.Version = "1.0 (tp-1)"
	app.Description = "Martini CLI\n     #####\n      ###\n       #\n       |\n       |\n     -----\n"
	app.Commands = []cli.Command{
		{
			Name:    "setup",
			Aliases: []string{"s"},
			Usage:   "Setup wizard. This will create the database schema and setup file. Should only be used the server itself.",
			Action: func(c *cli.Context) error {
				return setup.SetupWizard()
			},
		},
		{
			Name:    "connect",
			Aliases: []string{"c"},
			Usage:   "Connect and save config file",
			Action: func(c *cli.Context) error {
				var rerr error

				po := core.NewPrintOptionsFromCLIContext(c)
				conn := core.NewConnection(&po, c.String("server"), c.String("token"), c.String("username"), c.String("password"), c.Bool("ignoreSelfSignedCertificate"), c.String("renewtoken"), c.Int64("renewlifetime"), c.Int64("renewserverskew"))
				rerr = conn.Auth(func() string {
					pw := ""
					fmt.Print("Type in the admin password: ")
					userPasswordByte, errp := terminal.ReadPassword(int(syscall.Stdin))
					for errp != nil || len(string(userPasswordByte)) < 3 {
						fmt.Println()
						fmt.Print("Password can not be empty and must be min 2 char:")
						userPasswordByte, errp = terminal.ReadPassword(int(syscall.Stdin))
					}
					pw = string(userPasswordByte)
					return pw
				}, true)
				if rerr == nil {
					fmt.Println("Authenticated, saving")
					var cc core.ClientConfig
					cc.Token = conn.GetToken()
					cc.Server = conn.GetServer()
					cc.Username = conn.GetUsername()
					cc.Lifetime = conn.GetLifetime()
					cc.Renew = conn.GetRenew()
					cc.ServerSkew = conn.GetServerSkew()
					jstext, err := json.Marshal(cc)
					if err == nil {
						err = ioutil.WriteFile(cfile, jstext, os.FileMode(0640))
					} else {
						log.Printf("Unable to save config %v", err)
					}
				} else {
					log.Printf("Got error %s authenticating", rerr)
				}
				return rerr
			},
		},
		*commands.GetTenantCommands(),
		*commands.GetInstanceCommands(),
		*commands.GetJobCommands(),
		*commands.GetConfigCommands(),
		*commands.GetLicenseCommands(),
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server, s",
			Value:  serverDefault,
			Usage:  "API server",
			EnvVar: "MARTINICLI_SERVER",
		},
		cli.StringFlag{
			Name:   "username, u",
			Value:  usernameDefault,
			Usage:  "API Username",
			EnvVar: "MARTINICLI_USERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			Value:  "",
			Usage:  "API Password",
			EnvVar: "MARTINICLI_PASSWORD",
		},
		cli.StringFlag{
			Name:   "token, t",
			Value:  tokenDefault,
			Usage:  "API Token",
			EnvVar: "MARTINICLI_TOKEN",
		},
		cli.BoolFlag{
			Name:   "ignoreSelfSignedCertificate, i",
			Usage:  "Ignore Self Signed Certificate",
			EnvVar: "MARTINICLI_IGNORESSC",
		},
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "Be verbose",
			EnvVar: "MARTINICLI_VERBOSE",
		},
		cli.BoolFlag{
			Name:   "json",
			Usage:  "Pass json instead of printing",
			EnvVar: "MARTINICLI_JSON",
		},
		cli.StringFlag{
			Name:  "renewtoken",
			Value: renewtokenDefault,
			Usage: "renew token (internal)",
		},
		cli.Int64Flag{
			Name:  "renewlifetime",
			Value: renewlifetimeDefault,
			Usage: "lifetime of renewal token",
		},
		cli.Int64Flag{
			Name:  "renewserverskew",
			Value: renewserverskewDefault,
			Usage: "skew in clocks between server and client",
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
