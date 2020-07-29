## Go

+ 好用的RPC包

+ 内存安全

  + 垃圾回收

+ 类型安全

  

## 线程

+ goroutine
  + 通过channel同步
  + 也可使用mutex

+ I/O并发
+ 并行

> 异步(事件驱动)：只有一个循环，单核，单线程，高性能
>
> I/O多路复用是事件驱动的基础

### 线程挑战

+ 竞争Race
  + 用mutex解决
+ 协作Coordination
  + 用channel实现
  + sync.cond
  + waitGroup
+ 死锁Deadlock

