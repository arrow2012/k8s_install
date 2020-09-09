package cert

import (
	"fmt"
	"k8s_install/common/config"
	"k8s_install/common/utils"
	"os"
	"path/filepath"
)

type KubeConfig struct {
	KUBE_USER string
	KUBE_CERT string
	KUBE_CERT_KEY string
	KUBE_CONFIG string
}
var (
	log = config.Logger
)




func GenerateKubeConfig()  {
	k1:=&KubeConfig{
		"system:kube-controller-manager",
		config.Cert_dir+"sa.pem",
		config.Cert_dir+"sa-key.pem",
		config.Kube_config_dir+"controller-manager.kubeconfig",
	}
	k2:=&KubeConfig{
		"system:kube-scheduler",
		config.Cert_dir+"kube-scheduler.pem",
		config.Cert_dir+"kube-scheduler-key.pem",
		config.Kube_config_dir+"scheduler.kubeconfig",
	}
	k3:=&KubeConfig{
		"kubernetes-admin",
		config.Cert_dir+"admin.pem",
		config.Cert_dir+"admin-key.pem",
		config.Kube_config_dir+"admin.kubeconfig",
	}
	var arr = []*KubeConfig{
		k1,
		k2,
		k3,
	}
	for _,i :=range arr{
		cmd := fmt.Sprintf("kubectl config set-cluster %s --certificate-authority=%s --embed-certs=true --server=https://%s:6443 --kubeconfig=%s",config.K8s_cluster_name,filepath.Join(config.Cert_dir,config.Ca_cert), config.K8s_master_host[0],i.KUBE_CONFIG)
		fmt.Println(fmt.Sprintf("执行命令 %s",cmd))
		_,err := utils.ExecCmd(cmd,"",os.Environ())
		utils.CheckErrExit(err)

		cmd = fmt.Sprintf("kubectl config set-credentials %s --client-certificate=%s --client-key=%s --embed-certs=true --kubeconfig=%s ",i.KUBE_USER,i.KUBE_CERT,i.KUBE_CERT_KEY,i.KUBE_CONFIG)
		fmt.Println(fmt.Sprintf("执行命令 %s",cmd))
		_,err = utils.ExecCmd(cmd,"",os.Environ())
		utils.CheckErrExit(err)

		cmd = fmt.Sprintf("kubectl config set-context %s@%s --cluster=%s --user=%s --kubeconfig=%s",i.KUBE_USER,config.K8s_cluster_name,config.K8s_cluster_name,i.KUBE_USER,i.KUBE_CONFIG)
		fmt.Println(fmt.Sprintf("执行命令 %s",cmd))
		_,err = utils.ExecCmd(cmd,"",os.Environ())
		utils.CheckErrExit(err)

		cmd = fmt.Sprintf("kubectl config use-context %s@%s --kubeconfig=%s",i.KUBE_USER,config.K8s_cluster_name,i.KUBE_CONFIG)
		fmt.Println(fmt.Sprintf("执行命令 %s",cmd))
		_,err = utils.ExecCmd(cmd,"",os.Environ())
		utils.CheckErrExit(err)
	}


}