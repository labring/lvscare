# lvscare release note
```
lvscare care --vs 10.103.97.12:6443 --rs 192.168.0.2:6443 --rs 192.168.0.3:6443 --rs 192.168.0.4:6443 \
--health-path / --health-schem http
```

# Test
Clean your environment:
```
ip link del dev sealyun-ipvs
ipvsadm -C
```

start some nginx as realserver
```
docker run -p 8081:80 --name nginx1 -d nginx
docker run -p 8082:80 --name nginx2 -d nginx
docker run -p 8083:80 --name nginx3 -d nginx
```
```
lvscare care --vs 10.103.97.12:6443 --rs 127.0.0.1:8081 --rs 127.0.0.1:8082 --rs 127.0.0.1:8083 \
--health-path / --health-schem http
```

check ipvs rules:
```
ipvsadm -Ln
```

delete a nginx
```
docker stop nginx1
```

check ipvs again:
```
ipvsadm -Ln
```
