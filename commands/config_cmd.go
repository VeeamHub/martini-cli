package commands

import (
	"github.com/VeeamHub/martini-cli/config"
	"github.com/VeeamHub/martini-cli/core"
	"github.com/urfave/cli"
)

func GetConfigCommands() *cli.Command {
	return &cli.Command{
		Name:    "config",
		Aliases: []string{},
		Usage:   "config management",
		Subcommands: []cli.Command{
			{
				Name:    "broker",
				Aliases: []string{"b"},
				Usage:   "broker config",
				Subcommands: []cli.Command{
					{
						Name:  "add",
						Usage: "add port",
						Action: func(c *cli.Context) error {
							err := ValidateArray([]ValidString{
								ValidString{c.String("port"), "port", "[0-9]+"},
							})
							if err == nil {
								po := core.NewPrintOptionsFromCLIContext(c)
								conn := core.NewConnectionFromCLIContext(&po, c)
								err = conn.Auth(nil, false)
								if err == nil {
									err = config.BrokerAddPort(conn, c.String("port"))
								}
							}
							return err
						},
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:  "port, p",
								Value: "",
								Usage: "port",
							},
						},
					},
					{
						Name:  "delete",
						Usage: "delete port",
						Action: func(c *cli.Context) error {
							err := ValidateArray([]ValidString{
								ValidString{c.String("port"), "port", "[0-9]+"},
							})
							if err == nil {
								po := core.NewPrintOptionsFromCLIContext(c)
								conn := core.NewConnectionFromCLIContext(&po, c)

								err = conn.Auth(nil, false)
								if err == nil {
									err = config.BrokerDeletePort(conn, c.String("port"))
								}
							}
							return err
						},
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:  "port, p",
								Value: "",
								Usage: "port",
							},
						},
					},
					{
						Name:  "list",
						Usage: "list ports",
						Action: func(c *cli.Context) error {

							po := core.NewPrintOptionsFromCLIContext(c)
							conn := core.NewConnectionFromCLIContext(&po, c)

							err := conn.Auth(nil, false)
							if err == nil {
								portlist, rerr := config.BrokerList(conn)
								if rerr == nil {
									for _, p := range portlist.PortList {
										po.Println(p.Port)
									}
								} else {
									err = rerr
								}
							}
							return err

						},
					},
				},
			},
		},
	}
}
