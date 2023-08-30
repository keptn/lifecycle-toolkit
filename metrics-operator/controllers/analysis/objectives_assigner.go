package analysis

import metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"

type ITaskAssigner interface {
	AssignTasks(tasks []metricsapi.Objective, numWorkers int) [][]metricsapi.Objective
}

type TaskAssigner struct {
	numWorkers int
	tasks      []metricsapi.Objective
}

func (ta TaskAssigner) AssignTasks() map[int][]metricsapi.Objective {
	totalTasks := len(ta.tasks)
	taskMap := make(map[int][]metricsapi.Objective, ta.numWorkers)
	if ta.numWorkers > 0 {
		tasksPerWorker := totalTasks / ta.numWorkers
		start := 0
		for i := 0; i < ta.numWorkers && start < totalTasks; i++ {
			end := start + tasksPerWorker
			if i < totalTasks%ta.numWorkers {
				end++ // distribute the remainder tasks
			}
			taskMap[i+1] = ta.tasks[start:end]
			start = end
		}
	}
	return taskMap
}
