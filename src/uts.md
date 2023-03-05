# UTS

- <https://en.wikipedia.org/wiki/Linux_namespaces>

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
    hostname outside
    hostname
    ```

- pane2

    ```bash
    # here also the hostname changes to newhost
    hostname

    unshare --uts /bin/bash
    hostname

    hostname inside
    bash # update PS1 value
    ```

- pane1

    ```bash
    # the hostname doesn't take affect
    hostname

    # like the inception movie
    lsns --type uts
    ```
