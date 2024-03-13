create database if not exists testdb;

use testdb;

create table if not exists test(
    id int primary key auto_increment not null,
    `key` varchar(32) default '' not null,
    val int default 0 not null
)engine='innodb' charset='utf8mb4';

insert into test(`key`, val) values("111", 111), ("222", 222), ("2333", 2333);


create user canal identified by 'canal';
grant select, replication slave, replication client on *.* TO 'canal'@'%';
flush privileges ;