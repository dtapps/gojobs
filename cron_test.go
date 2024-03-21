package gojobs

import (
	"log"
	"testing"
)

func TestCron(t *testing.T) {
	c := NewCronWithSeconds(WithCronLog())
	t.Log("c.option.log", c.option.log)
	c.AddFunc("@every 1s", func() {
		//fmt.Println("every 1s")
	})
	c.AddFunc("@every 2s", func() {
		//fmt.Println("every 2s")
	})
	c.Start()
	log.Println(c.List())
	c.AddFunc("@every 10s", func() {
		c.ListShow()
	})

	select {}
}
