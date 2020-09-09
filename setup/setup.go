package setup

import (
	"k8s_install/common/config"
	"k8s_install/ssh"
	"fmt"
)

func Setup() {
	// 发送证书到远程服务器
	for _,i :=range	ssh.K8sMasterHost{
		i.CmdList = fmt.Sprintf("rm -fr %s %s",config.Etcd_dataDir,config.Kube_config_dir)
		ssh.SSHcommand(i)
	}
	for _,i :=range	ssh.K8snodeHost{
		i.CmdList = fmt.Sprintf("rm -fr %s",config.Kube_config_dir)
		ssh.SSHcommand(i)
	}
	for _,i :=range	ssh.K8sAllHost{
		i.CmdList = fmt.Sprintf("mkdir -p %s",config.K8s_certDir)
		ssh.SSHcommand(i)
		i.CmdList = fmt.Sprintf("ifconfig cni0 down && ip link delete cni0")
		ssh.SSHcommand(i)
	}
}
