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


## 打包智能合约



### 在 foo27 上打包

```shell

 peer chaincode package -n chentaocc -p github.com/hyperledger/fabric/chaincode/go/chaincode_example02 -v 0 -s -S -i "AND('OrgA.admin')" ccpack.out

```



### 在 peer0  上再一次对包签名

```shell
peer  chaincode signpackage ./channel-artifacts/ccpack.out ./channel-artifacts/signedccpack.out

```


### 在 barOrg1 节点上安装

在 cli 容器中用如下命令来安装签名过的只能合约

```shell
 peer chaincode install -n chentaocc -v 1.0 -p ./channel-artifacts/signedccpack.out
```

直接指定当前目录下签名的CDS包，得到如下错误

```
2019-11-15 06:22:12.456 UTC [msp] GetDefaultSigningIdentity -> DEBU 043 Obtaining default signing identity
2019-11-15 06:22:12.456 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 044 Using default escc
2019-11-15 06:22:12.456 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 045 Using default vscc
Error: error getting chaincode code chentaocc: path to chaincode does not exist: /opt/gopath/src/channel-artifacts/signedccpack.out

```

按照书中所写，“它必须位于用户的GOPATH的源码树中，　例如$GOPATH/src/sacc”

#### 是要把打包后的signedccpack.out部署到 cli容器 还是 peer节点的容器 GOPATH中？

##### 尝试放置到 cli容器中

```
cp ./channel-artifacts/signedccpack.out /opt/gopath/src/

peer chaincode install -n chentaocc -v 1.0 -p signedccpack.out

2019-11-15 06:42:29.708 UTC [msp] GetDefaultSigningIdentity -> DEBU 043 Obtaining default signing identity
2019-11-15 06:42:29.709 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 044 Using default escc
2019-11-15 06:42:29.709 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 045 Using default vscc
2019-11-15 06:42:29.817 UTC [chaincode.platform.golang] getCodeFromFS -> DEBU 046 getCodeFromFS signedccpack.out
Error: error getting chaincode code chentaocc: error getting chaincode package bytes: Error getting code code does not exist File /opt/gopath/src/signedccpack.out is not dir

```

仍然提示类似的错误



##### 尝试放置到 peer 容器中

仍然获得错误


```

peer chaincode install -n chentaocc -v 1.0 -p signedccpack.out

2019-11-15 07:06:03.045 UTC [msp] GetDefaultSigningIdentity -> DEBU 043 Obtaining default signing identity
2019-11-15 07:06:03.045 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 044 Using default escc
2019-11-15 07:06:03.045 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 045 Using default vscc
2019-11-15 07:06:03.186 UTC [chaincode.platform.golang] getCodeFromFS -> DEBU 046 getCodeFromFS signedccpack.out
Error: error getting chaincode code chentaocc: error getting chaincode package bytes: Error getting code code does not exist File /opt/gopath/src/signedccpack.out is not dir

```

#### 直接在 cli 容器中安装成功：   peer chaincode install ./channel-artifacts/signedccpack.out

```

 peer chaincode list --installed

 # get response
 2019-11-15 08:02:26.589 UTC [msp.identity] Sign -> DEBU 045 Sign: digest: 7F0EE977EF810723452FE1EDAEB7EB685C430A2E905B93F9A1FDF1F2242CD3FE
Get installed chaincodes on peer:
Name: chentaocc, Version: 0, Path: github.com/hyperledger/fabric/chaincode/go/chaincode_example02, Id: fec2908042ff87465f7ab0d7448d0e7a7f9a31742bd05d61d5b70795d2015383

# join channel
peer channel join -b ./channel-artifacts/chentaochannel.block

# query
peer chaincode query -C chentaochannel -n chentaocc -c '{"Args":["query","A"]}'

# get query error

Error: endorsement failure during query. response: status:500 message:"cannot retrieve package for chaincode chentaocc/1.0, error open /var/hyperledger/production/chaincodes/chentaocc.1.0: no such file or directory"

```

##### 上面的错误, 是不是在安装的时候没有指定名称和版本引起的？

```
root@bea1dd05b42d:/var/hyperledger/production/chaincodes# ls
chentaocc.0
```

从 peer 的容器目录下我们查看 可以看到 只有 chentaocc.0 这个合约，而 没有chentaocc.1.0 ， 所以是不是我们在安装CDS打包后的合约，要指定版本呢？


###### 指定版本重试：

