### How to enable debug oauth-proxy

1.Using privileged scc to enable management-ingress pod can be run as root
```
  annotations:
    openshift.io/scc: privileged
```
```
    securityContext:
      allowPrivilegeEscalation: true
      runAsNonRoot: false
      runAsUser: 0
      seccompProfile:
        type: Unconfined
```
Note: need to add serviceaccount into privileged scc `oc get scc privileged -oyaml`
```
supplementalGroups:
  type: RunAsAny
users:
- system:admin
- system:serviceaccount:openshift-infra:build-controller
- open-cluster-management:management-ingress-4e133
volumes:
```

2. rebuild oauth-proxy by enabling `dlv`
```
FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.16-openshift-4.8 AS builder
WORKDIR  /go/src/github.com/openshift/oauth-proxy
RUN go get github.com/go-delve/delve/cmd/dlv
COPY . .
RUN go build -gcflags="all=-N -l" .

FROM registry.ci.openshift.org/ocp/builder:rhel-8-base-openshift-4.8
COPY --from=builder /go/src/github.com/openshift/oauth-proxy/oauth-proxy /usr/bin/oauth-proxy
COPY --from=builder /go/bin/dlv /dlv
ENTRYPOINT ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/usr/bin/oauth-proxy"]
```

3. modifiy oauth-proxy command to using `dlv`
```
  containers:
  - command:
    - ./dlv
    - --listen=:40000
    - --headless=true
    - --api-version=2
    - --accept-multiclient
    - exec
    - /usr/bin/oauth-proxy
    - --
    - --provider=openshift
    - --upstream=https://management-ingress.open-cluster-management.svc:8443
    - --upstream-ca=/etc/tls/ca/service-ca.crt
    - --https-address=:9443
    - --client-id=multicloudingress
    - --client-secret=multicloudingresssecret
    - --pass-user-bearer-token=true
    - --pass-access-token=true
    - --scope=user:full
    - '-openshift-delegate-urls={"/": {"resource": "projects", "verb": "list"}}'
    - --skip-provider-button=true
    - --cookie-secure=true
    - --cookie-expire=300s
    - --cookie-refresh=30s
    - --tls-cert=/etc/tls/private/tls.crt
    - --tls-key=/etc/tls/private/tls.key
    - --cookie-secret=dlRCUFJvcUxOMUdWRnlIZg==
    - --openshift-ca=/etc/pki/tls/cert.pem
    - --openshift-ca=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    image: quay.io/clyang82/oauth-proxy:debug
```

4. using `port-forward` to forward port in local
`oc port-forward management-ingress-4e133-57d495ffb4-ktsbc 40000:40000`
Use VSCode to debug. the launch.json can be
```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Connect to server",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "port": 40000,
            "host": "127.0.0.1"
        }
    ]
}
```