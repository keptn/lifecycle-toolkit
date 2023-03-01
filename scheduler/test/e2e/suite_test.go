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

package e2e

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	testv1alpha3 "github.com/keptn/lifecycle-toolkit/scheduler/test/e2e/fake/v1alpha3"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// nolint:gci
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/v2 to learn more about Ginkgo.
var (
	cfg       *rest.Config
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
	k8sClient client.Client
	wg        *sync.WaitGroup
)

func TestE2EScheduler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scheduler Suite")
}

func logErrorIfPresent(err error) {
	if err != nil {
		GinkgoLogr.Error(err, "Something went wrong while cleaning up the test environment")
	}
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
	// apiServerArgs.Append("disable-admission-plugins", "TaintNodesByCondition", "Priority")
	apiServerArgs.Append("runtime-config", "api/all=true")

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	// +kubebuilder:scaffold:scheme
	err = testv1alpha3.AddToScheme(kscheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	k8sClient, err = client.New(cfg, client.Options{Scheme: kscheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
	wg = &sync.WaitGroup{}

	go func() {
		defer GinkgoRecover()
		time.Sleep(30 * time.Second) // wait for test to start
		wg.Wait()
		fmt.Println("SUITE FINISHED")
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
	f, err := os.Create("report.E2E-scheduler")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	fmt.Fprintf(f, "%s \n", time.Now().UTC())

	for _, specReport := range report.SpecReports {
		if specReport.FullText() != "" {
			fmt.Fprintf(f, "%s, ", specReport.ContainerHierarchyTexts[1])
			fmt.Fprintf(f, "%s%s ", specReport.ContainerHierarchyTexts[2], specReport.LeafNodeText)
			switch specReport.State {
			case types.SpecStatePassed:
				fmt.Fprintf(f, "%s\n", "✓")
			case types.SpecStateFailed:
				fmt.Fprintf(f, "%s\n", "✕")
			default:
				fmt.Fprintf(f, "%s\n", specReport.State)
			}
		}
	}
	f.Close()

})
