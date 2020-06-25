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

在多对多的关系和连接已常规用在关系数据库时，文档数据库和NoSQL重启了辩论：如何最好地在数据库中表示多对多关系。

IBM的信息管理系统（IMS）使用了一个相当简单的数据模型，称为层次模型（hierarchical model），它与文档数据库使用的JSON模型有一些惊人的相似之处【2】。它将所有数据表示为嵌套在记录 中的记录树

同文档数据库一样，IMS能良好处理一对多的关系，但是很难应对多对多的关系，并且不支持连接

开发人员必须决定是否复制（非规范化）数据或手动解决从一个记录到另一个记录的引用

这些二十世纪六七十年代的问题与现在开发人员遇到的文档数据库问题非常相似

那时人们提出了各种不同的解决方案来解决层次模型的局限性。其中最突出的两个是关系模型（relational model）（它变成了SQL，统治了世界）和网络模型（network model）（最初很受关注，但最终变得冷门）。

### 网络模型

也被称为*CODASYL*模型

网络模型种记录之间的链接不是外键，而更像指针

> 主键：一列（或一组列），其值能够唯一区分表 中每个行。

> 外键：某个表中的一列，它包含另一个表 的主键值，定义了两个表之间的关系。

访问记录的方法是沿着访问路径，类似于链表

无论是分层还是网络模型，如果没有所需数据的路径，必须浏览手写数据库查询代码

更改数据模型是很难的

### 关系模型

没有嵌套结构和复杂的访问路径

### 文档和关系数据库的融合

大多数关系数据库(MySQL除外)都支持XML

关系模型和文档模型混合是未来数据库一条很好的路线

## 数据查询语言

SQL是一种*声明式*查询语言

IMS和CODASYL使用*命令式*代码查询

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

可以将那些众所周知的算法运用到这些图上：例如，汽车导航系统搜索道路网络中两点之间 的最短路径，PageRank可以用在网络图上来确定网页的流行程度，从而确定该网页在搜索结 果中的排名。

在刚刚给出的例子中，图中的所有顶点代表了相同类型的事物（人，网页或交叉路口）。不 过，图并不局限于这样的同类数据：同样强大地是，图提供了一种一致的方式，用来在单个 数据存储中存储完全不同类型的对象。例如，Facebook维护一个包含许多不同类型的顶点和边的单个图：顶点表示人，地点，事件，签到和用户的评论;边缘表示哪些人是彼此的朋友， 哪个签到发生在何处，谁评论了哪条消息，谁参与了哪个事件，等等

有几种不同但相关的方法用来构建和查询图表中的数据。在本节中，我们将讨论属性图模型 （由Neo4j，Titan和InfiniteGraph实现）和三元组存储（triple-store）模型（由Datomic， AllegroGraph等实现）。我们将查看图的三种声明式查询语言：Cypher，SPARQL和 Datalog。除此之外，还有像Gremlin 【36】这样的图形查询语言和像Pregel这样的图形处理 框架（见第10章）。

### 属性图

在属性图模型中，每个顶点（vertex）包括：

+ 唯一的标识符 
+ 一组 出边（outgoing edges） 
+ 一组 入边（ingoing edges） 
+ 一组属性（键值对）

每条 边（edge） 包括： 

+ 唯一标识符 
+ 边的起点/尾部顶点（tail vertex）
+ 边的终点/头部顶点（head vertex）
+ 描述两个顶点之间关系类型的标签
+ 一组属性（键值对）

可以将图存储看作由两个关系表组成：一个存储顶点，另一个存储边，如例2-2所示（该模式 使用PostgreSQL json数据类型来存储每个顶点或每条边的属性）。头部和尾部顶点用来存储 每条边；如果你想要一组顶点的输入或输出边，你可以分别通 过 head_vertex 或 tail_vertex 来查询 edges 表。

```plsql
CREATE TABLE vertices (
    vertex_id INTEGER PRIMARY KEY,
    properties JSON
);
CREATE TABLE edges (
    edge_id INTEGER PRIMARY KEY,
    tail_vertex INTEGER REFERENCES vertices (vertex_id),
    head_vertex INTEGER REFERENCES vertices (vertex_id),
    label TEXT,
    properties JSON
);
CREATE INDEX edges_tails ON edges (tail_vertex);
CREATE INDEX edges_heads ON edges (head_vertex);
```

