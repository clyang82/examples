#!/bin/bash

TOKEN=`cat /var/run/secrets/kubernetes.io/serviceaccount/token`

curl -k     -X DELETE -H "Authorization: Bearer $TOKEN"     -H 'Accept: application/json'     -H 'Content-Type: application/json' https://kubernetes.default:443/api/v1/namespaces/open-cluster-management-observability/configmaps/demo
