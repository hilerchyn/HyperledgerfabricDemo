| server | count |
| ---- | ---- |
| zookeeper |  3 |
| kafka  | 4 |
| orderer | 3 |
| peer  | 2 |


kafka最小数量为4，为了满足容错的最小节点数。4个代理可以容错一个代理崩溃，即一个代理(排序节点还是对等节点？)停止服务，
channel仍然可以继续创建、读写。


# 环境准备

1. Docker
2. DockerCompose
3. Golang

# 拉取Docker镜像

不同的服务器分别拉取相应的镜像


# 生成必要文件

## CA文件

```shell
./bin/cryptogen generate --config=./crypto-config.yaml
```

## 创世区块

```shell
./bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
```

启动排序节点要用到


## Channel初始

```shell
./bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID mychannel
```

peer 节点启动后需要创建的频道文件



## 频道

### 创建

mychannel.tx 是频道初始时生成的

```shell
peer channel create -o orderer0.7shu.co:7050 -c mychannel -t 50s -f ./channel-artifacts/mychannel.tx
```

生成 频道初始区块  mychannel.block 

### 加入

mychannel.block 是上一步生成的

```shell
peer channel join -b mychannel.block 
```

## 智能合约 (链码)

### 安装

```shell
 peer chaincode install -n mycc -p github.com/hyperledger/fabric/chaincode/go/chaincode_example02 -v 1.0
```

### 实例化


```shell
peer chaincode instantiate -o orderer0.7shu.co:7050 -C mychannel -n mycc -c '{"Args":["init","A","10","B","10"]}' -P "OR ('Org1MSP.member','Org2MSP.member')" -v 1.0
```

注意此处的 -v 指定的版本为 1.0  与安装时指定的版本一致

### 查询

```shell
peer chaincode query -C mychannel -n mycc -c '{"Args":["query","A"]}'
```


## MSP 证书会过期吗？

```shell
2019-11-13 06:09:13.890 UTC [msp.identity] newIdentity -> DEBU 033 Creating identity instance for cert -----BEGIN CERTIFICATE-----
MIICDTCCAbSgAwIBAgIRAM7VdBZke0/mYHEhgPIjug0wCgYIKoZIzj0EAwIwazEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xFTATBgNVBAoTDG9yZzEuN3NodS5jbzEYMBYGA1UEAxMPY2Eub3Jn
MS43c2h1LmNvMB4XDTE5MTExMjA4MjYwMFoXDTI5MTEwOTA4MjYwMFowVzELMAkG
A1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFu
Y2lzY28xGzAZBgNVBAMMEkFkbWluQG9yZzEuN3NodS5jbzBZMBMGByqGSM49AgEG
CCqGSM49AwEHA0IABIM3fKgG3Vl0JXmEL/x9vAS4ZMvKiej/bxEXj1vXgA4mNWSG
lGBaezgF41KT9helj/vIXFqqBUc93pesfUqWAHWjTTBLMA4GA1UdDwEB/wQEAwIH
gDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQkMCKAINvDix4euRZAizf72ltp0ZJ3kmN2
r++jqvo60+A7pYOOMAoGCCqGSM49BAMCA0cAMEQCIFMu9WYjIyceYNW3kinMa1x7
ULKvBHKoLSCX8AOjdJN2AiAuJOh1KF721Pypr/BKmeNCF81HdQUVRmZmHlUT6G33
4w==
-----END CERTIFICATE-----
2019-11-13 06:09:13.890 UTC [msp] setupSigningIdentity -> DEBU 034 Signing identity expires at 2029-11-09 08:26:00 +0000 UTC
```

"Signing identity expires at" 指明了签名标识有实效时间 为10年。 这个实效时间是不是可以手动指定？


## 一个peer的docker 通过 docker-compose rm -sf 删除后，所有的数据都删除了！！

mychannel.block 丢失，则无法加入频道



## 重新创建、加入频道，安装、实例化 频道的智能和约



关键是  -channelID  这个参数，用来指定后续创建的 channel 名称

Error: got unexpected status: BAD_REQUEST -- initializing configtx manager failed: bad channel ID: channel ID '7shuchannel' contains illegal characters

