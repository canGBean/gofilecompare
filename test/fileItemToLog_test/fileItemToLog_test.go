package fileItemToLog_test

import (
	"comparetool/tools"
	"testing"
)

func TestFileItemToLog(t *testing.T) {
	logpath := tools.CreateFileItemsLogFile()
	f := tools.FileHashInfo{}
	f.HashKey = "123123123"
	f.FilePath = "D:/123/12/31/231/23/111.txt"

	tools.FileItemToLog(logpath, f)
	f.HashKey = "aaa"
	f.FilePath = "D:/123/12/31/231/23/111.txt"
	tools.FileItemToLog(logpath, f)
}
