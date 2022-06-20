# Big Storage System

## GFS

+ 专为顺序的访问大文件设计

  而不是需要随机访问的小文件

+ 允许弱一致性

  它对数据容错能力不需要像银行系统那样高

 ### Master Data

+ 两个表
  + filename -> chunk handle数组(non volatile ， 反射(reflected)到磁盘上)

  + handle  

    ​		-> chunk服务器列表(volatile，不用写入磁盘)

    ​			chunk的版本号(non volatile，写入磁盘)

    ​				版本号只在master认为没有primary时才增加

    ​			哪个是primary(volatile)

    ​			租约到期(lease expiration)(volatile)

日志和检查点在磁盘上

需要追加日志，写入磁盘的情况：

+ 任何文件达到64MB边界，需要创建一个新chunk
+ 由于新primary被指派，需要改变版本号

> 为什么用日志而不是数据库？
>
> 数据库用b-tree
>
> 增加(appending)日志会让这种操作快一点

### read

1. 客户端把文件名和偏移发给master

2. master发回handle和chunk服务器列表

   ​	客户端会尝试最近的服务器

   ​	客户端会缓存已发送的请求

3. 客户端和chunk服务器通信，发送一个handle和偏移，chunk服务器返回文件

### write

+ 没有primary

  ​	master选择一个chunk服务器作为primary，其他为secondary

  ​	master将版本号递增，并将其写入磁盘

  ​	告诉chunk服务器谁是primary，谁是secondary

  ​	primary和secondary将信息写入各自的磁盘，若系统重启，他们必须拿这个版本号向master报告

  > 若系统重启后，master发现报告的版本号高于自己的版本号，则会采用该版本号(说明master在分配版本号后，持久化之前崩了)

  ​	master给primary一个租约，60秒内它是primary(这个机制保证不会有两个primary)

  ​	master告诉client谁是primary，谁是secondary

  ​	client将要追加的数据以某种顺序给primary和所有secondaries

  ​	primary和所有secondaries将该数据写入一个临时位置

  ​	primary和所有secondaries都收到数据并回复yes

  ​	client发送一条消息给primary，说：你和所有secondaries节点都收到了我要追加的数据

  ​	primary会以某种顺序（假设同时有多条client请求），将数据追加到文件

  ​	确保块有足够空间，primary将client的数据写入当前块末尾，并告诉所有secondary也将其写入末尾相同偏移量处

  ​	如果secondary回复primary all ‘yes’，primary回复client 'success'

  ​		如果有一个没有回复'yes'，primary回复client ‘no’

  ​			如果某个secondary挂了，master可能会定义新的primary和secondary并增加版本号

  ​	若client收到‘no’，则应重新发起整个append过程

  ​		若client因某个secondary挂掉而收到'no'，并且在重发append请求之前挂掉

+ 有primary

### "spilit brain"

一种网络分区故障，master不能和primary通信，但primary可以和其他client通信

若此时简单地重新设置一个primary，就会导致出现两个primary

为避免为同一个chunk错误地指定两个primary，master在指定一个primary时建立一个租约(lease)，primary和master都知道租约的事件，若到期，primary就停止执行client的请求

如果master不能和primary建立联系，则会等待租约到期后，指定一个新的primary



### 强一致系统

两阶段提交

+ primary询问secondary是否有能力处理请求
+ 得到肯定答复后primary发送请求让secondary执行