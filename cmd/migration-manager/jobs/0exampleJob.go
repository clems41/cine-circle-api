package jobs

import (
	"gorm.io/gorm"
)

type ExampleJob struct{}

func (job ExampleJob) JobName() string {
	return "#XXXXX - Example - Job"
}

func (job ExampleJob) Job(tx *gorm.DB) (err error) {
	return
}

func (job ExampleJob) IsJobDone(DB *gorm.DB) bool {
	return true
}

func (job ExampleJob) PreventJobRerun() bool {
	return true
}
