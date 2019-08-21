package main

import (
	"fmt"
	"sync"
)

var workers = 9
var w_callNo = 0
var mu_workers sync.Mutex
var cond_workerBack = sync.NewCond(new(sync.Mutex))

func callWorker() {
	mu_workers.Lock()
	var callerNo = w_callNo
	w_callNo++
	fmt.Println(fmt.Sprintf("【%d】呼唤工人", callerNo))
	if workers <= 0 {
		// 等待工人回来消息
		cond_workerBack.L.Lock()
		mu_workers.Unlock()
		fmt.Println(fmt.Sprintf("【%d】等待中", callerNo))
		cond_workerBack.Wait()
		cond_workerBack.L.Unlock()
		fmt.Println(fmt.Sprintf("【%d】重新叫号", callerNo))
		callWorker()
	} else {
		workers--
		fmt.Println(">>>>>>>>> 给", callerNo, "派出了一个工人：现在有工人:", workers)
		mu_workers.Unlock()
	}
}

func workerBack() {
	mu_workers.Lock()
	cond_workerBack.L.Lock()
	defer mu_workers.Unlock()
	workers++
	// fmt.Println("<<<<<<<<< 一个工人回来了:现在有工人:", workers)
	cond_workerBack.Signal()
	cond_workerBack.L.Unlock()
}
