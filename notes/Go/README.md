## Go

+ 好用的RPC包

+ 内存安全

  + 垃圾回收

+ 类型安全

  

## 线程

+ goroutine
  + 通过channel同步
  + 也可使用mutex

+ I/O并发
+ 并行

> 异步(事件驱动)：只有一个循环，单核，单线程，高性能
>
> I/O多路复用是事件驱动的基础

### 线程挑战

+ 竞争Race

  + 用mutex解决

+ 协作Coordination

  + 用channel实现
  + sync.cond
  + waitGroup

+ 死锁Deadlock

  

## Tutorial

### 包，变量和函数

+ 包

  每个Go程序都是由包构成

  程序从main包开始运行

+ 导入

  ```go
  import (
  	"fmt"
  	"sync"
  )
  ```

+ 导出名Exported names

  + 一个已大写字母开头的名字

  + 导入一个包时，只能引用其中已导出的名字

    ```go
    math.Pi
    ```

+ 函数

  可接收没有参数或多个参数

  ```go
  func add(x int, y int) int {
  	return x + y
  }
  ```

  + 函数的两个或多个已命名形参类型相同时，除最后一个类型外，其他可省略

  ```go
  func add(x, y int) int {
  	return x + y
  }
  ```

  + 任意数量返回值

    ```go
    func swap(x, y string) (string, string) {
    	return y, x
    }
    ```

    ```go
    func main() {
    	a, b := swap("hello", "world")
    	fmt.Println(a, b)
    }
    ```

  + 命名返回值

    ```go
    func split(sum int) (x, y int) {
    	x = sum * 4 / 9
    	y = sum - x
    	return
    }
    ```

    视作命名在函数顶部的变量

    具有一定意义，可作为文档使用

    没有参数的return 返回已命名的返回值

    仅用在短函数中，长函数中回影响代码可读性

+ 变量

  var语句用于声明一个变量列表，类型在最后

  ```go
  var i int
  ```

  可在包或函数级别

  ```go
  var c, python, java bool
  
  func main() {
  	var i int
  	fmt.Println(i, c, python, java)
  }
  ```

  变量初始化

  ​	可以包含初始值，每个变量对应一个

  ​	若初始化值已存在，可省略类型；变量会从初始值中获得类型

  短变量声明

  ​	在函数中:=可在类型明确的地方代替var声明

  ​	在函数外的每个语句必须以关键字开始(var func 等)，因此`:=`结构**不能**在函数外使用

### 流程控制

### struct

### slice

### map

### 方法和接口

### 并发

