package test

import (
	. "github.com/onsi/ginkgo/v2"
	"time"
)

type Metric struct {
	creationTime            []time.Time
	succededAppVersionCount int
}

const LOAD = 100

var _ = Describe("[Feature:Performance] Load KeptnScheduler", Ordered, func() {

})
