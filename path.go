package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// execPath returns the executable path.
func AppFilePath() string {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return ""
	}
	return path
}

func AppPath() string {
	return filepath.Dir(AppFilePath())
}

// execPath returns the executable path.
func AppDir() string {
	return filepath.Base(AppPath())
}

// TEST
// return current file path including file name
func CurFilePath() string {
	//info, file, line, ok := runtime.Caller(3)
	_, file, _, _ := runtime.Caller(3)
	//path, _ := filepath.Split(file)
	return file
}

// TODO: 有缺陷
// get current file path without file name
func _cur_path() string {
	//info, file, line, ok := runtime.Caller(3)
	_, file, _, _ := runtime.Caller(1) // level 3
	fmt.Println("pa", file, CurFilePath())
	path, _ := path.Split(file)
	return path
}

// could not call _cur_path directly to get path. it must keep in 3 level
func CurPath() string {
	_, file, _, _ := runtime.Caller(2) // level 3
	path, _ := path.Split(file)
	return path
}

// get current file dir name
func CurDirName() string {
	_, file, _, _ := runtime.Caller(3) // level 3
	path, _ := path.Split(file)
	return filepath.Base(path)
}

// 获得文件信息,如果文件存在不出错且不是文件夹
func FileExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	//fmt.Println("FileIsDir", info.IsDir())
	return !info.IsDir()
}

func DirExists(dir string) bool {
	d, e := os.Stat(dir) // 返回FileInfo
	switch {
	case e != nil:
		return false
	case !d.IsDir(): // 是文件返回false
		return false
	}

	return true
}

func FilePathToPath(src string) string {
	return strings.Replace(src, `\`, `/`, -1)
}

func OpenDir(dir string) (filenames, dirs []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	for _, info := range files {
		if !info.IsDir() {
			filenames = append(filenames, info.Name())
		} else {
			dirs = append(dirs, info.Name())
		}
	}
	return
}

func GetFileLists(dir string) []string {
	filenames, _ := OpenDir(dir)
	var fs []string
	for _, file := range filenames {
		fs = append(fs, file)
	}
	return fs
}

func JoinURL(org, src string) string {
	// 当字符串里有["/"]时,保留它.某些URL /search/和/search不同
	url := path.Join(org, src)
	//log.Println(url, org, src)
	if strings.HasSuffix(src, "/") && src != "/" {
		return url + "/"
	} else {
		return url
	}
	///	Debug(url, org, src)
	return ""
}
