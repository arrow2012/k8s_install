package main

import (
	"fmt"
	"io"
	"k8s_install/addons"
	"k8s_install/bootstrap"
	"k8s_install/cert"
	"k8s_install/common/config"
	"k8s_install/etcd"
	"k8s_install/getBin"
	"k8s_install/master"
	"k8s_install/node"
	"k8s_install/setup"
	"os"
)


var (
	 log = config.Logger
)

func main() {
	//if runtime.GOOS == "darwin" {
	//	copy("file/bin/darwin/cfssl","/usr/local/bin/cfssl")
	//	copy("file/bin/darwin/cfssljson","/usr/local/bin/cfssljson")
	//}
	//if runtime.GOOS == "linux" {
	//	copy("file/bin/linux/cfssl","/usr/local/bin/cfssl")
	//	copy("file/bin/linux/cfssljson","/usr/local/bin/cfssljson")
	//}
	//
	//
	////下载解压二进制文件
	setup.Setup()
	getBin.Task()
	//生成证书
	cert.Task()
	//部署etcd集群
	etcd.Task()
	master.GenerateMasterconf()
	master.InitMaster()
	master.HealthNodePort()
	bootstrap.ExecBootstrap()
	node.DeployNode()
	addons.DeplyAddons()
	log.Info("部署完成")
}


func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}