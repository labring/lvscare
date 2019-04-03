[![Build Status](https://cloud.drone.io/api/badges/fanux/LVScare/status.svg)](https://cloud.drone.io/fanux/LVScare)
# LVScare
A lightweight LVS baby care, support ipvs health check

## Feature
If ipvs real server is unavilible, remove it, if real server return to normal, add it back.  This is useful for kubernetes master HA.

## Quick Start
```
lvscare create --vs 10.103.97.12:6443 --rs 192.168.0.2:6443 --rs 192.168.0.3:6443 --rs 192.168.0.4:6443 
```
Then kubeadm join can use `10.103.97.12:6443` instead real masters.

Run lvscare as a static pod on every kubernetes node.
```
lvscare care --vs 10.103.97.12:80 --rs 192.168.0.2:6443 --rs 192.168.0.3:6443 --rs 192.168.0.4:6443 -t 5s
```
* -t every 5s check the real server port
    * --health-port "6443" if not return 200 OK, remove the realserver
    * --health-path "/healthz" if not return 200 OK, remove the realserver
