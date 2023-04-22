package compile_service

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/golang/groupcache/lru"
)

var cache *lru.Cache

func init() {
	cache = lru.New(100)
}

func search(session string) (string, bool) {
	runOut, ok := cache.Get(session)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%s", runOut), true
}

// return "compileOut" "runOUt" "session"
func Compile(src string) (string, string, string) {
	// 计算哈希，写入文件
	h := fnv.New32a()
	h.Write([]byte(src))
	srcHash := h.Sum32()

	session := fmt.Sprintf("%d", srcHash)

	res, ok := search(session)
	if ok {
		return "✨ find on cache!", res, session
	}

	// compile file structure
	//
	// compile_file
	//  └─ session1
	//  ├─ session2
	//  │   └─ session.pl
	//  └─ session3
	cmdFile := fmt.Sprintf("%d.pi", session)
	// mkdir
	curDir := fmt.Sprintf("compile_file/%d", session)
	os.Mkdir(curDir, os.ModePerm)
	// write on file
	fileName := fmt.Sprintf("compile_file/%d/%s", session, cmdFile)

	err := ioutil.WriteFile(fileName, []byte(src), os.ModePerm)
	if err != nil {
		log.Fatalf("file write err:%s", err)
	}

	file2Name := fmt.Sprintf("compile_file/%d/Kagari.toml", session)
	tomlSrc := fmt.Sprintf("entry = '%s'\nproject = '%d'", cmdFile, session)
	err = ioutil.WriteFile(file2Name, []byte(tomlSrc), os.ModePerm)
	if err != nil {
		log.Fatalf("file write err:%s", err)
	}

	// defer 回调，删除文件
	defer func(curDir string) {
		cmd := exec.Command("rm", "-rf", curDir)
		if err := cmd.Start(); err != nil { //开始执行命令
			log.Fatalf("file del err:%s", err)
			return
		}
	}(curDir)

	// 编译运行，返回结果
	complieCmd := exec.Command("plc", cmdFile)
	complieCmd.Dir = curDir

	compileOut, err := complieCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s failed with %s", complieCmd, err)
	}

	runCmd := exec.Command(fmt.Sprintf("./out"))
	runCmd.Dir = curDir
	runOut, err := runCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s failed with %s", runCmd, err)
	}

	cache.Add(session, runOut)
	return string(compileOut), string(runOut), session

}
