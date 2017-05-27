package word_filter

import (
	"path"
	"os/exec"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	//默认脏词库
	DefaultDictDataPath string `yaml:"DefaultDictDataPath"`
	//自定义脏词库
	UserDictDataPath string `yaml:"UserDictDataPath"`

	//TCP监听地址
	TCPAddr string `yaml:"TCPAddr"`
}

func NewOptions(configPath string) *Options{
	var err error
	workDir, err := workDir()
	if err != nil {
		panic(err)
	}
	if len(configPath) == 0 {
		configPath = path.Join(workDir, "../../conf/app.yaml")
	}
	bConf, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	options := &Options{}
	err = yaml.Unmarshal(bConf, options)
	if err != nil {
		panic(err)
	}
	return options
}

// return app work dir.
func workDir() (string, error) {
	appPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	appPath, err = filepath.Abs(appPath)
	if err != nil {
		return "", err
	}
	// Note: we don't use path.Dir here because it does not handle case
	//	which path starts with two "/" in Windows: "//psf/Home/..."
	appPath = strings.Replace(appPath, "\\", "/", -1)

	i := strings.LastIndex(appPath, "/")
	if i == -1 {
		return appPath, nil
	}
	return appPath[:i], nil
}

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}