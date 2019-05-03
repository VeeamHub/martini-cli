package commands

import (
	"github.com/tdewin/martini-cli/config"
	"github.com/tdewin/martini-cli/core"
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
								conn := core.NewConnectionFromCLIContext(c)
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
								conn := core.NewConnectionFromCLIContext(c)
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
				},
			},
		},
	}
}
