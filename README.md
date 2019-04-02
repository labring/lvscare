# LVScare
A lightweight LVS baby care, support ipvs health check

## Feature
If ipvs real server is unavilible, remove it, if real server return to normal, add it back.  This is useful for kubernetes master HA.

## Quick Start
```
lvscare --vs 10.103.97.12:80 --rs 192.168.0.2:6443 --rs 192.168.0.3:6443 --rs 192.168.0.4:6443
```
