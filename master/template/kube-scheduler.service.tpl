[Unit]
Description=Kubernetes Scheduler
Documentation=https://github.com/kubernetes/kubernetes
After=network.target

[Service]
ExecStart=/usr/local/bin/kube-scheduler \
  --leader-elect=true \
  --kubeconfig=/etc/kubernetes/scheduler.kubeconfig \
  --bind-address=0.0.0.0 \
  --logtostderr=false \
  --log-dir=/var/log/kubernetes/kube-scheduler \
  --v=2
Restart=always
RestartSec=10s

[Install]
WantedBy=multi-user.target
