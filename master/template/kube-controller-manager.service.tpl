[Unit]
Description=Kubernetes Controller Manager
Documentation=https://github.com/kubernetes/kubernetes
After=network.target

[Service]
ExecStart=/usr/local/bin/kube-controller-manager \
  --allocate-node-cidrs=true \
  --authentication-kubeconfig=/etc/kubernetes/controller-manager.kubeconfig \
  --authorization-kubeconfig=/etc/kubernetes/controller-manager.kubeconfig \
  --bind-address=0.0.0.0 \
  --client-ca-file=/etc/kubernetes/pki/ca.pem \
  --cluster-signing-cert-file=/etc/kubernetes/pki/ca.pem \
  --cluster-signing-key-file=/etc/kubernetes/pki/ca-key.pem \
  --kubeconfig=/etc/kubernetes/controller-manager.kubeconfig \
  --leader-elect=true \
  --cluster-cidr={{ .PodCIDR }} \
  --service-cluster-ip-range={{ .SvcCIDR }} \
  --requestheader-client-ca-file=/etc/kubernetes/pki/ca.pem \
  --service-account-private-key-file=/etc/kubernetes/pki/sa-key.pem \
  --root-ca-file=/etc/kubernetes/pki/ca.pem \
  --use-service-account-credentials=true \
  --controllers=*,bootstrapsigner,tokencleaner \
  --experimental-cluster-signing-duration=86700h \
  --feature-gates=RotateKubeletClientCertificate=true \
  --node-monitor-period=5s \
  --node-monitor-grace-period=2m \
  --pod-eviction-timeout=1m \
  --logtostderr=false \
  --log-dir=/var/log/kubernetes/kube-controller-manager \
  --v=2


Restart=always
RestartSec=10s

[Install]
WantedBy=multi-user.target
