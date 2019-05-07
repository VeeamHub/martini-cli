package commands

import (
	"fmt"

	"github.com/tdewin/martini-cli/licensemgmt"

	"github.com/tdewin/martini-cli/core"
	"github.com/tdewin/martini-cli/tenant"
	"github.com/urfave/cli"
)

func GetLicenseCommands() *cli.Command {
	return &cli.Command{
		Name:    "license",
		Aliases: []string{"l"},
		Usage:   "license management",
		Subcommands: []cli.Command{
			{
				Name:    "listusers",
				Aliases: []string{"lu"},
				Usage:   "list users",
				Action: func(c *cli.Context) error {
					err := ValidateOrArray([][]ValidString{
						[]ValidString{
							ValidString{c.String("id"), "id (for tenant)", "."},
							ValidString{c.String("name"), "name (for tenant)", "."},
						},
					})
					if err == nil {
						po := core.NewPrintOptionsFromCLIContext(c)
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
									for _, id := range ids {
										licensearray, e := licensemgmt.ListUsers(conn, id)
										if e == nil {

											for _, l := range licensearray {
												po.Printf("%-15s | %-20s | %-35s | %25s\n", idtoname[id], l.OrganizationName, l.Name, l.LastBackupDate)
											}
											po.MarshalPrintJSON(licensearray)

										} else {
											err = fmt.Errorf("### Error Tenant ID %s %v", id, e)
										}
									}
								}
							} else {
								err = rerr
							}
						}
					}
					return err
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Value: "all",
						Usage: "Id of tenant",
					},
					cli.StringFlag{
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
					err := ValidateOrArray([][]ValidString{
						[]ValidString{
							ValidString{c.String("id"), "id (for tenant)", "."},
							ValidString{c.String("name"), "name (for tenant)", "."},
						},
					})
					if err == nil {
						po := core.NewPrintOptionsFromCLIContext(c)
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
									for _, id := range ids {
										licensearray, e := licensemgmt.ListInfo(conn, id)
										if e == nil {
											for _, l := range licensearray {

												po.Printf("%-15s | %-25s | %6d | %6d\n", idtoname[id], l.OrgName, l.LicensedUsers, l.NewUsers)
											}
											po.MarshalPrintJSON(licensearray)
										} else {
											err = fmt.Errorf("### Error Tenant ID %s %v", id, e)
										}
									}
								}
							} else {
								err = rerr
							}
						}
					}
					return err
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Value: "all",
						Usage: "Id of tenant",
					},
					cli.StringFlag{
						Name:  "name, n",
						Value: "",
						Usage: "Name of tenant",
					},
				},
			},
		},
	}
}
