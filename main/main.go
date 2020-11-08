package main

import (
	"comparetool/tools"
	"fmt"
	"github.com/jessevdk/go-flags"
	"sync"
	"time"
)

type InitCommand struct {
	CliWgNum           int    `short:"n" long:"num" description:"set groutine num" default:"1"`
	CliTraverseRootDir string `short:"p" long:"traverseRootDir" description:"required,travers directory path"  required:"true"`
	CliFileHashRule    string `short:"r" long:"hashRule" description:"filePath+fileName+fileSize+LastModificationTime more infomation read README.md"  default:"1110"`
}

type CompareCommand struct {
	CliFilelog1       string `short:"f" long:"filelog1" description:"required,the first fileinfo log file path"  required:"true"`
	CliRepPathPrefix1 string `short:"y" long:"prefix1" description:"replace fileinfo path  prefix in the first filelog"  default:""`
	CliFilelog2       string `short:"s" long:"filelog2" description:"required,the second fileinfo log file path"  required:"true"`
	CliRepPathPrefix2 string `short:"z" long:"prefix2" description:"replace fileinfo path  prefix in the seconde filelog"  default:""`
	CliWgNum          int    `short:"n" long:"num" description:"set groutine num" default:"1"`
}

func (p *InitCommand) Execute(args []string) error {
	logname := tools.CreateFileItemsLogFile("fileinfo ")
	fmt.Printf("create logfile: %s done,begin traverse files \n", logname)
	start := time.Now()

	var wg sync.WaitGroup
	writeCh := tools.AsyncFileItemService(p.CliTraverseRootDir, &wg, p.CliFileHashRule)
	for i := 0; i < p.CliWgNum; i++ {
		wg.Add(1)
		tools.FileItemDataReceiver(logname, writeCh, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("all files done,logname: %s. time-consuming：%s \n", logname, elapsed)
	return nil
}

func (p *CompareCommand) Execute(args []string) error {
	start := time.Now()
	fmt.Println("begin load log file to memery \n")
	filelogMap1 := tools.LoadLogToMem(p.CliFilelog1)
	filelogMap2 := tools.LoadLogToMem(p.CliFilelog2)
	elapsed := time.Since(start)
	fmt.Printf("load log files done. time-consuming：%s \n", elapsed)

	fmt.Println("begin compare files hash \n")
	start = time.Now()
	logname := tools.CreateFileItemsLogFile("filedifferent ")
	edp := tools.ExecDifferenceParam{filelogMap1, p.CliRepPathPrefix1, filelogMap2, p.CliRepPathPrefix2, logname, p.CliWgNum}
	tools.ExecDifferenceFileHashInfo(edp)
	elapsed = time.Since(start)
	fmt.Printf("compare files done.logname: %s. time-consuming：%s \n", logname, elapsed)
	return nil
}

type Option struct {
	Init    InitCommand    `command:"init"`
	Compare CompareCommand `command:"compare"`
}

/**
按照输入命令，生成指定目录下的所有文件的日志,比如：
输入：gofilecompare.exe init "D:\xxx"
输出：yyyyMMDD HH:mm:ss.log
遍历后输出的格式为
{hashkey:filehash,filepath:"xxx\xxx\xxxx.txt"}

对比两个日志
输入:gofilecompare.exe comp "yyyyMMDD HH:mm:ss1.log" "yyyyMMDD HH:mm:ss2.log"
输出：compyyyyMMDD HH:mm:ss.log
对比后输出的日志格式为
｛data:[{FileHashInfo:FileHashInfo,FileInfo:FileInfo},{FileHashInfo:FileHashInfo,FileInfo:FileInfo}]｝
*/
func main() {
	var opt Option
	flags.Parse(&opt)
}
