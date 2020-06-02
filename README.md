# all-about-mit-6.824

> 仅供个人学习的记录使用,请独立完成所有lab!
>
> All solutions written by github user [Qasak](https://qasak.github.io/).    

## lab1

MapReduce

实现一个分布式MapReduce，它由两个程序组成，master和worker。只有一个master进程，一个或多个worker进程并行执行。在现实的系统中，workers会在一堆不同的机器上运行，但是对于这个lab，你会在一台机器上运行他们。workers将通过RPC与masker交互。每个worker进程将向masker进程请求任务，从一个或多个文件中读取任务的输入，执行任务，并将任务的输出写入一个或多个文件。如果一个worker在一段合理的时间内没有完成任务（对于这个lab，用10秒），masker应该感知，并把相同的任务交给另一个worker。



## lab2

replication for fault-tolerance using Raft

## lab3

fault-tolerant key/value store

## lab4

sharded key/value store