package helper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"unsafe"
)

var runningDir string

// RunningDir 获取当前执行文件所在的目录
func RunningDir() string {
	if runningDir == "" {
		runningDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	return runningDir
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// EachFileCallback 遍历指定目录的文件名，并应用回调函数
func EachFileCallback(path string, cb func(name string)) {
	files, _ := ioutil.ReadDir(path)
	for _, fi := range files {
		if !fi.IsDir() {
			cb(fi.Name())
		}
	}
}

// Str2bytes 字串转为bytes
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Bytes2str bytes转为字串
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
