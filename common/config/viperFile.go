package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitViper(path, file, t string) (v *viper.Viper) {
	var runtime_viper = viper.New()
	runtime_viper.SetConfigType(t)
	runtime_viper.AddConfigPath(path)
	runtime_viper.SetConfigName(file)
	err := runtime_viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	return runtime_viper
}

var (
	Ver      = InitViper("conf", "cfg", "toml")

	//K8s_MasterHost  = Ver.GetStringSlice("k8s.master_host")



	//local
	Local_certDir=Ver.GetString("local.cert_dir")
	Local_binDir=Ver.GetString("local.bin_dir")
	SSH_key  = Ver.GetString("local.key")

	//ETCD
	Etcd_nodeName=Ver.GetString("etcd.node_name_prefix")
	Etcd_dataDir = Ver.GetString("etcd.data_dir")
	Etcd_confDir = Ver.GetString("etcd.conf_dir")
	Etcd_confFile = Ver.GetString("etcd.conf_file")
	Etcd_clusterIp = Ver.GetStringSlice("etcd.cluster_ip")
	Etcd_host_sshport = Ver.GetInt("etcd.ssh_port")
	Etcd_host_sshuser = Ver.GetString("etcd.ssh_username")
	Etcd_certDir= Ver.GetString("etcd.cert_dir")
	Etcd_clientsh = Ver.GetString("etcd.client_sh")

	//k8s
	K8s_kube_version = Ver.GetString("k8s.kube_version")
	K8s_cluster_name  = Ver.GetString("k8s.cluster_name")
	K8s_master_host=Ver.GetStringSlice("k8s.master_host")
	K8s_node_host=Ver.GetStringSlice("k8s.node_host")
	K8s_host_sshport = Ver.GetInt("k8s.ssh_port")
	K8s_host_sshuser = Ver.GetString("k8s.ssh_username")
	K8s_certDir = Ver.GetString("k8s.cert_dir")
	K8s_apiserver_service_unit = Ver.GetString("k8s.apiserver_service_unit")
	K8s_controllermanager_service_unit = Ver.GetString("k8s.controllermanager_service_unit")
	K8s_scheduler_service_unit = Ver.GetString("k8s.scheduler_service_unit")

	K8s_SvcCIDR = Ver.GetString("k8s.SvcCIDR")
	K8s_PodCIDR = Ver.GetString("k8s.PodCIDR")

	K8s_ClusterDns = Ver.GetString("k8s.ClusterDns")
	K8s_ClusterDomain = Ver.GetString("k8s.ClusterDomain")
	K8s_ServiceNodePortRange =Ver.GetString("k8s.ServiceNodePortRange")
	K8s_ServiceUnitDir = Ver.GetString("k8s.service_unit_dir")


	K8s_apiserver_logdir= Ver.GetString("k8s.apiserver_logdir")
	K8s_scheduler_logdir= Ver.GetString("k8s.scheduler_logdir")
	K8s_controllerManger_logdir= Ver.GetString("k8s.controllerManger_logdir")



	Ca_cert = Ver.GetString("k8s.ca_cert")
	Cert_dir = Ver.GetString("k8s.cert_dir")
	Kube_config_dir = Ver.GetString("k8s.kube_config_dir")
	Kube_addons_config_dir = Ver.GetString("k8s.kube_addons_config_dir")

	Cni_binary = Ver.GetString("docker.cni_binary")
	Pause_image = Ver.GetString("docker.pause_image")
	Flannel_image = Ver.GetString("flannel.flannel_image")
	Flannel_healthzPort = Ver.GetInt("flannel.healthzPort")

	Kube_Image  = Ver.GetString("docker.kube_image")
	Etcd_Image  = Ver.GetString("docker.etcd_image")
	Kubeproxy_Image  = Ver.GetString("docker.kube_proxy_image")
	Metrics_Image = Ver.GetString("docker.metrics_image")
	CoreDns_Image = Ver.GetString("docker.coredns_image")

)
