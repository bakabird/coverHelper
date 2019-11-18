package main

import (
	"fmt"
	"sync"
)

var downloadTasks = make(map[string]int)
var mu_downloadTasks sync.Mutex
var downloadCond = make(map[string]*sync.Cond)

const actionStatus_notExist = -1
const actionStatus_doing = 0
const actionStatus_complete = 1

func DAM_LOCK() {
	mu_downloadTasks.Lock()
}
func DAM_UNLOCK() {
	mu_downloadTasks.Unlock()
}

func DAM_complete(name string) {
	DAM_LOCK()
	downloadTasks[name] = actionStatus_complete
	DAM_UNLOCK()
}
func DAM_doing(name string) {
	if !DAM_isExist(name) {
		downloadTasks[name] = actionStatus_doing
	} else {
		fmt.Println("[警告]任务", name, "已经存在，不允许重新执行！")
	}
}

func DAM_isComplete(name string) bool {
	return downloadTasks[name] == actionStatus_complete
}
func DAM_isDoing(name string) bool {
	return downloadTasks[name] == actionStatus_doing
}
func DAM_isExist(name string) bool {
	_, exist := downloadTasks[name]
	return exist
}

func DAM_condWait(name string) {
	DAM_getCond(name).L.Lock()
	DAM_getCond(name).Wait()
	DAM_getCond(name).L.Unlock()
}

func DAM_getCond(name string) *sync.Cond {
	_, exist := downloadCond[name]
	if !exist {
		downloadCond[name] = sync.NewCond(new(sync.Mutex))
	}
	return downloadCond[name]
}
