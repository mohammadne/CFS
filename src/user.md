# user

## resources

- <https://kubernetes.io/docs/concepts/workloads/pods/user-namespaces/>

- <https://www.redhat.com/sysadmin/building-container-namespaces>

- <https://blog.quarkslab.com/digging-into-linux-namespaces-part-2.html>

## demo

```bash
# get your current user information
id

unshare --user bash
id

unshare --map-root-user bash
id
```
