package gojobs

import "testing"

func TestSpec(t *testing.T) {
	t.Log("每隔10秒执行一次：", GetSpecSeconds(10))
	t.Log("每隔10秒执行一次：", GetFrequencySeconds(10))

	t.Log("每隔60秒执行一次：", GetSpecSeconds(60))
	t.Log("每隔60秒执行一次：", GetFrequencySeconds(60))

	t.Log("每隔10分钟执行一次：", GetSpecMinutes(10))
	t.Log("每隔10分钟执行一次：", GetFrequencyMinutes(10))

	t.Log("每隔60分钟执行一次：", GetSpecMinutes(60))
	t.Log("每隔60分钟执行一次：", GetFrequencyMinutes(60))

	t.Log("每天n点执行一次：", GetSpecHour(10))
	t.Log("每天n点执行一次：", GetFrequencyHour(10))
}
