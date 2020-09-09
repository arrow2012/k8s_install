package addons

import (
	"k8s_install/common/config"
	"k8s_install/common/utils"
	"k8s_install/ssh"
	"os"
	"fmt"
	"strconv"
)

var (
	log = config.Logger
)

func DeplyAddons()  {

	os.Setenv("PodCIDR",config.K8s_PodCIDR)
	os.Setenv("KUBE_VERSION",config.K8s_kube_version)
	os.Setenv("KUBE_APISERVER","https://"+config.K8s_master_host[0]+":6443")
	os.Setenv("CLUSTER_NAME",config.K8s_cluster_name)
	os.Setenv("ClusterDns",config.K8s_ClusterDns)
	os.Setenv("ClusterDomain",config.K8s_ClusterDomain)
	os.Setenv("kube_proxy_image",config.Kubeproxy_Image)
	os.Setenv("metrics_image",config.Metrics_Image)
	os.Setenv("metrics_image",config.Metrics_Image)

	//generate_kubeproxyconf :=fmt.Sprintf("/bin/bash generate_kubeproxy_yml.sh")
	//_,err := utils.ExecCmd(generate_kubeproxyconf,"addons/",os.Environ())
	//if err != nil { log.Fatal(err.Error()) }

	log.Info("执行脚本 files/kubeconfig.sh")
	generate_kubeconfig :=fmt.Sprintf("/bin/bash files/kubeconfig.sh ${KUBE_APISERVER}")
	_,err := utils.ExecCmd(generate_kubeconfig,"addons/",os.Environ())
	if err != nil { log.Fatal(err.Error()) }

	m1host := ssh.K8sMasterHost[0]
	//putfilecmd := fmt.Sprintf("/usr/bin/rsync -avh /etc/kubernetes/kube-proxy.yml %s:/etc/kubernetes/addons/",m1host.Host)
	//_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
	//if err != nil { log.Fatal(err.Error()) }
	//m1host.CmdList = "kubectl apply -f /etc/kubernetes/addons/kube-proxy.yml"
	//ssh.SSHcommand(m1host)

	//os.Setenv("addon_dir",config.Kube_addons_config_dir)
	//os.Setenv("flannel_img",config.Flannel_image)
	//generate_flannel_yml :=fmt.Sprintf("/bin/bash generate_flannel_yml.sh")
	//_,err = utils.ExecCmd(generate_flannel_yml,"addons/",os.Environ())
	//if err != nil { log.Fatal(err.Error()) }

	//putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh ${addon_dir}/kube-flannel.yml %s:${addon_dir}/kube-flannel.yml",m1host.Host)
	//_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
	//if err != nil { log.Fatal(err.Error()) }
	//m1host.CmdList = fmt.Sprintf("kubectl apply -f %s/kube-flannel.yml",config.Kube_addons_config_dir)
	//ssh.SSHcommand(m1host)


	log.Info("执行脚本 generate_coredns_conf")
	//coredns server
	generate_coredns_yml :=fmt.Sprintf("/bin/bash generate_coredns_conf.sh")
	_,err = utils.ExecCmd(generate_coredns_yml,"addons/",os.Environ())
	if err != nil { log.Fatal(err.Error()) }

	putfilecmd := fmt.Sprintf("/usr/bin/rsync -avh /tmp/coredns.yml %s:/etc/kubernetes/addons/",m1host.Host)
	_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
	if err != nil { log.Fatal(err.Error()) }
	m1host.CmdList = "kubectl apply -f /etc/kubernetes/addons/coredns.yml"
	ssh.SSHcommand(m1host)

	//metrics server
	putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /tmp/metrics_server.yml %s:/etc/kubernetes/addons/",m1host.Host)
	_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
	if err != nil { log.Fatal(err.Error()) }
	m1host.CmdList = "kubectl apply -f /etc/kubernetes/addons/metrics_server.yml"
	ssh.SSHcommand(m1host)


	//生成flanneld的kubeconfig

	generate_flanneld_conf :=fmt.Sprintf("/bin/bash flannelkubeconfig.sh ${KUBE_APISERVER}")
	_,err = utils.ExecCmd(generate_flanneld_conf,"addons/files/",os.Environ())
	if err != nil { log.Fatal(err.Error()) }


	for _,i :=range ssh.K8sAllHost{
		os.Setenv("nodename",i.Host)
		os.Setenv("flannel_healthzport",strconv.Itoa(config.Flannel_healthzPort))
		//发送kube-proxy文件
		putfilecmd := fmt.Sprintf("/usr/bin/rsync -avh /usr/local/bin/kube-proxy %s:/usr/local/bin/",i.Host)
		_,err := utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		//发送kube-proxy.kubeconfig
		putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /etc/kubernetes/kube-proxy.kubeconfig %s:/etc/kubernetes/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { log.Fatal(err.Error()) }


		//生成kube-proxy.service 并发送

		generate_kubeproxyservice :=fmt.Sprintf("/bin/bash generate_kubeproxy_service.sh")
		_,err = utils.ExecCmd(generate_kubeproxyservice,"addons/",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /tmp/kube-proxy.service %s:/lib/systemd/system/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		//生成kube-proxy.conf 并发送
		generate_kubeproxyconf :=fmt.Sprintf("/bin/bash generate_kubeproxy_conf.sh")
		_,err = utils.ExecCmd(generate_kubeproxyconf,"addons/",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		//发送kube-proxy.kubeconfig
		putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /tmp/kube-proxy.conf %s:/etc/kubernetes/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		i.CmdList = "mkdir -p {/etc/cni/net.d,/run/flannel,/etc/kube-flannel/}"
		ssh.SSHcommand(i)

		putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /usr/local/bin/flanneld %s:/usr/local/bin/flanneld",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		//发送flanneld配置
		putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /etc/kubernetes/flanneld.kubeconfig %s:/etc/kubernetes/flanneld.kubeconfig",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }



		//生成flannel.service 并发送

		generate_flannel_service :=fmt.Sprintf("/bin/bash generate_flannel_service.sh")
		_,err = utils.ExecCmd(generate_flannel_service,"addons/",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /tmp/flanneld.service %s:/lib/systemd/system/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /tmp/10-flannel.conflist %s:/etc/cni/net.d/",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /tmp/net-conf.json %s:/etc/kube-flannel",i.Host)
		_,err = utils.ExecCmd(putfilecmd,"",os.Environ())
		if err != nil { log.Fatal(err.Error()) }

		//启动flannel
		i.CmdList = "systemctl daemon-reload && systemctl restart flanneld && sleep 5 && systemctl status flanneld |grep Active |grep running"
		ssh.SSHcommand(i)

		//启动kube-proxy
		i.CmdList = "systemctl daemon-reload && systemctl restart kube-proxy && sleep 5 && systemctl status kube-proxy |grep Active |grep running"
		ssh.SSHcommand(i)
	}

}