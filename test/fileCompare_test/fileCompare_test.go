package fileCompare_test

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)
import (
	compareTools "comparetool/tools"
)

func AllDifferenceFileHash(filelogMap1 map[string]string, filelogMap2 map[string]string) []compareTools.FileHashInfo {
	var retValue [] compareTools.FileHashInfo
	fileHashInfo := compareTools.FileHashInfo{}
	for k, v := range filelogMap1 {
		if _, ok := filelogMap2[k]; ok {
			continue
		} else {
			fileHashInfo.HashKey = k
			fileHashInfo.FilePath = v
			retValue = append(retValue, fileHashInfo)
		}
	}
	for k, v := range filelogMap2 {
		if _, ok := filelogMap1[k]; ok {
			continue
		} else {
			fileHashInfo.HashKey = k
			fileHashInfo.FilePath = v
			retValue = append(retValue, fileHashInfo)
		}
	}
	return retValue
}

func getFileInfo(pathPrefix string, fileHashInfo compareTools.FileHashInfo) compareTools.FileInfo {
	fileInfoTemp, err := os.Stat(fileHashInfo.FilePath)
	fileInfo := compareTools.FileInfo{}
	if err != nil {
		panic(errors.New("获取文件:"+fileHashInfo.FilePath+"信息失败"))
	}
	fileInfo.FileName = fileInfoTemp.Name();
	fileInfo.FileFullPath = fileHashInfo.FilePath
	fileInfo.FileModification = strings.Replace(fileInfoTemp.ModTime().String()[:27], ":", "-", -1)
	fileInfo.FileReplacePrefix = pathPrefix
	fileInfo.FilePathWithOutPrefix = strings.Replace(fileHashInfo.FilePath,pathPrefix,"",1)
	fileInfo.FileSize = strconv.FormatInt(fileInfoTemp.Size(), 10)
	return fileInfo
}


func FileItemToLog(outlogpath string,fileInfo compareTools.FileInfo) {
	logfile, err := os.OpenFile(outlogpath, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("创建日志失败", err)
	}
	defer logfile.Close()

	write := bufio.NewWriter(logfile)
	if b, err := json.Marshal(&fileInfo); err == nil {
		write.WriteString(string(b))
		write.WriteString("\n")
	}
	write.Flush()
}

func CombineDifferenceFileHash(wg *sync.WaitGroup, inCh chan compareTools.FileHashInfo,pathPrefix string,outlogpath string) {
	wg.Add(1)
	go func() {
		for {
			if data, ok := <-inCh; ok {
				fileInfo:=getFileInfo(pathPrefix,data)
				FileItemToLog(outlogpath,fileInfo)
			} else {
				break
			}
		}
		wg.Done()
	}()
}

func TestFileCompare(t *testing.T) {
	//input
	filelogMap1 := compareTools.LoadLogToMem("D:\\workspace\\wspace4go\\gofilecompare\\fileinfo 2020-10-20 13-22-04.log")

	filelogMap2 := compareTools.LoadLogToMem("D:\\workspace\\wspace4go\\gofilecompare\\fileinfo 2020-10-20 13-24-03.log")



	repPathPrefix1:="D://data1//"
	repPathPrefix2:="D://data0//"

	compareTools.ExecDifferenceFileHashInfo(filelogMap1,repPathPrefix1,filelogMap2,repPathPrefix2)

	//retCh1 := make(chan compareTools.FileHashInfo, 1)
	//retCh2 := make(chan compareTools.FileHashInfo, 1)
	//
	//
	////init
	//logpath := compareTools.CreateFileItemsLogFile("filedifferent ")
	//var wg sync.WaitGroup
	//compareTools.DifferenceFileHash(filelogMap1, filelogMap2, retCh1, &wg)
	//compareTools.DifferenceFileHash(filelogMap2, filelogMap1, retCh2, &wg)
	//CombineDifferenceFileHash(&wg, retCh1,repPathPrefix1,logpath)
	//CombineDifferenceFileHash(&wg, retCh2,repPathPrefix2,logpath)
	//CombineDifferenceFileHash(&wg, retCh2,repPathPrefix2,logpath)
	//wg.Wait()

}
