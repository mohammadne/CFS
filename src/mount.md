# mount

## resources

- <https://www.redhat.com/sysadmin/mount-namespaces>

- <https://www.bleepingcomputer.com/tutorials/introduction-to-mounting-filesystems-in-linux/>

- <https://superuser.com/questions/367595/what-is-the-difference-between-filesystemdevice-location-and-mounted-point-bo>

## demo1

```bash
# system mount points
mount

# report file system disk space usage
df
```

## demo2

- pane1

```bash
# we are on the ubuntu
lsb_release -a

cd

# download alpine minirootsf (https://alpinelinux.org/downloads/)
curl -O https://dl-cdn.alpinelinux.org/alpine/v3.17/releases/x86_64/alpine-minirootfs-3.17.2-x86_64.tar.gz

mkdir -p $HOME/alpine-rootfs

tar xzvf -C $HOME/alpine-rootfs

# create this file at the root of alpine-rootfs
touch $HOME/alpine-rootfs/HOST_UBUNTU_ROOT_FS
```

- pane2

```bash
unshare --uts --pid --fork chroot $HOME/alpine-rootfs /bin/bash -c "hostname inside; /bin/sh"

ls -lah

apt # not work
apk # works!

ping google.com
cat /etc/resolv.conf
echo "nameserver 1.1.1.1" >> /etc/resolv.conf
echo "nameserver 8.8.8.8" >> /etc/resolv.conf

# change the pid space
# mount -t type device dir
#
# nosuid: Do not allow set-user-identifier or set-group-identifier bits to take effect
# nodev: Do not interpret character or block special devices on the file system
# noexec: Do not allow execution of any binaries on the mounted file system.
mount -t proc proc ./proc -o nosuid,nodev,noexec
```

- pane 1

```bash
# not works
lsns

umount $HOME/alpine-rootfs/proc
mount -t proc proc /proc
```

- pane 2

```bash
mount -t proc proc ./proc -o nosuid,nodev,noexec
mkdir -p /mnt/test
mount --bind /usr/bin/ /mnt/test

# all binaries are here
ls -lah /mnt/test
```

- pane 1

```bash
# /mnt/test is also here
mount

# we can see child namespace mount points take affect on the host
# so we have to use mount namespace
mount | grep alpine-rootfs
umount $HOME/alpine-rootfs/proc
umount $HOME/alpine-rootfs/mnt/test
mount | grep alpine-rootfs
```

- pane 2

```bash
exit

# here we use --mount to seperate mount points of the processes inside the container
# also we can run this with --net (ip netns add)
unshare --uts --pid --mount --fork chroot $HOME/alpine-rootfs /bin/bash -c "hostname inside; /bin/sh"

mount

mkdir -p /mnt/test
mount --bind /usr/bin/ /mnt/test

mount
```

- pane 1

```bash
mount

# like the inception movie
lsns --type mnt
```
