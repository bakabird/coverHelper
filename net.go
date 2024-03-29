package main

import (
	"database/sql"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
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
	// defer conn.Close()

	fmt.Println(">>> 开始处理连接")
	const originHub = "./origin"
	const jpgHub = "./jpg"

	var taskRlt []string
	var tasks = 0
	var cond_tasks = sync.NewCond(new(sync.Mutex))
	var mu_conn sync.Mutex

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
			fmt.Println(`暂时接受到字符串：`, rcvStr)
			if err != nil {
				fmt.Println(`[STEP2]在接受连接输入时出现错误：`, err)
				// BIRDTODO: 如果执行出错应该返回一个错误给连接端
				return
			} else {
				coverUrlStrings += rcvStr[:n]
				if len(coverUrlStrings) >= urlsStrLen {
					coverUrls := strings.Split(coverUrlStrings, " ")
					tasks = len(coverUrls)
					taskRlt = make([]string, tasks)
					for index, coverUrl := range coverUrls {
						callWorker()
						go dealSave(db, index, coverUrl, taskRlt, cond_tasks, &tasks, &mu_conn, conn)
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
	dealReturn(conn, strings.Join(taskRlt, " "))
}

func dealSave(db *sql.DB, index int, coverUrl string, taskRlt []string, cond_tasks *sync.Cond, tasks *int, mu_conn *sync.Mutex, conn net.Conn) {
	defer workerBack()
	defer taskFinish(tasks, cond_tasks)
	defer (func() {
		(*mu_conn).Lock()
		conn.Write([]byte("W"))
		(*mu_conn).Unlock()
	})()

	isBilibCover := isBilibCover(coverUrl)
	if isBilibCover {
		isExist := db_muIsExist(db, coverUrl)
		if isExist {
			whenIsExist(db, index, coverUrl, taskRlt)
		} else {
			fmt.Println(`遇到没保存过图片的URL`, coverUrl)
			DAM_LOCK()
			if DAM_isExist(coverUrl) {
				if DAM_isDoing(coverUrl) {
					DAM_UNLOCK()
					whenOtherWorking(db, index, coverUrl, taskRlt)
				} else if DAM_isComplete(coverUrl) {
					DAM_UNLOCK()
					whenJobDone(db, index, coverUrl, taskRlt)
				} else {
					DAM_UNLOCK()
					whenError(index, taskRlt)
				}
			} else {
				DAM_doing(coverUrl)
				DAM_UNLOCK()
				whenWorking(db, index, coverUrl, taskRlt)
			}
		}
	} else {
		fmt.Println(`该封面URL不需要被转译`)
		taskRlt[index] = coverUrl
	}
}

func downloadPic(coverUrl string) string {
	urlParseRes, err := url.Parse(coverUrl)
	if err != nil {
		fmt.Println(`解析url时遇到问题：`, coverUrl, err)
		panic(err)
	}

	makeTempSubHostDir(urlParseRes.Host)

	base := filepath.Base(urlParseRes.Path)
	coverPath := filepath.Join(coverFileTempHub, "./"+urlParseRes.Host+"/", base)
	err = Download(coverUrl, coverPath)
	if err != nil {
		fmt.Println(`下载时遇到错误：`, coverUrl, err)
		panic(err)
	}
	fmt.Println("下载成功：", coverPath)
	return coverPath
}

func saveLocalCover(coverUrl string, img image.Image) string {
	urlParseRes, err := url.Parse(coverUrl)
	if err != nil {
		fmt.Println(`解析url时遇到问题：`, coverUrl, err)
		panic(err)
	}

	makeSubHostDir(urlParseRes.Host)

	base := filepath.Base(urlParseRes.Path)
	newFilePath := filepath.Join(coverFileHub, "./"+urlParseRes.Host+"/", base)
	newFileExt := filepath.Ext(newFilePath)

	newFile, err := os.Create(newFilePath)
	defer newFile.Close()
	if err != nil {
		fmt.Println(`新建文件时遇到错误：`, newFilePath, err)
		panic(err)
	}

	if newFileExt == ".png" {
		_ = png.Encode(newFile, img)
		fmt.Println("保存PNG文件成功", newFilePath)
	} else {
		_ = jpeg.Encode(newFile, img, nil)
		fmt.Println("保存文件成功", newFilePath)
	}

	return newFilePath
}

func makeSubHostDir(hostname string) {
	err := os.MkdirAll(path.Join(coverFileHub, "./"+hostname+"/"), os.ModePerm)
	if err != nil {
		fmt.Println(`域名子文件夹创建失败`)
		panic(err)
	}
}

func makeTempSubHostDir(hostname string) {
	err := os.MkdirAll(path.Join(coverFileTempHub, "./"+hostname+"/"), os.ModePerm)
	if err != nil {
		fmt.Println(`域名子文件夹创建失败`)
		panic(err)
	}
}

func isBilibCover(coverUrl string) bool {
	urlRes, _ := url.Parse(coverUrl)
	return strings.Index(urlRes.Host, "hdslb.com") > -1
}

func whenIsExist(db *sql.DB, index int, coverUrl string, taskRlt []string) {
	fmt.Println(`该URL的图片已经保存好了`, coverUrl)
	path := db_muGet(db, coverUrl)
	taskAddPath(taskRlt, index, path)
}

func whenOtherWorking(db *sql.DB, index int, coverUrl string, taskRlt []string) {
	fmt.Println(`任务由其它的工人完成，等待...`)
	DAM_condWait(coverUrl)

	path := db_muGet(db, coverUrl)
	taskAddPath(taskRlt, index, path)
	fmt.Println(`任务由其它的工人完成了！可喜可贺可喜可贺`)
}

func whenJobDone(db *sql.DB, index int, coverUrl string, taskRlt []string) {
	path := db_muGet(db, coverUrl)
	taskAddPath(taskRlt, index, path)
	fmt.Println(`任务已经由其它的工人完成过了`)
}

func whenError(index int, taskRlt []string) {
	taskRlt[index] = ""
	fmt.Println(`DAM 执行出错了 遇到了奇怪的 actionStats`)
}

func whenWorking(db *sql.DB, index int, coverUrl string, taskRlt []string) {
	coverPath := downloadPic(coverUrl)

	img, err := resizeImg(coverPath, 320, 200)
	if err != nil {
		fmt.Println(`调整图片大小时遇到错误：`, coverPath, err)
		return
	}
	fmt.Println("调整大小成功")

	newFilePath := saveLocalCover(coverUrl, img)

	// 在数据库中保存 coverUrl-path 对
	db_muSavePair(db, coverUrl, newFilePath)

	DAM_complete(coverUrl)
	DAM_getCond(coverUrl).Broadcast()
	// 添加对 taskRlt
	taskAddPath(taskRlt, index, newFilePath)
}

func taskAddPath(taskRlt []string, index int, coverPath string) {
	taskRlt[index] = "http://" + path.Join(*STATIC_COVER_HUB, coverPath)
}

func dealReturn(conn net.Conn, taskRlt string) {
	fmt.Println(`>>> 开始准备返回数据`)

	// -> OK
	conn_back_ok(conn)

	// <- R(Rlt)
	revBuf := make([]byte, 1024)
	n, err := conn.Read(revBuf)
	rcvStr := string(revBuf[:n])
	if err != nil {
		fmt.Println(`在等待连接返回 R 时出现错误：`, err)
	} else if rcvStr[0] == 'R' {
		// -> R{length of taskRlt}
		fmt.Println(`返回数据的长度L(len)`, len(taskRlt))
		conn.Write([]byte("L" + strconv.Itoa(len(taskRlt))))

		revBuf = make([]byte, 1024)
		n, err = conn.Read(revBuf)
		rcvStr = string(revBuf[:n])
		if err != nil {
			fmt.Println(`在等待连接返回 C(continue) 时出现错误：`, err)
		} else if rcvStr[0] == 'C' {
			// -> taskRlt
			fmt.Println(`返回数据：`, taskRlt)
			conn.Write([]byte(taskRlt))

			// <- O(OVER)
			revBuf = make([]byte, 1024)
			n, err = conn.Read(revBuf)
			rcvStr = string(revBuf[:n])
			if rcvStr[0] != 'O' {
				fmt.Println(`没能正常收到连接返回的 O`, rcvStr)
			}
		}
	} else {
		fmt.Println(`在准备返回数据过程中出错了，发送OK之后收到信息：`, rcvStr)
	}
}
