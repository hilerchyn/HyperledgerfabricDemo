# ansible


## 准备

将控制主机的 .ssh/id_rsa.pub 内容添加到远程主机的 .ssh/authorized_keys 中，实现免密登录

##  基本操作

测试链接状态

```
ansible -i ./ansible all -u vagrant -m ping
```

运行命令

```
ansible -i ./ansible all -u vagrant -a "docker ps" --become
```


## 启动集群中的Docker

```
// 启动
 ansible -i ./ansible zookeeper -u vagrant -a "/usr/local/bin/docker-compose -f /home/vagrant/docker-compose.yaml up -d" --become
 ansible -i ./ansible kafka -u vagrant -a "/usr/local/bin/docker-compose -f /home/vagrant/docker-compose.yaml up -d" --become
 ansible -i ./ansible orderer -u vagrant -a "/usr/local/bin/docker-compose -f /home/vagrant/docker-compose.yaml up -d" --become
 ansible -i ./ansible activepeers -u vagrant -a "/usr/local/bin/docker-compose -f /home/vagrant/docker-compose.yaml up -d" --become

// 停止
 ansible -i ./ansible all -u vagrant -a "/usr/local/bin/docker-compose -f /home/vagrant/docker-compose.yaml stop" --become

// 停止并删除
 ansible -i ./ansible all -u vagrant -a "/usr/local/bin/docker-compose -f /home/vagrant/docker-compose.yaml rm -sf" --become
```