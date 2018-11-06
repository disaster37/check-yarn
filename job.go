package main

import (
	"fmt"
	"github.com/disaster37/go-nagios"
	"github.com/disaster37/go-yarn-rest/client"
	"gopkg.in/urfave/cli.v1"
	"strconv"
	"time"
)

// Perform a node check
func checkJobs(c *cli.Context) error {

	monitoringData := nagiosPlugin.NewMonitoring()

	// Check global parameters
	err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), nagiosPlugin.STATUS_UNKNOWN)
	}

	// Check current parameters
	if c.Int("finished-since") == 0 {
		return cli.NewExitError("The parameter --finished-since can't be 0", nagiosPlugin.STATUS_UNKNOWN)
	}

	// Get Ambari connection
	yarnClient := client.New(yarnUrl, yarnLogin, yarnPassword)
	yarnClient.DisableVerifySSL()

	// Get failed jobs
	dateTime := time.Now().Add(time.Duration(c.Int("finished-since")) * -1 * time.Hour)
	filters := map[string]string{
		"finishedTimeBegin": strconv.FormatInt(dateTime.UnixNano()/1000000, 10),
		"states":            "FAILED",
	}
	if c.String("queue-name") != "" {
		filters["queue"] = c.String("queue-name")
	}
	if c.String("user-name") != "" {
		filters["user"] = c.String("user-name")
	}

	// Check node alertes
	jobs, err := yarnClient.Applications(filters)
	if err != nil {
		monitoringData.AddMessage("Somethink wrong when try to check jobs on Yarn: %v", err)
		monitoringData.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		monitoringData.ToSdtOut()
	}

	monitoringData, err = computeState(jobs, monitoringData)
	if err != nil {
		return err
	}

	monitoringData.ToSdtOut()
	return nil

}
