package main

import (
	"database/sql"
	"fmt"
	"image/jpeg"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func startListen(db *sql.DB) {
	ln, err := net.Listen("tcp", "127.0.0.1:25005")
	fmt.Println("开始监听")
	defer ln.Close()

	if err != nil {
		// handle error
		fmt.Println(`在开启监听过程中遇到问题`, err)
		return
	}

	for {
		conn, err := ln.Accept()
		fmt.Println(">>> 接受到访问")
		if err != nil {
			// handle error
			fmt.Println(`等待访问过程中遇到问题`, err)
			conn.Close()
			continue
		}

		callMaster()
		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *sql.DB) {
	defer masterBack()
	defer conn.Close()

	fmt.Println(">>> 开始处理连接")
	const originHub = "./origin"
	const jpgHub = "./jpg"

	var tasks = 0
	var cond_tasks = sync.NewCond(new(sync.Mutex))

	// step.1 接受命令，并了解将发送的字符长度
	revBuf := make([]byte, 1024)
	n, err := conn.Read(revBuf)
	rcvStr := string(revBuf[:n])

	if err != nil {
		fmt.Println(`[STEP1]在接受连接输入时出现错误：`, err)
	} else if rcvStr[0] == 'C' {
		urlsStrLen, _ := strconv.Atoi(rcvStr[1:])
		fmt.Println(">>> [STEP1]接受到C命令，即将发送的URLS字符串共有字符数量：", urlsStrLen)
		// G = GOON
		conn.Write([]byte("G"))
		fmt.Println("<<< [STEP1]返回G命令，示意继续")
		var coverUrlStrings = ""
		for {
			n, err := conn.Read(revBuf)
			rcvStr := string(revBuf[:n])
			if err != nil {
				fmt.Println(`[STEP2]在接受连接输入时出现错误：`, err)
				return
			} else {
				coverUrlStrings += rcvStr[:n]
				if len(coverUrlStrings) >= urlsStrLen {
					coverUrls := strings.Split(coverUrlStrings, " ")
					tasks = len(coverUrls)
					for _, coverUrl := range coverUrls {
						callWorker()
						go dealSave(db, coverUrl, cond_tasks, &tasks)
					}
					break
				} else {
					fmt.Println(`>>> [STEP2]已读取字符串 `, len(coverUrlStrings), ` ,继续读取字符串`)
				}
			}
		}
	} else {
		fmt.Println(">>> [STEP1]未知的命令")
	}

	waitTasksFinish(&tasks, cond_tasks)
	conn_back_ok(conn)
}

func dealSave(db *sql.DB, coverUrl string, cond_tasks *sync.Cond, tasks *int) {
	defer workerBack()
	defer taskFinish(tasks, cond_tasks)
	mu_db.Lock()
	isExist := db_IsExist(db, coverUrl)
	mu_db.Unlock()
	if isExist {
		fmt.Println(`该URL的图片已经保存好了`, coverUrl)
	} else {
		fmt.Println(`遇到没保存过图片的URL`, coverUrl)
		DAM_LOCK()
		if DAM_isExist(coverUrl) {
			if DAM_isDoing(coverUrl) {
				DAM_UNLOCK()
				fmt.Println(`任务由其它的工人完成，等待...`)
				DAM_getCond(coverUrl).L.Lock()
				DAM_getCond(coverUrl).Wait()
				DAM_getCond(coverUrl).L.Unlock()
				fmt.Println(`任务由其它的工人完成了！可喜可贺可喜可贺`)
			} else if DAM_isComplete(coverUrl) {
				DAM_UNLOCK()
				fmt.Println(`任务已经由其它的工人完成过了`)
			} else {
				DAM_UNLOCK()
				fmt.Println(`DAM 执行出错了 遇到了奇怪的 actionStats`)
			}
		} else {
			DAM_doing(coverUrl)
			DAM_UNLOCK()
			urlParseRes, err := url.Parse(coverUrl)
			if err != nil {
				fmt.Println(`解析url时遇到问题：`, coverUrl, err)
				return
			}

			base := filepath.Base(urlParseRes.Path)
			coverPath := filepath.Join(coverFileTempHub, base)

			err = Download(coverUrl, coverPath)
			if err != nil {
				fmt.Println(`下载时遇到错误：`, coverUrl, err)
				return
			}
			fmt.Println("下载成功：", coverPath)

			img, err := resizeImg(coverPath, 640, 400)
			if err != nil {
				fmt.Println(`调整图片大小时遇到错误：`, coverPath, err)
				return
			}
			fmt.Println("调整大小成功")

			jpgPath := path.Join(coverFileHub, base)
			jpgout, err := os.Create(jpgPath)
			defer jpgout.Close()
			if err != nil {
				fmt.Println(`新建文件时遇到错误：`, jpgPath, err)
				return
			}

			_ = jpeg.Encode(jpgout, img, nil)
			fmt.Println("保存 jpg 成功", jpgPath)

			// 在数据库中保存 coverUrl-path 对
			mu_db.Lock()
			db_SavePair(db, coverUrl, jpgPath)
			mu_db.Unlock()

			DAM_LOCK()
			DAM_complete(coverUrl)
			DAM_UNLOCK()
			DAM_getCond(coverUrl).Broadcast()
		}
	}
}
