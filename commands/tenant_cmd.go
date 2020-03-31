package commands

import (
	"fmt"
	"syscall"

	"github.com/VeeamHub/martini-cli/core"
	"github.com/VeeamHub/martini-cli/instance"
	"github.com/VeeamHub/martini-cli/tenant"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

type MartiniTenantWithInstance struct {
	Tenant    tenant.MartiniTenant
	Instances []instance.MartiniInstance
}

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
					err := ValidateArray([]ValidString{
						ValidString{c.String("tenant"), "tenant", "."},
						ValidString{c.String("email"), "email", `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`},
						ValidString{c.String("port"), "port", `[0-9]*`},
					})
					po := core.NewPrintOptionsFromCLIContext(c)
					rs := MartiniTenantWithInstance{}
					if err == nil {

						conn := core.NewConnectionFromCLIContext(&po, c)
						err = conn.Auth(nil, false)
						if err == nil {
							t := tenant.MartiniTenant{c.String("tenant"), c.String("email"), "-1", "", "-1"}
							err = t.Create(conn)
							if t.Password != "" {
								po.Printf("Password tenant : %s", t.Password)
							}

							rs.Tenant = t
							if c.String("fqdn") != "" {
								if t.Id != "-1" && t.Id != "" {
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
									i := instance.MartiniInstance{Name: fmt.Sprintf("%s-%s", t.Name, c.String("fqdn")), TenantId: t.Id, Type: "Manual", Status: "-1", Location: c.String("location"), Hostname: c.String("fqdn"), Port: c.String("port"),
										Username: c.String("username"), Password: pw}
									err = i.Create(conn)
									if err == nil {
										i.Password = ""
										rs.Instances = append(rs.Instances, i)
									}
								} else {
									err = fmt.Errorf("Tenant creation did not yield tenant id")
								}
							}

						}
					}

					return po.MarshalPrintJSONError(rs, err)
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
						Name:  "port",
						Value: "4443",
						Usage: "FQDN port",
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
					cli.StringFlag{
						Name:  "location",
						Value: "unknown",
						Usage: "Username instance",
					},
				},
			},
			{
				Name:    "deploy",
				Aliases: []string{"d"},
				Usage:   "deploy tenant (will create a new installation instead of just adding it to martini)",
				Subcommands: []cli.Command{
					{
						Name:    "amazon",
						Aliases: []string{"a"},
						Usage:   "deploy an amazon EC2 container",
						Action: func(c *cli.Context) error {
							po := core.NewPrintOptionsFromCLIContext(c)
							rs := tenant.MartiniTenant{Id: "-1", Name: "Not updated"}
							err := ValidateArray([]ValidString{
								ValidString{c.String("tenant"), "tenant", "."},
								ValidString{c.String("email"), "email", `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`},
								ValidString{c.String("region"), "region", "."},
							})
							if err == nil {
								conn := core.NewConnectionFromCLIContext(&po, c)
								err = conn.Auth(nil, false)
								if err == nil {
									t := tenant.MartiniTenant{c.String("tenant"), c.String("email"), "-1", "", "-1"}
									err = t.Create(conn)
									if t.Password != "" {
										po.Printf("Password for tenant is %s", t.Password)
									}
									rs = t

									if err == nil {
										t := instance.NewAWSConfig(t.Id, c.String("region"))
										rid, rerr := t.Deploy(conn)
										_ = rid
										if rerr != nil {
											err = rerr
										}
									}
								}

							}
							return po.MarshalPrintJSONError(rs, err)
						},
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:  "region",
								Value: "",
								Usage: "AWS Region",
							},
						},
					},
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
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all tenants",
				Action: func(c *cli.Context) error {
					po := core.NewPrintOptionsFromCLIContext(c)
					conn := core.NewConnectionFromCLIContext(&po, c)
					err := conn.Auth(nil, false)
					if err == nil {
						tenants, err := tenant.List(conn)
						if err == nil {

							for i := 0; i < 10; i++ {
								po.Print("##########")
							}

							for _, t := range tenants {
								po.Printf("\n| %5s | %20s | %30s | %15s |", t.Id, t.Name, t.Email, t.Registered)
							}
							po.Print("\n")
							for i := 0; i < 10; i++ {
								po.Print("##########")
							}
							po.Print("\n")

							po.MarshalPrintJSON(tenants)
						}
					}
					return err
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"x"},
				Usage:   "delete a tenant (does not delete the instances)",
				Action: func(c *cli.Context) error {
					po := core.NewPrintOptionsFromCLIContext(c)
					rs := core.ReturnStatus{Status: "not deleted"}
					err := ValidateArray([]ValidString{
						ValidString{c.String("id"), "id (for tenant)", "."},
					})
					if err == nil {

						conn := core.NewConnectionFromCLIContext(&po, c)
						err = conn.Auth(nil, false)
						if err == nil {
							err = tenant.Delete(conn, c.String("id"))
							rs.Status = "Deleted from db tenant"
							rs.Id = c.String("id")
						}
					}
					return po.MarshalPrintJSONError(rs, err)
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
