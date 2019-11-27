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