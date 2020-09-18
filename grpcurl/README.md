`grpcurl` is used to deploy `grpcurl` to interacte with the catalog server.

1. build `grpcurl` image

```
docker build . -t clyang/grpcurl
```

2. deploy `grpcurl.yaml`

3. access the pod to run

```grpcurl -plaintext acm-custom-registry:50051 list api.Registry
api.Registry.GetBundle
api.Registry.GetBundleForChannel
api.Registry.GetBundleThatReplaces
api.Registry.GetChannelEntriesThatProvide
api.Registry.GetChannelEntriesThatReplace
api.Registry.GetDefaultBundleThatProvides
api.Registry.GetLatestChannelEntriesThatProvide
api.Registry.GetPackage
api.Registry.ListBundles
api.Registry.ListPackages
```

```grpcurl -plaintext acm-custom-registry:50051  api.Registry/ListPackages
{
  "name": "advanced-cluster-management"
}
```

```
grpcurl -plaintext -d '{"name":"advanced-cluster-management"}' acm-custom-registry:50051 api.Registry/GetPackage
{
  "name": "advanced-cluster-management",
  "channels": [
    {
      "name": "release-2.1",
      "csvName": "advanced-cluster-management.v2.1.0"
    }
  ],
  "defaultChannelName": "release-2.1"
}
```

```
grpcurl -plaintext -d '{"pkgName":"advanced-cluster-management","channelName":"release-2.1"}' acm-custom-registry:50051 api.Registry/GetBundleForChannel
```
