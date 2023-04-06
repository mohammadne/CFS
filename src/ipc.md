# IPC

## resources

- <https://medium.com/@boutnaru/linux-namespaces-ipc-namespace-927f01cbcf3d>

## demo

```bash
# ipcmk create IPC objects
ipcmk -M 10 # create a shared-memory
ipcmk -Q # create a message-queue

# list IPC objects
ipcs

sudo unshare --ipc

# list IPC objects on the new container
ipcs
exit

# list IPC objects again on the host
ipcs
```
