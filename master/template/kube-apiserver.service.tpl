[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/kubernetes/kubernetes
After=network.target

[Service]
ExecStart=/usr/local/bin/kube-apiserver \
  --authorization-mode=Node,RBAC \
  --enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeClaimResize,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota,Priority,PodPreset \
  --advertise-address={{ .Host }} \
  --bind-address={{ .Host }}  \
  --insecure-port=0 \
  --secure-port=6443 \
  --allow-privileged=true \
  --audit-log-maxage=30 \
  --audit-log-maxbackup=3 \
  --audit-log-maxsize=100 \
  --audit-log-path=/var/log/audit.log \
  --storage-backend=etcd3 \
  --etcd-cafile=/etc/kubernetes/pki/ca.pem \
  --etcd-certfile=/etc/kubernetes/pki/etcd-client.pem \
  --etcd-keyfile=/etc/kubernetes/pki/etcd-client-key.pem \
  --etcd-servers={{ .EtcdServers }} \
  --event-ttl=1h \
  --enable-bootstrap-token-auth \
  --client-ca-file=/etc/kubernetes/pki/ca.pem \
  --kubelet-https \
  --kubelet-client-certificate=/etc/kubernetes/pki/kubelet-client.pem \
  --kubelet-client-key=/etc/kubernetes/pki/kubelet-client-key.pem \
  --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname \
  --runtime-config=api/all=true,settings.k8s.io/v1alpha1=true \
  --service-cluster-ip-range={{ .SvcCIDR }} \
  --service-node-port-range={{ .ServiceNodePortRange }} \
  --service-account-key-file=/etc/kubernetes/pki/sa.pub \
  --tls-cert-file=/etc/kubernetes/pki/kube-apiserver.pem \
  --tls-private-key-file=/etc/kubernetes/pki/kube-apiserver-key.pem \
  --requestheader-client-ca-file=/etc/kubernetes/pki/ca.pem \
  --requestheader-username-headers=X-Remote-User \
  --requestheader-group-headers=X-Remote-Group \
  --requestheader-allowed-names=front-proxy-client \
  --requestheader-extra-headers-prefix=X-Remote-Extra- \
  --proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.pem \
  --proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client-key.pem \
  --logtostderr=false \
  --log-dir=/var/log/kubernetes/kube-apiserver \
  --v=2

Restart=on-failure
RestartSec=10s
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
