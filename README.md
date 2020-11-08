# 说明
本工具主要实现2个功能:
1. 根据指定路径,遍历改路径下的所有文件（包含子文件夹下的文件）,最终生成json格式的多行日志文件。在生成的日志文件中,并不会输出文件夹及子文件夹的内容。输出格式为自定义的文件hash值与文件路径,并以换行符\n结束。日志输出格式如下:
```log
...
{"hashKey":"cb15034aae795b651b8654d81bb97e53","filePath":"D:\\aaa\\ccc\\ddd.txt"}
{"hashKey":"fd65c5509a1a4074d61ea8ae3bc9c301","filePath":"D:\\aaa\\ccc\\dd11.txt"}
...
```
2. 比较2个功能1中生成的日志文件,并生成这2个日志中的全差集数据文件。

# 数据格式
1.  使用遍历功能（功能1）,在遍历文件夹路径后,生成的文件名称格式为"fileinfo-yyyy-MM-DD-HH-mm-ss.log"
2.  使用对比功能（功能2）,在比较2个日志文件后,生成的文件名格式为"filedifferent--yyyy-MM-DD-HH-mm-ss.log"
3. 使用遍历功能（功能1）,生成的日志中文件属性含义为:
```go
/**
HashKey  - 指定格式的hashkey值,根据传入参数生成的md5数据,用来比较使用
FilePath - 当前文件的全路径
*/
type FileHashInfo struct {
	HashKey  string `json:"hashKey"`
	FilePath string `json:"filePath"`
}
```
4.  使用对比功能（功能2）,生成的日志中文件属性含义为:
```go
/**
FileName - 文件名称
FileSize - 文件大小,bytes
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
```

# 使用方法及命令说明
## 初始命令参数
在使用时,首先需要指定初始命令参数,可选参数为init与compare。其中init是指使用功能1,compare是指使用功能2,在win10-64bit环境下,执行如下:
```dos
>fct.exe -h
Usage:
  fct.exe [OPTIONS] <compare | init>

Help Options:
  /?          Show this help message
  /h, /help   Show this help message

Available commands:
  compare
  init
```
##  init命令说明
当执行init命令时,主要有3个可选参数,如下:
```dos
>fct.exe init -h
Usage:
  fct.exe [OPTIONS] init [init-OPTIONS]

Help Options:
  /?                         Show this help message
  /h, /help                  Show this help message

[init command options]
      /n, /num:              set groutine num (default: 1)
      /p, /traverseRootDir:  required,travers directory path
      /r, /hashRule:         filePath+fileName+fileSize+LastModificationTime
                             more infomation read README.md (default: 1110)
```
* /num参数 :短命令为/n或-n 。用来指定生成输出文件时使用的协程数量。默认为1。可选参数。
* /p参数:短命令为/p或-p 。用来指定需要遍历的文件夹。必选参数。
* /r参数: 短命令为/p或-p。用来指定输出文件中,hashKey属性的组成规则。这里的hashkey组成规则包含“文件路径+文件名称+文件大小+最后修改时间”每个属性如果使用则为1,如果不使用则为零。目前工具中支持3中hashkey的生成规则。说明如下:
> 1. 如果以“文件名称+文件路径+大小”做为该文件的hashkey,则 -r的值为1110,也是工具中默认的hashkey生成模式
>2. 如果以“文件名称+文件路径”做为该文件的hashkey,则-r的值为1100
>3. 如果以“文件名称+文件路径+大小+最后修改时间”做为该文件的hashkey,则-r的值为1111
当执行以下命令后,会在当前路径生成一个d:\A目录中的所有文件的日志
```dos
>fct.exe init -p d:\A -n 5 -r 1111
```

