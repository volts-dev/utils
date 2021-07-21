package utils

import (
	"errors"
	"fmt"
	"io"

	"os"
	"path/filepath"

	"strings"
)

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

//filter: jpg|go|png
func CopyDir(source string, dest string, filter string) (err error) {
	if source == dest {
		return errors.New("Can not copy to same dir !")
	}
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer, filter)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// 过滤文件格式
			lIndex := strings.Index(strings.ToLower(filter), strings.ToLower(filepath.Ext(sourcefilepointer))[1:]) // 不包含"."号
			if lIndex == -1 {
				// perform copy
				err = CopyFile(sourcefilepointer, destinationfilepointer)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

	}
	return
}
