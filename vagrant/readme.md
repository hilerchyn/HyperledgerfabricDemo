# IP tables

workbench: 192.168.31.100


zookeeper1: 192.168.31.201
zookeeper2: 192.168.31.202
zookeeper3: 192.168.31.203

kafka1: 192.168.31.211
kafka2: 192.168.31.212
kafka3: 192.168.31.213
kafka4: 192.168.31.214

orderer0: 192.168.31.220
orderer1: 192.168.31.221
orderer2: 192.168.31.222

peer0Org1: 192.168.31.110
peer1Org1: 192.168.31.111

peer0Org2: 192.168.31.120
peer1Org2: 192.168.31.121
foo27Org2: 192.168.31.122
barOrg2: 192.168.31.123

# 自建registry

官方文档 https://docs.docker.com/registry/

用官方提供的 registry 镜像(https://hub.docker.com/_/registry)，

最新版为2.7.1 （可以尝试第三方 Harbor，据说CNCF云平台用的就是这个）

registry 主要有 v1 和 v2 两个版本，v1基本已经弃用了。

*自建registry 的目的是提高镜像pull速度，统一管理，便于后续版本变更的统一管理*


## 启动
<hr/>

```
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

加上  --restart=always 参数后，重启docker服务时，该container会随之重启

## pull HyperledgerFabric 镜像
<hr/>


```
[vagrant@localhost ~]$ docker pull hyperledger/fabric-baseos
Using default tag: latest
latest: Pulling from hyperledger/fabric-baseos
Digest: sha256:85c85420ff06973532069faf4c3b88de0c8e96a73139c9937c7f82a69c876138
Status: Image is up to date for hyperledger/fabric-baseos:latest
docker.io/hyperledger/fabric-baseos:latest

```

注意上面最后一行是是拉取镜像的全名称


## 给镜像打标签，使其指向自建registry
<hr/>


```
 docker image tag docker.io/hyperledger/fabric-baseos registry:5000/tao-baseos
```

命令中的registry指向的是自建服务器的域名，已在 /etc/hosts 中设置，内容如下：

```
[vagrant@localhost ~]$ cat /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
192.168.31.101 registry
```

## 将镜像推到自建 registry
<hr/>

```
docker push registry:5000/tao-baseos
```

输出如下错误：

```
[vagrant@localhost ~]$ docker push registry:5000/tao-baseos
The push refers to repository [registry:5000/tao-baseos]
Get https://registry:5000/v2/: http: server gave HTTP response to HTTPS client
```

如字面意思，客户端用的是https协议，服务器端却给的是http的应答。

编辑docker配置文件，如果不存在则创建

```
sudo emacs /etc/docker/daemon.json
```

输入如下内容，

```
{"insecure-registries":["registry:5000"]}
```

指明改registry服务insecure，即不使用https协议

重启 docker

```
docker container stop registry
sudo systemctl restart docker
docker container start registry
```

启动时加上 --restart=always就不需要 stop和start了

重新push

```
[vagrant@localhost ~]$ docker push registry:5000/tao-baseos
The push refers to repository [registry:5000/tao-baseos]
f20608eea0b0: Pushed
541f67914ac3: Pushed
2db44bce66cd: Pushed
latest: digest: sha256:1e09c0d5d32d2ef3b157b271fd6ec5ae8c4ba9ae891a53137204e578531ff8a6 size: 948
```