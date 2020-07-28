## MapReduce: Simplified Data Processing on Large Clusters  

MapReduce: 简化大型集群上的数据处理

### 概要

MapReduce是用于处理和生成大型数据集的编程模型和相关实现。用户指定一个映射函数，该函数处理键/值对以生成一组中间键/值

一个reduce函数，它合并与同一中间键相关联的所有中间值。许多现实世界中的任务都可以在这个模型中表达，如本文所示。

用这种函数式编写的程序会自动并行化，并在大型商用计算机集群上执行。运行时系统负责对输入数据进行分区、在一组机器上调度程序的执行、处理机器故障以及管理所需的机器间通信等细节。这使得没有并行和分布式系统经验的程序员能够轻松地利用大型分布式系统的资源。

我们的MapReduce实现运行在大型商用计算机集群上，并且具有高度可扩展性：

典型的MapReduce计算在数千台机器上处理许多TB的数据。程序员发现该系统易于使用：数百个MapReduce程序已经实现，每天在Google集群上执行的MapReduce作业超过1000个

### 1 简介

在过去的五年里，作者和谷歌的许多其他人已经实现了数百个特殊目的计算：

处理大量原始数据的计算，

例如爬取文档、web请求日志等，以计算各种派生数据，例如

索引、web文档图结构，各种表示、每个被爬取主机的页面数量的摘要，在给定的一天内一组最频繁的查询

大多数这样的计算在概念上是直接的。然而，输入的数据通常很大，为了在合理的时间内完成计算，必须将计算分布在成百上千台机器上。如何并行化计算、分发数据和处理失败等问题，使原本简单的计算变得模糊，而处理这些问题需要大量复杂代码。

作为对这种复杂性的反应，我们设计了一种新的抽象，它允许我们表达我们试图执行的简单计算，但隐藏了库中并行化、容错、数据分布和负载均衡的混乱细节。我们的抽象源于Lisp和许多其他函数语言中的map和reduce原语。我们意识到，我们的大多数计算都涉及对输入中的每个逻辑“记录”应用映射(`map`)操作，以便计算一组中间键/值对，然后对共享同一键的所有值应用`reduce`操作，以便适当地组合派生的数据。我们使用带有用户指定的`map`和`reduce`操作的函数模型，可以轻松地并行化大型计算并使用重用(re-execution)作为容错的主要机制

这项工作的主要贡献是一个简单而强大的接口，它可以实现大规模计算的自动并行和分布，再加上这个接口的实现，可以在大型商用PC机群上实现高性能。

第2节描述了基本编程模型并给出了几个示例。第3节描述了为我们基于集群的计算环境定制的MapReduce接口的实现。第4节描述了我们发现有用的编程模型的一些改进。第5节介绍了我们对各种任务的实现的性能度量。第6节探讨了MapReduce在Google中的使用，包括我们使用它作为重写产品索引系统的基础的经验。第7节讨论相关工作和未来工作



### 2 编程模型

计算采用一组输入键/值对，并生成一组输出键/值对。MapReduce库的用户将计算表示为两个函数：Map和Reduce。

Map由用户编写，接受一个输入对并生成一组中间键/值对。MapReduce库将与同一中间键$I$关联的所有中间值组合在一起，并将它们传递给Reduce函数。

Reduce函数也是由用户编写的，它接受中间键$I$和该键的一组值。它将这些值合并在一起形成一个可能更小的值集。通常每次Reduce调用只产生零或一个输出值。中间值通过迭代器提供给用户的reduce函数。这使我们能够处理太大而无法放入内存的值列表

#### 2.1 例子

考虑计算一个大型文档集合中每个单词出现的次数的问题。用户将编写类似以下伪代码的代码：

```c++
map(String key, String value):
    // key: document name
    // value: document contents
    for each word w in value:
    	EmitIntermediate(w, "1");

reduce(String key, Iterator values):
    // key: a word
    // values: a list of counts
    int result = 0;
    for each v in values:
    	result += ParseInt(v);
    Emit(AsString(result));
```

map函数发出(emit)每个单词加上一个相关的出现次数（在这个简单的例子中是“1”）。reduce函数将为特定单词发出的所有计数相加