```
peer chaincode install -v 1.0 ./channel-artifacts/signedccpack.out

2019-11-15 08:21:17.542 UTC [grpc] HandleSubConnStateChange -> DEBU 042 pickfirstBalancer: HandleSubConnStateChange: 0xc0000a38d0, READY
2019-11-15 08:21:17.542 UTC [msp] GetDefaultSigningIdentity -> DEBU 043 Obtaining default signing identity
Error: chaincode version 1.0 does not match version 0 in packages


```

提示指定的版本与包中的版本不匹配。

###### 删除 peer 容器中的智能合约，再安装试试：

与上一步的错误仍然一样


###### 我们删除节点的所有容器，重新操作一遍试试 （）

与上一步的错误仍然一样

说明用 peer chaincode  package 打包的 CDS 需要指定版本


##### 重新带版本号打包


```
Usage:
  peer chaincode package [flags]

Flags:
  -s, --cc-package                  create CC deployment spec for owner endorsements instead of raw CC deployment spec
  -c, --ctor string                 Constructor message for the chaincode in JSON format (default "{}")
  -h, --help                        help for package
  -i, --instantiate-policy string   instantiation policy for the chaincode
  -l, --lang string                 Language the chaincode is written in (default "golang")
  -n, --name string                 Name of the chaincode
  -p, --path string                 Path to chaincode
  -S, --sign                        if creating CC deployment spec package for owner endorsements, also sign it with local MSP
  -v, --version string              Version of the chaincode specified in install/instantiate/upgrade commands

Global Flags:
      --cafile string                       Path to file containing PEM-encoded trusted certificate(s) for the ordering endpoint
      --certfile string                     Path to file containing PEM-encoded X509 public key to use for mutual TLS communication with the orderer endpoint
      --clientauth                          Use mutual TLS when communicating with the orderer endpoint
      --connTimeout duration                Timeout for client to connect (default 3s)
      --keyfile string                      Path to file containing PEM-encoded private key to use for mutual TLS communication with the orderer endpoint
  -o, --orderer string                      Ordering service endpoint
      --ordererTLSHostnameOverride string   The hostname override to use when validating the TLS connection to the orderer.
      --tls                                 Use TLS when communicating with the orderer endpoint
      --transient string                    Transient map of arguments in JSON encoding


```

我们看到package命令对 -v 参数的描述， 它是与 install/instantiate/upgrade 命令中指定的版本相匹配的，所以我们在打包的是偶要指定一个版本号

```

# foo27Org2 打包
peer chaincode package -n chentaocc -p github.com/hyperledger/fabric/chaincode/go/chaincode_example02 -v 1.0 -s -S -i "AND('OrgA.admin')" ./channel-artifacts/chentaocc.package

# peer0Org1 签名
peer  chaincode signpackage ./channel-artifacts/chentaocc.package ./channel-artifacts/chentaocc.package.signed

# barOrg2 安装，操作成功
peer chaincode install ./channel-artifacts/chentaocc.package.signed

# barOrg2 加入频道
 peer channel join -b ./channel-artifacts/chentaochannel.block

# barOrg2 查询
peer chaincode query -C chentaochannel -n chentaocc -c '{"Args":["query","A"]}'

# 出错了
Error: endorsement failure during query. response: status:500 message:"Instantiation policy mismatch for cc chentaocc/1.0"


```

###### 出错的原因是不是因为 打包的 -i 参数指定的  "AND('OrgA.admin')" 实例化 背书策略失败引起的？

因为 configtx.yaml 和  crypto-config.yaml 中没有 OrgA 的相关配置


 peer chaincode package -n chentaocc -p github.com/hyperledger/fabric/chaincode/go/chaincode_example02 -v 1.0 -s -S -i "AND('Org2.admin')" ./channel-artifacts/chentaocc.package


bar 属于 Org2 所以， -i 参数指定用 Org2 的amdin来实例化


peer  chaincode signpackage ./channel-artifacts/chentaocc.package ./channel-artifacts/chentaocc.package.signed

签名


重新安装CDS

```
peer chaincode install ./channel-artifacts/chentaocc.package.signed

peer channel join -b ./channel-artifacts/chentaochannel.block

peer chaincode query -C chentaochannel -n chentaocc -c '{"Args":["query","A"]}'


Error: endorsement failure during query. response: status:500 message:"Instantiation policy mismatch for cc chentaocc/1.0"


```

仍然错误, "Instantiation policy mismatch for cc chentaocc/1.0"  (实列化策略不匹配), 即新的peer在第一次查询时需要实例化(instantiate)刚安装智能合约
但是只能合约在打包时,指定了实列化策略,而当前节点不符合该策略,那么我们用  -i "AND('Org2.member')"  来测试一下, 原因是 当前节点不是 admin.


