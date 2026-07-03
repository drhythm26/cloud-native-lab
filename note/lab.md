## 01 阔容 replicaCount

操作：
改变 `podinfo/values.yaml` 中`replicaCount` 由3->5
现象：

```sh                             
> kubectl get pods -n podinfo -w
NAME                       READY   STATUS    RESTARTS   AGE
podinfo-5cfbcc5d5c-2m49j   1/1     Running   0          65m
podinfo-5cfbcc5d5c-4kq8f   1/1     Running   0          65m
podinfo-5cfbcc5d5c-5nsrc   1/1     Running   0          18s
podinfo-5cfbcc5d5c-bs4l9   1/1     Running   0          65m
podinfo-5cfbcc5d5c-srgrm   1/1     Running   0          18s

> kubectl describe application podinfo -n argocd
....
Events:
  Type    Reason              Age   From                           Message
  ----    ------              ----  ----                           -------
  Normal  OperationStarted    17m   argocd-application-controller  Initiated automated sync to '4ebe26e6fa669901c6be8eaebe3083b9072af15a'
  Normal  ResourceUpdated     17m   argocd-application-controller  Updated sync status: Synced -> OutOfSync
  Normal  OperationCompleted  17m   argocd-application-controller  Sync operation to 4ebe26e6fa669901c6be8eaebe3083b9072af15a succeeded
  Normal  ResourceUpdated     17m   argocd-application-controller  Updated sync status: OutOfSync -> Synced
  Normal  ResourceUpdated     17m   argocd-application-controller  Updated health status: Healthy -> Progressing
  Normal  ResourceUpdated     17m   argocd-application-controller  Updated health status: Progressing -> Healthy
```
其中新增了`2`个`pod`，`kube describe`的`Events`也很清楚 `Synced -> OutOfSync -> Synced`


## 
