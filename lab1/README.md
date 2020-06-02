## MapReduce

我们已经给了你一个开始的小部分代码。master和worker的“main”流程 main/mrmaster.go/和main/mrworker/go；不要更改这些文件。把你的实现放在mr/master.go，mr/worker.go，和mr/rpc.go.

下面是如何在word count MapReduce应用程序上运行代码。首先，确保word count插件是新构建的：

```shell
$ go build -buildmode=plugin ../mrapps/wc.go
```

在main 文件夹，运行master

```shell
$ rm mr-out*
$ go run mrmaster.go pg-*.txt
```

给mrmaster.go的参数pg-*.txt是输入文件，每个文件对应于一个“split”

是一个Map任务的输入。

在一个或多个其他窗口中，运行一些worker

```shell
$ go run mrworker.go wc.so
```

当worker和master完成后，看看mr out-*中的输出。完成lab后，输出文件的有序对(sorted union)应与顺序输出匹配，如下所示：

```shell
$ cat mr-out-* | sort | more
A 509
ABOUT 2
ACT 8
...
```

我们在main/test-mr.sh中为您提供了一个测试脚本。测试检查wc和indexer MapReduce应用程序在给定pg-xxx.txt文件文件作为输入时是否生成正确的输出。测试还检查您的实现是否并行运行Map和Reduce任务，以及您的实现是否从运行任务时崩溃的工作进程中恢复。

如果现在运行测试脚本，它将挂起，因为master永远不会完成

```shell
$ cd ~/6.824/src/main
$ sh test-mr.sh
*** Starting wc test.
```

您可以在mr/master.go的Done函数中将ret:=false更改为true，一遍让master立即退出。然后：

```shell
$ sh ./test-mr.sh
*** Starting wc test.
sort: No such file or directory
cmp: EOF on mr-wc-all
--- wc output is not the same as mr-correct-wc.txt
--- wc test: FAIL
$
```

测试脚本希望看到名为mr-out-X的文件中的输出，每个reduce任务一个。mr的空白mr/master.go以及mr/worker.go.不要生成那些文件（或者做很多其他事情），这样测试就失败了。

完成后，测试脚本输出应如下所示：

```shell
$ sh ./test-mr.sh
*** Starting wc test.
--- wc test: PASS
*** Starting indexer test.
--- indexer test: PASS
*** Starting map parallelism test.
--- map parallelism test: PASS
*** Starting reduce parallelism test.
--- reduce parallelism test: PASS
*** Starting crash test.
--- crash test: PASS
*** PASSED ALL TESTS
$
```

您还将看到Go-RPC包中的一些错误，这些错误看起来像

```shell
2019/12/16 13:27:09 rpc.Register: method "Done" has 1 input parameters; needs exactly three
```

忽略这些信息

### 一些规则:

