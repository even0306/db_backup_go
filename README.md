# 使用go语言编写的数据库备份工具

* 可以对linux端mysql数据库进行备份并压缩（后续可能会支持更多数据库类型）
* 可以发送到异机（windows需要开启ssh服务，用户需要设置密码）
* 可以设置删除n天前的旧备份（要查出所有备份并根据时间排序，保留最新的n份。因为一天的备份不可同时存在多份，所以不会出现所有的最新n份为同一天备份的情况）
* 删除旧备份可开启删除异机备份
* 会记录备份日志和异常日志
* 支持amd64和arm64构架