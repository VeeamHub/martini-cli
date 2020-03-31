package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/VeeamHub/martini-cli/instance"
	"github.com/VeeamHub/martini-cli/licensemgmt"

	"github.com/VeeamHub/martini-cli/core"
	"github.com/VeeamHub/martini-cli/tenant"
	"github.com/urfave/cli"
)

type MartiniLicenseUserInstance struct {
	TenantId     string                           `json:"tenantid"`
	TenantName   string                           `json:"tenantname"`
	InstanceId   string                           `json:"instanceid"`
	InstanceName string                           `json:"instancename"`
	LicenseUsers []licensemgmt.MartiniLicenseUser `json:"licenseusers"`
}

type MartiniLicenseInfoInstance struct {
	TenantId     string                           `json:"tenantid"`
	TenantName   string                           `json:"tenantname"`
	InstanceId   string                           `json:"instanceid"`
	InstanceName string                           `json:"instancename"`
	LicenseInfos []licensemgmt.MartiniLicenseInfo `json:"licenseinfos"`
}

func GetLicenseCommands() *cli.Command {
	return &cli.Command{
		Name:    "license",
		Aliases: []string{"l"},
		Usage:   "license management",
		Subcommands: []*cli.Command{
			{
				Name:    "listusers",
				Aliases: []string{"lu"},
				Usage:   "list users",
				Action: func(c *cli.Context) error {
					po := core.NewPrintOptionsFromCLIContext(c)
					allLicenseUserInstances := []MartiniLicenseUserInstance{}
					err := ValidateOrArray([][]ValidString{
						[]ValidString{
							ValidString{c.String("id"), "id (for tenant)", "."},
							ValidString{c.String("name"), "name (for tenant)", "."},
						},
					})
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
										instances, rerr := instance.List(conn, id)
										if rerr == nil {
											for _, instance := range instances {
												licensearray, e := licensemgmt.ListUsers(conn, instance.Id)
												if e == nil {

													for _, l := range licensearray {
														po.Printf("%-15s | %-6s | %-20s | %-35s | %25s\n", idtoname[id], instance.Id, l.OrganizationName, l.Name, l.LastBackupDate)

													}
													allLicenseUserInstances = append(allLicenseUserInstances, MartiniLicenseUserInstance{id, idtoname[id], instance.Id, instance.Name, licensearray})

												} else {
													allerrors = append(allerrors, fmt.Sprintf("### Error Instance id %s %v", instance.Id, e))
												}
											}
										} else {
											allerrors = append(allerrors, fmt.Sprintf("### Error Listinge tenant id %s %v", id, rerr))
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
					return po.MarshalPrintJSONError(allLicenseUserInstances, err)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "all",
						Usage: "Id of tenant",
					},
					&cli.StringFlag{
						Name:  "name, n",
						Value: "",
						Usage: "Name of tenant",
					},
				},
			},
			{
				Name:    "listinfo",
				Aliases: []string{"li"},
				Usage:   "list info",
				Action: func(c *cli.Context) error {
					allLicenseInfoInstances := []MartiniLicenseInfoInstance{}
					po := core.NewPrintOptionsFromCLIContext(c)
					err := ValidateOrArray([][]ValidString{
						[]ValidString{
							ValidString{c.String("id"), "id (for tenant)", "."},
							ValidString{c.String("name"), "name (for tenant)", "."},
						},
					})
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
										instances, rerr := instance.List(conn, id)
										if rerr == nil {
											for _, instance := range instances {
												licensearray, e := licensemgmt.ListInfo(conn, instance.Id)
												if e == nil {
													for _, l := range licensearray {

														po.Printf("%-15s | %-6s | %-25s | %6d | %6d\n", idtoname[id], instance.Id, l.OrgName, l.LicensedUsers, l.NewUsers)
														allLicenseInfoInstances = append(allLicenseInfoInstances, MartiniLicenseInfoInstance{id, idtoname[id], instance.Id, instance.Name, licensearray})
													}

												} else {

													allerrors = append(allerrors, fmt.Sprintf("### Error Instance id %s %v", instance.Id, e))
												}
											}
										} else {
											allerrors = append(allerrors, fmt.Sprintf("### Error Listinge tenant id %s %v", id, rerr))
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
					return po.MarshalPrintJSONError(allLicenseInfoInstances, err)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "all",
						Usage: "Id of tenant",
					},
					&cli.StringFlag{
						Name:  "name, n",
						Value: "",
						Usage: "Name of tenant",
					},
				},
			},
		},
	}
}
