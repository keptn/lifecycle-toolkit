/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/keptn-sandbox/lifecycle-controller/scheduler/pkg/klcpermit"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/component-base/cli"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"os"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	command := app.NewSchedulerCommand(
		app.WithPlugin(klcpermit.PluginName, klcpermit.New),
	)

	code := cli.Run(command)
	os.Exit(code)
}
