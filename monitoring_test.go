package main

import (
	"github.com/disaster37/go-nagios"
	"github.com/disaster37/go-yarn-rest/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeState(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)

	// when no job
	monitoringData := nagiosPlugin.NewMonitoring()
	monitoringData, err := computeState(make([]client.ApplicationInfo, 0, 0), monitoringData)
	assert.NoError(t, err)
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 1, len(monitoringData.Messages()))

	// When one job
	monitoringData = nagiosPlugin.NewMonitoring()
	job1 := client.ApplicationInfo{
		Id:           "id",
		Name:         "name",
		Queue:        "queue",
		FinishedTime: 0,
	}
	listJobs := make([]client.ApplicationInfo, 0, 2)
	listJobs = append(listJobs, job1)
	monitoringData, err = computeState(listJobs, monitoringData)
	assert.NoError(t, err)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 2, len(monitoringData.Messages()))

	// When more than one job
	monitoringData = nagiosPlugin.NewMonitoring()
	job2 := client.ApplicationInfo{
		Id:           "id2",
		Name:         "name2",
		Queue:        "queue2",
		FinishedTime: 0,
	}
	listJobs = append(listJobs, job2)
	monitoringData, err = computeState(listJobs, monitoringData)
	assert.NoError(t, err)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 3, len(monitoringData.Messages()))

}
