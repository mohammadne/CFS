# CFS

> CFS stands for Container From Scratch

here we aimed to demonstrate linux namespaces and a sample code show how docker uses namespaces to create containerized environments.

## resources

- <https://medium.com/swlh/build-containers-from-scratch-in-go-part-1-namespaces-c07d2291038b>

- <https://blog.quarkslab.com/digging-into-linux-namespaces-part-1.html>
- <https://blog.quarkslab.com/digging-into-linux-namespaces-part-2.html>

## Golang Implementation

```bash
hostname=devopshobbies
bash

alias cfs='go run main.go'
```

### UTS

```go
func run() {
  fmt.Printf("Running %v \n", os.Args[2:])

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    must(cmd.Run())
}
```

- on the container

    ```bash
    # see the current hostname
    hostname

    # change the hostname
    hostname container
    ```

- on the host

    ```bash
    # see the current hostname
    hostname
    # the container hostname affects the host
    ```

#### add the uts flag

- on the container

    ```bash
    # see the current hostname
    hostname

    # change the hostname
    hostname container
    ```

- on the host

    ```bash
    # see the current hostname
    hostname
    ```

- on the container

    ```bash
    # change the hostname
    hostname container

    # to take effect
    /bin/bash
    ```

#### add the sethostname

we do this in order when a new shell will be started, the new hostname takes effect.

```bash
ps -aux | tail -n10 
```

### PID

#### only add pid clone flag

- on the container

    ```bash
    ps -aux
    # shows all the processes on the machine
    ```

#### add chroot

- on the host

    ```bash
    # ubuntu
    curl -O http://cdimage.ubuntu.com/ubuntu-base/releases/20.04/release/ubuntu-base-20.04.3-base-amd64.tar.gz
    mkdir -p /tmp/ubuntu-rootfs
    tar xzvf ubuntu-base-20.04.3-base-amd64.tar.gz -C /tmp/ubuntu-rootfs
    touch /tmp/ubuntu-rootfs/CONTAINER_ROOT_FS

    # alpine
    curl -O https://dl-cdn.alpinelinux.org/alpine/v3.17/releases/x86_64/alpine-minirootfs-3.17.2-x86_64.tar.gz
    mkdir -p /tmp/alpine-rootfs
    tar xzvf alpine-minirootfs-3.17.2-x86_64.tar.gz -C /tmp/alpine-rootfs
    touch /tmp/alpine-rootfs/CONTAINER_ROOT_FS
    ```

- on the container

    ```bash
    # check the package manager
    apk
    apk add curl

    ping 8.8.8.8
    ping google.com
    echo 'nameserver 8.8.8.8' > /etc/resolv.conf

    apk add curl
    curl

    PS1="\u@\h:\# "

    # run the sleep process
    sleep 100
    ```

- on the host

    ```bash
    ps -C sleep

    ls /proc/7840

    # see where the root of the new process is
    ls -l /proc/7840/root
    ```

- on the container

    ```bash
    ps -a
    # will be failed

    ls /proc
    # not exists the /proc psudo filesystem
    ```

#### add mount and unmount section

- on the container

    ```bash
    ps
    # works

    # see the mount points on the container
    mount
    ```

- on the host

    ```bash
    mount | grep /tmp/container-rootfs
    ```

### MNT

#### add CLONE_NEWNS (for mount)

- on the container

    ```bash
    # see the mount points of the container
    mount
    ```

- on the host

    ```bash
    # you can't see the container mounts
    mount
    ```

- on the container

    ```bash
    sleep 100
    ```

- on the host

    ```bash
    ps -C sleep
    
    # but the host is aware of the mount points of the container but it doesn't clutter up the output
    cat /proc/8029/mounts
    ```
