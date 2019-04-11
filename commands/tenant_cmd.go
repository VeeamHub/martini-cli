package commands

import (
	"fmt"
	"log"
	"syscall"

	"github.com/tdewin/martini-cli/core"
	"github.com/tdewin/martini-cli/tenant"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

//seperation without using global variables
func GetTenantCommands() *cli.Command {
	return &cli.Command{
		Name:    "tenant",
		Aliases: []string{"t"},
		Usage:   "tenant management",
		Subcommands: []cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "create tenant",
				Action: func(c *cli.Context) error {
					conn := core.NewConnectionFromCLIContext(c)
					err := conn.Auth()
					if err == nil {
						pw := c.String("password")
						if pw == "" {
							fmt.Print("Enter tenant server password: ")
							dbbytePassword, errp := terminal.ReadPassword(int(syscall.Stdin))
							for errp != nil || len(string(dbbytePassword)) < 3 {
								fmt.Println()
								fmt.Print("Password can not be empty (min 3 char):")
								dbbytePassword, errp = terminal.ReadPassword(int(syscall.Stdin))
							}
							pw = string(dbbytePassword)
						}

						t := tenant.MartiniTenant{c.String("tenant"), c.String("email"), c.String("fqdn"), c.String("username"), pw, "-1"}

						err = t.Create(conn)
						if err != nil {
							log.Println("Error ", err)
						}
					}
					return err
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "tenant, t",
						Value: "",
						Usage: "Tenant name",
					},
					cli.StringFlag{
						Name:  "email, e",
						Value: "",
						Usage: "Email",
					},
					cli.StringFlag{
						Name:  "fqdn, f",
						Value: "",
						Usage: "FQDN instance",
					},
					cli.StringFlag{
						Name:  "username, u",
						Value: "",
						Usage: "Username instance",
					},
					cli.StringFlag{
						Name:  "password, p",
						Value: "",
						Usage: "Password instance",
					},
				},
			},
			{
				Name:    "deploy",
				Aliases: []string{"d"},
				Usage:   "deploy tenant (will create a new installation instead of just adding it to martini",
				Action: func(c *cli.Context) error {
					conn := core.NewConnectionFromCLIContext(c)
					err := conn.Auth()
					if err == nil {
						pw := c.String("password")
						if pw == "" {
							fmt.Print("Enter tenant server password: ")
							dbbytePassword, errp := terminal.ReadPassword(int(syscall.Stdin))
							for errp != nil || len(string(dbbytePassword)) < 3 {
								fmt.Println()
								fmt.Print("Password can not be empty (min 3 char):")
								dbbytePassword, errp = terminal.ReadPassword(int(syscall.Stdin))
							}
							pw = string(dbbytePassword)
						}

						t := tenant.MartiniTenant{c.String("tenant"), c.String("email"), c.String("fqdn"), c.String("username"), pw, "-1"}

						err = t.Deploy(conn)
						if err != nil {
							log.Println("Error ", err)
						}
					}
					return err
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "tenant, t",
						Value: "",
						Usage: "Tenant name",
					},
					cli.StringFlag{
						Name:  "email, e",
						Value: "",
						Usage: "Email",
					},
					cli.StringFlag{
						Name:  "fqdn, f",
						Value: "",
						Usage: "FQDN instance",
					},
					cli.StringFlag{
						Name:  "username, u",
						Value: "",
						Usage: "Username instance",
					},
					cli.StringFlag{
						Name:  "password, p",
						Value: "",
						Usage: "Password instance",
					},
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all tenants",
				Action: func(c *cli.Context) error {
					conn := core.NewConnectionFromCLIContext(c)
					err := conn.Auth()
					if err == nil {
						tenants, err := tenant.List(conn)
						if err == nil {
							fmt.Print("#####################################################################")
							for _, t := range tenants {
								fmt.Printf("\n| %5s | %15s | %25s | %30s | %25s |", t.Id, t.Name, t.Email, t.Instancefqdn, t.Instanceusername)
							}
							fmt.Println("\n#####################################################################")
						} else {
							fmt.Println(err)
						}
					}
					return err
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"x"},
				Usage:   "delete a tenant",
				Action: func(c *cli.Context) error {
					conn := core.NewConnectionFromCLIContext(c)
					err := conn.Auth()
					if err == nil {
						err := tenant.Delete(conn, c.String("id"))
						if err != nil {
							fmt.Println(err)
						}
					}
					return err
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "Id of tenant",
					},
				},
			},
			//more commands indent here
		},
	}
}
