package main

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/altsrc"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/urfave/cli"
	"os"
)

var debug bool
var yarnUrl string
var yarnLogin string
var yarnPassword string

func main() {

	// Logger setting
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)

	// CLI settings
	app := cli.NewApp()
	app.Usage = "Check Yarn jobs from Ressource Manager API"
	app.Version = "develop"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from `FILE`",
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "yarn-url",
			Usage:       "The full Knox URL to access on Ressource Manager API",
			EnvVar:      "YARN_URL",
			Destination: &yarnUrl,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "yarn-login",
			Usage:       "The Knox login",
			EnvVar:      "YARN_LOGIN",
			Destination: &yarnLogin,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "yarn-password",
			Usage:       "The Knox password",
			EnvVar:      "YARN_PASSWORD",
			Destination: &yarnPassword,
		}),
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Display debug output",
			Destination: &debug,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "check-jobs",
			Usage: "Check the jobs state",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "finished-since",
					Usage: "How many time in hour the jobs are finished",
					Value: 24,
				},
				cli.StringFlag{
					Name:  "queue-name",
					Usage: "The queue name where you should to checks the jobs",
				},
				cli.StringFlag{
					Name:  "user-name",
					Usage: "The user name where you should to checks the jobs",
				},
				cli.BoolFlag{
					Name: "fix-bug-2.7",
					Usage: "Use it to fix bug on Yarn 2.7 when you filter failed job on queue name",
				},
			},
			Action: checkJobs,
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.String("config") != "" {
			before := altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("config"))
			return before(c)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Check the global parameter
func manageGlobalParameters() error {
	if debug == true {
		log.SetLevel(log.DebugLevel)
	}

	if yarnUrl == "" {
		return errors.New("You must set --yarn-url parameter")
	}

	return nil
}
