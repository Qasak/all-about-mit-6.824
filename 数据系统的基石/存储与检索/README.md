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

相比散列索引的日志段，sstable有如下优势：

1. 合并段简单高效：类似归并排序

   有相同键时，保留最新段的值

2. 找特定键时，不需要保存所有键的索引(因为是有序的)假设你正在内存中寻找键 handiwork ， 但是你不知道段文件中该关键字的确切偏移量。
   然而， 你知道 handbag 和 handsome 的偏移， 而且由于排序特性， 你知道 handiwork
   必须出现在这两者之间。 这意味着您可以跳到 handbag 的偏移位置并从那里扫描， 直到
   您找到 handiwork （ 或没找到， 如果该文件中没有该键） 。  

3. 由于读请求往往需要扫描请求范围内的 多个key value对，可以考虑将这些记录保存到一个块中并在写磁盘之前将其压缩（如图 3-51二j:i 阴影区域所示）。然后稀 疏内存索引的每个条目指向压缩块的开头。除了节省磁盘空间，压缩还减少了 I/O带宽的占用。  

### 构建和维护SSTable

在磁盘上维护排序结构是可行的（参阅本章后面的“ B-trees”）  不过，将其保存在内存中更容易。内存排序有很多广为人知的树状数据结构，例如红黑树或AVL树［2］。使用这些数据结构，可以按任意顺序插入键并以排序后的顺序读取它们 。  

存储引擎的基本工作流程：

+ 当写入时，将其添加到内存中的平衡树数据结构中（｛列如红黑树）。这个内存中
  的树有时被称为内存表。  
+ 当内存表大于某个阈值（通常为几兆字节）时，将其作为 SSTable文件写入磁
  盘。由于树已经维护了按键排序的 key - value对， 写磁盘可以比较高效。新的
  SSTable文件成为数据库的最新部分。当 SSTable写磁盘的同时 ，写入可以继续添
  加到一个新的内存表实例 。  
+ 为了提供读取请求， 首先尝试在内存表中找到关键字， 然后在最近的磁盘段中， 然后在
  下一个较旧的段中找到该关键字  
+ 有时会在后台运行合并和压缩过程以组合段文件并丢弃覆盖或删除的值。  

这个方案效果很好。 它只会遇到一个问题： 如果数据库崩溃， 则最近的写入（ 在内存表中，
但尚未写入磁盘） 将丢失。 为了避免这个问题， 我们可以在磁盘上保存一个单独的日志， 每
个写入都会立即被附加到磁盘上， 就像在前一节中一样。 该日志不是按排序顺序， 但这并不
重要， 因为它的唯一目的是在崩溃后恢复内存表。 每当内存表写出到SSTable时， 相应的日志
都可以被丢弃。  

### 用SSTables制作LSM树

LSM树：Log-structured Merge-Tree日志结构合并树

这里描述的算法本质上是LevelDB 【6】 和RocksDB 【7】 中使用的关键值存储引擎库， 被设
计嵌入到其他应用程序中。 除此之外， LevelDB可以在Riak中用作Bitcask的替代品。 在
Cassandra和HBase中使用了类似的存储引擎【8】 ， 这两种引擎都受到了Google的Bigtable
文档【9】 （ 引入了SSTable和memtable） 的启发  

最初这种索引结构是由Patrick O'Neil等人描述的。 在日志结构合并树（ 或LSM树） 【10】 的
基础上， 建立在以前的工作上日志结构的文件系统【11】 。 基于这种合并和压缩排序文件原
理的存储引擎通常被称为LSM存储引擎。  

Lucene是Elasticsearch和Solr使用的一种全文搜索的索引引擎， 它使用类似的方法来存储它
的词典【12,13】 。 全文索引比键值索引复杂得多， 但是基于类似的想法： 在搜索查询中给出
一个单词， 找到提及单词的所有文档（ 网页， 产品描述等） 。 这是通过键值结构实现的， 其
中键是单词（ 关键词（ term） ） ， 值是包含单词（ 文章列表） 的所有文档的ID的列表。 在
Lucene中， 从术语到发布列表的这种映射保存在SSTable类的有序文件中， 根据需要在后台
合并【14】 。  



### 性能优化

与往常一样， 大量的细节使得存储引擎在实践中表现良好。 例如， 当查找数据库中不存在的
键时， LSM树算法可能会很慢： 您必须检查内存表， 然后将这些段一直回到最老的（ 可能必
须从磁盘读取每一个） ， 然后才能确定键不存在。 为了优化这种访问， 存储引擎通常使用额
外的Bloom过滤器【15】 。(布隆过滤器是内存高效数据结构，用于近似计算集合的内容， 如果数据库中不存在某个键，它能很快告诉你结果， 从而为不存在的键节省许多不必要的磁盘读取操作  )  

即使有许多微妙的东西， LSM树的基本思想 —— 保存一系列在后台合并的SSTables —— 简
单而有效。 即使数据集比可用内存大得多， 它仍能继续正常工作。 由于数据按排序顺序存
储， 因此可以高效地执行范围查询（ 扫描所有高于某些最小值和最高值的所有键） ， 并且因
为磁盘写入是连续的， 所以LSM树可以支持非常高的写入吞吐量。  

### B树

刚才讨论的日志结构索引正处在逐渐被接受的阶段， 但它们并不是最常见的索引类型。 使用
最广泛的索引结构在1970年被引入【17】 ， 不到10年后变得“无处不在”【18】 ， B树经受了时
间的考验。 在几乎所有的关系数据库中， 它们仍然是标准的索引实现， 许多非关系数据库也
使用它们 。

像SSTable一样， B- tree保留按键排序的 k-v对，这样可以实现高效的k-v查找和区间查询。但相似仅此而已 ： B-tree本质上具有非常不同的设计理念。  







