## MapReduce: Simplified Data Processing on Large Clusters  

MapReduce: 简化大型集群上的数据处理

### 概要

MapReduce是用于处理和生成大型数据集的编程模型和相关实现。用户指定一个***map***函数，该函数处理键/值对以生成一组中间键/值

一个***reduce***函数，它合并与同一中间键相关联的所有中间值。许多现实世界中的任务都可以在这个模型中表达，如本文所示。

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

![img](https://github.com/Qasak/all-about-distributed-system/blob/master/notes/mapreduce/%E5%9B%BE1-%E6%89%A7%E8%A1%8C%E6%A6%82%E8%A7%88-execution%20overview.png)

图1显示了我们实现中MapReduce操作的总体流程。当用户程序调用`MapReduce`函数时，将发生以下操作序列（图1中编号的标签对应于下面列表中的数字）：

1. 用户程序中的MapReduce库首先将输入文件分成M个部分，通常每段16兆字节到64兆字节（MB）（由用户通过可选参数控制）. 然后它会在一组机器上启动程序的许多副本. 



2. 这个程序的其中一个副本是特别的-master. 其余的都是master指派的worker. 有M个map任务和R个reduce任务要分配. master选择空闲的worker并为每个worker分配一个map任务或reduce任务. 



3. 分配了map任务的worker读取相应的输入拆分的内容. 它从输入数据中解析出键/值对，并将每个对传递给用户定义的***Map***函数. 由***Map***函数产生的中间键/值对被缓冲在内存中. 



4. 缓冲对被周期性地写入本地磁盘，由分区函数划分成R个区域. 本地磁盘上这些缓冲对的位置被传回主机，主机负责将这些位置转发给reduce workers. 



5. 当主机通知reduce worker这些位置时，它使用远程过程调用(RPC)从map worker的本地磁盘读取缓冲数据. 当reduce worker读取了所有中间数据时，它将按中间键对其进行排序，以便将同一键的所有出现组合在一起. 之所以需要排序，是因为通常有许多不同的键映射到同一个reduce任务. 如果中间数据量太大而无法放入内存，则使用外部排序. 



6. reduce worker迭代已排序的中间数据，对于遇到的每个唯一中间键，它将键和相应的中间值集传递给用户的***reduce***函数. Reduce函数的输出被附加到这个Reduce分区的最终输出文件中



7. 当所有map任务和reduce任务都已完成时，主程序将唤醒用户程序. 此时，用户程序中的`MapReduce`调用返回到用户代码. 

成功完成后，mapreduce执行的输出将在R个输出文件中可用（每个reduce任务一个，文件名由用户指定）。通常，用户不需要将这些R输出文件合并到一个文件中–他们通常将这些文件作为输入传递给另一个MapReduce调用，或者从另一个能够处理划分为多个文件的输入的分布式应用程序中使用这些文件。

#### 3.2 Master 数据结构

master有数个数据结构。对于每个map任务和reduce任务，它存储状态（空闲、正在进行或已完成）以及worker机的标识（对于非空闲任务）。

master是将中间文件区域的位置从map任务传播到reduce任务的中转人。

因此，对于每个完成的map任务，master存储由map任务生成的R中间文件区域的位置和大小。当map任务完成时，会收到对此位置和大小信息的更新。信息将以增量方式推送到具有正在进行的reduce任务的工人

#### 3.3 容错

由于MapReduce库的设计目的是帮助使用成百上千台机器处理大量数据，因此该库必须优雅地容忍机器故障

**worker 故障**

master定期对每个worker进行ping检查。如果在一定时间内没有从worker接收到响应，则主进程会将该worker标记为失败。由该worker完成的任何map任务都会重置回其初始空闲状态，因此有资格在其他worker上进行调度。类似地，对于失败的worker正在进行的任何map任务或reduce任务也将重置为空闲并有资格重新安排。

失败时会重新执行已完成的map任务，因为它们的输出存储在发生故障的计算机的本地磁盘上，因此无法访问。完成的reduce任务不需要重新执行，因为它们的输出存储在全局文件系统中。

当一个map任务首先由worker A执行，然后由worker B执行（因为a失败），所有执行reduce任务的worker都会收到重新执行的通知。任何尚未从worker A读取数据的reduce任务都将从worker B读取数据。

MapReduce对大规模的worker故障具有弹性。例如，在一次MapReduce操作期间，运行集群上的网络维护导致一次80台机器的组在几分钟内无法访问。MapReduce主机只需重新执行无法访问的工作机所做的工作，并继续向前推进，最终完成MapReduce操作。

**Master 故障**

很容易让master周期性写入上述master数据结构的检查点。如果master任务终止，则可以从上一个检查点状态启动新副本。但是，考虑到只有一个主节点，它不太可能失败；因此，如果主节点失败，我们当前的实现将中止MapReduce计算。客户机可以检查这种情况，如果需要，可以重试MapReduce操作

**出现故障时的语义**

当用户提供的***map***和***reduce***运算符是其输入值的确定函数时，我们的分布式实现产生的输出与整个程序的无故障顺序执行所产生的输出相同

我们依靠原子提交(atomic commits)map和reduce任务的输出来实现这个属性。

每个正在进行的任务将其输出写入私有临时文件。reduce任务生成一个这样的文件，map任务生成R个这样的文件（每个reduce任务一个）。map任务完成后，worker向主服务器发送一条消息，并在消息中包含R临时文件的名称。如果主机接收到已完成的map任务的完成消息，它将忽略该消息。否则，它会在主数据结构中记录R个文件的名称

当reduce任务完成时，reduce worker将其临时输出文件自动重命名为最终输出文件。如果在多台计算机上执行相同的reduce任务，则将对同一最终输出文件执行多个重命名调用。我们依赖底层文件系统提供的原子重命名操作来保证最终文件系统状态只包含reduce任务一次执行所产生的数据

我们的***map***和***reduce***运算符绝大多数都是确定性的，而且在这种情况下，我们的语义相当于一个顺序执行，这使得程序员很容易对他们的程序行为进行推理。

当***map***和/或***reduce***运算符不确定时，我们提供了较弱但仍然合理的语义。

在存在非确定性运算符的情况下，特定reduce任务$R_1$的输出相当于由非确定性程序的顺序执行所产生的$R_1$的输出。

然而，不同reduce任务$R_2$的输出可以对应于非确定性程序的不同顺序执行所产生的$R_2$的输出。(However, the output for a different reduce task R2 may correspond to the output for R2 produced by a different sequential execution of the non deterministic program.)

考虑map任务$M$和reduce任务$R_1$和$R_2$。设$e(R_i)$是所提交的$R_i$的执行（只有一个这样的执行）。语义较弱的原因是$e(R_1)$可能已经读取了$M$的一次执行所产生的输出，而$e(R_2)$可能已经读取了$M$的另一次执行所产生的输出。

#### 3.4 局部性

在我们的计算环境中，网络带宽是一种相对稀缺的资源。我们利用输入数据（由GFS[8]管理）存储在组成集群的机器的本地磁盘上，从而节省了网络带宽。GFS将每个文件分成64MB块，并在不同的计算机上存储每个块的多个副本（通常为3个副本）。MapReduce master考虑输入文件的位置信息，并尝试在包含相应输入数据副本的计算机上调度map任务。否则，它会尝试在任务输入数据的副本附近调度map任务（例如，在与包含数据的计算机位于同一网络交换机上的工作机上）。当在集群中的大部分worker上运行大型MapReduce操作时，大多数输入数据都是在本地读取的，并且不消耗网络带宽。

#### 3.5 任务粒度

我们将map阶段细分为M个片段，并将reduce阶段细分为R个片段，如上所述。理想情况下，M和R应该远大于工作机器的数量。让每个worker执行许多不同的任务可以提高动态负载平衡，并在worker失败时加速恢复：它完成的许多map任务可以分散到所有其他worker计算机上。在我们的实现中，M和R的大小有实际的限制，因为主机必须做出O（M+R）调度决策，并如上所述在内存中保持O（M*R）状态。（然而，内存使用的常量因素很小：状态的O（M*R）段由每个map task/reduce任务对大约一个字节的数据组成）此外，R常常受到用户的限制，因为每个reduce任务的输出都在一个单独的输出文件中。在实践中，我们倾向于选择M，以便每个单独的任务大约有16mb到64mb的输入数据（因此上面描述的局部性优化是最有效的），并且我们将R设为期望使用的工作机数量的一个小倍数。我们经常使用2000台工人机器，以M=200000和R=5000执行MapReduce计算。

此外，R常常受到用户的限制，因为每个reduce任务的输出都会在一个单独的输出文件中结束。在实践中，我们倾向于选择M，以便每个单独的任务大约有16mb到64mb的输入数据（因此上面描述的局部性优化是最有效的），并且我们将R设为期望使用的worker机数量的一个小倍数。我们经常使用2000台worker机器，以M=200000和R=5000执行MapReduce计算。

#### 3.6 备份任务

使MapReduce操作所花费的总时间变长的一个常见原因是“掉队者”：一台需要非常长时间来完成计算中最后几个map或reduce任务之一的机器。

掉队者的出现有很多原因。例如，一台有坏磁盘的机器可能会遇到频繁的可纠正错误，导致其读取性能从30 MB/s降低到1 MB/s。群集调度系统可能已在计算机上计划了其他任务，导致其由于CPU、内存、本地磁盘或网络带宽的竞争而更慢地执行MapReduce代码。最近我们遇到的一个问题是机器初始化代码中的一个错误，导致处理器缓存被禁用：受影响机器上的计算速度减慢了100倍以上。

我们有一个总的机制来缓解掉队者的问题。当MapReduce操作即将完成时，主服务器将调度剩余*正在进行*的任务的备份执行。

每当主或备份执行完成时，任务都会标记为已完成。我们已经对该机制进行了优化，使其通常只将操作使用的计算资源增加不超过百分之几。

我们发现这大大缩短了完成大型MapReduce操作的时间。例如，当备份任务机制被禁用时，第5.3节中描述的排序程序需要花费44%的时间才能完成。

### 4 改进

