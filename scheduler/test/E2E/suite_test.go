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

package E2E

import (
	"context"
	"fmt"
	testv1alpha1 "github.com/keptn/lifecycle-toolkit/scheduler/test/E2E/fake/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	clientset1 "k8s.io/client-go/kubernetes"
	kscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/scheduler"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
	"time"
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

	go func() {
		defer GinkgoRecover()
		<-ctx.Done()
		fmt.Println("FINISHED")
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
	fmt.Fprintf(f, "%s \n", time.Now().UTC())

	for _, specReport := range report.SpecReports {
		if specReport.FullText() != "" {
			fmt.Fprintf(f, "%s, ", specReport.ContainerHierarchyTexts[1])
			fmt.Fprintf(f, "%s%s | %s\n", specReport.ContainerHierarchyTexts[2], specReport.LeafNodeText, specReport.State)

		}
	}
	f.Close()
})
