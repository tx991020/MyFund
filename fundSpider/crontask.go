package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron/v3"
	"github.com/tx991020/utils"
)

type wrappedLogger struct {
	log *logs.BeeLogger
}

func NewWrappedLogger() *wrappedLogger {
	return &wrappedLogger{log: logs.GetBeeLogger()}
}

// Info logs routine messages about cron's operation.
func (wl *wrappedLogger) Info(msg string, keysAndValues ...interface{}) {
	wl.log.Info(msg)
}

// Error logs an error condition.
func (wl *wrappedLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	wl.log.Error(err.Error())
}

var (
	logg = NewWrappedLogger()
	c    = cron.New(cron.WithChain(
		cron.SkipIfStillRunning(logg),
		cron.Recover(logg),
	))
)

func SyncTasks() {

	c.AddFunc("3 0 * * *", func() {
		list := ReqFundList(fundList)
		if len(list) > 0 {
			data = strings.Split(utils.Time2Date(time.Now().AddDate(0, 0, -1)), "-")
			table = fmt.Sprintf("fund%s%s", data[1], data[2])

			list := ReqFundList(fundList)
			tmpPath := fmt.Sprintf(`./files/tmp%s.csv`, table)
			fundPath := fmt.Sprintf(`./files/%s.csv`, table)
			if len(list) > 0 {
				err := ReadCSVFund(fundDetail, tmpPath, fundPath, list)
				if err != nil {
					logs.Error("下载失败", err)
					return
				}
				err = CopyToPG(fundPath, table, createTable)
				if err != nil {
					logs.Error("导入失败", err)
					return
				}
			}
			logs.Info("定时任务爬取fund完成")
		}

	})
	//定时推送比特币，以太坊价格
	c.AddFunc("0 9 * * *", func() {
		_, err := exec.Command("/bin/bash", "-c", "/root/run.sh").Output()
		if err != nil {
			fmt.Println(err)
		}
		//定时推送基金

	})
	//c.AddFunc("10 0 * * *", func() {
	//	data = strings.Split(utils.Time2Date(time.Now().AddDate(0, 0, -1)), "-")
	//	pTable := fmt.Sprintf("rank%s%s", data[1], data[2])
	//	paihang := fmt.Sprintf(`./files/rank%s%s.csv`, data[1], data[2])
	//	PostDayRank(dayRank, paihang)
	//	err := CopyToPG(paihang, pTable,rankTable)
	//	if err != nil {
	//		logs.Error("导入失败", err)
	//		return
	//	}
	//	logs.Info("定时任务爬取rank完成")
	//})
	c.Start()
}
