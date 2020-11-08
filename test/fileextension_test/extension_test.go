package fileextension_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)


type FileHashInfo struct {
	HashKey  string `json:"hashKey"`
	FilePath string `json:"filePath"`
}
type FileInfo struct {
	FileName              string `json:"fileName"`
	FileSize              string `json:"fileSize"`
	FileModification      string `json:"fileModification"` //yyyy-MM-DD HH:mm:ss
	FileFullPath          string `json:"fileFullPath"`
	FilePathWithOutPrefix string `json:"filePathWithOutPrefix"`
	FileReplacePrefix     string `json:"fileReplacePrefix"`
}

func (p *FileHashInfo) FileItemToLog(){
}
//
//func (p *FileInfo) FileItemToLog() {
//
//}

func writeFileItemToLog(outlogpath string, p interface{}) {
	write:=openLog(outlogpath)
	if b, err := json.Marshal(&p); err == nil {
		//write.WriteString(string(b))
		//write.WriteString("\n")
		fmt.Println(string(b))
	}
	write.Flush()
}




func openLog(outlogpath string) *bufio.Writer{
	logfile, err := os.OpenFile(outlogpath, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("创建日志失败", err)
	}
	defer logfile.Close()
	return bufio.NewWriter(logfile)
}

func TestExtension(t *testing.T) {
	filehash:=new(FileHashInfo)
	filehash.FilePath="D:\\"
	filehash.HashKey="123123"
	writeFileItemToLog("D:\\",filehash)
	//f:=new(FileInfo)
	//FileItemToLog1("D:\\",f)
	fileInfo:=new(FileInfo)
	fileInfo.FileSize="123"
	fileInfo.FileName="2323"
	writeFileItemToLog("D:\\111",fileInfo)

}
