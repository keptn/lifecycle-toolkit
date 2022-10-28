package test

import (
	"fmt"
	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	keptncontroller "github.com/keptn/lifecycle-controller/operator/controllers/common"
	"github.com/keptn/lifecycle-controller/operator/controllers/keptnapp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//const (
//	MinPodStartupMeasurements = 30
//	TotalPodCount             = 100
//)
//
//var _ = Describe("Job E2E Test", func() {
//	It("[Feature:Performance] Schedule Density Job", func() {
//		// add 20 CRD with 100 print jobs as pre deployment
//		// create pod watch for jobs
//		// verify status succeeded for all
//		// here there is an example collecting metrics
//
//		err := waitDensityCRDReady(context, pg, TotalPodCount)
//		checkError(context, err)
//
//		By(fmt.Sprintf("Scheduling additional %d Pods to measure startup latencies", latencyPodsIterations*nodeCount))
//
//		createTimes := make(map[string]metav1.Time, 0)
//		nodeNames := make(map[string]string, 0)
//		scheduleTimes := make(map[string]metav1.Time, 0)
//		runTimes := make(map[string]metav1.Time, 0)
//		watchTimes := make(map[string]metav1.Time, 0)
//
//		var mutex sync.Mutex
//		checkPod := func(p *v1.Pod) {
//			mutex.Lock()
//			defer mutex.Unlock()
//			defer GinkgoRecover()
//
//			if p.Status.Phase == v1.PodRunning {
//				if _, found := watchTimes[p.Name]; !found {
//					watchTimes[p.Name] = metav1.Now()
//					createTimes[p.Name] = p.CreationTimestamp
//					nodeNames[p.Name] = p.Spec.NodeName
//					var startTime metav1.Time
//					for _, cs := range p.Status.ContainerStatuses {
//						if cs.State.Running != nil {
//							if startTime.Before(&cs.State.Running.StartedAt) {
//								startTime = cs.State.Running.StartedAt
//							}
//						}
//					}
//					if startTime != metav1.NewTime(time.Time{}) {
//						runTimes[p.Name] = startTime
//					} else {
//						fmt.Println("Pod  is reported to be running, but none of its containers is", p.Name)
//					}
//				}
//			}
//		}
//
//		additionalPodsPrefix := "density-latency-pod"
//		stopCh := make(chan struct{})
//
//		nsName := namespace
//		_, controller := cache.NewInformer(
//			&cache.ListWatch{
//				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
//					options.LabelSelector = labels.SelectorFromSet(labels.Set{"type": additionalPodsPrefix}).String()
//					obj, err := context.kubeclient.CoreV1().Pods(nsName).List(options)
//					return runtime.Object(obj), err
//				},
//				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
//					options.LabelSelector = labels.SelectorFromSet(labels.Set{"type": additionalPodsPrefix}).String()
//					return context.kubeclient.CoreV1().Pods(nsName).Watch(options)
//				},
//			},
//			&v1.Pod{},
//			0,
//			cache.ResourceEventHandlerFuncs{
//				AddFunc: func(obj interface{}) {
//					p, ok := obj.(*v1.Pod)
//					if !ok {
//						fmt.Println("Failed to cast observed object to *v1.Pod.")
//					}
//					Expect(ok).To(Equal(true))
//					go checkPod(p)
//				},
//				UpdateFunc: func(oldObj, newObj interface{}) {
//					p, ok := newObj.(*v1.Pod)
//					if !ok {
//						fmt.Println("Failed to cast observed object to *v1.Pod.")
//					}
//					Expect(ok).To(Equal(true))
//					go checkPod(p)
//				},
//			},
//		)
//
//		go controller.Run(stopCh)
//
//		for latencyPodsIteration := 0; latencyPodsIteration < latencyPodsIterations; latencyPodsIteration++ {
//			podIndexOffset := latencyPodsIteration * nodeCount
//			fmt.Println("Creating  latency pods in range ", nodeCount, podIndexOffset+1, podIndexOffset+nodeCount)
//
//			watchTimesLen := len(watchTimes)
//
//			var wg sync.WaitGroup
//			wg.Add(nodeCount)
//
//			cpuRequest := *resource.NewMilliQuantity(1, resource.DecimalSI)
//			memRequest := *resource.NewQuantity(1, resource.DecimalSI)
//
//			rcNameToNsMap := map[string]string{}
//			for i := 1; i <= nodeCount; i++ {
//				name := additionalPodsPrefix + "-" + strconv.Itoa(podIndexOffset+i)
//				nsName := context.namespace
//				rcNameToNsMap[name] = nsName
//				go createRunningPodFromRC(&wg, context, name, imageutils.GetPauseImageName(), additionalPodsPrefix, cpuRequest, memRequest)
//				time.Sleep(200 * time.Millisecond)
//			}
//			wg.Wait()
//
//			By("Waiting for all Pods begin observed by the watch...")
//			waitTimeout := 10 * time.Minute
//			for start := time.Now(); len(watchTimes) < watchTimesLen+nodeCount; time.Sleep(10 * time.Second) {
//				if time.Since(start) < waitTimeout {
//					fmt.Println("Timeout reached waiting for all Pods being observed by the watch.")
//				}
//			}
//
//			By("Removing additional replication controllers")
//			deleteRC := func(i int) {
//				defer GinkgoRecover()
//				name := additionalPodsPrefix + "-" + strconv.Itoa(podIndexOffset+i+1)
//				deleteReplicationController(context, name)
//			}
//			workqueue.ParallelizeUntil(con.TODO(), 25, nodeCount, deleteRC)
//		}
//		close(stopCh)
//
//		nsName = context.namespace
//		//time.Sleep(1 * time.Minute) // sleep to be added for large number of pods
//		selector := fields.Set{
//			"involvedObject.kind":      "Pod",
//			"involvedObject.namespace": nsName,
//			"source":                   "kube-batch",
//		}.AsSelector().String()
//		options := metav1.ListOptions{FieldSelector: selector}
//		schedEvents, _ := context.kubeclient.CoreV1().Events(nsName).List(options)
//		for k := range createTimes {
//			for _, event := range schedEvents.Items {
//				if event.InvolvedObject.Name == k {
//					scheduleTimes[k] = event.FirstTimestamp
//					break
//				}
//			}
//		}
//
//		scheduleLag := make([]PodLatencyData, 0)
//		startupLag := make([]PodLatencyData, 0)
//		watchLag := make([]PodLatencyData, 0)
//		schedToWatchLag := make([]PodLatencyData, 0)
//		e2eLag := make([]PodLatencyData, 0)
//
//		for name, create := range createTimes {
//			sched, ok := scheduleTimes[name]
//			if !ok {
//				fmt.Println("Failed to find schedule time for ", name)
//				missingMeasurements++
//			}
//			run, ok := runTimes[name]
//			if !ok {
//				fmt.Println("Failed to find run time for", name)
//				missingMeasurements++
//			}
//			watch, ok := watchTimes[name]
//			if !ok {
//				fmt.Println("Failed to find watch time for", name)
//				missingMeasurements++
//			}
//			node, ok := nodeNames[name]
//			if !ok {
//				fmt.Println("Failed to find node for", name)
//				missingMeasurements++
//			}
//			scheduleLag = append(scheduleLag, PodLatencyData{Name: name, Node: node, Latency: sched.Time.Sub(create.Time)})
//			startupLag = append(startupLag, PodLatencyData{Name: name, Node: node, Latency: run.Time.Sub(sched.Time)})
//			watchLag = append(watchLag, PodLatencyData{Name: name, Node: node, Latency: watch.Time.Sub(run.Time)})
//			schedToWatchLag = append(schedToWatchLag, PodLatencyData{Name: name, Node: node, Latency: watch.Time.Sub(sched.Time)})
//			e2eLag = append(e2eLag, PodLatencyData{Name: name, Node: node, Latency: watch.Time.Sub(create.Time)})
//		}
//
//		sort.Sort(LatencySlice(scheduleLag))
//		sort.Sort(LatencySlice(startupLag))
//		sort.Sort(LatencySlice(watchLag))
//		sort.Sort(LatencySlice(schedToWatchLag))
//		sort.Sort(LatencySlice(e2eLag))
//
//		PrintLatencies(scheduleLag, "worst create-to-schedule latencies")
//		PrintLatencies(startupLag, "worst schedule-to-run latencies")
//		PrintLatencies(watchLag, "worst run-to-watch latencies")
//		PrintLatencies(schedToWatchLag, "worst schedule-to-watch latencies")
//		PrintLatencies(e2eLag, "worst e2e latencies")
//
//		//// Capture latency metrics related to pod-startup.
//		podStartupLatency := &PodStartupLatency{
//			CreateToScheduleLatency: ExtractLatencyMetrics(scheduleLag),
//			ScheduleToRunLatency:    ExtractLatencyMetrics(startupLag),
//			RunToWatchLatency:       ExtractLatencyMetrics(watchLag),
//			ScheduleToWatchLatency:  ExtractLatencyMetrics(schedToWatchLag),
//			E2ELatency:              ExtractLatencyMetrics(e2eLag),
//		}
//
//		fmt.Println(podStartupLatency.PrintJSON())
//
//		dir, err := os.Getwd()
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println(dir)
//
//		filePath := path.Join(dir, "MetricsForE2ESuite_"+time.Now().Format(time.RFC3339)+".json")
//		if err := ioutil.WriteFile(filePath, []byte(podStartupLatency.PrintJSON()), 0644); err != nil {
//			fmt.Errorf("error writing to %q: %v", filePath, err)
//		}
//
//	})
//})

const LOAD = 100

var _ = Describe("[Feature:Performance] Load KeptnAppController", Ordered, func() {
	var (
		apps         []*klcv1alpha1.KeptnApp //Shelf is declared here
		appVersions  []*klcv1alpha1.KeptnAppVersion
		spanRecorder *sdktest.SpanRecorder
		tracer       *otelsdk.TracerProvider
	)
	BeforeAll(func() {
		//setup once
		By("Waiting for Manager")
		Eventually(func() bool {
			return k8sManager != nil
		}).Should(Equal(true))

		spanRecorder = sdktest.NewSpanRecorder()
		tracer = otelsdk.NewTracerProvider(otelsdk.WithSpanProcessor(spanRecorder))

		controllers := []keptncontroller.Controller{&keptnapp.KeptnAppReconciler{
			Client:   k8sManager.GetClient(),
			Scheme:   k8sManager.GetScheme(),
			Recorder: k8sManager.GetEventRecorderFor("test-app-controller"),
			Log:      GinkgoLogr,
			Tracer:   tracer.Tracer("test-app-tracer"),
		}}
		setupManager(controllers)
	})

	BeforeEach(func() {

		for i := 0; i < LOAD; i++ {
			instance := &klcv1alpha1.KeptnApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("app-%d", i),
					Namespace: "default",
				},
				Spec: klcv1alpha1.KeptnAppSpec{
					Version: "1.2.3",
					Workloads: []klcv1alpha1.KeptnWorkloadRef{
						{
							Name:    "app-wname",
							Version: "2.0",
						},
					},
				},
			}
			apps = append(apps, instance)
			Expect(k8sClient.Create(ctx, instance)).Should(Succeed())

		}
	})

	AfterEach(func() {
		for i, app := range apps {
			// Remember to clean up the cluster after each test
			deleteAppInCluster(app)
			deleteAppVersionInCluster(appVersions[i])
		}
	})
	It("should create the app version CR", func() {
		for _, app := range apps {
			appVersions = append(appVersions, assertResourceUpdated(app))
			//ResetSpanRecords()
		}
	})
})
