# 存储与检索

> Wer Ordnung hält, ist nur zu faul zum Suchen.
> (把东西放规矩，就不用找了)
> —German proverb  

数据库本质上做两件事：插入数据时，保存；

查询时，返回。

如何存储？如何返回？

数据库优化：针对特定工作负载：需要了解存储引擎的底层机制

> 事务型工作负载和分析性负载的存储引擎在优化时非常不同。

### 存储引擎：

两大类存储引擎：日志结构(log-structured)存储引擎和面向页面(page-oriented)的存储引擎(例如B树)

## 驱动数据库的数据结构

世界上最简单的数据库可以用两个Bash函数实现  

> grep (global search regular expression(RE) and print out the line,全局搜索正则表达式并打印行)

```bash
#!/bin/bash
db_set () {
	echo "$1,$2" >> database
} 

b_get () {
	grep "^$1," database | sed -e "s/^$1,//" | tail -n 1
}
```

底层的存储格式非常简单： 一个文本文件， 每行包含一条逗号分隔的键值对（ 忽略转义问题
的话， 大致与CSV文件类似） 。 每次对 db_set 的调用都会向文件末尾追加记录， 所以更新
键的时候旧版本的值不会被覆盖 —— 因而查找最新值的时候， 需要找到文件中键最后一次出
现的位置（ 因此 db_get 中使用了 tail -n 1 。 )  

db_set 函数对于极其简单的场景其实有非常好的性能， 因为在文件尾部追加写入通常是非常
高效的。 与 db_set 做的事情类似， 许多数据库在内部使用了日志（ log） ， 也就是一个 仅追
加（ append-only） 的数据文件。  

真正的数据库有更多的问题需要处理（ 如并发控制， 回收
磁盘空间以避免日志无限增长， 处理错误与部分写入的记录） ， 但基本原理是一样的。 日志
极其有用， 我们还将在本书的其它部分重复见到它好几次  

> 日志（ log） 这个词通常指应用日志： 即应用程序输出的描述发生事情的文本。 本书在更普遍的意义下使用日志这一词： *一个仅追加的记录序列*。 它可能压根就不是给人类看
> 的， 使用二进制格式， 并仅能由其他程序读取。  

另一方面， 如果这个数据库中有着大量记录， 则这个 db_get 函数的性能会非常糟糕。 每次你
想查找一个键时， db_get 必须从头到尾扫描整个数据库文件来查找键的出现。 用算法的语言
来说， 查找的开销是 O(n) ： 如果数据库记录数量 n 翻了一倍， 查找时间也要翻一倍。 这就
不好了。  

### 哈希表索引

![img](https://github.com/Qasak/distributed-system/blob/master/%E6%95%B0%E6%8D%AE%E7%B3%BB%E7%BB%9F%E7%9A%84%E5%9F%BA%E7%9F%B3/%E6%95%B0%E6%8D%AE%E6%A8%A1%E5%9E%8B%E4%B8%8E%E6%9F%A5%E8%AF%A2%E8%AF%AD%E8%A8%80/hash_map0.png)

键值存储与在大多数编程语言中可以找到的字典（ dictionary） 类型非常相似， 通常字典都是
用散列映射（ hash map） （ 或哈希表（ hash table） ） 实现的。   

### python的字典实现

+ 开放寻址法(open addressing)

+ 解释：所有的元素都存放在散列表里，当产生哈希冲突时，通过一个探测函数计算出下一个候选位置，如果下一个获选位置还是有冲突，那么不断通过探测函数往下找，直到找个一个空槽来存放待插入元素。

+ 字典中的一个key-value键值对元素称为entry（也叫做slots），对应到Python内部是PyDictEntry

+ PyDictEntry：

  ```c
  typedef struct {
      /* Cached hash code of me_key.  Note that hash codes are C longs.
       * We have to use Py_ssize_t instead because dict_popitem() abuses
       * me_hash to hold a search finger.
       */
      Py_ssize_t me_hash;
      PyObject *me_key;
      PyObject *me_value;
  } PyDictEntry;
  ```

  

+ PyDictObject: PyDictEntry对象的集合

  ```c
  typedef struct _dictobject PyDictObject;
  struct _dictobject {
      PyObject_HEAD
      Py_ssize_t ma_fill;  /* # Active + # Dummy */
      Py_ssize_t ma_used;  /* # Active */
  
      /* The table contains ma_mask + 1 slots, and that's a power of 2.
       * We store the mask instead of the size because the mask is more
       * frequently needed.
       */
      Py_ssize_t ma_mask;
  
      /* ma_table points to ma_smalltable for small tables, else to
       * additional malloc'ed memory.  ma_table is never NULL!  This rule
       * saves repeated runtime null-tests in the workhorse getitem and
       * setitem calls.
       */
      PyDictEntry *ma_table;
      PyDictEntry *(*ma_lookup)(PyDictObject *mp, PyObject *key, long hash);
      PyDictEntry ma_smalltable[PyDict_MINSIZE];
  };
  ```

  

+ 其中的PyObject_HEAD:

  ```c
      Py_ssize_t ob_refcnt;
      struct _typeobject *ob_type;
  /*
  ob_refcnt，引用记数
  ob_type，类型对象的指针
  */
  ```

  

### 用hash map索引磁盘上的数据

假设我们的数据存储是一个追加写入(appending to)的文件

最简单的索引(indexing)策略是：keep an in-memory hash map where every key is mapped to a byte offset in the data file  —the location at which the value can be found  

Bitcask(Riak的默认存储引擎)就是这样实现的

像Bitcask这样的存储引 擎非常适合每个键的值频繁更新的场景。

实现中需要考虑的重要问题：

+ 文件格式
  + CSV不是日志的最佳格式，更快更简单的方法是二进制格式，首先以字节为单位来记录字符串长度，之后跟上原始字符串(不需要转义)
+ 删除记录
  + 如果要删除键和它关联的值，必须在数据文件中追加一个特殊的删除记录(有时称为墓碑)。合并日志段时，一旦发现墓碑标记，则丢弃这个已删除键的所有值
+ 崩溃恢复
  + 如果数据库重新启动，则内存中的hash map将丢失。可以从头到尾读整个段文件。但这样很慢，Bitcask将每个段的hash map快照存储在磁盘上，可以更快地加载到内存中，以此加快恢复速度。
+ 部分写入的记录
  + 数据库随时可能崩溃，包括将记录追加到日志的过程。Bitcask文件包括校验值，这样可发现损坏部分并丢弃
+ 并发控制
  + 写入以严格的先后顺序追加到日志，通常的实现是只有一个写线程。数据文件段是追加的，并且是不可变的，所以可以被多个线程同时读取。

追加的日志乍看很浪费空间，为什么不用原地更新，新值覆盖旧值？结果证明追加式设计非常不错：

+ 追加和分段合并主要是顺序写，它比随机写快得多
+ 如果段文件是追加的或不可变的，并发和崩溃恢复要简单得多。
+ 合并旧段可避免随时间推移数据文件出现碎片化问题。

哈希表索引的局限性：

+ 哈希表必须全放入内存，键太多时不行。
+ 区间查询效率不高。例如，不能简单地支持扫描kitty00000, kitty99999 区间
  内的所有键，只能采用逐个查找的方式查询每一个键  

## SSTables和LSM-Tree

排序字符串表(sorted string table)：SSTable

