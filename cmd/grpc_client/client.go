package main

import (

	"google.golang.org/grpc"
	"github.com/k0kubun/pp"
	"golang.org/x/net/context"

	pb "github.com/cyuxlif/go-word-filter/pb/word_filter"
)

const (
	// Address gRPC服务地址
	Addr = "127.0.0.1:7890"
)

func main(){
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
	Text.Text = `操你妈
操你妈妈
操你大爷
操你妹|1|我傻x
操你菊花|1
卖毒品|3
你妈了个逼|2|
波多野结衣全集
fuck you
fuck your mother`
	_, err =  c.EditUserDict(context.Background(), Text)
	if err != nil {
		pp.Println(err)
	}

	t, err := c.GetUserDict(context.Background(), &pb.Empty{})
	if err != nil {
		pp.Println(err)
	}
	pp.Println(t.Text)


	Text.Text = `有了 gRPC， 我们可以一次性的在一个 .proto 文件中定义服务并使用任何支持它的操你妹语言去实现客户端和服务器，反过来，
	它们可以在各种环境中，你妈了个逼从Google的服务器到你自己的平板电脑—— gRPC 帮你解决了性吧春暖花开不同语言及环境间通信的复杂性.
	使用 protocol buffers 还能获得波多野结衣全集其他好处，包括高效的序列号，简单的 IDL 以及容易进行操你妈接口更新`
	r, err := c.FindKeyWords(context.Background(), Text)
	if err != nil {
		pp.Println(err)
	}
	pp.Println(r.KeyWords)
}
