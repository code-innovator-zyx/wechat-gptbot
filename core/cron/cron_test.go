package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2024/6/20 13:27
* @Package:
 */

func Test_Cron(t *testing.T) {
	// 创建一个新的 Cron 调度器
	c := cron.New()

	// 启动 Cron 调度器
	c.Start()

	// 在主协程中运行一段时间，以便我们可以添加和观察任务的执行
	go func() {

		// 动态添加一个任务
		addTask(c, "32 13 * * *", func() { fmt.Println("每分钟执行一次的新任务") })

	}()

	// 为了演示，保持主协程运行 10 分钟
	time.Sleep(5 * time.Second)

	// 停止 Cron 调度器 (非强制)
	c.Stop()
}

// 添加任务的辅助函数
func addTask(c *cron.Cron, spec string, cmd func()) {
	id, err := c.AddFunc(spec, cmd)
	if err != nil {
		fmt.Printf("添加任务失败: %v\n", err)
		return
	}
	fmt.Printf("任务添加成功，任务 ID: %d\n", id)
	fmt.Println(c.Entry(id).Job)
}
