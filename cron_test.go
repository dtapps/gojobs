package gojobs

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"log"
	"testing"
)

func TestCron1(t *testing.T) {

	// 创建一个cron实例 精确到秒
	crontab := NewCron()

	log.Println(crontab)

	err := crontab.AddJobByFunc("1", "*/1 * * * * *", func() {
		log.Println("哈哈哈哈")
	})
	if err != nil {
		fmt.Printf("添加任务时出错：%s", err)
		return
	}

	err = crontab.AddJobByFunc("2", "*/2 * * * * *", func() {
		log.Println("啊啊啊啊")
	})
	if err != nil {
		fmt.Printf("添加任务时出错：%s", err)
		return
	}

	crontab.Start()
	select {}
}

func TestCron2(t *testing.T) {
	i := 0
	s := gocron.NewScheduler()
	s.Every(5).Seconds().Do(func() {
		i++
		log.Println("execute per 5 seconds", i)
	})
	<-s.Start()
}