```

# peer0Org1 打包

peer chaincode package -n chentaocc -p github.com/hyperledger/fabric/chaincode/go/chaincode_example02 -v 1.0 -s -S -i "AND('Org2.member')" ./channel-artifacts/chentaocc.package

# foo27Org2 签名

peer  chaincode signpackage ./channel-artifacts/chentaocc.package ./channel-artifacts/chentaocc.package.signed

# barOrg2 执行新节点的所有操作

peer chaincode install ./channel-artifacts/chentaocc.package.signed

peer channel join -b ./channel-artifacts/chentaochannel.block

peer chaincode query -C chentaochannel -n chentaocc -c '{"Args":["query","A"]}'


Error: endorsement failure during query. response: status:500 message:"Instantiation policy mismatch for cc chentaocc/1.0"


```


查看代码

```
#core/common/ccprovider/ccprovider.go

func CheckInstantiationPolicy(name, version string, cdLedger *ChaincodeData) error {
	ccdata, err := GetChaincodeData(name, version)
	if err != nil {
		return err
	}

	// we have the info from the fs, check that the policy
	// matches the one on the file system if one was specified;
	// this check is required because the admin of this peer
	// might have specified instantiation policies for their
	// chaincode, for example to make sure that the chaincode
	// is only instantiated on certain channels; a malicious
	// peer on the other hand might have created a deploy
	// transaction that attempts to bypass the instantiation
	// policy. This check is there to ensure that this will not
	// happen, i.e. that the peer will refuse to invoke the
	// chaincode under these conditions. More info on
	// https://jira.hyperledger.org/browse/FAB-3156
	if ccdata.InstantiationPolicy != nil {
		if !bytes.Equal(ccdata.InstantiationPolicy, cdLedger.InstantiationPolicy) {
			return fmt.Errorf("Instantiation policy mismatch for cc %s/%s", name, version)
		}
	}

	return nil
}


```


!bytes.Equal(ccdata.InstantiationPolicy, cdLedger.InstantiationPolicy)

此处实例化策略是与分类账本原有的实例化策略做的对比, 那么我们在打包的时候-i 参数指定的策略也应该与 最初的实例化策略一致才能通过,下面测试:

```

# peer0Org1 打包

peer chaincode package -n chentaocc -p github.com/hyperledger/fabric/chaincode/go/chaincode_example02 -v 1.0 -s -S -i  "OR ('Org1MSP.member','Org2MSP.member')" ./channel-artifacts/chentaocc.package

# foo27Org2 签名

peer  chaincode signpackage ./channel-artifacts/chentaocc.package ./channel-artifacts/chentaocc.package.signed

# barOrg2 执行新节点的所有操作

peer chaincode install ./channel-artifacts/chentaocc.package.signed
#  peer chaincode install -n chentaocc -p github.com/hyperledger/fabric/chaincode/go/chaincode_example02 -v 1.0

peer channel join -b ./channel-artifacts/chentaochannel.block

peer chaincode query -C chentaochannel -n chentaocc -c '{"Args":["query","A"]}'


Error: endorsement failure during query. response: status:500 message:"Instantiation policy mismatch for cc chentaocc/1.0"


# 实例化出错, 创建者的MSP(Org2MSP)与实例化的MSP(Org2MSP)不一致
peer chaincode instantiate -o orderer0.7shu.co:7050 -C chentaochannel -n chentaocc -c '{"Args":["init","A","10","B","15"]}' -P "OR ('Org1MSP.member','Org2MSP.member')" -v 1.0

2019-11-15 10:36:12.297 UTC [msp.identity] Sign -> DEBU 04c Sign: plaintext: 0AA0070A6C08031A0C089CFFB9EE0510...324D53500A04657363630A0476736363
2019-11-15 10:36:12.297 UTC [msp.identity] Sign -> DEBU 04d Sign: digest: 7F6B4AC669C6A87D07FAE87FE9C9E080E2DEFCAC8ADA0C043AE39BEBB4D72370
Error: error endorsing chaincode: rpc error: code = Unknown desc = access denied: channel [chentaochannel] creator org [Org2MSP]


```


查看日志,在查询过程中提示 admin 的 tls  证书目录不存在,继续调查


## 智能合约本地开发环境的搭建


### 克隆　fabric-samples 仓库

