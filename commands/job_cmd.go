package commands

import (
	"fmt"

	"github.com/tdewin/martini-cli/core"
	"github.com/tdewin/martini-cli/job"
	"github.com/tdewin/martini-cli/tenant"
	"github.com/urfave/cli"
)

func GetJobCommands() *cli.Command {
	return &cli.Command{
		Name:    "job",
		Aliases: []string{"j"},
		Usage:   "job management",
		Subcommands: []cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list jobs",
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
										jobarray, e := job.List(conn, id)
										if e == nil {
											po.Println("### Tenant ID", id, "(", idtoname[id], ")")
											for _, job := range jobarray {
												po.Println(job.Id, job.Name, job.LastRun, job.LastStatus)
											}
											po.MarshalPrintJSON(jobarray)
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
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "start job",
				Action: func(c *cli.Context) error {
					err := ValidateOrArray([][]ValidString{
						[]ValidString{
							ValidString{c.String("id"), "id (for tenant)", "."},
							ValidString{c.String("name"), "name (for tenant)", "."},
						}, []ValidString{
							ValidString{c.String("jobid"), "jobid", "."},
							ValidString{c.String("jobname"), "name (for job)", "."},
						},
					})
					if err == nil {
						po := core.NewPrintOptionsFromCLIContext(c)

						conn := core.NewConnectionFromCLIContext(&po, c)

						err = conn.Auth(nil, false)
						if err == nil {
							id := c.String("id")
							if c.String("name") != "" {
								id, err = tenant.Resolve(conn, c.String("name"))
							}
							if err == nil {
								jobid := c.String("jobid")
								jobname := c.String("jobname")
								if jobname != "" {
									jobid, err = job.Resolve(conn, id, jobname)
								}
								if err == nil {
									err = job.Start(conn, id, jobid)
								}
							}
						}
					}
					return err
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "Id of tenant",
					}, cli.StringFlag{
						Name:  "name, n",
						Value: "",
						Usage: "Name of tenant",
					},
					cli.StringFlag{
						Name:  "jobid, j",
						Value: "",
						Usage: "Id of job",
					},
					cli.StringFlag{
						Name:  "jobname, o",
						Value: "",
						Usage: "Name of Job",
					},
				},
			},
		},
	}
}
