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

  ​	在函数外的每个语句必须以关键字开始(var，func 等)，因此`:=`结构**不能**在函数外使用

+ 基本类型

  ```go
  bool
  string
  int int8 int16 int32 int64
  uint uint8 uint16 uint32 uint64 uintptr
  byte // uint8的别名
  rune // int32的别名
       // 表示一个Unicode 码点(Unicode code point)
       //(也就是编号，Unicode行列中相交的点)
  float32 float64
  complex64 complex128
  ```

  int, uint 和 uintptr 在 32 位系统上通常为 32 位宽，在 64 位系统上则为 64 位宽。 当你需要一个整数值时应使用 int 类型，除非你有特殊的理由使用固定大小或无符号的整数类型。

  ```go
  var (
  	ToBe   bool       = false
  	MaxInt uint64     = 1<<64 - 1
  	z      complex128 = cmplx.Sqrt(-5 + 12i)
  )
  
  func main() {
  	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
  	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
  	fmt.Printf("Type: %T Value: %v\n", z, z)
  }
  ```

+ 零值

  没有明确初始值的变量声明会被赋予零值

  ​	数值类型 0

  ​	布尔类型 false

  ​	字符串 ""

+ 类型转换

  `T(v)`将v值转换为类型T

  ```go
  var i int = 42
  var f float64 = float64(i)
  var u uint = uint(f)
  
  i := 42
  f := float64(i)
  u := uint(f)
  ```

  与 C 不同的是，Go 在不同类型的项之间赋值时需要显式转换

+ 类型推导

  声明一个变量而不指定其类型时，变量的类型由*右值*(等号右边的值)推到得出

  右值声明了类型：

  ```go
  var i int
  j := i // j 也是一个 int
  ```

  右值为数值常量：

  ```go
  i := 42           // int
  f := 3.142        // float64
  g := 0.867 + 0.5i // complex128
  ```

+ 常量

  与变量类似，但用`const`关键字

  不能用`:=`

  数值常量：

  ​	是高精度的值

  ​	一个未指定类型的常量由上下文决定

  ```go
  const (
  	// 将 1 左移 100 位来创建一个非常大的数字
  	// 即这个数的二进制是 1 后面跟着 100 个 0
  	Big = 1 << 100
  	// 再往右移 99 位，即 Small = 1 << 1，或者说 Small = 2
  	Small = Big >> 99
  )
  ```

  

  

### 流程控制

#### for

Go仅有的一种循环

```go
for i := 0; i < 10; i++ {
    sum += i
}
```

注意：和 C、Java、JavaScript 之类的语言不同，Go 的 for 语句后面的三个构成部分外没有小括号， 大括号 { } 则是必须的



初始化语句和后置语句是可选的。

```go
sum := 1
for sum < 1000 {
    sum += sum
}
```



如果省略循环条件，该循环就不会结束

```go
for {
}
```





#### if

#### switch

#### defer



### struct

### slice

### map

### 方法和接口

### 并发

