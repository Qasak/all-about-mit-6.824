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



