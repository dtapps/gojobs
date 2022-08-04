package gojobs

import "testing"

func TestSpec(t *testing.T) {
	t.Log("每隔10秒执行一次：", GetSeconds(10).Spec())
	t.Log("每隔10秒执行一次：", GetSeconds(10).Frequency())

	t.Log("每隔60秒执行一次：", GetSeconds(60).Spec())
	t.Log("每隔60秒执行一次：", GetSeconds(60).Frequency())

	t.Log("每隔5分钟执行一次：", GetMinutes(5).Spec())
	t.Log("每隔5分钟执行一次：", GetMinutes(5).Frequency())

	t.Log("每隔10分钟执行一次：", GetMinutes(10).Spec())
	t.Log("每隔10分钟执行一次：", GetMinutes(10).Frequency())

	t.Log("每隔60分钟执行一次：", GetMinutes(60).Spec())
	t.Log("每隔60分钟执行一次：", GetMinutes(60).Frequency())

	t.Log("每天n点执行一次：", GetHour(10).Spec())
	t.Log("每天n点执行一次：", GetHour(10).Frequency())

	t.Log("每隔n小时执行一次：", GetHourInterval(1).Spec())
	t.Log("每隔n小时执行一次：", GetHourInterval(1).Frequency())
}
