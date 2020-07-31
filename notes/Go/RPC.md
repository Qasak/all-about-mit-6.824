

rpc包通过网络或其他I/O连接提供对对象导出方法的访问。使对象的名称可见。注册后，可以远程访问对象的导出方法。服务器可以注册不同类型的多个对象（服务），但注册同一类型的多个对象是错误的。

只有满足这些条件的方法才可用于远程访问；其他方法将被忽略：

-方法的类型是导出的。

-方法被导出。

-该方法有两个参数，都是导出（或内置）类型。

-方法的第二个参数是指针。

-方法具有返回类型错误。



实际上，该方法必须形如

```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

其中T1和T2可以通过encoding/gob进行封送。即使使用不同的编解码器，这些要求也适用。（将来，这些要求可能会对自定义编解码器有所软化。）

方法的第一个参数表示调用方提供的参数；第二个参数表示要返回给调用方的结果参数。方法的返回值（如果为非nil）将作为一个字符串返回，客户端产看是否是被 errors.New创建。如果返回错误，则不会将reply参数发送回客户端。

服务器可以通过调用ServeConn来处理单个连接上的请求。更典型的是，它将创建一个网络侦听器并调用Accept，对于HTTP侦听器，则调用HandleHTTP和http.Serve。 

客户端希望用服务建立连接，然后在连接上调用NewClient。方便的函数Dial（DialHTTP）对原始网络连接（HTTP连接）执行这两个步骤。生成的客户端对象有两个方法，Call和Go，它们指定要调用的服务和方法、一个包含参数的指针和一个接收结果参数的指针。

Call方法等待远程调用完成，而Go方法异步启动调用，并使用Call结构的Done通道发出完成信号。

除非设置了显式编解码器，否则将使用包encoding/gob来传输数据。

这里有一个简单的例子。服务器希望导出Arith类型的对象：

```go
package server

import "errors"

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
```

服务器调用(HTTP服务)

```go
arith := new(Arith)
rpc.Register(arith)
rpc.HandleHTTP()
l, e := net.Listen("tcp", ":1234")
if e != nil {
	log.Fatal("listen error:", e)
}
go http.Serve(l, nil)
```

此时，客户机可以看到一个带有方法的“Arith”服务算术乘法“和”算术除法". 要调用一个，客户端首先拨打服务器：

```go
client, err := rpc.DialHTTP("tcp", serverAddress + ":1234")
if err != nil {
	log.Fatal("dialing:", err)
}
```

然后远程调用

```go
// Synchronous call
args := &server.Args{7,8}
var reply int
err = client.Call("Arith.Multiply", args, &reply)
if err != nil {
	log.Fatal("arith error:", err)
}
fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
```

或

```go
// Asynchronous call
quotient := new(Quotient)
divCall := client.Go("Arith.Divide", args, quotient, nil)
replyCall := <-divCall.Done	// will be equal to divCall
// check errors, print, etc.
```

服务器实现通常会为客户机提供一个简单的、类型安全的包装器。



net/rpc包已冻结，不接受新功能。