# UTS

## resources

- <https://en.wikipedia.org/wiki/Linux_namespaces>

- <https://www.redhat.com/sysadmin/uts-namespace>

## demo

- pane1

    ```bash
    hostname
    ```

- pane2

    ```bash
    hostname
    ```

- pane1

    ```bash
    hostname uts1
    hostname
    ```

- pane2

    ```bash
    # here also the hostname changes to uts1
    hostname

    unshare --uts /bin/bash
    hostname

    hostname uts2
    bash # update PS1 value
    ```

- pane1

    ```bash
    # the hostname doesn't take affect
    hostname

    # list all uts namespaces
    lsns --type uts
    ```
