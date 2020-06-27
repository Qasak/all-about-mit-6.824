# 数据编码与演化



> Everything changes and nothing stands still  

## 数据编码格式

许多编程语言都内建了将内存对象编码为字节序列的支持。 例如， Java
有 java.io.Serializable 【1】 ， Ruby有 Marshal 【2】 ， Python有 pickle 【3】 等等。 许多
第三方库也存在， 例如 Kryo for Java 【4】 。  

使用语言内置的编码方案通常不是个好主意，除非只是为了临时尝试。  

### JSON, XML, 二进制变体

