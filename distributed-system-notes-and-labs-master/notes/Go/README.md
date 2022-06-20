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

一个结构体就是一组字段(field)

结构体字段用点号访问

```go
package main

import "fmt"

type Vertex struct {
	X int
	Y int
}

func main() {
	v := Vertex{1, 2}
	v.X = 4
	fmt.Println(v.X)
}

```

结构体指针

> 如果我们有一个指向结构体的指针 `p`，那么可以通过 `(*p).X` 来访问其字段 `X`。不过这么写太啰嗦了，所以语言也允许我们使用隐式间接引用，直接写 `p.X` 就可以。

```go
v := Vertex{1, 2}
p := &v
p.X = 1e9
```



结构体文法

```go
type Vertex struct {
	X, Y int
}

var (
	v1 = Vertex{1, 2}  // 创建一个 Vertex 类型的结构体
	v2 = Vertex{X: 1}  // Y:0 被隐式地赋予
	v3 = Vertex{}      // X:0 Y:0
	p  = &Vertex{1, 2} // 创建一个 *Vertex 类型的结构体（指针）
    				  // 特殊的前缀 & 返回一个指向结构体的指针。
)
```







### slice切片

切片是一个具有三项内容的描述符， 包含一个指向（ 数组内部） 数据的指针、 长度以及容量，在这三项被初始化之前， 该切片为 nil

#### range

+ for 循环的range形式可以遍历slice或map

+ 当使用 for 循环遍历切片时，每次迭代都会返回两个值。第一个值为当前元素的下标，第二个值为该下标所对应元素的一份*复制*(copy)。

  ```go
  var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}
  	for i, v := range pow {
  		fmt.Printf("2**%d = %d\n", i, v)
  	}
  ```

  

### map

### 方法和接口

#### 方法

Go没有类，不过你可以为结构体类型定义方法

方法就是一类带特殊的`接收者`参数的函数

方法接收者在它自己的参数列表内，位于func关键字和方法名之间

```go
package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
	v := Vertex{3, 4}
	fmt.Println(v.Abs())
}
```

> 此例中，`Abs` 方法拥有一个名为 `v`，类型为 `Vertex` 的接收者。

方法即函数，下面的函数和上面的方法功能是一样的

```go
func Abs(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
```



可以为非结构体类型声明方法

```go
package main

import (
	"fmt"
	"math"
)

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func main() {
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())
}

```

> 接收者的类型定义和方法声明必须在同一包内；不能为内建类型声明方法

+ 指针接收者

  对于类型T，接收者的类型可以用\*T。(此外，T不能是像*int这样的指针)

  ```go
  func (v *Vertex) Scale(f float64) {
  	v.X = v.X * f
  	v.Y = v.Y * f
  }
  ```

  > 指针接收者的方法可以修改接收者指向的值

  方法与指针重定向

  带指针参数的函数必须接受一个指针

  而以指针为接收者的方法被调用时，接收者既能为值又能为指针

  这是因为：当v是个值而非指针时，Go 会将语句 v.Scale(5) 解释为 (&v).Scale(5)。

  ```go
  package main
  
  import "fmt"
  
  type Vertex struct {
  	X, Y float64
  }
  
  func (v *Vertex) Scale(f float64) {
  	v.X = v.X * f
  	v.Y = v.Y * f
  }
  
  func ScaleFunc(v *Vertex, f float64) {
  	v.X = v.X * f
  	v.Y = v.Y * f
  }
  
  func main() {
  	v := Vertex{3, 4}
  	v.Scale(2)
  	ScaleFunc(&v, 10)
  
  	p := &Vertex{4, 3}
  	p.Scale(3)
  	ScaleFunc(p, 8)
  
  	fmt.Println(v, p)
  }
  ```

  



#### 接口

*接口类型*是由一组`方法签名`定义的集合

An *interface type* is defined as a set of `method signatures`.

接口类型的变量可以保存任何实现了这些方法的值。



### 并发

```go
var done sync.WaitGroup
for _, u := range urls {
    done.Add(1)
    u2 := u
    go func() {
        defer done.Done()
        ConcurrentMutex(u2, fetcher, f)
    }()
    //go func(u string) {
    //	defer done.Done()
    //	ConcurrentMutex(u, fetcher, f)
    //}(u)
}
done.Wait()
```

> 为了既调用ConcurrentMutex又调用Done

#### 信道channel

用make创建

默认值是零， 表示不带缓冲的或*同步*的信道。  

无缓冲信道在通信时会同步交换数据， 它能确保（ 两个 Go 程的） 计算处于确定状态。  

```go
ch:=make(chan int)
```

可以带缓冲

```go
ch := make(chan int, 100)
```

仅当信道的缓冲区填满后，向其发送数据时才会阻塞。当缓冲区为空时，接受方会阻塞。











### 内存

### new 分配

#### new

new(T) 会为类型为 T 的新项分配已置零的内存空间， 并返回它的地址  



每当获取一个复合字面(composite literals)的地址时， 都将为一个新的实例分配内存  

以字段: 值 对的形式明确地标出元素， 初始化
字段时就可以按任何顺序出现， 未给出的字段值将赋予零值。   

```go
func NewFile(fd int, name string) *File {
    if fd < 0 {
    	return nil
    } 
    return &File{fd: fd, name: name}
}
```



#### make

它只用于创建切片、映射和信道并返回类型为 T（ 而非 *T ） 的一个已初始化 （ 而非置零） 的值  

出现这种用差异的原因在于， 这三种类型**本质上为引用**数据类型， 它们在使用前必须初始化。  

```go
make([]int, 10, 100)
make([]int, 10)
```

会分配一个具有 100 个 int 的数组空间， 接着创建一个长度为 10， 容量为 100 并指向该数组
中前 10 个元素的切片结构  



(生成切片时， 其容量可以省略 )

#### copy

```go
func main() {

    slice := []int{0, 1, 2, 3, 4}
    slice2 := slice[1:4]

    slice4 := make([]int, len(slice2))

    copy(slice4, slice2)

    fmt.Printf("slice %v, slice4 %v \n", slice, slice4)
    slice[1] = 1111
    fmt.Printf("slice %v, slice4 %v \n", slice, slice4)
}
```

slice4是从slice2中copy生成，slice和slice4底层的匿名数组是不一样的。因此修改他们不会影响彼此

### 数组arrary

数组是值。 将一个数组赋予另一个数组会复制其所有元素。
特别地， 若将某个数组传入某个函数， 它将接收到该数组的一份副本而非指针。
数组的大小是其类型的一部分。 类型 [10]int 和 [20]int 是不同的。  

数组为值的属性很有用， 但代价高昂； 若你想要 C 那样的行为和效率， 你可以传递一个指向
该数组的指针。  

```go
func Sum(a *[3]float64) (sum float64) {
    for _, v := range *a {
    	sum += v
    } 
    return
} 
array := [...]float64{7.0, 8.5, 9.1}
x := Sum(&array) // 注意显式的取址操作
```

