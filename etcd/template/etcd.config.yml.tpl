name: '{{ .NodeName  }}'
data-dir: {{ .DataDir }}
wal-dir: {{ .Waldir }}
snapshot-count: 5000
heartbeat-interval: 100
election-timeout: 1000
quota-backend-bytes: 0
listen-peer-urls: 'https://{{ .Host }}:2380'
listen-client-urls: 'https://{{ .Host }}:2379,http://127.0.0.1:2379'
max-snapshots: 3
max-wals: 5
cors:
initial-advertise-peer-urls: 'https://{{ .Host }}:2380'
advertise-client-urls: 'https://{{ .Host }}:2379'
discovery:
discovery-fallback: 'proxy'
discovery-proxy:
discovery-srv:
initial-cluster: '{{ .Init_cluster }}'

initial-cluster-token: 'etcd-k8s-cluster'
initial-cluster-state: 'new'
strict-reconfig-check: false
enable-v2: false
enable-pprof: true
proxy: 'off'
proxy-failure-wait: 5000
proxy-refresh-interval: 30000
proxy-dial-timeout: 1000
proxy-write-timeout: 5000
proxy-read-timeout: 0
client-transport-security:
  ca-file: '/etc/kubernetes/pki/ca.pem'
  cert-file: '/etc/kubernetes/pki/etcd.pem'
  key-file: '/etc/kubernetes/pki/etcd-key.pem'
  client-cert-auth: true
  trusted-ca-file: '/etc/kubernetes/pki/ca.pem'
  auto-tls: true
peer-transport-security:
  ca-file: '/etc/kubernetes/pki/ca.pem'
  cert-file: '/etc/kubernetes/pki/etcd.pem'
  key-file: '/etc/kubernetes/pki/etcd-key.pem'
  peer-client-cert-auth: true
  trusted-ca-file: '/etc/kubernetes/pki/ca.pem'
  auto-tls: true
debug: false
logger: zap
log-outputs: [stderr]
force-new-cluster: false

auto-compaction-mode: periodic
auto-compaction-retention: "1"

