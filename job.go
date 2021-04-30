package main

import (
	"fmt"
	"strconv"
	"time"
	"strings"
	"regexp"

	"github.com/disaster37/go-nagios"
	"github.com/disaster37/go-yarn-rest/client"
	"github.com/urfave/cli"
)

// Perform a node check
func checkJobs(c *cli.Context) error {

	monitoringData := nagiosPlugin.NewMonitoring()

	// Check global parameters
	err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), nagiosPlugin.STATUS_UNKNOWN)
	}

	if len(c.StringSlice("state")) == 0 {
		return cli.NewExitError("The parameter --state can't be empty", nagiosPlugin.STATUS_UNKNOWN)
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
	}

	for _, param := range c.StringSlice("state") {
		keyValue := strings.Split(param, "=")
		if len(keyValue) != 2 {
			return cli.NewExitError("The parameter --state can content key=value", nagiosPlugin.STATUS_UNKNOWN)
		}
		filters[keyValue[0]] = keyValue[1]
	}

	if !c.Bool("fix-bug-2.7") && c.String("queue-name") != "" {
		filters["queue"] = c.String("queue-name")
	}
	if c.String("user-name") != "" {
		filters["user"] = c.String("user-name")
	}

	// Check node alertes
	jobsTmp, err := yarnClient.Applications(filters)
	if err != nil {
		monitoringData.AddMessage("Somethink wrong when try to check jobs on Yarn: %v", err)
		monitoringData.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		monitoringData.ToSdtOut()
	}

	// Check if filter must be run after grab data
	var jobs []client.ApplicationInfo
	if (c.Bool("fix-bug-2.7") && c.String("queue-name") != "") || c.String("job-name") != "" {
		jobs = make([]client.ApplicationInfo, 0)
		r, err := regexp.Compile(c.String("job-name"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), nagiosPlugin.STATUS_UNKNOWN)
		}
		for _, job := range jobsTmp {
			isKeep := true
			
			// Filter job name
			if c.String("job-name") != "" {
				if !r.MatchString(job.Name) {
					isKeep = false
				}
			}

			// Filter queue name
			if c.Bool("fix-bug-2.7") && c.String("queue-name") != "" && job.Queue != c.String("queue-name") {
				isKeep = false
			}

			// Compute
			if isKeep {
				jobs = append(jobs, job)
			}
		}
	} else {
		jobs = jobsTmp
	}

	monitoringData, err = computeState(jobs, monitoringData)
	if err != nil {
		return err
	}

	monitoringData.ToSdtOut()
	return nil

}
