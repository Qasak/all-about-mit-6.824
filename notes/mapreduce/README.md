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

