## intro
a demo to explain MySQL CDC(Capture Data Change) with canal and golang sdk

## 配置 MySQL
先检查对应的配置项是否正确设置：
```sql
show variables like 'log_bin';        # 要求是ON
show variables like 'binlog_format';  # 要求是ROW
```

## 给Canal创建 MySQL 账号并授权
```sql
create user canal identified by 'canal';
grant select, replication slave, replication client on *.* TO 'canal'@'%';
flush privileges ;
```

## 使用docker-compose启动数据库和canal
```shell
make network # create network 
make up      # start 
# make down  # stop 
```


## 运行go程序
```shell
make run
```

## bug排查
进入canal的docker容器，看里面的logs