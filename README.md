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



