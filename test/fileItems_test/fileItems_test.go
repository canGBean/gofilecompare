package fileItems_test

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

type FileHashInfo struct {
	hashKey  string
	filePath string
}
type FileInfo struct {
	fileName         string
	fileSize         string
	fileModification string //yyyy-MM-DD
}

func AsyncFileItemService(rootpath string) chan []FileHashInfo {
	retCh := make(chan []FileHashInfo, 1)
	fileHashInfoArray := [] FileHashInfo{}
	go func() {
		ret := FileItemsDataProducer(rootpath, retCh, fileHashInfoArray)
		fmt.Println("returned result",ret)
		retCh <- ret
		fmt.Println("service exited")
	}()
	return retCh
}

func FileItemsDataProducer(rootpath string, fch chan []FileHashInfo, fileHashInfoArray []FileHashInfo) []FileHashInfo {
	dir, err := ioutil.ReadDir(rootpath)
	if err != nil {
		panic(errors.New("ReadDir error"))
	}
	pthSep := string(os.PathSeparator)
	Md5Inst := md5.New()
	var hashkey string
	fileHashInfo := FileHashInfo{}
	for _, itm := range dir {
		if itm.IsDir() {
			newPath := rootpath + pthSep + itm.Name()
			fmt.Println(newPath, "----------DIR")
			fileHashInfoArray = FileItemsDataProducer(newPath, fch, fileHashInfoArray)
		} else {
			fileHashInfo.filePath = rootpath + pthSep + itm.Name()
			hashkey = rootpath + pthSep + itm.Name() + strconv.Itoa(itm.ModTime().Second())
			Md5Inst.Write([]byte(hashkey))
			hashResult := Md5Inst.Sum([]byte(""))
			fileHashInfo.hashKey = hex.EncodeToString(hashResult)
			fileHashInfoArray = append(fileHashInfoArray, fileHashInfo)
			fmt.Println(fileHashInfo, "----------FIlE")
		}
	}
	return fileHashInfoArray
}

func TestFileItems(t *testing.T) {
	retCh := AsyncFileItemService("D:\\data0")
	//从channel中取出返回值
	fmt.Println(<-retCh)
}
