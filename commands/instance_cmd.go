package commands

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/VeeamHub/martini-cli/instance"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/VeeamHub/martini-cli/core"

	"github.com/VeeamHub/martini-cli/tenant"
	"github.com/urfave/cli"
)

type MartiniInstanceTenant struct {
	TenantId   string `json:"tenantid"`
	TenantName string `json:"tenantname"`

	Instances []instance.MartiniInstance `json:"instances"`
}

func GetInstanceCommands() *cli.Command {
	return &cli.Command{
		Name:    "instance",
		Aliases: []string{"i"},
		Usage:   "instance management",
		Subcommands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list",
				Action: func(c *cli.Context) error {
					po := core.NewPrintOptionsFromCLIContext(c)
					err := ValidateArray([]ValidString{
						ValidString{c.String("id"), "tenant id", "[0-9]*"},
						ValidString{c.String("name"), "name (for tenant)", "[a-zA-Z0-9-_]*"},
					})
					allinstances := []MartiniInstanceTenant{}
					if err == nil {
						conn := core.NewConnectionFromCLIContext(&po, c)
						err = conn.Auth(nil, false)
						if err == nil {
							idtoname, nametoid, rerr := tenant.Mappings(conn)
							if rerr == nil {
								id := c.String("id")
								ids := []string{}
								if c.String("name") != "" {
									tenantid, ok := nametoid[c.String("name")]
									if ok {
										ids = append(ids, tenantid)
									} else {
										err = fmt.Errorf("Could not find tenant %s", c.String("name"))
									}
								} else if id == "all" {
									po.Verbose("All tenants selected, requesting all tenants from server")

									for id := range idtoname {
										ids = append(ids, id)
									}
								} else {
									ids = append(ids, id)
								}

								if err == nil {

									var allerrors []string

									for _, id := range ids {
										instancearray, e := instance.List(conn, id)
										if e == nil {
											po.Println("### Tenant ID", id, "(", idtoname[id], ")")

											for _, instance := range instancearray {
												po.Println(instance.Id, instance.Name, instance.Hostname, instance.Type, instance.Location)
											}
											t := MartiniInstanceTenant{id, idtoname[id], instancearray}
											allinstances = append(allinstances, t)
										} else {
											allerrors = append(allerrors, fmt.Sprintf("### Error Tenant ID %s %v", id, e))
										}
									}

									if c.Bool("showorphans") {
										instancearray, e := instance.ListOrphans(conn)
										if e == nil {
											po.Println("####### Orphans")

											for _, instance := range instancearray {
												po.Println(instance.Id, instance.Name, instance.Hostname, instance.Type, instance.Location)
											}
											t := MartiniInstanceTenant{"<orphan>", "<orphan>", instancearray}
											allinstances = append(allinstances, t)
										} else {
											allerrors = append(allerrors, fmt.Sprintf("### Error Orphans %v", e))
										}
									}
									if len(allerrors) > 0 {
										err = errors.New(strings.Join(allerrors, "\n"))
									}

								}
							} else {
								err = rerr
							}

						}
					}
					return po.MarshalPrintJSONError(allinstances, err)

				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "all",
						Usage: "Id of tenant",
					},
					&cli.BoolFlag{
						Name:  "showorphans",
						Usage: "by default orphans are not shown",
					},
					&cli.StringFlag{
						Name:  "name, n",
						Value: "",
						Usage: "Name of tenant",
					},
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"x"},
				Usage:   "delete an instance (will clean it up)",
				Action: func(c *cli.Context) error {
					po := core.NewPrintOptionsFromCLIContext(c)
					rs := core.ReturnStatus{Status: "not deleted"}

					err := ValidateArray([]ValidString{
						ValidString{c.String("id"), "id (for instance)", "."},
					})
					if err == nil {

						conn := core.NewConnectionFromCLIContext(&po, c)
						err = conn.Auth(nil, false)
						if err == nil {
							err = instance.Delete(conn, c.String("id"))
							rs.Status = "Set for deletion in db"
							rs.Id = c.String("id")
						}
					}
					return po.MarshalPrintJSONError(rs, err)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "Id of tenant",
					},
				},
			},
			{
				Name:    "assign",
				Aliases: []string{""},
				Usage:   "(re)assign a tenant to an id",
				Action: func(c *cli.Context) error {
					po := core.NewPrintOptionsFromCLIContext(c)
					rs := core.ReturnStatus{Status: "not assigned"}

					err := ValidateOrArray([][]ValidString{
						[]ValidString{
							ValidString{c.String("newtenantid"), "newtenantid", "."},
							ValidString{c.String("newtenantname"), "newtenantname", "."},
						},
						[]ValidString{
							ValidString{c.String("id"), "id (for instance)", "."},
						},
					})
					if err == nil {

						conn := core.NewConnectionFromCLIContext(&po, c)
						err = conn.Auth(nil, false)
						if err == nil {
							tenantid := c.String("newtenantid")
							if tenantid == "" {
								tenantid, err = tenant.Resolve(conn, c.String("newtenantname"))
							}

							if err == nil {
								if tenantid != "-1" && tenantid != "" {
									err = instance.Assign(conn, c.String("id"), tenantid)
								}
							}
						}
					}
					return po.MarshalPrintJSONError(rs, err)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "newtenantid",
						Value: "",
						Usage: "Id of tenant",
					},
					&cli.StringFlag{
						Name:  "newtenantname, n",
						Value: "",
						Usage: "Name of the tenant",
					},
					&cli.StringFlag{
						Name:  "id, i",
						Value: "first",
						Usage: "Id of tenant instance; default is first which will just select the first tenant if you don't specify an instance id but you specify tenant",
					},
				},
			},
			{
				Name:    "broker",
				Aliases: []string{"b"},
				Usage:   "broker an rdp connection via the martini server to a tenant",
				Action: func(c *cli.Context) error {
					po := core.NewPrintOptionsFromCLIContext(c)
					var err error
					var bep instance.MartiniBrokerEndpoint

					iid := c.String("id")

					if iid == "first" {
						err = ValidateOrArray([][]ValidString{
							[]ValidString{
								ValidString{c.String("tenantid"), "tenantid", "."},
								ValidString{c.String("tenantname"), "tenantname", "."},
							},
						})
					} else {
						icheck := (ValidString{iid, "id (for instance)", "."})
						err = icheck.Validate()
					}

					if err == nil {
						conn := core.NewConnectionFromCLIContext(&po, c)
						err = conn.Auth(nil, false)
						if err == nil {
							if iid == "first" {
								tenantid := c.String("tenantid")
								if tenantid == "" {
									tenantid, err = tenant.Resolve(conn, c.String("tenantname"))
								}
								if err == nil {
									instances, rerr := instance.List(conn, tenantid)
									if rerr == nil {
										if len(instances) > 0 {
											iid = instances[0].Id
										} else {
											err = fmt.Errorf("Tenant has no connected instances")
										}
									} else {
										err = rerr
									}
								}
							}
							if err == nil {
								if iid != "" && iid != "first" {

									bep, err = instance.Broker(conn, iid, c.String("clientip"))
									if err == nil {
										po.Printf("Opened endpoint on %s (expecting ip %s)\n", bep.Port, bep.ExpectedClient)
									}
								} else {
									err = fmt.Errorf("Was not able to resolve the tenant name or find the first instance. Try to use the instance id instead")
								}
							}
						}
					}
					return po.MarshalPrintJSONError(bep, err)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tenantid",
						Value: "",
						Usage: "Id of tenant",
					},
					&cli.StringFlag{
						Name:  "tenantname, n",
						Value: "",
						Usage: "Name of the tenant",
					},
					&cli.StringFlag{
						Name:  "id, i",
						Value: "first",
						Usage: "Id of tenant instance; default is first which will just select the first tenant if you don't specify an instance id but you specify tenant",
					},
					&cli.StringFlag{
						Name:  "clientip, c",
						Value: "",
						Usage: "IP of your local break-out towards the server. If empty, the server will try to autodetect",
					},
				},
			},
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "create instance on existing tenant",
				Action: func(c *cli.Context) error {

					err := ValidateOrArray([][]ValidString{
						[]ValidString{
							ValidString{c.String("tenantid"), "tenantid", "."},
							ValidString{c.String("tenantname"), "tenantname", "."},
						}, []ValidString{
							ValidString{c.String("port"), "port", `[0-9]*`},
						},
					})

					po := core.NewPrintOptionsFromCLIContext(c)
					rs := core.ReturnStatus{Status: "Nothing done (check error)", Id: "-1", SubId: "-1"}
					if err == nil {

						conn := core.NewConnectionFromCLIContext(&po, c)
						err = conn.Auth(nil, false)
						if err == nil {
							tenantid := c.String("tenantid")
							tenantname := ""

							if tenantid == "" {
								tenantid, err = tenant.Resolve(conn, c.String("tenantname"))
								tenantname = c.String("tenantname")
							} else {
								tenantname, err = tenant.ReverseResolve(conn, tenantid)
							}
							if err == nil {
								if c.String("fqdn") != "" {
									if tenantid != "-1" && tenantid != "" {
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
										i := instance.MartiniInstance{Name: fmt.Sprintf("%s-%s", tenantname, c.String("fqdn")), TenantId: tenantid, Type: "Manual", Status: "-1", Location: c.String("location"), Hostname: c.String("fqdn"), Port: c.String("port"),
											Username: c.String("username"), Password: pw}
										err = i.Create(conn)
										if err == nil {
											rs.Status = "instance added to tenant"
											rs.Id = i.Id
										}
									} else {
										err = fmt.Errorf("Tenant creation did not yield tenant id")
									}
								}
							}

						}
					}

					return po.MarshalPrintJSONError(rs, err)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tenantid",
						Value: "",
						Usage: "Id of tenant",
					},
					&cli.StringFlag{
						Name:  "tenantname, n",
						Value: "",
						Usage: "Name of the tenant",
					},
					&cli.StringFlag{
						Name:  "fqdn, f",
						Value: "",
						Usage: "FQDN instance",
					},
					&cli.StringFlag{
						Name:  "port",
						Value: "4443",
						Usage: "FQDN port",
					},
					&cli.StringFlag{
						Name:  "username, u",
						Value: "",
						Usage: "Username instance",
					},
					&cli.StringFlag{
						Name:  "password, p",
						Value: "",
						Usage: "Password instance",
					},
					&cli.StringFlag{
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
				Subcommands: []*cli.Command{
					{
						Name:    "amazon",
						Aliases: []string{"a"},
						Usage:   "deploy an amazon EC2 container",
						Action: func(c *cli.Context) error {
							po := core.NewPrintOptionsFromCLIContext(c)
							rs := core.ReturnStatus{Status: "Nothing done (check error)", Id: "-1", SubId: "-1"}

							err := ValidateOrArray([][]ValidString{
								[]ValidString{
									ValidString{c.String("tenantid"), "tenantid", "."},
									ValidString{c.String("tenantname"), "tenantname", "."},
								}, []ValidString{
									ValidString{c.String("region"), "region", "."},
								},
							})

							if err == nil {
								conn := core.NewConnectionFromCLIContext(&po, c)
								err = conn.Auth(nil, false)
								if err == nil {
									tenantid := c.String("tenantid")

									if tenantid == "" {
										tenantid, err = tenant.Resolve(conn, c.String("tenantname"))
									}

									if err == nil {
										t := instance.NewAWSConfig(tenantid, c.String("region"))
										rid, rerr := t.Deploy(conn)
										if rerr == nil {
											rs.Id = rid
											rs.Status = "Deployement is started"
										}
									}
								}

							}
							return po.MarshalPrintJSONError(rs, err)
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "region",
								Value: "",
								Usage: "AWS Region",
							},
						},
					},
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tenantid",
						Value: "",
						Usage: "Id of tenant",
					},
					&cli.StringFlag{
						Name:  "tenantname, n",
						Value: "",
						Usage: "Name of the tenant",
					},
				},
			},

			//cmd ident here
		},
	}
}
