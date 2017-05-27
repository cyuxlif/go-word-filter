package word_filter

import (
	"testing"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/k0kubun/pp"
)

func TestDict_Load(t *testing.T) {
	sr := strings.NewReader(`
操你妈
操你妈妈
操你大爷
操你妹|1|操你姐
操你菊花|1
卖毒品|3
江泽民|2|
习近平|2
波多野结衣全集
fuck you
fuck your mother`)
	dict := NewDict(&Options{})
	dict.Load(sr)

	pp.Println(dict.Trie)
}

func TestDict_FindKeyWords(t *testing.T) {

	s := `我们的例子是一夜情交友聊天室一个简单的路由映射性吧春暖花开的应用，它允激情成人网络电视许客户端获取路由特性的信息，
	生成路由的总结，以及交互路由信息，如服务器和其他客户端的流量更新老挝国营赌场，务并使用任何支持它的语言去实现客户端和服务器。`
	opts := NewOptions("./conf/app.yaml")
	w := New(opts)
	keyWords := w.dict.FindKeyWords(s)

	assertKeyWords := []*KeyWord{
		&KeyWord{Word: "一夜情交友聊天室", Attr: 1},
		&KeyWord{Word: "性吧春暖花开", Attr: 1},
		&KeyWord{Word: "激情成人网络电视", Attr: 1},
		&KeyWord{Word: "老挝国营赌场", Attr: 1},
	}
	for k, keyWord := range keyWords {
		assert.Equal(t, keyWord, assertKeyWords[k])
	}
}

func BenchmarkDict_FindKeyWords(b *testing.B) {
	s := `有了 gRPC， 我们可以一次性的在一个 .proto 文件中定义服务并使用任何支持它的语言去实现客户端和服务器，反过来，
	它们可以在各种环境中，从Google的服务器到你自己的平板电脑—— gRPC 帮你解决了性吧春暖花开不同语言及环境间通信的复杂性.
	使用 protocol buffers 还能获得其他好处，包括高效的序列号，简单的 IDL 以及容易进行接口更新`
	opts := NewOptions("./conf/app.yaml")
	w := New(opts)
	w.dict.LoadDict()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.dict.FindKeyWords(s)
	}
}
