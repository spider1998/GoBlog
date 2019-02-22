package main

import (
	"Project/doit/app"
	"Project/doit/routing"
	"Project/doit/service"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	var err error
	//初始化NSQ日志收发配置
	err = service.Log.Boot()
	if err != nil {
		return
	}
	err = service.SLog.Boot()
	if err != nil {
		return
	}
	/*-----全局初始化-----*/
	err = app.Init()
	if err != nil {
		panic(err)
	}

	go service.App.CronRedis() //持久化存储redis数据

	go service.BackupApp.RunCronLoop() //同步定时备份

	var wg sync.WaitGroup
	wg.Add(3)
	/*-----监听信号并处理正常关闭服务-----*/
	go func() {
		defer wg.Done()
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		sig := <-ch
		signal.Ignore(syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		close(ch)
		app.Logger.Info().Interface("signal", sig).Msg("received signal.")
		err := app.ServerClose()
		if err != nil {
			panic(err)
		}
	}()

	/*-----启动路由服务-----*/
	go func() {
		defer wg.Done()
		err := routing.Run()
		if err != nil {
			panic(err)
		}
	}()

	// 等待线程关闭
	wg.Wait()
	app.Logger.Info().Msg("shutdown gracefully.")
}
