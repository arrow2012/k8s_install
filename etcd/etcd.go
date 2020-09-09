package etcd

import (
	"fmt"
	"k8s_install/common/config"
	"k8s_install/common/utils"
	"k8s_install/ssh"
	"os"
	"strconv"
	"strings"
	"text/template"

)

var (
	log = config.Logger
)


func initEtcd()  {
	for _,i :=range ssh.EtcdHost{
		i.CmdList = fmt.Sprintf("mkdir -p %s && chmod 700 %s",config.Etcd_dataDir,config.Etcd_dataDir)
		ssh.SSHcommand(i)

		putfilecmd :=fmt.Sprintf("/usr/bin/rsync -avh \"/usr/local/bin/etcd\" %s:/usr/local/bin/",i.Host)
		_,err := utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { panic(err) }
		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh \"/usr/local/bin/etcdctl\" %s:/usr/local/bin/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { panic(err) }

		fmt.Println(i.Host)
		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh etcd.service %s:/usr/lib/systemd/system/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"etcd/template/",os.Environ())
		if err != nil { panic(err) }

		//启动ETCD
		i.CmdList =fmt.Sprintf("source /etc/profile.d/etcd.sh && systemctl daemon-reload && systemctl enable etcd && systemctl restart etcd")
		ssh.SSHcommand(i)
	}
}

func GetEtcdEndpointsStr() string {
	var ClusterStr []string
	for _,host :=range config.Etcd_clusterIp {
		ClusterStr = append(ClusterStr, fmt.Sprintf("https://%s:2379",host))
	}
	return strings.Join(ClusterStr, ",")
}


func GenerateEtcdConf() {
	type EtcdConfig struct {
		Name        	string
		NodeName        string
		DataDir         string
		Waldir         	string
		Host     		string
		Init_cluster 		string
		Cluster_endpoints string
	}

	confFile, err := os.OpenFile(config.Etcd_confFile, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		utils.CheckErrExit(err)
	}

	clientshFile, err := os.OpenFile(config.Etcd_clientsh, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		utils.CheckErrExit(err)
	}


	var etcdconf []EtcdConfig
	var InitCluster []string
	var ClusterStr []string

	for index,host :=range config.Etcd_clusterIp {
		ec := EtcdConfig{
			NodeName:     config.Etcd_nodeName+strconv.Itoa(index+1),
			DataDir     : config.Etcd_dataDir,
			Waldir : config.Etcd_dataDir+"/wal",
			Host: host,
		}
		etcdconf = append(etcdconf, ec)
		InitCluster = append(InitCluster, fmt.Sprintf("%s=https://%s:2380",ec.NodeName,ec.Host))
		ClusterStr = append(ClusterStr, fmt.Sprintf("https://%s:2379",ec.Host))
	}

	InitClusterString := strings.Join(InitCluster, ",")
	endpointsStr := strings.Join(ClusterStr, ",")

	t, err := template.ParseGlob("etcd/template/*.tpl")
	if err != nil { panic(err) }

	GeClientsh := false

	for _,i :=range etcdconf {
		i.Init_cluster = InitClusterString
		os.Truncate(config.Etcd_confFile, 0)
		if err !=nil{
			log.Error(err.Error())
		}
		i.Cluster_endpoints = endpointsStr

		err = t.ExecuteTemplate(confFile,"etcd.config.yml.tpl",i)
		if err != nil { panic(err) }

		putfilecmd :=fmt.Sprintf("/usr/bin/rsync -a --rsync-path=\"mkdir -p %s && rsync\"  %s %s:%s/%s",config.Etcd_confDir,config.Etcd_confFile,i.Host,config.Etcd_confDir,config.Etcd_confFile)
		_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { panic(err) }


		if GeClientsh == false{
			err = t.ExecuteTemplate(clientshFile,"etcd.sh.tpl",i)
			if err != nil { panic(err) }
			GeClientsh = true
		}

		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh %s %s:/etc/profile.d/",config.Etcd_clientsh,i.Host)
		_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { panic(err) }
	}

	os.Remove(config.Etcd_confFile)
	if err !=nil{
		log.Error(err.Error())
	}
	os.Remove(config.Etcd_clientsh)
	if err !=nil{
		log.Error(err.Error())
	}
}


func Task()  {
	GenerateEtcdConf()
	initEtcd()
}
