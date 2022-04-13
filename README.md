# 使用 go 语言编写的数据库备份工具

## 功能点

* 可以对 linux 端 mysql 和 postgresql 数据库进行压缩备份（后续可能会支持更多数据库类型）
* 可以发送到异机（windows 需要开启 ssh 服务，用户需要设置密码）
* 可以设置删除 n 天前的旧备份（要查出所有备份并根据时间排序，保留最新的 n 份。因为一天的备份不可同时存在多份，所以不会出现所有的最新n份为同一天备份的情况）
* 删除旧备份可开启删除异机备份
* 会记录备份日志和异常日志
* 支持 amd64 和 arm64 构架

## 使用方式

* 在 config.json 配置文件中做相关配置，内有注释说明。在 dbs.txt 中写数据库名称，一行一个。当正向筛选时，会备份这里写到的所有数据库，当反向筛选时，将备份除这里的其他所有数据库。
* 配置完成后，放到 db_backup_go_xxx64 执行文件同目录，直接执行 db_backup_go_xxx64 即可

**注意：本软件不实现数据备份的主体功能，备份依赖于 mysqldump 和 pg_dump。如这些工具本身使用存在报错，那本程序也会报错**