# go-word-filter

基于trie树的敏感词过滤服务， 支持grpc调用


## 安装

```
$ go get github.com/cyuxlif/go-word-filter

$ cd $GOPATH/src/github.com/cyuxlif/go-word-filter/conf
$ cp app.yaml.example app.yaml
$ vim app.yaml #配置字典路径， 监听地址
$ cd $GOPATH/src/github.com/cyuxlif/go-word-filter/cmd/word_filter
$ go build
$ ./word_filter -c /yourpath/conf/app.yaml
```


## 使用

参考 $GOPATH/src/github.com/cyuxlif/go-word-filter/cmd/grpc_client

部分代码如：
```
	// 连接
	conn, err := grpc.Dial(Addr, grpc.WithInsecure())

	if err != nil {
		pp.Println(err)
	}

	defer conn.Close()

	// 初始化客户端
	c := pb.NewWordFilterClient(conn)

	// 调用方法
	Text := new(pb.Text)
	Text.Text = `有了 gRPC， 我们可以一次性的在一个 .proto 文件中定义服务并使用任何支持它的操你妹语言去实现客户端和服务器，反过来，
	它们可以在各种环境中，你妈了个逼从Google的服务器到你自己的平板电脑—— gRPC 帮你解决了性吧春暖花开不同语言及环境间通信的复杂性.
	使用 protocol buffers 还能获得波多野结衣全集其他好处，包括高效的序列号，简单的 IDL 以及容易进行操你妈接口更新`
	r, err := c.FindKeyWords(context.Background(), Text)
	if err != nil {
		pp.Println(err)
	}
	pp.Println(r.KeyWords)
```
将会打印
```
[]*word_filter.KeyWord{
  &word_filter.KeyWord{
    Word:    "操你妹",
    Attr:    1,
    Replace: "我傻x",
  },
  &word_filter.KeyWord{
    Word:    "你妈了个逼",
    Attr:    2,
    Replace: "",
  },
  &word_filter.KeyWord{
    Word:    "性吧春暖花开",
    Attr:    1,
    Replace: "",
  },
  &word_filter.KeyWord{
    Word:    "波多野结衣全集",
    Attr:    1,
    Replace: "",
  },
  &word_filter.KeyWord{
    Word:    "操你妈",
    Attr:    1,
    Replace: "",
  },
}

```

## 字典格式

这里的字典格式指的是data目录下的txt的字典格式
词典格式为文本格式：

关键词|词性（可选，默认是审核词 1代表替换词，2代表审核词， 3代表禁止词）|替换词（可选）

如:

操你妈|1 //替换词

你妈逼|1|你妹妹 //替换词

卖毒品|3 //禁止词

你妈了个逼|2 //审核词
