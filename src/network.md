# Network

## resources

- <https://dustinspecker.com/posts/how-do-kubernetes-and-docker-create-ip-addresses/>

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
cat /proc/sys/net/ipv4/ip_forward
echo 1 | tee /proc/sys/net/ipv4/ip_forward

# forward traffic from the virtual device to the physical device and vice versa
iptables --append FORWARD --in-interface veth_dustin --out-interface enp0s7 --jump ACCEPT
iptables --append FORWARD --in-interface enp0s7 --out-interface veth_dustin --jump ACCEPT
iptables-save

# MASQUERADE the IP
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
```