一些属性：

1. 任何顶点都可以有一条边连接到任何其他顶点。没有模式限制哪种事物可不可以关联。 

2. 给定任何顶点，可以高效地找到它的入边和出边，从而遍历图，即沿着一系列顶点的路 径前后移动。（这就是为什么例2-2在 tail_vertex 和 head_vertex 列上都有索引的原因。） 
3.  通过对不同类型的关系使用不同的标签，可以在一个图中存储几种不同的信息，同时仍 然保持一个清晰的数据模型。

这些特性为数据建模提供了很大的灵活性，如图2-5所示。图中显示了一些传统关系模式难以 表达的事情，例如不同国家的不同地区结构（法国有省和州，美国有不同的州和州），国中国的怪事（先忽略主权国家和国家错综复杂的烂摊子），不同的数据粒度（Lucy现在的住所 被指定为一个城市，而她的出生地点只是在一个州的级别）。

你可以想象延伸图还能包括许多关于Lucy和Alain，或其他人的其他更多的事实。例如，你可 以用它来表示食物过敏（为每个过敏源增加一个顶点，并增加人与过敏源之间的一条边来指 示一种过敏情况），并链接到过敏源，每个过敏源具有一组顶点用来显示哪些食物含有哪些 物质。然后，你可以写一个查询，找出每个人吃什么是安全的。图表在可演化性是富有优势 的：当向应用程序添加功能时，可以轻松扩展图以适应应用程序数据结构的变化。



## Cypher查询语言

Cypher式属性图的声明式查询语言，为Neo4j图数据库而发明

```cypher
CREATE
    (NAmerica:Location {name:'North America', type:'continent'}),
    (USA:Location {name:'United States', type:'country' }),
    (Idaho:Location {name:'Idaho', type:'state' }),
    (Lucy:Person {name:'Lucy' }),
    (Idaho) -[:WITHIN]-> (USA) -[:WITHIN]-> (NAmerica),
    (Lucy) -[:BORN_IN]-> (Idaho)
```

当图2-5的所有顶点和边被添加到数据库后，让我们提些有趣的问题：例如，找到所有从美国 移民到欧洲的人的名字。更确切地说，这里我们想要找到符合下面条件的所有顶点，并且返 回这些顶点的 name 属性：该顶点拥有一条连到美国任一位置的 BORN_IN 边，和一条连到欧洲 的任一位置的 LIVING_IN 边。

例2-4展示了如何在Cypher中表达这个查询。在MATCH子句中使用相同的箭头符号来查找图 中的模式： (person) -[:BORN_IN]-> () 可以匹配 BORN_IN 边的任意两个顶点。该边的尾节点 被绑定了变量 person ，头节点则未被绑定。

```cypher
MATCH
    (person) -[:BORN_IN]-> () -[:WITHIN*0..]-> (us:Location {name:'United States'}),
    (person) -[:LIVES_IN]-> () -[:WITHIN*0..]-> (eu:Location {name:'Europe'})
RETURN person.name

```



> 找到满足以下两个条件的所有顶点（称之为person顶点）：
>
> 1. person 顶点拥有一条到某个顶点的 BORN_IN 出边。从那个顶点开始，沿着一系 列 WITHIN 出边最终到达一个类型为 Location ， name 属性为 United States 的顶 点。
>
>     
>
> 2. person 顶点还拥有一条 LIVES_IN 出边。沿着这条边，可以通过一系列 WITHIN 出边 最终到达一个类型为 Location ， name 属性为 Europe 的顶点。 对于这样的 Person 顶点，返回其 name 属性。

执行这条查询可能会有几种可行的查询路径。这里给出的描述建议首先扫描数据库中的所有 人，检查每个人的出生地和居住地，然后只返回符合条件的那些人。 等价地，也可以从两个 Location 顶点开始反向地查找。假如 name 属性上有索引，则可以高 效地找到代表美国和欧洲的两个顶点。然后，沿着所有 WITHIN 入边，可以继续查找出所有在 美国和欧洲的位置（州，地区，城市等）。最后，查找出那些可以由 BORN_IN 或 LIVES_IN 入 边到那些位置顶点的人。 通常对于声明式查询语言来说，在编写查询语句时，不需要指定执行细节：查询优化程序会 自动选择预测效率最高的策略，因此你可以继续编写应用程序的其他部分。

