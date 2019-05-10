package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tdewin/martini-cli/core"
	"github.com/tdewin/martini-cli/instance"
	"github.com/tdewin/martini-cli/job"
	"github.com/tdewin/martini-cli/tenant"
	"github.com/urfave/cli"
)

type MartiniJobTenant struct {
	TenantId     string           `json:"tenantid"`
	TenantName   string           `json:"tenantname"`
	InstanceId   string           `json:"instanceid"`
	InstanceName string           `json:"instancename"`
	Jobs         []job.MartiniJob `json:"jobs"`
}

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

					var alljobs []MartiniJobTenant
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
										if rerr != nil {
											allerrors = append(allerrors, fmt.Sprintf("### Error Tenant ID listing id %s %v", id, rerr))
										} else {
											po.Println("### Tenant ID", id, "(", idtoname[id], ")")
											for _, instance := range instances {
												po.Println("### Instance ID", instance.Id, "(", instance.Name, ")")
												jobarray, e := job.List(conn, instance.Id)

												if e == nil {
													for _, job := range jobarray {
														po.Println(job.Id, job.Name, job.LastRun, job.LastStatus)
													}
													t := MartiniJobTenant{id, idtoname[id], instance.Id, instance.Name, jobarray}
													alljobs = append(alljobs, t)
												} else {
													allerrors = append(allerrors, fmt.Sprintf("### Error Tenant ID %s %v", id, e))
												}
											}
											po.Println("")
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
					return po.MarshalPrintJSONError(alljobs, err)
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
							ValidString{c.String("id"), "id (for instance)", "."},
						}, []ValidString{
							ValidString{c.String("jobid"), "jobid", "."},
							ValidString{c.String("jobname"), "name (for job)", "."},
						},
					})
					po := core.NewPrintOptionsFromCLIContext(c)
					rs := core.ReturnStatus{Status: "Not Started"}
					if err == nil {

						conn := core.NewConnectionFromCLIContext(&po, c)

						err = conn.Auth(nil, false)
						if err == nil {
							id := c.String("id")

							if err == nil {
								jobid := c.String("jobid")
								jobname := c.String("jobname")
								if jobname != "" {
									jobid, err = job.Resolve(conn, id, jobname)
								}
								if err == nil {
									err = job.Start(conn, id, jobid)
									if err == nil {
										rs.Id = id
										rs.SubId = jobid
										rs.Status = "Job Started"
									}
								}
							}
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