- 映射阶段应该将中间键(intermediate keys)划分到桶(buckets)`nReduce` reduce任务, nReduce` 是 `main/mrmaster.go` 传给 `MakeMaster()`的参数.
- worker 实现应该把 X'th 个reduce 任务输出放到  `mr-out-X`文件中.
- 一个 `mr-out-X` 文件应为每个Reduce函数输出一行. 行应该生成为 Go `"%v %v"`形式, 被key和value调用. 查看 `main/mrsequential.go` 中注释为"this is the correct format"的行. 如果您的实现偏离此格式太多，则测试脚本将失败。
- 你可以修改 `mr/worker.go`, `mr/master.go`, and `mr/rpc.go`. 你可以在测试时临时修改其他文件, 但请确保您的代码与原始版本兼容；我们将使用原始版本进行测试
- worker应该将中间Map输出放在当前目录中的文件中，稍后您的worker可以将它们作为Reduce任务的输入
- `main/mrmaster.go` 期待 `mr/master.go` 实现一个 `Done()` 方法，该方法在MapReduce作业(job)完成时返回true; 此时 `mrmaster.go` 将退出.
- 当作业完全完成时，应退出worker进程. 实现这一点的一个简单方法是使用`call()`返回值:如果worker进程无法联系master进程，则可以假定主进程已经退出，因为该作业已经完成，因此master进程也可以终止.取决与您的设计，您可能会发现拥有一个“请退出”伪任务（pseudo-task）（master可以将其交给workers）很有帮助。

### Hints

- 一种开始着手的办法是修改 `mr/worker.go`的 `Worker()` 将 RPC发送到master来请求一个任务。然后修改master，以一个尚未启动(as-yet unstarted)map任务的文件名响应。然后修改worker以读取该文件并调用application Map函数,如 `mrsequential.go`.
-  Map and Reduce 应用函数在运行时(run-time) 从`.so`文件中被Go plugin包加载.
- 如果您更改mr/目录中的任何内容，可能需要重新构建您使用的任何MapReduce插件.with something like `go build -buildmode=plugin ../mrapps/wc.go`
- 这个lab依赖于workers共享一个文件系统. 当所有工人运行在同一台机器上时，这很简单(straightforward), 但是如果worker运行在不同的机器上，则需要像GFS这样的全局文件系统
- 中间文件的合理命名约定(convention)是`mr-X-Y`，其中X是Map任务号，Y是reduce任务号。
- worker的map任务代码将需要一种方法来存储文件中的中间键/值对，以便在reduce任务期间正确读取。一种可能是使用Go的encoding/json包。要将键/值对写入JSON文件，请执行以下操作：

```go
  enc := json.NewEncoder(file)
  for _, kv := ... {
    err := enc.Encode(&kv)
```

读回文件(read such a file back)

```go
  dec := json.NewDecoder(file)
  for {
    var kv KeyValue
    if err := dec.Decode(&kv); err != nil {
      break
    }
    kva = append(kva, kv)
  }
```

- 您的worker的map部分可以使用`ihash(key)`函数(in `worker.go`)为给定的键选择reduce任务。
- 你可以从 `mrsequential.go` 搬运一些代码来读取Map输入文件,用来排序intermedate Map和Reduce之间的中间k/v对, 以及用来存储Reduce到文件的输出.
- master作为一个 RPC server, 将是并发的;不要忘记锁定共享数据.
- 使用 Go的race检测器, with `go build -race` and `go run -race`. `test-mr.sh` 有一条注释，说明如何为测试启用race检测器
- Workers 有些时候需要等待, e.g. 最后一个map完成后reduces才能开始. 一种可行的犯法是  workers周期性的向master请求work,  在每个请求之间休眠`time.Sleep()` . 另一种办法是 master 中相关的PRC处理程序(PRC handler)循环等待(have a loop that waits), with `time.Sleep()` or `sync.Cond`. Go在自己的线程中为每个RPC运行处理程序(handler)，因此一个处理程序正在等待不会阻止master处理其他RPCs。
- master 不能可靠地区分崩溃的workers, 存活但由于某种原因卡住/停滞的workers .你能做的最好的事情就是让master等上一段时间，然后放弃，把任务重新分配给另一个worker。对于这个lab，让master等10秒钟；之后master应该假设worker已经挂了（当然，可能没有）。
- 要测试崩溃恢复，可以使用`mrapps/crash.go` application plugin. 他在Map和Reduce函数中随机退出.
- 为了确保出线崩溃的情况下部分写入的文件不被观测到MapReduce 论文提到了一个trick：利用一个临时文件，一旦完成写入，就原子重命名它(using a tempororay file and atomically renaming it once it is completely written) 你可以使用 `ioutil.TempFile` 来创建一个临时文件，用 `os.Rename` 原子重命名它.
- `test-mr.sh` 运行 `mr-tmp`子目录中所有进程, 因此如果出现问题，并且您想查看中间文件或输出文件，请查看那里。