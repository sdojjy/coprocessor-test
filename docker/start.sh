#!/bin/bash

#docker run -d --name pd  -p 2379:2379   -p 2380:2380  pingcap/pd:latest  --client-urls="http://0.0.0.0:2379"
#docker run -d --name tikv -p 20160:20160  --link pd:pd --ulimit nofile=1000000:1000000   pingcap/tikv:latest --pd="pd:2379"
#docker run -d --name tidb -p 4000:4000   --link pd:pd -p 10080:10080 pingcap/tidb:latest  --path="pd:2379"

kubectl port-forward svc/demo-pd 2379:2379 -n tidb
kubectl port-forward svc/demo-tidb 4000:4000 -n tidb
kubectl port-forward svc/demo-tikv-peer 20160:20160 -n tidb
kubectl port-forward svc/demo-tikv-peer 20160:20160 -n tidb