```
git clone https://github.com/hyperledger/fabric-samples.git

# 切换到docker对应版本的分支下
git checkout v1.4.4  
```

.
├── balance-transfer
├── basic-network
├── bin
├── chaincode
├── chaincode-docker-devmode
├── ci
├── ci.properties
├── CODE_OF_CONDUCT.md
├── commercial-paper
├── config
├── CONTRIBUTING.md
├── docs
├── fabcar
├── first-network
├── high-throughput
├── interest_rate_swaps
├── Jenkinsfile
├── LICENSE
├── MAINTAINERS.md
├── off_chain_data
├── README.md
├── scripts
└── SECURITY.md



### 下载对应版本编译好的二进制文件

```

wget -c https://github.com/hyperledger/fabric/blob/v1.4.4/scripts/bootstrap.sh

bash ./bootstrap.sh -s 1.4.4


```

.
├── bin
│   ├── configtxgen
│   ├── configtxlator
│   ├── cryptogen
│   ├── discover
│   ├── idemixgen
│   ├── orderer
│   └── peer
└── config
    ├── configtx.yaml
    ├── core.yaml
    └── orderer.yaml



得到 hyperledger-fabric-linux-amd64-1.4.4.tar.gz, 解压缩后的内容拷贝到， fabric-samples 仓库的目录下


### 启动开发环境的docker


进入fabric-samples仓库的 chaincode-docker-devmode 目录


```
cd chaincode-docker-devmode
```

目录结构如下

.
├── chaincode
├── docker-compose-simple.yaml
├── msp
│   ├── admincerts
│   │   └── admincert.pem
│   ├── cacerts
│   │   └── cacert.pem
│   ├── keystore
│   │   └── key.pem
│   ├── signcerts
│   │   └── peer.pem
│   ├── tlscacerts
│   │   └── tlsroot.pem
│   └── tlsintermediatecerts
│       └── tlsintermediate.pem
├── myc.block
├── myc.tx
├── orderer.block
├── README.rst
└── script.sh

注意 myc.tx 和 myc.block 文件，对应的channel 名称为  myc ， 在后续安装和实例化智能合约时指定的channel 名字必须是 myc


```
docker-compose -f ./docker-compose-simple.yaml  up -d

[root@localhost chaincode-docker-devmode]# docker ps
CONTAINER ID        IMAGE                        COMMAND                  CREATED             STATUS              PORTS                                            NAMES
0af358c383f1        hyperledger/fabric-ccenv     "/bin/sh -c 'sleep 6…"   50 minutes ago      Up 50 minutes                                                        chaincode
552b24ad6fe6        hyperledger/fabric-tools     "/bin/bash -c ./scri…"   50 minutes ago      Up 50 minutes                                                        cli
2c167b34df85        hyperledger/fabric-peer      "peer node start --p…"   50 minutes ago      Up 50 minutes       0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp   peer
e415c523d336        hyperledger/fabric-orderer   "orderer"                50 minutes ago      Up 50 minutes       0.0.0.0:7050->7050/tcp                           orderer

```

chaincode 容器用于编译测试

cli 容器用于安装、实例化、查询智能合约


### 开发测试智能合约


确认peer节点没有任何智能合约

```

 docker exec -it peer bash


root@d8a94ed6aaab:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls /var/hyperledger/production/chaincodes/ -al
total 0
drwxr-xr-x 2 root root  6 Nov 28 09:26 .
drwxr-xr-x 5 root root 65 Nov 28 09:26 ..
root@d8a94ed6aaab:/opt/gopath/src/github.com/hyperledger/fabric/peer#

```

chaincode 节点编译，初始智能合约

