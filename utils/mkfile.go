package utils

import (
	"io"
	"os"
	"fmt"
	logger "k8s_install/log"

)


var (
	log = logger.MyLogger
)

func CreateDir(dir ...string)  {
	for _,v :=range dir{
		err :=os.MkdirAll(v,os.ModePerm)
		if err !=nil{
			fmt.Println(err)
		}
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func CheckErr(err error) {
	if err != nil {
		log.Error(err.Error())
	}
}

func CheckErrExit(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}


func WriteStrToFile(fileName, s string) {
	var f *os.File
	var err1 error
	if checkFileIsExist(fileName) { //如果文件存在
		log.Warn("删除文件")
		err1 = os.RemoveAll(fileName)
		if err1 != nil {
			log.Warn(fmt.Sprintf("删除错误: %s", err1.Error()))
		}
		log.Warn("创建文件")
		f, err1 = os.Create(fileName) //创建文件
	} else {
		log.Warn("文件不存在")
		f, err1 = os.Create(fileName) //创建文件
	}
	CheckErr(err1)
	_, err1 = io.WriteString(f, s) //写入文件(字符串)
	CheckErr(err1)
}
