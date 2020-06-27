# 数据编码与演化



> Everything changes and nothing stands still  

## 数据编码格式

许多编程语言都内置支持将内存中的对象编码为字节序列。例如， J av a有 j av a.
io. Serializable , Ruby有Marshal , Python有 pickle l3l等。此外，还有许多第
三方库，例如用于Java 的Kryo  。  

使用语言内置的编码方案通常不是个好主意，除非只是为了临时尝试。  