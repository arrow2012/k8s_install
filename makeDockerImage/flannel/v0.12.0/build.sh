#!/bin/sh
moduleName="flannel"
work_path=$(dirname $0)
cd ./${work_path}
work_path=$(pwd)
version=${work_path##*/}
docker build -t  registry.cn-shanghai.aliyuncs.com/itcam/${moduleName}:${version} .
docker login --username=yeliang4@gmail.com --password=itcam12345 registry.cn-shanghai.aliyuncs.com
docker push registry.cn-shanghai.aliyuncs.com/itcam/${moduleName}:${version}