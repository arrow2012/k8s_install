package master

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"k8s_install/common/config"
	"k8s_install/common/utils"
	"k8s_install/etcd"
	"k8s_install/ssh"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)


var (
	log = config.Logger
)

func InitMaster()  {

	putfilecmd :=fmt.Sprintf("/usr/bin/rsync -avh healthz-rbac.yml /etc/kubernetes/")
	_,err := utils.ExecCmd(putfilecmd,"master/template/",os.Environ())
	if err != nil { log.Error(err.Error()) }



	for _,i :=range ssh.K8sMasterHost{
		i.CmdList = fmt.Sprintf("mkdir -p {%s,~/.kube}",config.K8s_certDir)
		ssh.SSHcommand(i)

		putfilecmd :=fmt.Sprintf("/usr/bin/rsync -avh /usr/local/bin/kube* %s:/usr/local/bin/",i.Host)
		_,err := utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { panic(err) }

		i.CmdList = fmt.Sprintf("[ ! -d \"%s\" ] && mkdir -p %s || echo \"目录%s已经存在\"",config.K8s_apiserver_logdir,config.K8s_apiserver_logdir,config.K8s_apiserver_logdir)
		ssh.SSHcommand(i)
		i.CmdList = fmt.Sprintf("[ ! -d \"%s\" ] && mkdir -p %s || echo \"目录%s已经存在\"",config.K8s_scheduler_logdir,config.K8s_scheduler_logdir,config.K8s_scheduler_logdir)
		ssh.SSHcommand(i)
		i.CmdList = fmt.Sprintf("[ ! -d \"%s\" ] && mkdir -p %s || echo \"目录%s已经存在\"",config.K8s_controllerManger_logdir,config.K8s_controllerManger_logdir,config.K8s_controllerManger_logdir)
		ssh.SSHcommand(i)

		for _,f :=range []string{
			"admin.kubeconfig",
			"scheduler.kubeconfig",
			"controller-manager.kubeconfig",
			"healthz-rbac.yml"}{
			putfilecmd := fmt.Sprintf("/usr/bin/rsync -avh /etc/kubernetes/%s %s:/etc/kubernetes/",f,i.Host)
			_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
			if err != nil { panic(err) }

			if f == "admin.kubeconfig"{
				putfilecmd :=fmt.Sprintf("/usr/bin/rsync -avh %s/%s ~/.kube/config",config.Kube_config_dir,f)
				_,err := utils.ExecCmd(putfilecmd,"./",os.Environ())
				if err != nil { panic(err) }
				putfilecmd = fmt.Sprintf("/usr/bin/rsync -avh /etc/kubernetes/%s %s:~/.kube/config",f,i.Host)
				_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
				if err != nil { panic(err) }
			}
		}




		//启动k8s
		i.CmdList =fmt.Sprintf("systemctl daemon-reload && systemctl enable kube-apiserver.service kube-scheduler.service  kube-controller-manager.service  && systemctl restart kube-apiserver.service kube-scheduler.service  kube-controller-manager.service ")
		ssh.SSHcommand(i)
	}
}

