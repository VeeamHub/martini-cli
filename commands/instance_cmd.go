package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tdewin/martini-cli/instance"

	"github.com/tdewin/martini-cli/core"

	"github.com/tdewin/martini-cli/tenant"
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
		Subcommands: []cli.Command{
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
					cli.StringFlag{
						Name:  "id, i",
						Value: "all",
						Usage: "Id of tenant",
					},
					cli.BoolFlag{
						Name:  "showorphans",
						Usage: "by default orphans are not shown",
					},
					cli.StringFlag{
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
					cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "Id of tenant",
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
					cli.StringFlag{
						Name:  "tenantid",
						Value: "",
						Usage: "Id of tenant",
					},
					cli.StringFlag{
						Name:  "tenantname, n",
						Value: "",
						Usage: "Name of the tenant",
					},
					cli.StringFlag{
						Name:  "id, i",
						Value: "first",
						Usage: "Id of tenant instance; default is first which will just select the first tenant if you don't specify an instance id but you specify tenant",
					},
					cli.StringFlag{
						Name:  "clientip, c",
						Value: "",
						Usage: "IP of your local break-out towards the server. If empty, the server will try to autodetect",
					},
				},
			},
			//cmd ident here
		},
	}
}
