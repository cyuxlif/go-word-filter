package word_filter

import (
	"net"

	"google.golang.org/grpc"
	"golang.org/x/net/context"
	log "github.com/sirupsen/logrus"

	pb "github.com/cyuxlif/go-word-filter/pb/word_filter"
)

const VERSION = "v1.0.0"

type App struct {
	dict        *Dict
	tcpListener net.Listener
}

func New(opts *Options) *App {
	return &App{
		dict: NewDict(opts),
	}
}

func (app *App) Run(addr string) {
	//加载词典
	app.dict.LoadDict()

	listen, err := net.Listen("tcp", addr)
	app.tcpListener = listen
	if err != nil {
		log.WithError(err).WithField("addr", addr).Fatal("faild listen")
	}

	// 实例化grpc Server
	s := grpc.NewServer()

	// 注册Service
	pb.RegisterWordFilterServer(s, app)

	log.WithField("addr", addr).Info("listen on")
	s.Serve(app.tcpListener)
}

//rpc handler
//根据文本查找关键词
func (app *App) FindKeyWords(ctx context.Context, t *pb.Text) (*pb.FindKeyWordsRes, error) {
	log.WithField("text", t.Text).Debug("FindKeyWords")
	keyWords := app.dict.FindKeyWords(t.Text)
	outputKeyWords := make([]*pb.KeyWord, 0, len(keyWords))
	for _, v := range keyWords {
		outPutKeyWord := &pb.KeyWord{
			Word:    v.Word,
			Attr:    int32(v.Attr),
			Replace: v.Replace,
		}
		outputKeyWords = append(outputKeyWords, outPutKeyWord)
	}
	keyWordsRes := &pb.FindKeyWordsRes{
		KeyWords: outputKeyWords,
	}
	return keyWordsRes, nil
}

//rpc handler
//获取用户词典内容
func (app *App) GetUserDict(ctx context.Context, empty *pb.Empty) (*pb.Text, error) {
	b, err := app.dict.GetUserDict()
	if err != nil {
		return nil, err
	}
	return &pb.Text{Text: string(b)}, nil
}

//rpc handler
//修改用户词典内容
func (app *App) EditUserDict(ctx context.Context, t *pb.Text) (*pb.Empty, error) {
	err := app.dict.EditUserDict([]byte(t.Text))
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
