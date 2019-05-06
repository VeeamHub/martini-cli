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
							id := c.String("id")
							ids := []string{}
							if c.String("name") != "" {
								tenantid, rerr := tenant.Resolve(conn, c.String("name"))
								if rerr != nil {
									err = rerr
								} else {
									ids = append(ids, tenantid)
								}
							} else if id == "all" {
								po.Verbose("All tenants selected, requesting all tenants from server")
								tenants, rerr := tenant.List(conn)
								if rerr == nil {
									for _, t := range tenants {
										ids = append(ids, t.Id)
									}
								} else {
									err = rerr
								}
							} else {
								ids = append(ids, id)
							}
							if err == nil {
								for _, id := range ids {
									jobarray, e := job.List(conn, id)
									if e == nil {
										fmt.Println("### Tenant ID", id)
										for _, job := range jobarray {
											fmt.Println(job.Id, job.Name, job.LastRun, job.LastStatus)
										}
									} else {
										fmt.Println("### Error Tenant ID", id, e)
									}
								}
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