func GenerateMasterconf() {
	type MasterConf struct {
		Host string
		EtcdServers string
		SvcCIDR         string
		ServiceNodePortRange         	string
		PodCIDR		string
	}

	K8s_apiserver_service_unit, err := os.OpenFile(config.K8s_apiserver_service_unit, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		utils.CheckErrExit(err)
	}

	K8s_controllermanager_service_unit, err := os.OpenFile(config.K8s_controllermanager_service_unit, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		utils.CheckErrExit(err)
	}
	K8s_scheduler_service_unit, err := os.OpenFile(config.K8s_scheduler_service_unit, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		utils.CheckErrExit(err)
	}



	kubeep, err := os.OpenFile("kube-ep.yml", os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		utils.CheckErrExit(err)
	}



	var mconf []MasterConf

	for _,host :=range config.K8s_master_host {
		ec := MasterConf {
			Host:     host,
			EtcdServers     : etcd.GetEtcdEndpointsStr(),
			SvcCIDR:	config.K8s_SvcCIDR,
			ServiceNodePortRange: config.K8s_ServiceNodePortRange,
			PodCIDR: config.K8s_PodCIDR,
		}
		mconf = append(mconf, ec)
	}
	t, err := template.ParseGlob("master/template/*.tpl")
	if err != nil { panic(err) }

	for _,i :=range mconf {
		os.Truncate(config.K8s_apiserver_service_unit, 0)
		if err !=nil{
			log.Error(err.Error())
		}
		fmt.Println("---------------")
		fmt.Println(i.Host,i.SvcCIDR)
		fmt.Println("---------------")
		err = t.ExecuteTemplate(K8s_apiserver_service_unit,"kube-apiserver.service.tpl",i)
		if err != nil { panic(err) }

		putfilecmd :=fmt.Sprintf("/usr/bin/rsync -avh %s %s:%s/%s",config.K8s_apiserver_service_unit,i.Host,config.K8s_ServiceUnitDir,config.K8s_apiserver_service_unit)
		_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { panic(err) }

		//生成 scheduler
		os.Truncate(config.K8s_scheduler_service_unit, 0)
		if err !=nil{
			log.Error(err.Error())
		}
		err = t.ExecuteTemplate(K8s_scheduler_service_unit,"kube-scheduler.service.tpl",i)
		if err != nil { panic(err) }

		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh %s %s:%s/%s",config.K8s_scheduler_service_unit,i.Host,config.K8s_ServiceUnitDir,config.K8s_scheduler_service_unit)
		_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { panic(err) }

		//生成 controller-manager
		os.Truncate(config.K8s_controllermanager_service_unit, 0)
		if err !=nil{
			log.Error(err.Error())
		}
		err = t.ExecuteTemplate(K8s_controllermanager_service_unit,"kube-controller-manager.service.tpl",i)
		if err != nil { panic(err) }

		putfilecmd =fmt.Sprintf("/usr/bin/rsync -avh %s %s:%s/%s",config.K8s_controllermanager_service_unit,i.Host,config.K8s_ServiceUnitDir, config.K8s_controllermanager_service_unit)
		_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { panic(err) }
	}

	//生成 kube-ep.yml

	type Master struct {
		Host []string
	}
	m := &Master{
		Host: config.K8s_master_host,
	}
	os.Truncate("kube-ep.yml", 0)
	if err !=nil{
		log.Error(err.Error())
	}
	err = t.ExecuteTemplate(kubeep,"kube-ep.yml.tpl",m)
	if err != nil { panic(err) }

	putfilecmd :=fmt.Sprintf("/usr/bin/rsync -avh kube-ep.yml %s/kube-ep.yml",config.Kube_config_dir)
	_,err = utils.ExecCmd(putfilecmd,"./",os.Environ())
	if err != nil { panic(err) }

	os.Remove(config.Etcd_confFile)
	if err !=nil{
		log.Error(err.Error())
	}

	os.Remove(config.Etcd_clientsh)
	if err !=nil{
		log.Error(err.Error())
	}
}


func CheckApiServerHealth() bool {
	url := fmt.Sprintf("https://192.168.120.93:6443/healthz")
	method := "GET"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //InsecureSkipVerify参数值只能在客户端上设置有效
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, nil)
	for i:= 0;i<10;i++ {
		time.Sleep(5*time.Second)
		fmt.Println(fmt.Sprintf("开始第%d 次检测接口 %s ",i+1,url))
		if err != nil {
			continue
		}
		res, err := client.Do(req)
		if err != nil {
			continue
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if strings.Contains(string(body),"ok"){
			return true
		}else {
			continue
		}
	}
	return false

}

func HealthNodePort()  {

	if !CheckApiServerHealth() {
		log.Fatal("Api server 检测失败")
	}


	ApiserverReady := false
	for i:= 0;i<10;i++ {
		time.Sleep(5*time.Second)
		cmd := "kubectl get cs"
		fmt.Println(fmt.Sprintf("开始第%d 次检测 %s ",i+1,cmd))

		content,err := utils.ExecCmd(cmd,"",os.Environ())
		if err !=nil{
			continue
		}
		if !strings.Contains(content,"Unhealthy") {
			ApiserverReady = true
			break
		}
	}

	log.Info("cs 准备OK")
	if ApiserverReady == true {
		putfilecmd :=fmt.Sprintf("kubectl apply -f %s/kube-ep.yml && kubectl apply -f %s/healthz-rbac.yml || echo 执行失败",config.Kube_config_dir,config.Kube_config_dir)
		_ ,err := utils.ExecCmd(putfilecmd,"./",os.Environ())
		if err != nil { panic(err) }
	}else {
		log.Fatal("Api server not ready")
	}



}