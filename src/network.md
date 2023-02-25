# Network

## resources

- <https://dustinspecker.com/posts/how-do-kubernetes-and-docker-create-ip-addresses/>

- <https://linuxconfig.org/how-to-turn-on-off-ip-forwarding-in-linux>

## demo 1

```bash
# start HTTP server in host namespace
python3 -m http.server 8080

# create new network namespace
ip netns add apple_ns

# start HTTP server in apple_ns
ip netns exec apple_ns python3 -m http.server 8080

# list network namespaces
ip netns list
lsns --type=net

# Shows addresses assigned to all network interfaces
ip netns exec apple_ns ip address list

# start the loopback device up
ip netns exec apple_ns ip link set dev lo up

# Shows addresses assigned to all network interfaces
ip netns exec apple_ns ip address list

# request http server on apple_ns namespace
ip netns exec apple_ns curl localhost:8080

# create two virtual ethernet devices in the host network namespace
ip link add dev host_veth type veth peer name apple_veth
ip link list

# move the apple_veth device to the apple_ns network namespace
ip link set apple_veth netns apple_ns
ip netns exec apple_ns ip link list

# start host_veth device
ip link set dev host_veth up
ip link list

# assign ip address to host_veth
ip address add 10.0.0.10/24 dev host_veth
ip address list

# start apple_veth device
ip netns exec apple_ns ip link set dev apple_veth up

# assign ip address to apple_veth
ip netns exec apple_ns ip address add 10.0.0.11/24 dev apple_veth
ip netns exec apple_ns ip address list

# test ping connectivity
ping 10.0.0.10 -c 1
ping 10.0.0.11 -c 1
ip netns exec apple_ns ping 10.0.0.10 -c 1
ip netns exec apple_ns ping 10.0.0.11 -c 1

# request to apple_ns http server from host namespace
curl 10.0.0.11:8080

# request to host≈ http server from apple_ns namespace
ip netns exec apple_ns curl 10.0.0.10:8080

# use computer’s local IP instead of host_veth ip address (doesn't work)
ip netns exec apple_ns curl 192.168.64.4:8080
ip netns exec apple_ns ip route get 192.168.64.4

# enable apple_ns to make requests to computer’s local IP
ip netns exec apple_ns ip route show
ip netns exec apple_ns ip route add default via 10.0.0.10
ip netns exec apple_ns ip route show

# use computer’s local IP instead of host_veth ip address
ip netns exec apple_ns curl 192.168.64.4:8080

# talk to internet from our apple_ns (fails)
ip netns exec apple_ns ping 8.8.8.8 -c 1

# enable IP forwarding
sysctl net.ipv4.ip_forward
sysctl -w net.ipv4.ip_forward=1
sysctl -p

# forward traffic from the virtual device to the physical device and vice versa
iptables --append FORWARD --in-interface host_veth --out-interface enp0s7 --jump ACCEPT
iptables --append FORWARD --in-interface enp0s7 --out-interface host_veth --jump ACCEPT
iptables-save

# MASQUERADE the IPs (SNAT)
sudo iptables --append POSTROUTING --table nat --out-interface enp0s7 --jump MASQUERADE

# talk to internet from our apple_ns
ip netns exec apple_ns ping 8.8.8.8 -c 1

# test DNS (fails)
ip netns exec apple_ns ping google.com -c 1

# configure apple_ns DNS nameserver
mkdir -p /etc/netns/apple_ns
echo "nameserver 8.8.8.8" > /etc/netns/apple_ns/resolv.conf

# test DNS
ip netns exec apple_ns ping google.com -c 1

# cleanup
ip link delete dev host_veth
ip netns delete apple_ns
iptables --delete FORWARD --in-interface host_veth --out-interface enp0s7 --jump ACCEPT
iptables --delete FORWARD --in-interface enp0s7 --out-interface host_veth --jump ACCEPT
iptables --delete POSTROUTING --table nat --out-interface enp0s7 --jump MASQUERADE
```

