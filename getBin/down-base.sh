#!/bin/bash

kube_image="registry.cn-shanghai.aliyuncs.com/itcam/k8s:v1.18.6"
etcd_image="registry.cn-shanghai.aliyuncs.com/itcam/etcd:v3.4.10"
pause_image="registry.cn-shanghai.aliyuncs.com/itcam/pause-amd64:3.2"
cni_image="registry.cn-shanghai.aliyuncs.com/itcam/cni-plugins:v0.8.6"
flanneld_image="registry.cn-shanghai.aliyuncs.com/itcam/flannel:v0.12.0"
calico_image="registry.cn-shanghai.aliyuncs.com/itcam/pause-amd64:3.2"
docker_image_dir="/opt/image"
cni_binary="cni-plugins-linux-amd64-v0.8.6.tgz"


down_kube(){
    [ ! -f kubernetes-server-linux-amd64.tar.gz ] && {
        docker pull ${kube_image}
        docker run --rm -d --name kube ${kube_image} sleep 10
        docker cp kube:${docker_image_dir}/kubernetes-server-linux-amd64.tar.gz .
        tar -zxvf kubernetes-server-linux-amd64.tar.gz  --strip-components=3 -C /usr/local/bin kubernetes/server/bin/kube{let,ctl,-apiserver,-controller-manager,-scheduler,-proxy}
    } || :
}

down_etcd(){
    docker pull ${etcd_image}
    docker run --rm -d --name etcd ${etcd_image} sleep 10
    docker cp etcd:${docker_image_dir}/etcd /usr/local/bin
    docker cp etcd:${docker_image_dir}/etcdctl /usr/local/bin
}

down_flannel(){
    docker pull ${flanneld_image}
    docker run --rm --entrypoint sh -d --name flanneld ${flanneld_image} -c 'sleep 10'
    docker cp flanneld:${docker_image_dir}/flanneld /usr/local/bin/
    docker cp flanneld:${docker_image_dir}/mk-docker-opts.sh /usr/local/bin/
}

down_cni(){
    docker pull ${cni_image}
    docker run -d --rm --name cni ${cni_image} sleep 10
    docker cp cni:${docker_image_dir}/${cni_binary} /usr/local/src/
}

down_base(){
    down_kube
    down_etcd
    down_flannel
    down_cni
}

if [ "${#@}" -eq 1 ];then
    if [ "$1" != 'all' ];then
        down_$1
    else
        down_kube
        down_etcd
        down_flannel
        down_cni
    fi
else
    echo you must choose a type to download
    exit 0
fi
