package main

import (
	"fmt"
	"os"
)

/*
统计文件个数
传入一个目录，递归攻击当前目录中的所有文件个数
*/

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

	count, err := countFiles(dirPath)
	if err != nil {
		fmt.Println("获取文件个数失败")
		return
	}

	fmt.Println("文件个数：", count)

}

// 递归获取目录下的文件个数
func countFiles(dirPath string) (int, error) {
	// 打开目录
	dir, err := os.Open(dirPath)
	if err != nil {
		return 0, err
	}
	defer dir.Close()

	// 读取目录中的文件
	files, err := dir.Readdir(-1)
	if err != nil {
		return 0, err
	}

	// 遍历文件
	count := 0
	for _, file := range files {
		// 如果是目录，则递归获取文件个数
		if file.IsDir() {
			subCount, err := countFiles(dirPath + "/" + file.Name())
			if err != nil {
				return 0, err
			}
			count += subCount
		} else {
			// 排除一些隐藏文件
			if file.Name()[0] == '.' {
				continue
			}

			fmt.Println(file.Name())
			count++
		}
	}

	return count, nil
}
