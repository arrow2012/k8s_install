package bootstrap

import (
	"k8s_install/common/config"
	"k8s_install/common/utils"

	"os"
)

var (
	log = config.Logger
)

func ExecBootstrap() {
	os.Setenv("KUBE_APISERVER","https://"+config.K8s_master_host[0]+":6443")
	os.Setenv("CLUSTER_NAME",config.K8s_cluster_name)
	_,err := utils.ExecCmd("/bin/sh bootstrap.sh ${KUBE_APISERVER}","bootstrap/",os.Environ())
	if err !=nil {log.Fatal(err.Error())}
}
