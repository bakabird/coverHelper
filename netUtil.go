package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nfnt/resize"
)

func Download(coverUrl string, fPath string) error {
	// creat File
	newFile, err := os.Create(fPath)
	defer newFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	client := http.Client{Timeout: 900 * time.Second}
	resp, err := client.Get(coverUrl)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func resizeImg(sourcePath string, aimW, aimH int) (image.Image, error) {
	file, _ := os.Open(sourcePath)
	defer file.Close()

	ext := filepath.Ext(sourcePath)
	var img image.Image
	var err error
	if ext == ".png" {
		img, err = png.Decode(file)
	} else {
		img, _, err = image.Decode(file)
	}
	if err != nil {
		fmt.Println("Decode时出错", err.Error())
		return nil, err
	}

	m := resize.Resize(uint(aimW), uint(aimH), img, resize.Lanczos2)
	return m, nil
}

func conn_back_error(conn net.Conn, err error) {
	conn.Write([]byte("err:" + err.Error()))
}
func conn_back_ok(conn net.Conn) {
	conn.Write([]byte("OK"))
}

func printbackError(conn net.Conn, log string, err error) {
	fmt.Println(log)
	conn_back_error(conn, err)
}

func waitTasksFinish(tasks *int, cond *sync.Cond) {
	fmt.Println("开始等待任务完成")
	cond.L.Lock()
	for *tasks > 0 {
		cond.Wait()
		fmt.Println("有一个任务被完成了,剩余任务有", *tasks)
		cond.L.Unlock()
		cond.L.Lock()
	}
	cond.L.Unlock()
	fmt.Println("任务完成了")
}

func taskFinish(tasks *int, cond *sync.Cond) {
	// fmt.Println("有一个任务被完成了")
	cond.L.Lock()
	*tasks--
	cond.L.Unlock()
	cond.Signal()
}
