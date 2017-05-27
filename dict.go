package word_filter

import (
	"sync"
	"os"
	"bufio"
	"strings"
	"io"
	"unicode/utf8"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/Unknwon/com"
)

type AttrType uint8

const (
	//替换关键词
	REPLACE AttrType = iota + 1 //1
	//审核关键词
	REVIEW  //2
	//禁止关键词
	BAN  //3
)

const SPLIT_FLAG = "|"

//字典
type Dict struct {
	//默认脏词库
	DefaultDictDataPath string
	//自定义脏词库
	UserDictDataPath string

	MaxWordLength int
	//trie 树
	Trie *Trie
	sync.RWMutex
}

type Trie struct {
	Root *Node
}

type Node struct {
	KeyWord *KeyWord //节点末才有数据， 非节点末为nil
	Nodes   map[rune]*Node
}

type KeyWord struct {
	Word    string
	Attr    AttrType
	Replace string
}

func NewDict(opts *Options) *Dict {
	return &Dict{
		DefaultDictDataPath: opts.DefaultDictDataPath,
		UserDictDataPath:    opts.UserDictDataPath,
		Trie: &Trie{
			Root: &Node{
				Nodes: make(map[rune]*Node),
			},
		},
	}
}

//初始化字典 trie树
//词典格式为文本格式  关键词|词性（可选，默认是替换词）|替换词（可选）
//如:
// 彭丽媛          //替换词
// 江泽民|1        //替换词
// 你妈逼|1|你妹妹  //替换词
// 叫小姐|2        //审核词
// 卖毒品|3        //禁止词
func (d *Dict) LoadDict() {
	//载入默认词库 为必须
	defaultDictDataPath := d.DefaultDictDataPath
	log.WithField("dict", defaultDictDataPath).Info("load default dict")
	defaultDictFile, defaultDictErr := os.Open(defaultDictDataPath)
	if defaultDictErr != nil {
		log.WithField("dict", defaultDictDataPath).WithError(defaultDictErr).Fatal("failed load default dict")
	}
	defer defaultDictFile.Close()

	d.Load(defaultDictFile)

	//载入用户词库 非必须
	userDictDataPath := d.UserDictDataPath
	log.WithField("dict", userDictDataPath).Info("load user dict")
	userDictFile, userDictErr := os.OpenFile(userDictDataPath, os.O_RDWR|os.O_CREATE, 0666)
	if userDictErr != nil {
		log.WithError(userDictErr).WithField("dict", userDictFile).Fatal("failed load user dict")
	}
	defer userDictFile.Close()
	d.Load(userDictFile)
}

func (d *Dict) Load(reader io.Reader) {
	d.Lock()
	defer d.Unlock()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		lineText := scanner.Text()
		params := strings.Split(lineText, SPLIT_FLAG)

		if len(params) < 1 {
			continue
		}
		for k, param := range params {
			params[k] = strings.Trim(param, " \n\r\t")
		}
		switch len(params) {
		case 1:
			d.addKeyWord(&KeyWord{Word: params[0], Attr: REPLACE})
		case 2:
			attr, err := com.StrTo(params[1]).Uint8()
			if err != nil {
				log.WithError(err).WithField("text", lineText).Warn("failed load dict line")
				continue
			}
			d.addKeyWord(&KeyWord{Word: params[0], Attr: AttrType(attr)})
		case 3:
			attr, err := com.StrTo(params[1]).Uint8()
			if err != nil {
				log.WithError(err).WithField("text", lineText).Warn("failed load dict line")
				continue
			}
			d.addKeyWord(&KeyWord{Word: params[0], Attr: AttrType(attr), Replace: params[2]})
		}
	}

	if scanner.Err() != nil {
		log.WithError(scanner.Err()).Error("scan error")
	}
}

//添加trie树节点
func (d *Dict) addKeyWord(keyWord *KeyWord) bool {
	n := d.Trie.Root
	isNew := false
	for _, c := range keyWord.Word {
		if _, ok := n.Nodes[c]; !ok {
			n.Nodes[c] = &Node{Nodes:make(map[rune]*Node)}
			isNew = true
		}
		n = n.Nodes[c]
	}

	//词末节点 添加单词属性 和替换词
	if isNew {
		n.KeyWord = keyWord
		lenWord := utf8.RuneCountInString(keyWord.Word)
		if lenWord > d.MaxWordLength {
			d.MaxWordLength = lenWord
		}
	} else {
		//覆盖字典已有 keyWord
		if n.KeyWord != keyWord {
			n.KeyWord = keyWord
		}
	}
	return isNew
}

//查找节点
func (d *Dict) hasWord(wordRune []rune) (*KeyWord, bool) {
	n := d.Trie.Root
	for _, c := range wordRune {
		if _, ok := n.Nodes[c]; !ok {
			return nil, false
		}
		n = n.Nodes[c]
	}
	if n.KeyWord == nil {
		return nil, false
	}
	return n.KeyWord, true
}

//从文本中找出  所有关键词
//采用正向最大匹配算法
func (d *Dict) FindKeyWords(text string) ([]*KeyWord) {
	d.RLock()
	defer d.RUnlock()
	keyWords := make([]*KeyWord, 0)
	maxLen := d.MaxWordLength
	var start int
	lenText := utf8.RuneCountInString(text)
	textRune := []rune(text)
	for start < lenText {
	L:
		for end := min(start + maxLen, lenText); end > start; end-- {
			//匹配到  偏移匹配词的长度
			if keyWord, ok := d.hasWord(textRune[start:end]); ok {
				keyWords = append(keyWords, keyWord)
				start = end
				break L
			}
		}
		//没匹配到向后 偏移1
		start++
	}
	return keyWords
}

//获取用户自定义词库
func (d *Dict) GetUserDict() ([]byte, error) {
	userDictDataPath := d.UserDictDataPath
	userDictFile, userDictErr := os.Open(userDictDataPath)
	if userDictErr != nil {
		return nil, userDictErr
	}
	defer userDictFile.Close()
	return ioutil.ReadAll(userDictFile)
}

//更新用户词库
func (d *Dict) EditUserDict(userDictContent []byte) error {
	//载入用户词库
	userDictDataPath := d.UserDictDataPath
	userDictFile, userDictErr := os.OpenFile(userDictDataPath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
	if userDictErr != nil {
		log.WithError(userDictErr).WithField("dict", userDictFile).Error("failed open user dict")
		return userDictErr
	}
	defer userDictFile.Close()
	_, err := userDictFile.Write(userDictContent)
	if err != nil {
		log.WithError(err).Error("failed write userDict")
		return err
	}
	//reload dict
	d.ReloadDict()
	return nil
}

//重新加载词库
func (d *Dict) ReloadDict() {
	log.Info("reload dict")
	d.MaxWordLength = 0
	d.Trie = &Trie{
		Root: &Node{
			Nodes: make(map[rune]*Node),
		},
	}
	d.LoadDict()
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}