此外，用户编写代码，用输入和输出文件的名称以及可选的调优参数填充mapreduce规范对象。然后，用户调用`MapReduce`函数，将规范对象传递给它。用户的代码与`MapReduce`库（C++实现）链接在一起。附录A包含此示例的完整程序文本

```c++
#include "mapreduce/mapreduce.h"
// User’s map function
class WordCounter : public Mapper {
	public:
		virtual void Map(const MapInput& input) {
            const string& text = input.value();
            const int n = text.size();
            for (int i = 0; i < n; ) {
            // Skip past leading whitespace
            while ((i < n) && isspace(text[i]))
            	i++;
            
            // Find word end
            int start = i;
            while ((i < n) && !isspace(text[i]))
            	i++;
            if (start < i)
            Emit(text.substr(start,i-start),"1");
        }
    }
};
REGISTER_MAPPER(WordCounter);

// User’s reduce function
class Adder : public Reducer {
    virtual void Reduce(ReduceInput* input) {
        // Iterate over all entries with the
        // same key and add the values
        int64 value = 0;
        while (!input->done()) {
            value += StringToInt(input->value());
            input->NextValue();
    }
    // Emit sum for input->key()
    Emit(IntToString(value));
    }
};
REGISTER_REDUCER(Adder);

int main(int argc, char** argv) {
    ParseCommandLineFlags(argc, argv);
    MapReduceSpecification spec;
    // Store list of input files into "spec"
    for (int i = 1; i < argc; i++) {
        MapReduceInput* input = spec.add_input();
        input->set_format("text");
        input->set_filepattern(argv[i]);
        input->set_mapper_class("WordCounter");
    }
    // Specify the output files:
    // /gfs/test/freq-00000-of-00100
    // /gfs/test/freq-00001-of-00100
    // ...
    MapReduceOutput* out = spec.output();
    out->set_filebase("/gfs/test/freq");
    out->set_num_tasks(100);
    out->set_format("text");
    out->set_reducer_class("Adder");
    // Optional: do partial sums within map
    // tasks to save network bandwidth
    out->set_combiner_class("Adder");
    // Tuning parameters: use at most 2000
    // machines and 100 MB of memory per task
    spec.set_machines(2000);
    spec.set_map_megabytes(100);
    spec.set_reduce_megabytes(100);
    // Now run it
    MapReduceResult result;
    if (!MapReduce(spec, &result)) abort();
    // Done: ’result’ structure contains info
    // about counters, time taken, number of
    // machines used, etc.
    return 0;
}
```



#### 2.2 类型

尽管前面的伪代码是根据字符串输入和输出编写的，但从概念上讲，用户提供的map和reduce函数具有关联类型：

```c
map		(k1, v1)		->list(k2, v2)
reduce	(k2, list(v2))	 ->list(v2)
```



 例如，输入键和值是从与输出键和值不同的域中提取的。

此外，中间键和值与输出键和值来自同一个域。

我们的C++实现将字符串传递到用户定义的函数，并将其留给用户代码，以便在字符串和适当类型之间转换。

#### 2.3 更多例子

下面是一些有趣程序的简单示例

他们可以很容易地表示为MapReduce计算。

+ **分布式Grep**

  如果map函数与提供的模式匹配，那么它将发出(emit)一行。

  reduce函数是一个标识函数(identify function)，它只将提供的中间数据复制到输出

+ **URL访问频率计数**

  map函数处理网页请求和输出的日志`<URL, 1>`

  reduce函数将同一URL的所有值相加，并发出一个`<URL, total count>`对



+ **反向Web链接图**

  map函数输出`<target，source>`对，

  对于每个指向`target`URL，在其页面中发现的`source`

  reduce函数串联list并发送`<target, list(sorce)>`



+ **每个主机的术语向量**

  术语向量把出现在一个文档或一组文档中的最重要的单词概括为一个<word，frequency>对的列表。

  map函数为每个输入文档发出一个<hostname，term vector>对（主机名从文档的URL中提取）。

  reduce函数传递给定主机的所有每个文档项向量。它将这些术语向量相加，丢弃不常见的术语，然后发出最终的<hostname，term vector>对

