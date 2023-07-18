// 统计文件个数
// 传入一个目录，递归攻击当前目录中的所有文件个数
// 使用一个缓冲通道来进行并发处理。缓冲通道是为了控制并发量，防止并发过多导致内存溢出。
// 因为读取目录的时候，是一个IO操作，需要将它放到一个goroutine里面去运行

package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

// 创建一个files chan
var files_chan = make(chan string, 100)

// 控制通道的数量 改变这个参数可以明显的看到运行时间的变化，但是不宜过高
// 可以使用 ulimit -n 查看系统的最大打开文件数
var chan_limit = make(chan bool, 100)

var wg sync.WaitGroup

func main() {
	// 获取当前目录
	dirPath, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前目录失败")
		return
	}

	// 如果有参数传入，则使用参数作为目录
	if len(os.Args) > 1 {
		dirPath = os.Args[1]
	}

	start := time.Now()
	// 计算运行时间
	defer func() {
		fmt.Println("运行时间：", time.Since(start))
	}()

	wg.Add(1)
	go readDir(dirPath)

	go func() {
		wg.Wait()
		close(files_chan)
	}()

	// 文件个数
	var fileCount int

	// 给文件去重
	fileMap := make(map[string]int)

	// 读取chan中的数据
	for fileName := range files_chan {
		//fmt.Println(fileName)

		fileExt := path.Ext(fileName)
		fileName = strings.TrimSuffix(fileName, fileExt)
		fileCount++
		fileMap[fileName]++

	}

	fmt.Println("文件个数：", fileCount)

	// 去重后的文件个数
	fmt.Println("去重后的文件个数：", len(fileMap))

}

// 读取目录
func readDir(dirPath string) {
	// 限制goroutine数量
	chan_limit <- true
	defer func() {
		<-chan_limit
		wg.Done()
	}()

	// 打开目录
	dir, err := os.Open(dirPath)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	// 读取目录中的文件
	files, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		// 排除一些隐藏文件
		if file.Name()[0] == '.' {
			continue
		}

		// 如果是目录，则递归获取文件个数
		if file.IsDir() {
			// 写入通道
			wg.Add(1)
			go readDir(dirPath + "/" + file.Name())
		} else {
			files_chan <- file.Name()
		}
	}
}