### SQL中的图查询

例2-2建议在关系数据库中表示图数据。但是，如果把图数据放入关系结构中，我们是否也可 以使用SQL查询它？ 答案是肯定的，但有些困难。

在关系数据库中，你通常会事先知道在查询中需要哪些连接。 在图查询中，你可能需要在找到待查找的顶点之前，遍历可变数量的边。也就是说，连接的 数量事先并不确定。

在Cypher中，用 WITHIN * 0 非常简洁地表述了上述事实：“沿着 WITHIN 边，零次或多次”。它 很像正则表达式中的 * 运算符。 自SQL:1999，查询可变长度遍历路径的思想可以使用称为递归公用表表达式（ *WITH RECURSIVE* 语法）的东西来表示。例2-5显示了同样的查询 - 查找从美国移民到欧洲的人的姓名 - 在SQL使用这种技术（PostgreSQL，IBM DB2，Oracle和SQL Server均支持）来表述。但 是，与Cypher相比，其语法非常笨拙。 例2-5 与示例2-4同样的查询，在SQL中使用递归公用表表达式表示

```SQL
WITH RECURSIVE
    -- in_usa 包含所有的美国境内的位置ID
    in_usa(vertex_id) AS (
    SELECT vertex_id FROM vertices WHERE properties ->> 'name' = 'United States'
    UNION
    SELECT edges.tail_vertex FROM edges
        JOIN in_usa ON edges.head_vertex = in_usa.vertex_id
        WHERE edges.label = 'within'
    ),
    -- in_europe 包含所有的欧洲境内的位置ID
    in_europe(vertex_id) AS (
    SELECT vertex_id FROM vertices WHERE properties ->> 'name' = 'Europe'
    UNION
    SELECT edges.tail_vertex FROM edges
        JOIN in_europe ON edges.head_vertex = in_europe.vertex_id
        WHERE edges.label = 'within' ),
    -- born_in_usa 包含了所有类型为Person，且出生在美国的顶点
    born_in_usa(vertex_id) AS (
    SELECT edges.tail_vertex FROM edges
        JOIN in_usa ON edges.head_vertex = in_usa.vertex_id
        WHERE edges.label = 'born_in' ),
    -- lives_in_europe 包含了所有类型为Person，且居住在欧洲的顶点。
    lives_in_europe(vertex_id) AS (
    SELECT edges.tail_vertex FROM edges
        JOIN in_europe ON edges.head_vertex = in_europe.vertex_id
        WHERE edges.label = 'lives_in')
    SELECT vertices.properties ->> 'name'
    FROM vertices
        JOIN born_in_usa ON vertices.vertex_id = born_in_usa.vertex_id
        JOIN lives_in_europe ON vertices.vertex_id = lives_in_europe.vertex_id;
```

+ 首先，查找 name 属性为 United States 的顶点，将其作为 in_usa 顶点的集合的第一个 元素。 
+ 从 in_usa 集合的顶点出发，沿着所有的 with_in 入边，将其尾顶点加入同一集合，不断 递归直到所有 with_in 入边都被访问完毕。 
+ 同理，从 name 属性为 Europe 的顶点出发，建立 in_europe 顶点的集合。 
+ 对于 in_usa 集合中的每个顶点，根据 born_in 入边来查找出生在美国某个地方的人。 
+ 同样，对于 in_europe 集合中的每个顶点，根据 lives_in 入边来查找居住在欧洲的人。 
+ 最后，把在美国出生的人的集合与在欧洲居住的人的集合相交。

同一个查询，用某一个查询语言可以写成4行，而用另一个查询语言需要29行，这恰恰说明了 不同的数据模型是为不同的应用场景而设计的。选择适合应用程序的数据模型非常重要

### 三元组存储和SPARQL

三元组：(主语,谓语,宾语)

例如：(吉姆,喜欢,香蕉)

三元组的主语相当于图中的一个顶点。谓语宾语是下面两情况之一：

