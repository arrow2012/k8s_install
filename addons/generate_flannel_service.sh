#!/bin/bash

cat > /tmp/flanneld.service <<EOF
[Unit]
Description=Network fabric for containers
Documentation=https://github.com/coreos/flannel
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
Restart=always
RestartSec=5
KillMode=process
# This is needed because of this: https://github.com/coreos/flannel/issues/792
# Kubernetes knows the nodes by their FQDN so we have to use the FQDN
#Environment=NODE_NAME=my-node.foo.bar.com
# Note that we don't specify any etcd option. This is because we want to talk
# to the apiserver instead. The apiserver then talks to etcd on flannel's
# behalf.
Environment=NODE_NAME=${nodename}
ExecStart=/usr/local/bin/flanneld \
  --kube-subnet-mgr=true \
  --kubeconfig-file=/etc/kubernetes/flanneld.kubeconfig \
  --ip-masq=true \
  --iface=eth0 \
  --public-ip ${nodename} \
  --healthz-ip ${nodename} \
  --healthz-port ${flannel_healthzport} \
  --v=2

[Install]
WantedBy=multi-user.target
EOF


cat > /tmp/10-flannel.conflist <<EOF
{
  "name": "cbr0",
  "cniVersion": "0.3.1",
  "plugins": [
    {
      "type": "flannel",
      "delegate": {
        "hairpinMode": true,
        "isDefaultGateway": true
      }
    },
    {
      "type": "portmap",
      "capabilities": {
        "portMappings": true
      }
    }
  ]
}
EOF

cat > /tmp/net-conf.json <<EOF
{
  "Network": "${PodCIDR}",
  "Backend": {
    "Type": "vxlan"
  }
}
EOF