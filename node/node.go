package node

import (
	"fmt"
	"k8s_install/common/config"
	"k8s_install/common/utils"
	"k8s_install/ssh"
	"os"
	"time"
)

var (
	log = config.Logger
)

func DeployNode() {

	for _, i := range ssh.K8sAllHost {
		i.CmdList = "mkdir -p {/etc/kubernetes/pki,/etc/kubernetes/manifests,/var/lib/kubelet/,/opt/cni/bin/,/var/log/kubernetes/kubelet}"
		ssh.SSHcommand(i)

		//发送CNI文件
		putfilecmd :=fmt.Sprintf("/usr/bin/rsync -avh /usr/local/src/%s %s:/usr/local/src/",config.Cni_binary,i.Host)
		_,err := utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		//发送kubelet文件
		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh /usr/local/bin/kubelet %s:/usr/local/bin/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { log.Fatal(err.Error()) }


		//生成kubelet.service文件并且发送
		os.Setenv("node_host",i.Host)
		os.Setenv("pause_image",config.Pause_image)
		os.Setenv("ClusterDns",config.K8s_ClusterDns)
		os.Setenv("ClusterDomain",config.K8s_ClusterDomain)
		generate_kubelet_service := fmt.Sprintf("/bin/bash generate_kubelet_service.sh")
		_,err = utils.ExecCmd(generate_kubelet_service,"node/",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh kubelet.service %s:/lib/systemd/system/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"node/",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh kubelet-conf.yml %s:/etc/kubernetes/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"node/",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		//发送证书
		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh /etc/kubernetes/bootstrap.kubeconfig %s:/etc/kubernetes/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh /etc/kubernetes/pki/ca.pem %s:/etc/kubernetes/pki/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		//启动kubelet
		i.CmdList = "systemctl daemon-reload && systemctl restart kubelet && sleep 5 && systemctl status kubelet |grep Active |grep running"
		ssh.SSHcommand(i)
	}
	time.Sleep(5*time.Second)
	// master打上污点
	for _,i:=range  config.K8s_master_host{
		taintNodeCmd :=fmt.Sprintf("kubectl  taint nodes %s node-role.kubernetes.io/master=\"\":NoSchedule --overwrite",i)
		_,err := utils.ExecCmd(taintNodeCmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }
	}
}