```

# 进入容器
 docker exec -it chaincode bash

# 查看当前目录结构，所有的可测试的智能合约就在当前目录下，如果自己开发的话，也要将内容放置到当前目录下
root@1e6051626195:/opt/gopath/src/chaincode# ls -al
total 8
drwxrwxrwx 1 1000 1000 4096 Nov 28 07:48 .
drwxr-xr-x 4 root root   41 Nov 28 09:26 ..
drwxrwxrwx 1 1000 1000    0 Nov 28 07:39 abac
drwxrwxrwx 1 1000 1000    0 Nov 28 07:48 chaincode_example02
drwxrwxrwx 1 1000 1000 4096 Nov 28 07:48 fabcar
drwxrwxrwx 1 1000 1000    0 Nov 28 07:48 marbles02
drwxrwxrwx 1 1000 1000    0 Nov 28 07:48 marbles02_private
drwxrwxrwx 1 1000 1000    0 Nov 28 08:37 sacc
drwxrwxrwx 1 1000 1000    0 Nov 28 08:37 sacc

# 进入对应的智能合约目录下并编译
root@1e6051626195:/opt/gopath/src/chaincode# cd sacc
root@1e6051626195:/opt/gopath/src/chaincode/sacc# go build
root@1e6051626195:/opt/gopath/src/chaincode/sacc# ls -al
total 22944
drwxrwxrwx 1 1000 1000        0 Nov 28 09:34 .
drwxrwxrwx 1 1000 1000     4096 Nov 28 07:48 ..
-rwxrwxrwx 1 1000 1000 23484336 Nov 28 09:34 sacc
-rwxrwxrwx 1 1000 1000     2851 Nov 28 07:48 sacc.go


# 运行智能合约, 此处启动完，不能退出，因为在cli容器操作查询时要跟此容器的进程交互才能完成查询功能

root@1e6051626195:/opt/gopath/src/chaincode/sacc# CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=sacc:0 ./sacc
2019-11-28 09:47:56.215 UTC [shim] setupChaincodeLogging -> INFO 001 Chaincode log level not provided; defaulting to: INFO
2019-11-28 09:47:56.218 UTC [shim] setupChaincodeLogging -> INFO 002 Chaincode (build level: ) starting up ...
2019-11-28 09:47:56.222 UTC [bccsp] initBCCSP -> DEBU 001 Initialize BCCSP [SW]
2019-11-28 09:47:56.229 UTC [grpc] DialContext -> DEBU 002 parsed scheme: ""
2019-11-28 09:47:56.230 UTC [grpc] DialContext -> DEBU 003 scheme "" not registered, fallback to default scheme
2019-11-28 09:47:56.233 UTC [grpc] watcher -> DEBU 004 ccResolverWrapper: sending new addresses to cc: [{peer:7052 0  <nil>}]
2019-11-28 09:47:56.234 UTC [grpc] switchBalancer -> DEBU 005 ClientConn switching balancer to "pick_first"
2019-11-28 09:47:56.238 UTC [grpc] HandleSubConnStateChange -> DEBU 006 pickfirstBalancer: HandleSubConnStateChange: 0xc0003d4ad0, CONNECTING
2019-11-28 09:47:56.299 UTC [grpc] HandleSubConnStateChange -> DEBU 007 pickfirstBalancer: HandleSubConnStateChange: 0xc0003d4ad0, READY



```

此处启动完，不能退出，因为在cli容器操作查询时要跟此容器的进程交互才能完成查询功能

cli 容器安装测试

```
docker exec -it cli bash

root@cb5cf1bd622a:/opt/gopath/src/chaincodedev# ls
chaincode  docker-compose-simple.yaml  msp  myc.block  myc.tx  orderer.block  README.rst  script.sh


# 安装
peer chaincode install -p chaincodedev/chaincode/sacc -n sacc -v 0


2019-11-28 10:00:57.349 UTC [container] WriteFileToPackage -> DEBU 04a Writing file to tarball: src/chaincodedev/chaincode/sacc/sacc.go
2019-11-28 10:00:57.355 UTC [msp.identity] Sign -> DEBU 04b Sign: plaintext: 0AC4070A5C08031A0C08D9B5FEEE0510...7719FF040000FFFF996692F900120000
2019-11-28 10:00:57.355 UTC [msp.identity] Sign -> DEBU 04c Sign: digest: A5031D5B6B27C282F9753A481C68412F997513C620CD4ECB2EEFC61F8B44EFC2
2019-11-28 10:00:57.361 UTC [chaincodeCmd] install -> INFO 04d Installed remotely response:<status:200 payload:"OK" >

# peer 节点中已经
#
#root@d8a94ed6aaab:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls /var/hyperledger/production/chaincodes/ -al
#total 8
#drwxr-xr-x 2 root root   34 Nov 28 10:00 .
#drwxr-xr-x 5 root root   65 Nov 28 09:26 ..
#-rw-r--r-- 1 root root 1244 Nov 28 09:59 sacc.0
#-rw-r--r-- 1 root root 1244 Nov 28 10:00 sacc.1


# 实例化

peer chaincode instantiate -n sacc -v 0 -c '{"Args":["a","10"]}' -C myc

peer chaincode query -n sacc -c '{"Args":["query","a"]}' -C myc

peer chaincode invoke -n sacc -c '{"Args":["set","a", "200"]}' -C myc


```