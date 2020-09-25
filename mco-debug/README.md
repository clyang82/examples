This is designed to debug mco operator with your image. Please modify the image to your image in operator.yaml
1. `oc apply -f https://raw.githubusercontent.com/clyang82/examples/master/mco-debug/operator.yaml -n open-cluster-management`
2. `oc delete pod -lname=multicluster-observability-operator-debug -n open-cluster-management; oc delete deploy multicluster-observability-operator -n open-cluster-management`
