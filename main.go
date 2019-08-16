package main

import (
	"./util"
	"fmt"
	"strings"
	"time"
)

var RetryCount = 0

func main() {
	fmt.Println(util.HostName(),"检测器启动")
	//定时检测路由器是否在线
	go OnlineTimerExec()
	go LanCheckTimerExec()
	select {}
}

func OnlineTimerExec() {
	d := time.Duration(time.Minute * 3)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		select {

		case <-t.C:
			if util.IsOnline("baidu.com") == false {
				Retry()
			} else {
				RetryCount = 0
			}
			fmt.Println(time.Now())
		}
	}
}

func LanCheckTimerExec() {
	d := time.Duration(time.Minute * 3)
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			list := util.LanCompare()
			if len(list) > 0 {
				SendNotice(list)
			}
			fmt.Println(time.Now())

		}
	}
}
func Retry() {
	if (RetryCount > 5) {
		SendNotice(nil)
	} else {
		util.ReDial()
		RetryCount++
		fmt.Println("掉线")
	}
}
func SendNotice(list []string) {
	fmt.Println("掉线", list)
	util.SendMail(util.HostName(),strings.Join(list, ","))
}
