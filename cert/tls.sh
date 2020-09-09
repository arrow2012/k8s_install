#!/bin/bash


if [ $# == 0 ];then
    echo "带上k8s证书目录参数；"
    exit 1
fi

timenow=$(date +%Y%m%d-%H%M%S)
csrDir="../"

mkdir -p  ../file/cert/${timenow}
cd ../file/cert/${timenow}

chmod +x /usr/local/bin/cfssl
chmod +x /usr/local/bin/cfssljson

ls ../

CACertDir=ca/
k8s_cert_dir=$1
if [ ! -d ${k8s_cert_dir} ];then
  mkdir -p ${k8s_cert_dir} || exit 1
else
  echo "删除目录"
  rm -fr ${k8s_cert_dir}
  mkdir -p ${k8s_cert_dir} || exit 1
fi

echo "生成CA"
[ ! -d "ca/" ] && mkdir ca || echo "开始执行"
pwd
cfssl gencert -initca ${csrDir}/ca-csr.json | cfssljson -bare ca/ca || exit 1
echo "生成 api server服务端证书"
/usr/bin/rsync ca/*.pem ${k8s_cert_dir}  || exit 1

sleep 1
[ ! -d "kube-apiserver/" ] && mkdir kube-apiserver || echo "开始执行"
cfssl gencert -ca=${CACertDir}/ca.pem -ca-key=${CACertDir}/ca-key.pem -config=${csrDir}/ca-config.json -profile=server ${csrDir}/kube-apiserver-csr.json | cfssljson -bare kube-apiserver/kube-apiserver || exit 1
/usr/bin/rsync kube-apiserver/*.pem ${k8s_cert_dir} || exit 1
echo "生成 kubelet client证书"
sleep 1
[ ! -d "kubelet/" ] && mkdir kubelet || echo "开始执行"
cfssl gencert -ca=${CACertDir}/ca.pem -ca-key=${CACertDir}/ca-key.pem -config=${csrDir}/ca-config.json -profile=client ${csrDir}/kubelet-client-csr.json | cfssljson -bare kubelet/kubelet-client || exit 1
/usr/bin/rsync kubelet/*.pem ${k8s_cert_dir} || exit 1
echo "生成 etcd对等证书"
sleep 1
[ ! -d "etcd/" ] && mkdir etcd || echo "开始执行"
cfssl gencert -ca=${CACertDir}/ca.pem -ca-key=${CACertDir}/ca-key.pem -config=${csrDir}/ca-config.json -profile=peer ${csrDir}/etcd-csr.json | cfssljson -bare etcd/etcd || exit 1
echo "生成 etcd client证书"
sleep 1
[ ! -d "etcd/" ] && mkdir etcd || echo "开始执行"
cfssl gencert -ca=${CACertDir}/ca.pem -ca-key=${CACertDir}/ca-key.pem -config=${csrDir}/ca-config.json -profile=client ${csrDir}/etcd-csr.json | cfssljson -bare etcd/etcd-client || exit 1
/usr/bin/rsync etcd/*.pem ${k8s_cert_dir} || exit 1

echo "生成 admin client 证书"
sleep 1
[ ! -d "admin/" ] && mkdir admin || echo "开始执行"
cfssl gencert -ca=${CACertDir}/ca.pem -ca-key=${CACertDir}/ca-key.pem -config=${csrDir}/ca-config.json -profile=client ${csrDir}/admin-csr.json | cfssljson -bare admin/admin || exit 1
/usr/bin/rsync admin/*.pem ${k8s_cert_dir} || exit 1

echo "生成 kube-scheduler 证书"
sleep 1
[ ! -d "kube-scheduler/" ] && mkdir kube-scheduler || echo "开始执行"
cfssl gencert -ca=${CACertDir}/ca.pem -ca-key=${CACertDir}/ca-key.pem -config=${csrDir}/ca-config.json -profile=client ${csrDir}/kube-scheduler-csr.json | cfssljson -bare kube-scheduler/kube-scheduler || exit 1
/usr/bin/rsync  kube-scheduler/*.pem ${k8s_cert_dir} || exit 1
echo "生成 front-proxy-client client证书"
sleep 1
[ ! -d "front-proxy/" ] && mkdir front-proxy || echo "开始执行"
cfssl gencert -ca=${CACertDir}/ca.pem -ca-key=${CACertDir}/ca-key.pem -config=${csrDir}/ca-config.json -profile=client ${csrDir}/front-proxy-client-csr.json | cfssljson -bare front-proxy/front-proxy-client || exit 1
/usr/bin/rsync front-proxy/*.pem ${k8s_cert_dir} || exit 1

echo "生成 sa 证书和公私钥"
sleep 1

[ ! -d "sa/" ] && mkdir sa || echo "开始执行"
cfssl gencert -ca=${CACertDir}/ca.pem -ca-key=${CACertDir}/ca-key.pem -config=${csrDir}/ca-config.json -profile=client ${csrDir}/sa-csr.json | cfssljson -bare sa/sa || exit 1
openssl rsa -in sa/sa-key.pem -outform PEM -pubout -out sa/sa.pub || exit 1

/usr/bin/rsync sa/*.pem ${k8s_cert_dir} || exit 1
/usr/bin/rsync sa/*.pub ${k8s_cert_dir} || exit 1

echo "证书生成完成"



