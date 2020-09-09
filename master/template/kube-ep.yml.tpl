apiVersion: v1
kind: Service
metadata:
  namespace: kube-system
  name: kube-controller-manager
  labels:
    k8s-app: kube-controller-manager
spec:
#  selector:
#    component: kube-controller-manager
  type: ClusterIP
#  clusterIP: None
  ports:
  - name: http-metrics
    port: 10252
    targetPort: 10252
    protocol: TCP
---
apiVersion: v1
kind: Service
metadata:
  namespace: kube-system
  name: kube-scheduler
  labels:
    k8s-app: kube-scheduler
spec:
#  selector:
#    component: kube-scheduler
  type: ClusterIP
#  clusterIP: None
  ports:
  - name: http-metrics
    port: 10251
    targetPort: 10251
    protocol: TCP
---
apiVersion: v1
kind: Endpoints
metadata:
  labels:
    k8s-app: kube-controller-manager
  name: kube-controller-manager
  namespace: kube-system
subsets:
- addresses:
  {{- range $index, $host := .Host}}
  - ip: {{$host}}
  {{- end}}
  ports:
  - name: http-metrics
    port: 10252
    protocol: TCP
---
apiVersion: v1
kind: Endpoints
metadata:
  labels:
    k8s-app: kube-scheduler
  name: kube-scheduler
  namespace: kube-system
subsets:
- addresses:
  {{- range $index, $host := .Host}}
  - ip: {{$host}}
  {{- end}}

  ports:
  - name: http-metrics
    port: 10251
    protocol: TCP
