ETCD_CERET_DIR=/etc/kubernetes/pki/
ETCD_CA_FILE=ca.pem
ETCD_KEY_FILE=etcd-key.pem
ETCD_CERT_FILE=etcd.pem
ETCD_EP={{ .Cluster_endpoints }}
alias etcd_v2="etcdctl --cert-file ${ETCD_CERET_DIR}/${ETCD_CERT_FILE} \
              --key-file ${ETCD_CERET_DIR}/${ETCD_KEY_FILE}  \
              --ca-file ${ETCD_CERET_DIR}/${ETCD_CA_FILE}  \
              --endpoints $ETCD_EP"

alias etcd_v3="ETCDCTL_API=3 \
    etcdctl   \
   --cert ${ETCD_CERET_DIR}/${ETCD_CERT_FILE} \
   --key ${ETCD_CERET_DIR}/${ETCD_KEY_FILE} \
   --cacert ${ETCD_CERET_DIR}/${ETCD_CA_FILE} \
    --endpoints $ETCD_EP"

function etcd-ha(){
    etcd_v3 endpoint status --write-out=table
}
