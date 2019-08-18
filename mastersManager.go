package main

import (
	"sync"
)

var masters = 3
var m_callNo = 0
var mu_masters sync.Mutex
var cond_masterBack = sync.NewCond(new(sync.Mutex))

func callMaster() {
	mu_masters.Lock()
	// tCallNo := m_callNo
	m_callNo++
	// fmt.Println(tCallNo, "呼呼管理员...")
	if masters <= 0 {
		// 等待工人回来消息
		mu_masters.Unlock()
		cond_masterBack.L.Lock()
		// fmt.Println(tCallNo, "等待管理员...")
		cond_masterBack.Wait()
		cond_masterBack.L.Unlock()
		// fmt.Println(tCallNo, "重新叫号")
		callMaster()
	} else {
		masters--
		// fmt.Println(">>> >>> >>> 给", tCallNo, " 派出了一个管理者，现在有管理者:", masters)
		mu_masters.Unlock()
	}
}

func masterBack() {
	mu_masters.Lock()
	masters++
	// fmt.Println("<<< <<< <<< 一个管理者回来了，现在有管理者:", masters)
	mu_masters.Unlock()
	cond_masterBack.Signal()
}
