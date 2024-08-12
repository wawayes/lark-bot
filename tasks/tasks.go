package tasks

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Task interface {
	Run() error
	GetSchedule() string
}

type TaskManager struct {
	cron   *cron.Cron
	logger *logrus.Logger
	tasks  []Task
}

func NewTaskManager(logger *logrus.Logger) *TaskManager {
	return &TaskManager{
		cron:   cron.New(),
		logger: logger,
		tasks:  []Task{},
	}
}

func (tm *TaskManager) AddTask(task Task) {
	tm.tasks = append(tm.tasks, task)
	_, err := tm.cron.AddFunc(task.GetSchedule(), func() {
		if err := task.Run(); err != nil {
			tm.logger.Errorf("Task execution failed: %v", err)
		}
	})
	if err != nil {
		tm.logger.Errorf("Failed to add task: %v", err)
	}
}

func (tm *TaskManager) Start() {
	tm.cron.Start()
}

func (tm *TaskManager) Stop() {
	tm.cron.Stop()
}
