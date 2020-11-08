package tools

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

/**
FileName - 文件名称
FileSize - 文件大小
FileModification - 文件最后修改时间
FileFullPath - 文件完整路径
FilePathWithOutPrefix - 替换前缀后的文件全
FileReplacePrefix - 需要替换掉FileFullPath中的路径前缀,仅替换1次
HashKey - 指定格式的HashKey
*/
type FileInfo struct {
	FileName              string `json:"fileName"`
	FileSize              string `json:"fileSize"`
	FileModification      string `json:"fileModification"` //yyyy-MM-DD HH:mm:ss
	FileFullPath          string `json:"fileFullPath"`
	FilePathWithOutPrefix string `json:"filePathWithOutPrefix"`
	FileReplacePrefix     string `json:"fileReplacePrefix"`
	HashKey               string `json:"hashKey"`
}

type ExecDifferenceParam struct {
	FilelogMap1  map[string]string
	Prefix1      string
	FilelogMap2  map[string]string
	Prefix2      string
	Logname      string
	CoroutineNum int
}

/**
获取两个日志文件不同的内容
取filelogMap1里有而filelogMap2里没有的
*/
func DifferenceFileHash(filelogMap1 map[string]string, filelogMap2 map[string]string, fch chan FileHashInfo, wg *sync.WaitGroup) {
	var retValue [] FileHashInfo
	fileHashInfo := FileHashInfo{}
	wg.Add(1)
	go func() {
		for k, v := range filelogMap1 {
			if _, ok := filelogMap2[k]; ok {
				continue
			} else {
				fileHashInfo.HashKey = k
				fileHashInfo.FilePath = v
				retValue = append(retValue, fileHashInfo)
				fch <- fileHashInfo
			}
		}
		close(fch)
		wg.Done()
	}()
}

/**
执行对比文件
*/
func ExecDifferenceFileHashInfo(edp ExecDifferenceParam) {
	retCh1 := make(chan FileHashInfo, 1)
	retCh2 := make(chan FileHashInfo, 1)
	var wg sync.WaitGroup
	DifferenceFileHash(edp.FilelogMap1, edp.FilelogMap2, retCh1, &wg)
	DifferenceFileHash(edp.FilelogMap2, edp.FilelogMap1, retCh2, &wg)

	for i := 0; i < edp.CoroutineNum; i++ {
		combineDifferenceFileHash(&wg, retCh1, edp.Prefix1, edp.Logname)
		combineDifferenceFileHash(&wg, retCh2, edp.Prefix2, edp.Logname)
	}

	wg.Wait()
}

/**
获取文件信息
pathPrefix string : 目录前缀
fileHashInfo FileHashInfo:
*/
func GetFileInfo(pathPrefix string, fileHashInfo FileHashInfo) FileInfo {
	fileInfoTemp, err := os.Stat(fileHashInfo.FilePath)
	fileInfo := FileInfo{}
	if err != nil {
		panic(errors.New("获取文件:" + fileHashInfo.FilePath + "信息失败"))
	}
	fileInfo.FileName = fileInfoTemp.Name();
	fileInfo.FileFullPath = fileHashInfo.FilePath
	fileInfo.FileModification = strings.Replace(fileInfoTemp.ModTime().String()[:23], ":", "-", -1)
	if len(pathPrefix) != 0 {
		fileInfo.FileReplacePrefix = pathPrefix
		fileInfo.FilePathWithOutPrefix = strings.Replace(fileHashInfo.FilePath, pathPrefix, "", 1)
	}
	fileInfo.FileSize = strconv.FormatInt(fileInfoTemp.Size(), 10)
	fileInfo.HashKey = fileHashInfo.HashKey
	return fileInfo
}

/**
合并文件差集并写入日志
*/
func combineDifferenceFileHash(wg *sync.WaitGroup, inCh chan FileHashInfo, pathPrefix string, outlogpath string) {
	wg.Add(1)
	go func() {
		for {
			if data, ok := <-inCh; ok {
				fileInfo := GetFileInfo(pathPrefix, data)
				FileItemToLog(outlogpath, fileInfo)
			} else {
				break
			}
		}
		wg.Done()
	}()
}

func LoadLogToMem(logfilepath string) map[string]string {
	logfile, err := os.OpenFile(logfilepath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer logfile.Close()
	fileHashMap := map[string]string{}
	reader := bufio.NewReader(logfile)
	for {
		fileHashData, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		var fileHashInfo FileHashInfo
		if err := json.Unmarshal([]byte(fileHashData), &fileHashInfo); err != nil {
			panic("读取日志文件转换失败:" + logfilepath)
		}
		fileHashMap[fileHashInfo.HashKey] = fileHashInfo.FilePath
	}
	return fileHashMap
}
