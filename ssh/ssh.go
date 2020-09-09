package ssh

import (
	"bytes"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"k8s_install/common/config"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	log = config.Logger
)


type SSHResult struct {
	Host    string
	Success bool
	Result  string
}

type SSHHost struct {
	Host      string
	Port      int
	Username  string
	Password  string
	CmdList   string
	Key       string
}




type sshhost SSHHost

func connect(user, password, host, key string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	if key == "" {
		auth = append(auth, ssh.Password(password))
	} else {
		pemBytes, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}

		var signer ssh.Signer

		if password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
		}

		if err != nil {
			return nil, err
		}

		auth = append(auth, ssh.PublicKeys(signer))
	}

	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},

	}

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}


	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	return session, nil
}

func Dossh(username, password, host, key string, cmdlist []string, port, timeout int,  linuxMode bool) {
	chSSH := make(chan SSHResult)
	if linuxMode {
		go dossh_run(username, password, host, key, cmdlist, port, chSSH)
	} else {
		go dossh_session(username, password, host, key, cmdlist, port,  chSSH)
	}
	var res SSHResult

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		res.Host = host
		res.Success = false
		res.Result = ("SSH run timeout：" + strconv.Itoa(timeout) + " second.")
	}
	return
}

func dossh_session(username, password, host, key string, cmdlist []string, port int,  ch chan SSHResult) {
	session, err := connect(username, password, host, key, port)
	var sshResult SSHResult
	sshResult.Host = host

	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>", err.Error())
		ch <- sshResult
		return
	}
	defer session.Close()

	cmdlist = append(cmdlist, "exit")

	stdinBuf, _ := session.StdinPipe()

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt

	session.Stderr = &errbt
	err = session.Shell()
	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>", err.Error())
		ch <- sshResult
		return
	}
	for _, c := range cmdlist {
		c = c + "\n"
		stdinBuf.Write([]byte(c))
	}
	session.Wait()
	if errbt.String() != "" {
		sshResult.Success = false
		sshResult.Result = errbt.String()
		ch <- sshResult
	} else {
		sshResult.Success = true
		sshResult.Result = outbt.String()
		ch <- sshResult
	}

	return
}

func dossh_run(username, password, host, key string, cmdlist []string, port int,  ch chan SSHResult) {

	session, err := connect(username, password, host, key, port)
	var sshResult SSHResult
	sshResult.Host = host

	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>", err.Error())
		ch <- sshResult
		return
	}
	defer session.Close()


	cmdlist = append(cmdlist, "exit")
	newcmd := strings.Join(cmdlist, "&&")

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt

	session.Stderr = &errbt
	err = session.Run(newcmd)
	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>", err.Error())
		ch <- sshResult
		return
	}

	if errbt.String() != "" {
		sshResult.Success = false
		sshResult.Result = errbt.String()
		ch <- sshResult
	} else {
		sshResult.Success = true
		sshResult.Result = outbt.String()
		ch <- sshResult
	}

	return
}

func PutFile(host string,port int,username,key,src,dest string) error {
	sshHost := host
	sshUser := username
	sshPort := port
	sshKeyPath :=key

	fmt.Println(sshHost)
	fmt.Println(sshUser)
	fmt.Println(sshPort)
	fmt.Println(sshKeyPath)
	fmt.Println(dest)

	sshConfig, err := auth.PrivateKey(sshUser, sshKeyPath, ssh.InsecureIgnoreHostKey())
	if err !=nil{
		return err
	}

	scpClient := scp.NewClient(fmt.Sprintf("%s:%d", sshHost, sshPort), &sshConfig)

	err = scpClient.Connect()
	if err !=nil{
		return err
	}

	fileData, err := os.Open(src)
	if err !=nil{
		return err
	}
	scpClient.CopyFile(fileData, dest, "0655")
	defer scpClient.Session.Close()
	defer fileData.Close()
	return nil

}


func SSHcommand(s *SSHHost){
	sshHost := s.Host
	sshUser := s.Username
	sshPassword := s.Password
	sshType := "key"//password 或者 key
	sshKeyPath := s.Key
	sshPort := s.Port


	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second,//ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            sshUser,
		Config: 	ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"}},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
	}
	if sshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
	}



	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal(fmt.Sprintf("创建ssh client 失败 %s",err.Error()))
	}
	defer sshClient.Close()


	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal(fmt.Sprintf("创建ssh session 失败 %s",err.Error()))
	}

	defer session.Close()
	//执行远程命令
	fmt.Println("***************************************************************************************")
	combo,err := session.CombinedOutput(s.CmdList)
	fmt.Println(fmt.Sprintf("%v         ************************",time.Now()))
	color.White("[命令输入]# %s\n",s.CmdList)
	color.Yellow( "[命令输出]# %s\n",string(combo))
	fmt.Println("***************************************************************************************")
	if err != nil {
		red := color.New(color.FgRed)
		boldRed := red.Add(color.Bold)
		boldRed.Printf("fail: [%v]\n",s.Host)
		log.Fatal(err.Error())
	}else {
		color.Green("ok: [%v]",s.Host)
	}
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("find key's home dir failed %s", err.Error()))
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("ssh key file read failed %s", err.Error()))
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(fmt.Sprintf("ssh key signer failed %s", err.Error()))
	}
	return ssh.PublicKeys(signer)
}