## demo 2

```bash
# enable IP forwarding
sysctl net.ipv4.ip_forward
sysctl -w net.ipv4.ip_forward=1
sysctl -p

# # MASQUERADE the IPs (SNAT)
sudo iptables --append POSTROUTING --table nat --out-interface enp0s7 --jump MASQUERADE

# setup apple namespace
ip link add dev host1_veth type veth peer name apple_veth
ip link set dev host1_veth up
ip address add 10.0.0.10/24 dev host1_veth
ip netns add apple_ns
ip link set dev apple_veth netns apple_ns
ip netns exec apple_ns ip link set dev lo up
ip netns exec apple_ns ip link set dev apple_veth up
ip netns exec apple_ns ip address add 10.0.0.11/24 dev apple_veth
ip netns exec apple_ns ip route add default via 10.0.0.10
ip netns exec apple_ns python3 -m http.server 8080 &

# setup lemon namespace
ip link add dev host2_veth type veth peer name lemon_veth
ip link set dev host2_veth up
ip address add 10.0.0.20/24 dev host2_veth
ip netns add lemon_ns
ip link set dev lemon_veth netns lemon_ns
ip netns exec lemon_ns ip link set dev lo up
ip netns exec lemon_ns ip link set dev lemon_veth up
ip netns exec lemon_ns ip address add 10.0.0.21/24 dev lemon_veth
ip netns exec lemon_ns ip route add default via 10.0.0.20
ip netns exec lemon_ns python3 -m http.server 8080 &

# routing issues for forwarding packets between apple and lemon namespaces
ip link list

# so we use bridge interface to overcome that
ip link add dev host_bridge type bridge
ip address add 10.0.0.1/24 dev host_bridge
ip link set host_bridge up
ip address list

# connect virtual ethernets to the host_bridge
ip link set dev host1_veth master host_bridge
ip link set dev host2_veth master host_bridge

# update default routes
ip netns exec apple_ns ip route delete default via 10.0.0.10
ip netns exec apple_ns ip route add default via 10.0.0.1
ip netns exec lemon_ns ip route delete default via 10.0.0.20
ip netns exec lemon_ns ip route add default via 10.0.0.1

ip address delete 10.0.0.10/24 dev host1_veth
ip address delete 10.0.0.20/24 dev host2_veth

# check connectivity from host to namespaces
ping 10.0.0.11 -c 1
ping 10.0.0.21 -c 1

# check connectivity between namespaces (fails)
ip netns exec apple_ns ping 10.0.0.21 -c 1
ip netns exec lemon_ns ping 10.0.0.11 -c 1

# enabling a bridge to forward traffic from one veth to another veth
iptables --append FORWARD --in-interface host_bridge --out-interface host_bridge --jump ACCEPT

# check connectivity between namespaces
ip netns exec apple_ns ping 10.0.0.21 -c 1
ip netns exec lemon_ns ping 10.0.0.11 -c 1
ip netns exec apple_ns curl 10.0.0.21:8080
ip netns exec lemon_ns curl 10.0.0.11:8080

# add rules to forward traffic between host_bridge and enp0s7
iptables --append FORWARD --in-interface host_bridge --out-interface enp0s7 --jump ACCEPT
iptables --append FORWARD --in-interface enp0s7 --out-interface host_bridge --jump ACCEPT

# cleanup
ip link delete dev host_bridge
ip link delete dev host1_veth
ip link delete dev host2_veth
ip netns delete apple_ns
ip netns delete lemon_ns
iptables --delete FORWARD --in-interface host_bridge --out-interface enp0s7 --jump ACCEPT
iptables --delete FORWARD --in-interface enp0s7 --out-interface host_bridge --jump ACCEPT
iptables --delete POSTROUTING --table nat --out-interface enp0s7 --jump MASQUERADE
```
