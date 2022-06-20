# 数据编码与演化



> Everything changes and nothing stands still  

## 数据编码格式

许多编程语言都内建了将内存对象编码为字节序列的支持。 例如， Java
有 java.io.Serializable 【1】 ， Ruby有 Marshal 【2】 ， Python有 pickle 【3】 等等。 许多
第三方库也存在， 例如 Kryo for Java 【4】 。  

使用语言内置的编码方案通常不是个好主意，除非只是为了临时尝试。  

### JSON, XML, 二进制变体

JSON， XML和CSV是文本格式， 因此具有人类可读性（ 尽管语法是一个热门辩题） 。 除了表
面的语法问题之外， 它们也有一些微妙的问题：  

+ 数字的编码多有歧义之处。 XML和CSV不能区分数字和字符串（ 除非引用外部模式） 。
  JSON虽然区分字符串和数字， 但不区分整数和浮点数， 而且不能指定精度。
+ 当处理大量数据时， 这个问题更严重了。 例如， 大于$2^{53}$的整数不能在IEEE 754双
  精度浮点数中精确表示， 因此在使用浮点数（ 例如JavaScript） 的语言进行分析时， 这些
  数字会变得不准确。 Twitter上有一个大于$2^{53}$的数字的例子， 它使用一个64位的数
  字来标识每条推文。 Twitter API返回的JSON包含了两种推特ID， 一个JSON数字， 另一
  个是十进制字符串， 以此避免JavaScript程序无法正确解析数字的问题【10】 。
+ JSON和XML对Unicode字符串（ 即人类可读的文本） 有很好的支持， 但是它们不支持二
  进制数据（ 不带字符编码(character encoding)的字节序列） 。 二进制串是很实用的功
  能， 所以人们通过使用Base64将二进制数据编码为文本来绕开这个限制。 模式然后用于
  表示该值应该被解释为Base64编码。 这个工作， 但它有点hacky， 并增加了33％的数据
  大小。 XML 【11】 和JSON 【12】 都有可选的模式支持。 这些模式语言相当强大， 所以
  学习和实现起来相当复杂。 XML模式的使用相当普遍， 但许多基于JSON的工具嫌麻烦
  才不会使用模式。 由于数据的正确解释（ 例如数字和二进制字符串） 取决于模式中的信
  息， 因此不使用XML/JSON模式的应用程序可能需要对相应的编码/解码逻辑进行硬编
  码。
+ CSV没有任何模式， 因此应用程序需要定义每行和每列的含义。 如果应用程序更改添加
  新的行或列， 则必须手动处理该变更。 CSV也是一个相当模糊的格式（ 如果一个值包含
  逗号或换行符， 会发生什么？ ） 。 尽管其转义规则已经被正式指定【13】 ， 但并不是所
  有的解析器都正确的实现了标准。  

### 二进制编码

![img](https://github.com/Qasak/distributed-system/blob/master/%E6%95%B0%E6%8D%AE%E7%B3%BB%E7%BB%9F%E7%9A%84%E5%9F%BA%E7%9F%B3/%E7%BC%96%E7%A0%81%E4%B8%8E%E6%BC%94%E5%8C%96/messagepack-json.png)

### Thrift, Protocol Buffers

![img](https://github.com/Qasak/distributed-system/blob/master/%E6%95%B0%E6%8D%AE%E7%B3%BB%E7%BB%9F%E7%9A%84%E5%9F%BA%E7%9F%B3/%E7%BC%96%E7%A0%81%E4%B8%8E%E6%BC%94%E5%8C%96/Thrift%20BinaryProtocol.png)

### 数据流模式

