on v1.7.3
--rbca
kubectl create -f ServiceAccount.yaml            --only 1 once
kubectl create -f ClusterRoleBinding.yaml        --only 1 once 
kubectl create -f tomcat.yaml                    --2 items   1.serviceAccount: admin   2.serviceAccountName: admin
test.go                                          -- golang call rest api

if not use default namespace, you must add resources->limits->cpu and memory
