package tools

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
hashKey - 指定格式的hashkey值，根据传入参数生成的md5数据，用来比较使用
filePath - 当前文件的全路径
*/
type FileHashInfo struct {
	HashKey  string `json:"hashKey"`
	FilePath string `json:"filePath"`
}

/**
创建目录下所有文件的日志
日志格式e.g:
fileinfo-2020-10-20-11-37-02.log
*/
func CreateFileItemsLogFile(filenamePrefix string) string {
	//截取日期格式2020-10-20 11-37-02
	curtime := strings.Replace(time.Now().String()[:19], ":", "-", -1)
	outlogpath := filenamePrefix + curtime + ".log"
	return strings.Replace(outlogpath, " ", "-", -1)
}

/**
将FileHashInfo以json写入日志
传入需要转json的对象
*/
func FileItemToLog(outlogpath string, p interface{}) {
	logfile, err := os.OpenFile(outlogpath, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("创建日志失败", err)
	}
	defer logfile.Close()

	write := bufio.NewWriter(logfile)
	if b, err := json.Marshal(&p); err == nil {
		write.WriteString(string(b))
		write.WriteString("\n")
	}
	write.Flush()
}

func FileItemDataReceiver(outlogname string, fch chan FileHashInfo, wg *sync.WaitGroup) {
	go func() {
		for {
			if data, ok := <-fch; ok {
				//fmt.Println(data)
				FileItemToLog(outlogname, data)
			} else {
				break
			}
		}
		wg.Done()
	}()
}

func AsyncFileItemService(rootpath string, wg *sync.WaitGroup, fileHashRule string) chan FileHashInfo {
	retCh := make(chan []FileHashInfo, 1)
	fileHashInfoArray := [] FileHashInfo{}
	go func() {
		ret := FileItemsDataProducer(rootpath, rootpath, fileHashInfoArray, fileHashRule)
		//fmt.Println("returned result",ret)
		retCh <- ret
		//fmt.Println("service exited")
	}()
	fileHashInfoArray = <-retCh
	writeCh := make(chan FileHashInfo)
	wg.Add(1)
	go func() {
		for _, v := range fileHashInfoArray {
			writeCh <- v
		}
		close(writeCh)
		wg.Done()
	}()
	return writeCh
}

func FileItemsDataProducer(rootPrefix string, rootpath string, fileHashInfoArray []FileHashInfo, fileHashRule string) []FileHashInfo {
	dir, err := ioutil.ReadDir(rootpath)
	if err != nil {
		panic(errors.New("ReadDir error"))
	}
	pthSep := string(os.PathSeparator)

	var hashkey string
	fileHashInfo := FileHashInfo{}
	for _, itm := range dir {
		if itm.IsDir() {
			newPath := rootpath + pthSep + itm.Name()
			fileHashInfoArray = FileItemsDataProducer(rootPrefix, newPath, fileHashInfoArray, fileHashRule)
		} else {
			fileHashInfo.FilePath = rootpath + pthSep + itm.Name()
			switch fileHashRule {
			case "1100":
				temp_path := rootpath + pthSep + itm.Name()
				hashkey = strings.Replace(temp_path, rootPrefix, "", 1)
			case "1111":
				temp_path := rootpath + pthSep + itm.Name()
				temp_path = strings.Replace(temp_path, rootPrefix, "", 1)
				hashkey = temp_path + strconv.FormatInt(itm.Size(), 10) + strconv.Itoa(itm.ModTime().Second())
			//default 1110
			default:
				//将文件d:\1\2\3\4.txt 替换为1\2\3\4.txt 其中d:\为用户传入的根路径
				temp_path := rootpath + pthSep + itm.Name()
				temp_path = strings.Replace(temp_path, rootPrefix, "", 1)
				//替换掉跟路径的目录后的 路径及文件名及大小作为hash
				hashkey = temp_path + strconv.FormatInt(itm.Size(), 10)
				//fmt.Printf("rootpath: %s, itmName: %s, temppath:%s ,hashkey:%s \n",rootpath,itm.Name(),temp_path,hashkey)
			}

			Md5Inst := md5.New()
			Md5Inst.Write([]byte(hashkey))
			hashResult := Md5Inst.Sum(nil) //和Md5Inst.Sum([]byte(""))有什么区别？

			//fmt.Print("md5hashkey:"+hex.EncodeToString(hashResult)+"\n")
			fileHashInfo.HashKey = hex.EncodeToString(hashResult)
			fileHashInfoArray = append(fileHashInfoArray, fileHashInfo)
		}
	}
	return fileHashInfoArray
}
