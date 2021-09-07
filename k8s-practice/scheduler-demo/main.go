package main

import (
	"math/rand"
	"os"
	"time"

	"scheduler-demo/pkg/plugins"

	"k8s.io/component-base/logs"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := app.NewSchedulerCommand(app.WithPlugin(plugins.Name, plugins.New))

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
