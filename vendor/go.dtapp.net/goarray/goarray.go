package goarray

func Grouping() {

}

// TurnString []string 转 string
func TurnString(ss []string) (s string) {
	sl := len(ss)
	for k, v := range ss {
		if k+1 == sl {
			s = s + v
		} else {
			s = s + v + ","
		}
	}
	return s
}
