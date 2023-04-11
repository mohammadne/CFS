# Cgroup

## resources

- <https://book.hacktricks.xyz/linux-hardening/privilege-escalation/docker-security/cgroups>

- <https://book.hacktricks.xyz/linux-hardening/privilege-escalation/docker-security/namespaces/cgroup-namespace>

## demo

- create and run the script

    ```bash
    #!/bin/bash

    for i in {1..100}; do
        sleep 100 &
    done
    ```

- edit cgroups

    ```bash
    sudo mkdir -p /sys/fs/cgroup/cfs

    sudo echo 30 > /sys/fs/cgroup/cfs/pids.max
    sudo vim /sys/fs/cgroup/cfs/pids.max

    echo $$
    sudo echo 10761 > /sys/fs/cgroup/cfs/cgroup.procs
    ```

- run the script again