```shell
./bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/chentaochannel.tx -channelID chentaochannel
```



```shell
peer channel create -o orderer0.7shu.co:7050 -c chentaochannel -t 50s -f ./channel-artifacts/chentaochannel.tx
```


```shell
peer channel join -b chentaochannel.block 
```


```shell
 peer chaincode install -n chentaocc -p github.com/hyperledger/fabric/chaincode/go/chaincode_example02 -v 1.0
```


* 实例化智能合约非常耗时间，需要搞明白是什么原因？

```shell
peer chaincode instantiate -o orderer0.7shu.co:7050 -C chentaochannel -n chentaocc -c '{"Args":["init","A","10","B","15"]}' -P "OR ('Org1MSP.member','Org2MSP.member')" -v 1.0
```

```shell
peer chaincode query -C chentaochannel -n chentaocc -c '{"Args":["query","A"]}'
```


### 新加入频道的peer 如果没有安装智能合约就执行查询操作会遇到问题


如  foo27.org2.7shu.co 节点，

peer channel join -b chentaochannel.block 


未安装智能合约，进行下面操作，

peer chaincode query -C chentaochannel -n chentaocc -c '{"Args":["query","A"]}'


发生如下错误：

Error: endorsement failure during query. response: status:500 message:"cannot retrieve package for chaincode chentaocc/1.0, error open /var/hyperledger/production/chaincodes/chentaocc.1.0: no such file or directory" 


* 是否意味着 智能合约编译后的程序即放置在报错的文件夹内呢？

经过查找，该智能合约编译后，被放置在peer容器对应的目录中，即foo27.org2.7shu.co容器的  /var/hyperledger/production/chaincodes/chentaocc.1.0

* 查看容器列表(docker ps) 发现多了一个容器

```shell
CONTAINER ID        IMAGE                                                                                                   COMMAND                  CREATED             STATUS              PORTS                                        NAMES
1b82e8db3658        dev-foo27.org2.7shu.co-chentaocc-1.0-0d5dd5a7d9c5c4ea387de7b26b9558ae38c4df37e44c1521be56aec05c261288   "chaincode -peer.add…"   6 minutes ago       Up 6 minutes                                                     dev-foo27.org2.7shu.co-chentaocc-1.0
f0bb1c399236        hyperledger/fabric-tools                                                                                "/bin/bash"              2 hours ago         Up 22 minutes                                                    cli
196d6e00f707        hyperledger/fabric-peer                                                                                 "peer node start"        2 hours ago         Up 22 minutes       0.0.0.0:7051-7053->7051-7053/tcp             foo27.org2.7shu.co
c4af12231983        hyperledger/fabric-ca                                                                                   "sh -c 'fabric-ca-se…"   2 hours ago         Up 22 minutes       0.0.0.0:7054->7054/tcp                       ca
bb36cbfa8b4c        hyperledger/fabric-couchdb                                                                              "tini -- /docker-ent…"   2 hours ago         Up 22 minutes       4369/tcp, 9100/tcp, 0.0.0.0:5984->5984/tcp   couchdb
[root@localhost vagrant]# docker exe

```

名字：
dev-foo27.org2.7shu.co-chentaocc-1.0-0d5dd5a7d9c5c4ea387de7b26b9558ae38c4df37e44c1521be56aec05c261288   
推断为用来编译智能合约的

命令：
"chaincode -peer.add…" 
向节点发送编译后的智能合约

* 通过实际操作：

> 关闭一个节点  docker-compose stop 重新开启的时候，这个节点是未启动的，因为docker-compose.yaml 配置文件中没有相应的配置
> 当执行查询操作  peer chaincode query -C chentaochannel -n chentaocc -c '{"Args":["query","A"]}' 后，发现这个容器被启动了


cli节点的配置文件中的volumes 配置中作了如下映射

volumes:
      - /var/run/:/host/var/run/
      
由此可推断，在执行查询操作时，cli容器中的程序通过调用docker的接口创建的，且用于智能合约的验证。
正是因为在创建新容器要拉取新的docker镜像，所以导致新的节点在安装合约时和第一次查询时的等待时间特别长。

## 智能合约 参考  https://hyperledger-fabric.readthedocs.io/en/latest/chaincode.html
