package main

import (
	"github.com/disaster37/go-nagios"
	"github.com/disaster37/go-yarn-rest/client"
)

func computeState(jobs []client.ApplicationInfo, monitoringData *nagiosPlugin.Monitoring) (*nagiosPlugin.Monitoring, error) {

	monitoringData.AddPerfdata("nbJobFailed", len(jobs), "")
	if len(jobs) > 0 {
		monitoringData.SetStatus(nagiosPlugin.STATUS_CRITICAL)
		monitoringData.AddMessage("There are %d jobs in failed state", len(jobs))
	} else {
		monitoringData.AddMessage("All work fine!")
	}

	for _, job := range jobs {
		monitoringData.AddMessage("Job %s (%s) in queue %s failed at %s", job.Name, job.Id, job.Queue, job.FinishedDateTime())
	}

	return monitoringData, nil

}
