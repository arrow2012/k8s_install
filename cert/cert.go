package cert

import (
	"fmt"
	"k8s_install/common/config"
	"k8s_install/tlsutil"
	"k8s_install/common/utils"
	"os"
)
var (
	local_certdir = config.Ver.GetString("local.cert_dir")
)



func GenerateCertByCfssl()  {
	utils.ExecCmd(fmt.Sprintf("/bin/bash tls.sh %s",local_certdir),"cert/",nil)

	// 发送证书到远程服务器
	for _,i :=range	config.Etcd_clusterIp{
		os.Setenv("remote_ip",i)
		os.Setenv("cert_dir",config.Etcd_certDir)
		fmt.Println(os.Getenv("remote_ip"))
		utils.ExecCmd("/usr/bin/rsync -avpz --delete ./ ${remote_ip}:${cert_dir}/",config.Local_certDir,os.Environ())
	}

	for _,i :=range	config.K8s_master_host{
		os.Setenv("remote_ip",i)
		os.Setenv("cert_dir",config.Etcd_certDir)
		fmt.Println(os.Getenv("remote_ip"))
		utils.ExecCmd("/usr/bin/rsync -avpz --delete ./ ${remote_ip}:${cert_dir}/",config.Local_certDir,os.Environ())
	}

	os.Unsetenv("remote_ip")
	os.Unsetenv("cert_dir")
}


func GenerateCAFile()  {
	//生成CA私钥
	sn, err := tlsutil.GenerateSerialNumber()
	if err !=nil{
		fmt.Print(err)
	}
	signer,str, err := tlsutil.GeneratePrivateKey()
	if err !=nil{
		fmt.Print(err)
	}
	utils.WriteStrToFile("file/cert/cert-key.pem",str)

	//创建CA证书
	ca, err := tlsutil.GenerateCA(signer, sn, 36500, nil)
	if err !=nil{
		fmt.Print(err)
	}
	utils.WriteStrToFile("file/cert/cert.pem",ca)
}

func Task()  {

	GenerateCertByCfssl()
	GenerateKubeConfig()

}

