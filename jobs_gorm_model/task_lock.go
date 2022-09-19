package jobs_gorm_model

import (
	"context"
	"errors"
	"go.dtapp.net/gojobs"
)

type TaskLockOperation struct {
	task   Task           // 任务
	client *gojobs.Client // 实例
}

func (task Task) NewLock(c *gojobs.Client) (*TaskLockOperation, error) {
	if task.Id == 0 {
		return nil, errors.New("任务数据不正常")
	}
	return &TaskLockOperation{
		task:   task,
		client: c,
	}, nil
}

// LockId 上锁
func (tlo *TaskLockOperation) LockId(ctx context.Context) error {
	_, err := tlo.client.LockId(ctx, tlo.task)
	return err
}

// UnlockId 解锁
func (tlo *TaskLockOperation) UnlockId(ctx context.Context) error {
	return tlo.client.UnlockId(ctx, tlo.task)
}

// LockForeverId 永远上锁
func (tlo *TaskLockOperation) LockForeverId(ctx context.Context) error {
	_, err := tlo.client.LockForeverId(ctx, tlo.task)
	return err
}

// LockCustomId 上锁
func (tlo *TaskLockOperation) LockCustomId(ctx context.Context) error {
	_, err := tlo.client.LockCustomId(ctx, tlo.task)
	return err
}

// UnlockCustomId 解锁
func (tlo *TaskLockOperation) UnlockCustomId(ctx context.Context) error {
	return tlo.client.UnlockCustomId(ctx, tlo.task)
}

// LockForeverCustomId 永远上锁
func (tlo *TaskLockOperation) LockForeverCustomId(ctx context.Context) error {
	_, err := tlo.client.LockForeverCustomId(ctx, tlo.task)
	return err
}
