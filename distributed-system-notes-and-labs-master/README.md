# all-about-mit-6.824

> 仅供个人学习的记录使用,请独立完成所有lab!
>
> All solutions written by github user [Qasak](https://qasak.github.io/).    

## lab1

[MapReduce](https://github.com/Qasak/distributed-system/blob/master/lab1/README.md)

实现一个分布式MapReduce，它由两个程序组成，master和worker。

只有一个master进程，一个或多个worker进程并行执行。

在现实的系统中，workers会在一堆不同的机器上运行，但是对于这个lab，你会在一台机器上运行他们。

workers将通过RPC与masker交互。每个worker进程将向masker进程请求任务，从一个或多个文件中读取任务的输入，执行任务，并将任务的输出写入一个或多个文件。

如果一个worker在一段合理的时间内没有完成任务（对于这个lab，用10秒），masker应该感知，并把相同的任务交给另一个worker。



## lab2

replication for fault-tolerance using Raft

第一个构建容错k/v存储系统的lab.

Raft:一个复制状态机的协议

在下一个lab中，你将基于Raft构建一个k/v服务。

然后，你将在多个复制状态机上“切分”("shard")你的服务，以获得更高性能

---

复制服务通过在多个复制服务器上存储其状态(即数据)的完整副本来实现容错

复制允许服务继续运行，即使某些服务器出现故障（崩溃/网络中断/不稳定）

挑战在于：上述失败可能会导致副本保存不同的数据副本

---

Raft将客户端请求组织成一个成为log的序列，并确保所有副本服务器都看到相同的日志

每个副本按日志顺序执行客户端请求，并将它们应用于服务状态的本地副本

由于所有活动副本都看到相同的日志内容，所以它们都以相同的顺序执行相同的请求，因此继续具有相同的服务状态

如果服务器出现故障，但稍后恢复，Raft会负责更新其日志

只要至少大多数服务器都是活动的，并且可以相互通信，Raft将继续运行

如果没有足够多，Raft将不会有任何进展，但会在多数服务器能够再次通信后，尽快找到它中断的地方。

---

在这个lab中，您将把Raft作为一个Go对象类型和相关的方法实现，以便在更大的服务中用作一个模块

一组Raft实例通过RPC相互通信来维护复制的日志。

你的Raft接口将支持不确定的编号命令序列，也称为日志条目(log entries)

条目用索引号编号。具有给定索引的日志条目最终将被提交。

此时，您的Raft应该将日志条目发送到更大的服务器，以便执行。

---

您应该遵循扩展[Raft Paper](https://pdos.csail.mit.edu/6.824/papers/raft-extended.pdf)中的设计，特别注意图2。您将实现本文中的大部分内容，包括保存持久状态，并在节点失败后读取它，然后重新启动。

您将不会实现群集成员更改（第6节）。您将在稍后的实验室中实现日志压缩/快照(log compaction/snapshotting)（第7节）

---

您可能会发现本指南很有用，以及关于锁和并发结构的建议。

从更广泛的角度来看，可以看看Paxos、Chubby、Paxos Made Live、Spanner、Zookeeper、Harp、Viewstamped Replication,[Bolosky](http://static.usenix.org/event/nsdi11/tech/full_papers/Bolosky.pdf)。

## lab3

fault-tolerant key/value store

## lab4

sharded key/value store