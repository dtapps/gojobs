package gojobs

import (
	"testing"
)

func TestSpec(t *testing.T) {
	t.Log("每隔n秒执行一次：", GetSpecSeconds(10))
	t.Log("每隔n秒执行一次：", GetFrequencySeconds(10))

	t.Log("每隔n分钟执行一次：", GetSpecMinutes(10))
	t.Log("每隔n分钟执行一次：", GetFrequencyMinutes(10))

	t.Log("每天n点执行一次：", GetSpecHour(10))
	t.Log("每天n点执行一次：", GetFrequencyHour(10))
}
