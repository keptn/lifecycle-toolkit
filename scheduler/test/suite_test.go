/*
Copyright 2022.

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

package test

import (
	"context"
	"fmt"
	testv1alpha1 "github.com/keptn/lifecycle-toolkit/scheduler/test/fake/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	clientset1 "k8s.io/client-go/kubernetes"
	kscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/events"
	"k8s.io/kube-scheduler/config/v1beta3"
	"k8s.io/kubernetes/pkg/scheduler"
	schedapi "k8s.io/kubernetes/pkg/scheduler/apis/config"
	"k8s.io/kubernetes/pkg/scheduler/apis/config/scheme"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"k8s.io/kubernetes/pkg/scheduler/profile"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"strings"
	"testing"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/v2 to learn more about Ginkgo.
var (
	schedulingConfigFile = filepath.Join("..", "manifests", "permit", "scheduler-config.yaml")
	cfg                  *rest.Config
	testEnv              *envtest.Environment
	ctx                  context.Context
	cancel               context.CancelFunc
	k8sManager           ctrl.Manager
	k8sClient            client.Client
	testCtx              *testContext
	fw                   framework.Framework
)

type testContext struct {
	ClientSet          clientset1.Interface
	KubeConfig         *rest.Config
	InformerFactory    informers.SharedInformerFactory
	DynInformerFactory dynamicinformer.DynamicSharedInformerFactory
	Scheduler          *scheduler.Scheduler
	Ctx                context.Context
	CancelFn           context.CancelFunc
}

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scheduler Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.TODO())
	By("bootstrapping test environment")

	t := true
	testEnv = &envtest.Environment{
		UseExistingCluster: &t,
	}

	apiServerArgs := testEnv.ControlPlane.GetAPIServer().Configure()
	//apiServerArgs.Append("disable-admission-plugins", "TaintNodesByCondition", "Priority")
	apiServerArgs.Append("runtime-config", "api/all=true")

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	//+kubebuilder:scaffold:scheme
	err = testv1alpha1.AddToScheme(kscheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	k8sClient, err = client.New(cfg, client.Options{Scheme: kscheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
	//
	//testCtx = &testContext{}
	//testCtx.Ctx, testCtx.CancelFn = context.WithCancel(context.Background())
	//
	//cs := clientset1.NewForConfigOrDie(cfg)
	//testCtx.ClientSet = cs
	//testCtx.KubeConfig = cfg
	//
	//conf, err := NewDefaultSchedulerComponentConfig()
	//Expect(err).NotTo(HaveOccurred())
	//
	//conf.Profiles[0].SchedulerName = "keptn-scheduler"
	//conf.Profiles[0].Plugins.Permit.Enabled = []schedapi.Plugin{schedapi.Plugin{Name: klcpermit.PluginName}}
	//conf.Profiles[0].PluginConfig = append(conf.Profiles[0].PluginConfig, schedapi.PluginConfig{
	//	Name: klcpermit.PluginName,
	//})
	//
	//testCtx = initTestSchedulerWithOptions(
	//	testCtx,
	//	scheduler.WithProfiles(conf.Profiles...),
	//	scheduler.WithFrameworkOutOfTreeRegistry(fwkruntime.Registry{klcpermit.PluginName: klcpermit.New}),
	//)
	//syncInformerFactory(testCtx)
	//
	//k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
	//	Scheme: kscheme.Scheme,
	//})
	//Expect(err).ToNot(HaveOccurred())
	//fmt.Println("plugins", conf.Profiles[0].Plugins)
	go func() {
		defer GinkgoRecover()
		<-ctx.Done()
		//testCtx.Scheduler.Run(ctx) // in case you want to test

		//fw, err = fwkruntime.NewFramework(fwkruntime.Registry{klcpermit.PluginName: klcpermit.New}, &conf.Profiles[0], fwkruntime.WithKubeConfig(cfg),
		//	fwkruntime.WithClientSet(cs), fwkruntime.WithInformerFactory(testCtx.InformerFactory))
		//gexec.KillAndWait(4 * time.Second)

		// Teardown the test environment once controller is finished.
		// Otherwise, from Kubernetes 1.21+, teardown timeouts waiting on
		// kube-apiserver to return
		err := testEnv.Stop()
		Expect(err).ToNot(HaveOccurred())
	}()

})

func ignoreAlreadyExists(err error) error {
	if apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.custom")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		path := strings.Split(specReport.FileName(), "/")
		testFile := path[len(path)-1]
		if specReport.ContainerHierarchyTexts != nil {
			testFile = specReport.ContainerHierarchyTexts[0]
		}
		fmt.Fprintf(f, "%s %s | %s\n", testFile, specReport.LeafNodeText, specReport.State)
	}
	f.Close()
})

func initTestSchedulerWithOptions(testCtx *testContext, opts ...scheduler.Option) *testContext {
	testCtx.InformerFactory = scheduler.NewInformerFactory(testCtx.ClientSet, 0)
	if cfg != nil {
		dynClient := dynamic.NewForConfigOrDie(cfg)
		testCtx.DynInformerFactory = dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynClient, 0, v1.NamespaceAll, nil)
	}

	var err error
	eventBroadcaster := events.NewBroadcaster(&events.EventSinkImpl{
		Interface: testCtx.ClientSet.EventsV1(),
	})

	Expect(err).ToNot(HaveOccurred())

	opts = append(opts, scheduler.WithKubeConfig(cfg))
	testCtx.Scheduler, err = scheduler.New(
		testCtx.ClientSet,
		testCtx.InformerFactory,
		testCtx.DynInformerFactory,
		profile.NewRecorderFactory(eventBroadcaster),
		testCtx.Ctx.Done(),
		opts...,
	)
	Expect(err).ToNot(HaveOccurred())
	stopCh := make(chan struct{})
	eventBroadcaster.StartRecordingToSink(stopCh)

	return testCtx
}

func NewDefaultSchedulerComponentConfig() (schedapi.KubeSchedulerConfiguration, error) {
	var versionedCfg v1beta3.KubeSchedulerConfiguration
	scheme.Scheme.Default(&versionedCfg)
	conf := schedapi.KubeSchedulerConfiguration{}
	if err := scheme.Scheme.Convert(&versionedCfg, &conf, nil); err != nil {
		return schedapi.KubeSchedulerConfiguration{}, err
	}
	return conf, nil
}

func syncInformerFactory(testCtx *testContext) {
	testCtx.InformerFactory.Start(testCtx.Ctx.Done())
	if testCtx.DynInformerFactory != nil {
		testCtx.DynInformerFactory.Start(testCtx.Ctx.Done())
	}
	testCtx.InformerFactory.WaitForCacheSync(testCtx.Ctx.Done())
	if testCtx.DynInformerFactory != nil {
		testCtx.DynInformerFactory.WaitForCacheSync(testCtx.Ctx.Done())
	}
}
