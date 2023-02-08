// /*
// Copyright 2020 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */
package main

//
//import (
//	"context"
//	"fmt"
//	"io/ioutil"
//	"net"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"path/filepath"
//	"testing"
//
//	"github.com/google/go-cmp/cmp"
//	"github.com/spf13/pflag"
//
//	"k8s.io/kubernetes/cmd/kube-scheduler/app"
//	"k8s.io/kubernetes/cmd/kube-scheduler/app/options"
//	"k8s.io/kubernetes/pkg/scheduler/apis/config"
//	"k8s.io/kubernetes/pkg/scheduler/apis/config/testing/defaults"
//)
//
//func TestSetup(t *testing.T) {
//	// temp dir
//	tmpDir, err := ioutil.TempDir("", "scheduler-options")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer os.RemoveAll(tmpDir)
//
//	// https server
//	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(200)
//		w.Write([]byte(`{"metadata": {"name": "test"}}`))
//	}))
//	defer server.Close()
//
//	configKubeconfig := filepath.Join(tmpDir, "config.kubeconfig")
//	if err := ioutil.WriteFile(configKubeconfig, []byte(fmt.Sprintf(`
//apiVersion: v1
//kind: Config
//clusters:
//- cluster:
//    insecure-skip-tls-verify: true
//    server: %s
//  name: default
//contexts:
//- context:
//    cluster: default
//    user: default
//  name: default
//current-context: default
//users:
//- name: default
//  user:
//    username: config
//`, server.URL)), os.FileMode(0600)); err != nil {
//		t.Fatal(err)
//	}
//
//
//	// Coscheduling plugin config
//	coschedulingConfigFile := filepath.Join(tmpDir, "coscheduling.yaml")
//	if err := ioutil.WriteFile(coschedulingConfigFile, []byte(fmt.Sprintf(`
//apiVersion: kubescheduler.config.k8s.io/v1beta3
//kind: KubeSchedulerConfiguration
//clientConnection:
//  kubeconfig: "%s"
//profiles:
//- plugins:
//    multiPoint:
//      enabled:
//      - name: Coscheduling
//    queueSort:
//      disabled:
//      - name: PrioritySort
//    filter:
//      disabled:
//      - name: "*"
//    score:
//      disabled:
//      - name: "*"
//    preScore:
//      disabled:
//      - name: "*"
//  pluginConfig:
//  - name: Coscheduling
//    args:
//      permitWaitingTimeSeconds: 10
//`, configKubeconfig)), os.FileMode(0600)); err != nil {
//		t.Fatal(err)
//	}
//
//
//	testcases := []struct {
//		name            string
//		flags           []string
//		registryOptions []app.Option
//		wantPlugins     map[string]*config.Plugins
//	}{
//		{
//			name: "default config",
//			flags: []string{
//				"--kubeconfig", configKubeconfig,
//			},
//			wantPlugins: map[string]*config.Plugins{
//				"default-scheduler": defaults.ExpandedPluginsV1beta3,
//			},
//		},
//		{
//			name:            "single profile config - Coscheduling",
//			flags:           []string{"--config", coschedulingConfigFile},
//			registryOptions: []app.Option{app.WithPlugin(coscheduling.Name, coscheduling.New)},
//			wantPlugins: map[string]*config.Plugins{
//				"default-scheduler": {
//					MultiPoint: defaults.ExpandedPluginsV1beta3.MultiPoint,
//					QueueSort:  config.PluginSet{Enabled: []config.Plugin{{Name: coscheduling.Name}}},
//					Bind:       defaults.ExpandedPluginsV1beta3.Bind,
//					PreFilter: config.PluginSet{
//						Enabled: append(defaults.ExpandedPluginsV1beta3.PreFilter.Enabled, config.Plugin{Name: coscheduling.Name}),
//					},
//					PostFilter: config.PluginSet{
//						Enabled: append(defaults.ExpandedPluginsV1beta3.PostFilter.Enabled, config.Plugin{Name: coscheduling.Name}),
//					},
//					Permit: config.PluginSet{
//						Enabled: append(defaults.ExpandedPluginsV1beta3.Permit.Enabled, config.Plugin{Name: coscheduling.Name}),
//					},
//					Reserve: config.PluginSet{
//						Enabled: append(defaults.ExpandedPluginsV1beta3.Reserve.Enabled, config.Plugin{Name: coscheduling.Name}),
//					},
//					PreBind: defaults.ExpandedPluginsV1beta3.PreBind,
//					PostBind: config.PluginSet{
//						Enabled: append(defaults.ExpandedPluginsV1beta3.PostBind.Enabled, config.Plugin{Name: coscheduling.Name}),
//					},
//				},
//			},
//		},
//		// TODO: add a multi profile test.
//		// Ref: test "plugin config with multiple profiles" in
//		// https://github.com/kubernetes/kubernetes/blob/master/cmd/kube-scheduler/app/server_test.go
//	}
//
//	makeListener := func(t *testing.T) net.Listener {
//		t.Helper()
//		l, err := net.Listen("tcp", ":0")
//		if err != nil {
//			t.Fatal(err)
//		}
//		return l
//	}
//
//	for _, tc := range testcases {
//		t.Run(tc.name, func(t *testing.T) {
//			fs := pflag.NewFlagSet("test", pflag.PanicOnError)
//			opts := options.NewOptions()
//
//			nfs := opts.Flags
//			for _, f := range nfs.FlagSets {
//				fs.AddFlagSet(f)
//			}
//			if err := fs.Parse(tc.flags); err != nil {
//				t.Fatal(err)
//			}
//
//			// use listeners instead of static ports so parallel test runs don't conflict
//			opts.SecureServing.Listener = makeListener(t)
//			defer opts.SecureServing.Listener.Close()
//
//			ctx, cancel := context.WithCancel(context.Background())
//			defer cancel()
//			_, sched, err := app.Setup(ctx, opts, tc.registryOptions...)
//			if err != nil {
//				t.Fatal(err)
//			}
//
//			gotPlugins := make(map[string]*config.Plugins)
//			for n, p := range sched.Profiles {
//				gotPlugins[n] = p.ListPlugins()
//			}
//
//			if diff := cmp.Diff(tc.wantPlugins, gotPlugins); diff != "" {
//				t.Errorf("unexpected plugins diff (-want, +got): %s", diff)
//			}
//		})
//	}
//}