+ **Inverted index(倒排索引)**

  > 一个未经处理的数据库中，一般是以文档ID作为索引，以文档内容作为记录。
  > 而Inverted index 指的是将单词或记录作为索引，将文档ID作为记录，这样便可以方便地通过单词或记录查找到其所在的文档。

  map函数解析每个文档，并发出一个<word，document ID>对的序列。

  reduce函数接受给定单词的所有对，对相应的文档id进行排序并发出一个

  <word，list（document ID）>配对。

  所有输出对的集合形成一个简单的倒排索引。很容易增加这种计算来跟踪单词的位置。

+ **分布式排序**

  map函数从每个记录中提取键，并发出一个<key，record>对。

  reduce函数不变地发出所有对。该计算取决于第4.1节中描述的partitioning facilities和第4.2节中描述的排序属性。

  

### 3实现

MapReduce接口有许多不同的实现。正确的选择取决于环境。例如，一种实现可能适合于小型共享内存机器，另一种适用于大型NUMA多处理器，以及另一种更大的联网机器集合。

本节介绍了一个针对Google广泛使用的计算环境的实现：

通过交换式以太网连接在一起的大型商品PC集群[4]。在我们的环境中：

（1） 机器通常是双核x86处理器

​	运行Linux，每台机器有2-4GB内存。

（2） 通常使用商品网络硬件

​	100Mb/s或1Gb/s

​	但平均在整个平分带宽中要少得多。

（3） 集群由成百上千台机器组成，因此机器故障很常见。

（4） 存储由直接连接到单个机器的廉价IDE磁盘提供。内部开发的分布式文件系统[8]用于管理存储在这些磁盘上的数据。文件系统使用复制在不可靠的硬件上提供可用性和可靠性。

（5） 用户向调度系统提交作业。每项工作

由一组任务组成，并由调度程序映射

到群集中的一组可用计算机。

#### 3.1 执行概览

通过自动将输入数据划分为一组M个分割，*Map*调用分布在多台机器上。

输入拆分可以由不同的机器并行处理。

Reduce调用是通过使用一个分区函数(patitioning)（例如$hash(key)\mod R$）将中间密钥空间划分为R个片段来分布的。分区数（R）和分区函数由用户指定。

图1显示了我们实现中MapReduce操作的总体流程。当用户程序调用`MapReduce`函数时，将发生以下操作序列（图1中编号的标签对应于下面列表中的数字）：

1. 用户程序中的MapReduce库首先将输入文件分成M个部分，通常每段16兆字节到64兆字节（MB）（由用户通过可选参数控制）. 然后它会在一组机器上启动程序的许多副本. 



2. 这个程序的其中一个副本是特别的-大师. 其余的都是主人指派的工人. 有M个映射任务和R个reduce任务要分配. 主节点选择空闲的工人并为每个工人分配一个map任务或reduce任务. 



3. 分配了映射任务的工作人员读取相应的输入拆分的内容. 它从输入数据中解析出键/值对，并将每个对传递给用户定义的Map函数. 由Map函数产生的中间键/值对被缓冲在内存中. 



3. 缓冲对被周期性地写入本地磁盘，由分区函数划分成R个区域. 本地磁盘上这些缓冲对的位置被传回主机，主机负责将这些位置转发给reduce workers. 



3. 当主机通知reduce worker这些位置时，它使用远程过程调用从map worker的本地磁盘读取缓冲数据. 当reduce worker读取了所有中间数据时，它将按中间键对其进行排序，以便将同一键的所有出现组合在一起. 之所以需要排序，是因为通常有许多不同的键映射到同一个reduce任务. 如果中间数据量太大而无法放入内存，则使用外部排序. 



3. reduce worker迭代已排序的中间数据，对于遇到的每个唯一中间键，它将键和相应的中间值集传递给用户的reduce函数. Reduce函数的输出被附加到这个Reduce分区的最终输出文件中



7. 当所有映射任务和reduce任务都已完成时，主程序将唤醒用户程序. 此时，用户程序中的MapReduce调用返回到用户代码. 