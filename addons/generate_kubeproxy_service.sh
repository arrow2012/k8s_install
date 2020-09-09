#!/bin/bash

cat > /tmp/kube-proxy.service << EOF
[Unit]
Description=Kubernetes Kube Proxy
Documentation=https://github.com/kubernetes/kubernetes
After=network.target

[Service]
ExecStart=/usr/local/bin/kube-proxy \
  --hostname-override=${nodename} \
  --config=/etc/kubernetes/kube-proxy.conf \
  --logtostderr=false \
  --log-dir=/var/log/kubernetes/kube-proxy \
  --v=2

Restart=always
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF