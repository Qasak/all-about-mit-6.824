## 关系模型与文档模型

### NoSQL

并非特定技术，(Not only SQL)

+ 比关系型数据库更好的可扩展性，包括非常大的数据集或非常高的写入吞吐量
+ 开源
+ 关系模型不能很好支持特殊的查询
+ 更具多动态性与表现力的数据模型

### 对象关系不匹配

如果数据存储在关系表中，那么需要一个笨拙的转换层，处于应用程序代码中的对 象和表，行，列的数据库模型之间

> "阻抗不匹配":每个电路的输入和输出都有一定的阻抗（交流电阻）
>
> 将一个电路的输出连接到另一个电路的输入时，如果两个电路的输出和输入阻抗匹配， 则连接上的功率传输将被最大化。阻抗不匹配会导致信号反射及其他问题。

![img](https://github.com/Qasak/distributed-system/blob/master/%E6%95%B0%E6%8D%AE%E7%B3%BB%E7%BB%9F%E7%9A%84%E5%9F%BA%E7%9F%B3/%E6%95%B0%E6%8D%AE%E6%A8%A1%E5%9E%8B%E4%B8%8E%E6%9F%A5%E8%AF%A2%E8%AF%AD%E8%A8%80/linkedin0.png)

使用关系型模式来表示领英简介

整个简介可以通过一 个唯一的标识符` user_id `来标识

像 first_name 和 last_name 这样的字段每个用户只出现一 次，所以可以在User表上将其建模为列

但是，大多数人在职业生涯中拥有多于一份的工 作，人们可能有不同样的教育阶段和任意数量的联系信息。从用户到这些项目之间存在一对 多的关系，可以用多种方式来表示：

+ 传统SQL模型（SQL：1999之前）中，最常见的规范化表示形式是将职位，教育和联系 信息放在单独的表中，对User表提供外键引用，如上图。

+ 后续的SQL标准增加了对结构化数据类型和XML数据的支持;这允许将多值数据存储在单 行内，并支持在这些文档内查询和索引。这些功能在Oracle，IBM DB2，MS SQL Server 和PostgreSQL中都有不同程度的支持。JSON数据类型也得到多个数据库的支 持，包括IBM DB2，MySQL和PostgreSQL 。
+ 第三种选择是将职业，教育和联系信息编码为JSON或XML文档，将其存储在数据库的文 本列中，并让应用程序解析其结构和内容。这种配置下，通常不能使用数据库来查询该 编码列中的值。

```json
{
    "user_id": 251,
    "first_name": "Bill",
    "last_name": "Gates",
    "summary": "Co-chair of the Bill & Melinda Gates... Active blogger.",
    "region_id": "us:91",
    "industry_id": 131,
    "photo_url": "/p/7/000/253/05b/308dd6e.jpg",
    "positions": [
        {
            "job_title": "Co-chair",
            "organization": "Bill & Melinda Gates Foundation"
        },
        {
            "job_title": "Co-founder, Chairman",
            "organization": "Microsoft"
        }
    ],
    "education": [
        {
            "school_name": "Harvard University",
            "start": 1973,
            "end": 1975
        },
        {
            "school_name": "Lakeside School, Seattle",
            "start": null,
            "end": null
        }
    ],
    "contact_info": {
        "blog": "http://thegatesnotes.com",
        "twitter": "http://twitter.com/BillGates"
    }
}

```



对于一个像简历这样自包含文档的数据结构而言，JSON表示是非常合适的：如上。 JSON比XML更简单。*面向文档*的数据库（如MongoDB ，RethinkDB ， CouchDB 和Espresso）支持这种数据模型

JSON表示比图中的多表模式具有更好的局部性（locality）。如果在前面的关系型示例中 获取简介，那需要执行多个查询（通过 user_id 查询每个表），或者在User表与其下属表之 间混乱地执行多路连接。而在JSON表示中，所有相关信息都在同一个地方，一个查询就足够 了。

从用户简介文件到用户职位，教育历史和联系信息，这种一对多关系隐含了数据中的一个树 状结构，而JSON表示使得这个树状结构变得明确

![img](https://github.com/Qasak/distributed-system/blob/master/%E6%95%B0%E6%8D%AE%E7%B3%BB%E7%BB%9F%E7%9A%84%E5%9F%BA%E7%9F%B3/%E6%95%B0%E6%8D%AE%E6%A8%A1%E5%9E%8B%E4%B8%8E%E6%9F%A5%E8%AF%A2%E8%AF%AD%E8%A8%80/json0.png)

一对多关系构建了一个树结构

### 多对一和一对多的关系

region_id 和 industry_id 是以ID，而不是纯字符串“Greater Seattle Area”和“Philanthropy”的形式给出的。为什么？

存储ID还是文本字符串，这是个 副本（duplication） 问题

使用ID的好处是，ID对人类没有任何意义，因而永远不需要改变：ID可以保持不变，即使它 标识的信息发生变化。任何对人类有意义的东西都可能需要在将来某个时候改变——如果这 些信息被复制，所有的冗余副本都需要更新。这会导致写入开销，也存在不一致的风险（一些副本被更新了，还有些副本没有被更新）。

去除此类重复是数据库 规范化 （normalization） 的关键思想。

> 规范化：如果重复存储了可以存储在一个地方的值，就是不规范化的

对这些数据进行规范化需要多对一的关系（许多人生活在一个特定的地区，许多 人在一个特定的行业工作），这与文档模型不太吻合。

在关系数据库中，通过ID来引用其他 表中的行是正常的，因为连接很容易。在文档数据库中，一对多树结构没有必要用连接，对 连接的支持通常很弱

> RethinkDB支持连接，MongoDB不支持连接，ChouchDB只支持预先声明的试图

如果数据库本身不支持连接，则必须在应用程序代码中通过对数据库进行多个查询来模拟连接。

此外，即便应用程序的最初版本适合无连接的文档模型，随着功能添加到应用程序中，数据 会变得更加互联。例如，考虑一下对简历例子进行的一些修改：

组织和学校作为实体

在前面的描述中， organization （用户工作的公司）和 school_name （他们学习的地方）只 是字符串。也许他们应该是对实体的引用呢？然后，每个组织，学校或大学都可以拥有自己 的网页（标识，新闻提要等）。每个简历可以链接到它所提到的组织和学校，并且包括他们 的图标和其他信息（参见图2-3，来自LinkedIn的一个例子）。



![img]()

公司名不仅是字符串，还是一个指向公司实体的链接



推荐

推荐应该拥有作者个人简介的引用。

### 文档数据库是否重蹈覆辙？

回顾关系模型和网络模型的辩论

### 网络模型

网络模型种记录之间的链接不是外键，而更像指针

访问记录的方法是沿着访问路径，类似于链表

无论是分层还是网络模型，如果没有所需数据的路径，必须浏览手写数据库查询代码

更改数据模型是很难的

### 关系模型

没有嵌套结构和复杂的访问路径

### 文档和关系数据库的融合

大多数关系数据库(MySQL除外)都支持XML

关系模型和文档模型混合是未来数据库一条很好的路线

## 数据查询语言

SQL是一种声明式查询语言

IMS和CODASYL使用命令式代码查询

SQL相当有限的功能性为数据库提供了更多自动优化的空间

声明式语言往往适合并行执行

命令代码很难在多个内核和多个机器之间并行化，因为它制定了指令必须以特定顺序执行（它仅指定结果的模式，不指定用于确定结果的算法）

### Web上的声明式查询

在Web浏览器中，使用声明式CSS样式比使用JavaScript命令式地操作样式要好得多。类似 地，在数据库中，使用像SQL这样的声明式查询语言比使用命令式查询API要好得多

### MapReduce查询

一些NoSQL数据存储(包括MongoDB,CouchDB)支持有限形式的MapReduce作为多个文档中执行只读查询的机制

简述MongoDB使用的模型：

MapReduce既不是一个声明式查询语言，也不是一个完全命令式的查询API，而是处于两者之间：查询的逻辑用代码片段表示，这些代码片段会被处理框架重复性调用。它基于`map` `reduce`函数

eg：假设你是一名海洋生物学家，每当你看到海洋中的动物 时，你都会在数据库中添加一条观察记录。现在你想生成一个报告，说明你每月看到多少鲨 鱼。

在PostgreSQL中：

```sql
SELECT
    date_trunc('month', observation_timestamp) AS observation_month,
    sum(num_animals) AS total_animals
FROM observations
WHERE family = 'Sharks'
GROUP BY observation_month;
```

date_trunc('month'，timestamp) 函数用于确定包含 timestamp 的日历月份，并返回代表该月 份开始的另一个时间戳。换句话说，它将时间戳舍入成最近的月份

同样的查询用MongoDB的MapReduce功能可以如下表述：

```javascript
db.observations.mapReduce(function map() {
    var year = this.observationTimestamp.getFullYear();
    var month = this.observationTimestamp.getMonth() + 1;
    emit(year + "-" + month, this.numAnimals);
},
function reduce(key, values) {
    return Array.sum(values);
},
{
    query: {
        family: "Sharks"
    },
    out: "monthlySharkReport"
});
```

+ 可以声明式地指定只考虑鲨鱼种类的过滤器（这是一个针对MapReduce的特定于 MongoDB的扩展）。
+ 将 this 设置为文档对象,每个匹配查询的文档都会调用一次JavaScript函数 map
+ map 函数发出一个键（包括年份和月份的字符串，如 "2013-12" 或 "2014-1" ）和一个值 （该观察记录中的动物数量）。
+ map 发出的键值对按键来分组。对于具有相同键（即，相同的月份和年份）的所有键值 对，调用一次 reduce 函数。
+ reduce 函数将特定月份内所有观测记录中的动物数量相加。
+ 将最终的输出写入到 monthlySharkReport 集合中。

例如，假设`observations`集合中包含这两个文档：

```js
{
    observationTimestamp: Date.parse( "Mon, 25 Dec 1995 12:34:56 GMT"),
    family: "Sharks",
    species: "Carcharodon carcharias",
    numAnimals: 3
{
}
    observationTimestamp: Date.parse("Tue, 12 Dec 1995 16:17:18 GMT"),
    family: "Sharks",
    species: "Carcharias taurus",
    numAnimals: 4
}

```

对每个文档都会调用一次 map 函数，结果将是 emit("1995-12",3) 和 emit("1995-12",4) 。随 后，以 reduce("1995-12",[3,4]) 调用 reduce 函数，将返回 7 。

map和reduce函数在功能上有所限制：它们必须是纯函数，这意味着它们只使用传递给它们 的数据作为输入，它们不能执行额外的数据库查询，也不能有任何副作用。这些限制允许数 据库以任何顺序运行任何功能，并在失败时重新运行它们。然而，map和reduce函数仍然是 强大的：它们可以解析字符串，调用库函数，执行计算等等。

> 纯函数：
>
> - 它应始终返回相同的值。不管调用该函数多少次，无论今天、明天还是将来某个时候调用它。
> - 自包含（不使用全局变量）。
> - 它不应修改程序的状态或引起副作用（修改全局变量）。

MapReduce是一个相当底层的编程模型，用于计算机集群上的分布式执行。像SQL这样的更 高级的查询语言可以用一系列的MapReduce操作来实现（见第10章），但是也有很多不使用 MapReduce的分布式SQL实现。请注意，SQL中没有任何内容限制它在单个机器上运行，而 MapReduce在分布式查询执行上没有垄断权。



能够在查询中使用JavaScript代码是高级查询的一个重要特性，但这不限于MapReduce，一 些SQL数据库也可以用JavaScript函数进行扩展

MapReduce的一个可用性问题是，必须编写两个密切合作的JavaScript函数，这通常比编写 单个查询更困难。此外，声明式查询语言为查询优化器提供了更多机会来提高查询的性能。 基于这些原因，MongoDB 2.2添加了一种叫做聚合管道的声明式查询语言的支持【9】。用这 种语言表述鲨鱼计数查询如下所示：



```json
db.observations.aggregate([
    { $match: { family: "Sharks" } },
    { $group: {
    _id: {
        year: { $year: "$observationTimestamp" },
        month: { $month: "$observationTimestamp" }
    },
    totalAnimals: { $sum: "$numAnimals" } }}
]);

```

聚合管道语言与SQL的子集具有类似表现力，但是它使用基于JSON的语法而不是SQL的英语 句子式语法; 这种差异也许是口味问题。这个故事的寓意是NoSQL系统可能会发现自己意外地 重新发明了SQL，尽管带着伪装。

## 图数据模型

多对多关系是不同数据模型之间具有区别性的重要特征。

如果你的应用程 序大多数的关系是一对多关系（树状结构化数据），或者大多数记录之间不存在关系，那么 使用文档模型是合适的。

但是，要是多对多关系在你的数据中很常见呢？关系模型可以处理多对多关系的简单情况， 但是随着数据之间的连接变得更加复杂，将数据建模为图形显得更加自然。

一个图由两种对象组成：顶点（vertices）（也称为节点（nodes） 或实体（entities））， 和边（edges）（ 也称为关系（relationships）或弧 （arcs） ）。多种数据可以被建模为 一个图形。典型的例子包括：

社交图谱 顶点是人，边指示哪些人彼此认识。 

网络图谱 顶点是网页，边缘表示指向其他页面的HTML链接。 

公路或铁路网络 顶点是交叉路口，边线代表它们之间的道路或铁路线。









