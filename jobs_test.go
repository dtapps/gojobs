package gojobs

import "testing"

func TestSpec(t *testing.T) {
	t.Log(GetSpecSeconds(10))
	t.Log(GetFrequencySeconds(10))

	t.Log(GetSpecMinutes(1))
	t.Log(GetFrequencyMinutes(1))
	t.Log(GetSpecMinutes(10))
	t.Log(GetFrequencyMinutes(10))
	t.Log(GetSpecMinutes(30))
	t.Log(GetFrequencyMinutes(30))

	t.Log(GetSpecHour(10))
	t.Log(GetFrequencyHour(10))
}
