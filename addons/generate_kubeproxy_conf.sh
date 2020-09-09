#!/bin/bash

cat > /tmp/kube-proxy.conf << EOF
apiVersion: kubeproxy.config.k8s.io/v1alpha1
bindAddress: ${nodename}
clientConnection:
    acceptContentTypes: ""
    burst: 10
    contentType: application/vnd.kubernetes.protobuf
    kubeconfig: /etc/kubernetes/kube-proxy.kubeconfig
    qps: 5
clusterCIDR: "${PodCIDR}"
configSyncPeriod: 15m0s
conntrack:
    max: null
    maxPerCore: 32768
    min: 131072
    tcpCloseWaitTimeout: 1h0m0s
    tcpEstablishedTimeout: 24h0m0s
enableProfiling: false
healthzBindAddress: ${nodename}:10256
hostnameOverride: ${nodename}
iptables:
    masqueradeAll: true
    masqueradeBit: 14
    minSyncPeriod: 0s
    syncPeriod: 30s
ipvs:
    excludeCIDRs: null
    minSyncPeriod: 0s
    scheduler: ""
    syncPeriod: 30s
kind: KubeProxyConfiguration
metricsBindAddress: ${nodename}:10249
mode: "ipvs"
nodePortAddresses: null
oomScoreAdj: -999
portRange: ""
resourceContainer: /kube-proxy
udpIdleTimeout: 250ms
EOF