##  compare命令说明
compare命令为比较init命令生成的日志文件,这里支持2个文件的比较,并根据hashkey筛选出2个日志中的差集文件。这里比较时需要2个日志文件的hashkey生成方式相同,否则没有比较意义。compare命令参数如下:
```dos
>fct.exe compare -h
Usage:
  fct.exe [OPTIONS] compare [compare-OPTIONS]

Help Options:
  /?                  Show this help message
  /h, /help           Show this help message

[compare command options]
      /f, /filelog1:  required,the first fileinfo log file path
      /y, /prefix1:   replace fileinfo path  prefix in the first filelog
      /s, /filelog2:  required,the second fileinfo log file path
      /z, /prefix2:   replace fileinfo path  prefix in the seconde filelog
```
* /filelog1参数 :短命令为/f或-f 。用来指定第1个输入的日志路径,必选
* /prefix1参数 :短命令为/y或-y 。用来指定替换第1个日志中的FilePath 的路径前缀,并将值赋予FilePathWithOutPrefix 属性
* /filelog2参数 :短命令为/s或-s 。用来指定第2个输入的日志路径,必选
* /prefix2参数 :短命令为/z或-z。用来指定替换第2个日志中的FilePath 的路径前缀,并将值赋予FilePathWithOutPrefix 属性

### 举例说明:
1. 在init阶段生成了2个日志结果文件分别为fileinfo-log1.log与fileinfo-log2.log内容如下:
fileinfo-log1.log
```log
{"hashKey":"cb15034aae795b651b8654d81bb97e53","filePath":"D:\\aaa\\ccc\\AAA.txt"}
{"hashKey":"fd65c5509a1a4074d61ea8ae3bc9c301","filePath":"D:\\aaa\\ccc\\BBB.txt"}
```
fileinfo-log2.log
```log
{"hashKey":"cb15034aae795b651b8654d81bb97e53","filePath":"D:\\aaa\\bbb\\AAA.txt"}
{"hashKey":"dfc832a597b4539c9eaae5a084b41672","filePath":"D:\\aaa\\bbb\\CCC.txt"}
```
2.此时需要比较这两个log的内容,在不替换指定路径前缀的时,可使用命令如下:
```dos
>fct.exe compare -f fileinfo-log1.log -s fileinfo-log2.log
```
则输出的结果为:
```
{"fileName":"BBB.txt","fileSize":"123","fileModification":"2016-09-26 23-58-40.920","fileFullPath":"D:\\aaa\\ccc\\BBB.txt","filePathWithOutPrefix":"","fileReplacePrefix":"","hashKey":"fd65c5509a1a4074d61ea8ae3bc9c301"}

{"fileName":"CCC.txt","fileSize":"666","fileModification":"2016-02-23 20-18-40.920","fileFullPath":"D:\\aaa\\bbb\\CCC.txt","filePathWithOutPrefix":"","fileReplacePrefix":"","hashKey":"dfc832a597b4539c9eaae5a084b41672"}
```
3. 此时需要比较这两个log的内容,并替换掉内容为:
> * fileinfo-log1.log中filePath的"D:\\aaa"
> * fileinfo-log2.log中filePath的"D:\\aaa\\bbb"

可使用命令如下:
```dos
fct.exe compare -f fileinfo-log1.log -y D:\\aaa -s fileinfo-log2.log -z D:\\aaa\\bbb
```
则输出的结果为:
```
{"fileName":"BBB.txt","fileSize":"123","fileModification":"2016-09-26 23-58-40.920","fileFullPath":"D:\\aaa\\ccc\\BBB.txt","filePathWithOutPrefix":"\\ccc\\BBB.txt","fileReplacePrefix":"D:\\aaa","hashKey":"fd65c5509a1a4074d61ea8ae3bc9c301"}

{"fileName":"CCC.txt","fileSize":"666","fileModification":"2016-02-23 20-18-40.920","fileFullPath":"D:\\aaa\\bbb\\CCC.txt","filePathWithOutPrefix":"CCC.txt","fileReplacePrefix":"D:\\aaa\\bbb","hashKey":"dfc832a597b4539c9eaae5a084b41672"}
```
