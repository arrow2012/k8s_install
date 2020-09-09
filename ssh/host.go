package ssh

import "k8s_install/common/config"


func GetClusterHosts(host []string,port int,user,key string) []*SSHHost {
	var hosts []*SSHHost
	for _,i :=range host{
		s := &SSHHost{
			Host: i,
			Port: port,
			Username: user,
			Key: key,
		}
		hosts = append(hosts, s)
	}
	return hosts
}

var (
	 K8sMasterHost  = GetClusterHosts(config.K8s_master_host,config.K8s_host_sshport,config.K8s_host_sshuser,config.SSH_key)
	 K8snodeHost  = GetClusterHosts(config.K8s_node_host,config.K8s_host_sshport,config.K8s_host_sshuser,config.SSH_key)
	 EtcdHost  = GetClusterHosts(config.Etcd_clusterIp,config.Etcd_host_sshport,config.Etcd_host_sshuser,config.SSH_key)
	 K8sAllHost = append(K8sMasterHost,K8snodeHost...)
)