1. 原始数据类型中的值，例如字符串或数字。在这种情况下，三元组的谓语和宾语相当于 主语顶点上的属性的键和值。例如， (lucy, age, 33) 就像属性 {“age”：33} 的顶点 lucy。
2. 图中的另一个顶点。在这种情况下，谓语是图中的一条边，主语是其尾部顶点，而宾语 是其头部顶点。例如，在 (lucy, marriedTo, alain) 中主语和宾语 lucy 和 alain 都是顶 点，并且谓语 marriedTo 是连接他们的边的标签。



```turtle
@prefix : <urn:example:>.
_:lucy a :Person.
_:lucy :name "Lucy".
_:lucy :bornIn _:idaho.
_:idaho a :Location.
_:idaho :name "Idaho".
_:idaho :type "state".
_:idaho :within _:usa.
_:usa a :Location
_:usa :name "United States"
_:usa :type "country".
_:usa :within _:namerica.
_:namerica a :Location
_:namerica :name "North America"
_:namerica :type :"continent"
```

图的顶点：` _：someName`

当谓语表示边时，该宾语是一 个顶点，如` _:idaho :within _:usa` 。当谓语是一个属性时，该宾语是一个字符串，如 `_:usa :name "United States"`

更简洁的写法（分号隔开，说明关于同一个主语的多个事情）：

```turtle
@prefix : <urn:example:>.
_:lucy a :Person; :name "Lucy"; :bornIn _:idaho.
_:idaho a :Location; :name "Idaho"; :type "state"; :within _:usa
_:usa a :Loaction; :name "United States"; :type "country"; :within _:namerica.
_:namerica a :Location; :name "North America"; :type "continent".

```

### 语义网络

### RDF(资源描述框架)数据模型

上面的Turtle语言是一种用于RDF数据的人类可读格式

Apache Jena可以根据需要在不同的RDF格式之间自动转换

```xml
<rdf:RDF
    xmlns="urn:example:"
    xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <Location rdf:nodeID="idaho">
        <name>Idaho</name>
        <type>state</type>
        <within>
            <Location rdf:nodeID="usa">
                <name>United States</name>
                <type>country</type>
                <within>
                    <Location rdf:nodeID="namerica">
                        <name>North America</name>
                        <type>continent</type>
                    </Location>
                </within>
            </Location>
        </within>
    </Location>
    <Person rdf:nodeID="lucy">
        <name>Lucy</name>
        <bornIn rdf:nodeID="idaho"/>
    </Person>
</rdf:RDF>

```

### SPARQL查询语言

```SPARQL
PREFIX : <urn:example:>
SELECT ?personName WHERE {
    ?person :name ?personName.
    ?person :bornIn / :within* / :name "United States".
    ?person :livesIn / :within* / :name "Europe".
}

```

sparql和cypher(sparql的变量以?开头):

```
(person) -[:BORN_IN]-> () -[:WITHIN*0..]-> (location) # Cypher
?person :bornIn / :within* ?location. # SPARQL
```

因为RDF不区分属性和边，而只是将它们作为谓语，所以可以使用相同的语法来匹配属性。 在下面的表达式中，变量 usa 被绑定到任意具有值为字符串 "United States" 的 name 属性的 顶点：

```
(usa {name:'United States'}) # Cypher
?usa :name "United States". # SPARQL
```



### Datalog

三元组写成：谓语(主语，宾语)

```Datalog
name(namerica, 'North America').
type(namerica, continent).
name(usa, 'United States').
type(usa, country).
within(usa, namerica).
name(idaho, 'Idaho').
type(idaho, state).
within(idaho, usa).
name(lucy, 'Lucy').
born_in(lucy, idaho).

```

Datalog是Prolog的子集

## 小结

在历史上，数据最开始被表示为一棵大树（层次数据模型），但是这不利于表示多对多的关 系，所以发明了关系模型来解决这个问题

最近，开发人员发现一些应用程序也不适合采用 关系模型

。新的非关系型“NoSQL”数据存储在两个主要方向上存在分歧：

1. 文档数据库的应用场景是：数据通常是自我包含的，而且文档之间的关系非常稀少。 
2. 图形数据库用于相反的场景：任意事物都可能与任何事物相关联。

文档数据库和图数据库有一个共同点，那就是它们通常不会为存储的数据强制一个模式，