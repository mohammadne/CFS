# user

## resources

- <https://kubernetes.io/docs/concepts/workloads/pods/user-namespaces/>

- <https://www.redhat.com/sysadmin/building-container-namespaces>

- <https://blog.quarkslab.com/digging-into-linux-namespaces-part-2.html>

## demo

- pane 1

    ```bash
    # get your current user information
    id

    lsns --type user
    ```

- pane 2

    ```bash
    # without remapping
    unshare --user bash

    # get current process user namespace
    readlink /proc/$$/ns/user

    # 56534 comes from /proc/sys/kernel/overflowuid
    id
    ```

- pane 1

    ```bash
    # get your current user information
    id

    # get all user namespaces
    lsns --type user
    ```

- pane 2

    ```bash
    exit

    # with remapping
    unshare --map-root-user bash # unshare -Ur bash
    id

    # get remapping (ID-inside-ns ID-outside-ns range)
    
    # $$ is the PID of the bash
    cat /proc/$$/uid_map 
    
    # the self refers to the cat command itself (ls -l /proc/self/exe)
    # every time you run the ls command, you'll get a new process ID
    cat /proc/self/uid_map 

    # create a file
    touch temp.txt
    ls -l
    ```

- pane 1

    ```bash
    # get ubuntu user-id and group-id
    sudo cat /etc/passwd | grep ubuntu

    # permissions of temp.txt
    ls -l
    ```

- pane 2

    ```bash
    exit
    unshare --user bash
    echo $$
    ```

- pane 1

    ```bash
    echo "0 1000 65335" | sudo tee /proc/1168/uid_map
    echo "0 1000 65335" | sudo tee /proc/1168/gid_map
    ```

- pane 2

    ```bash
    cat /proc/1168/uid_map
    cat /proc/1168/gid_map

    touch hello
    ```

- pane 1

    ```bash
    ls -l hello
    ```
