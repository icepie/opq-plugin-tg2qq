package conf

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/astaxie/beego/logs"

	"gopkg.in/yaml.v2"

	"github.com/spf13/viper"
)

// OPQBotConfig OPQ Bot 配置
type OPQBotConfig struct {
	Url   string
	QQ    int64
	Group int64
}

type ProxyConfig struct {
	Enable bool
	Url    string
}

// TGBotConfig Telegram Bot 配置
type TGBotConfig struct {
	Token  string
	ChatID string
	Proxy  ProxyConfig
}

// Config 基础配置
type Config struct {
	OPQBot OPQBotConfig
	TGBot  TGBotConfig
}

// ProConf 新建实例
var ProConf = new(Config)

// initConfig 初始化配置
func initConfig(cfpath string) error {

	b, err := yaml.Marshal(ProConf)
	if err != nil {
		return err
	}

	f, err := os.Create(cfpath)
	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString(string(b))

	return nil
}

// init 初始化函数
func init() {

	// 取项目地址
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	var pathsep string
	if runtime.GOOS == "windows" {
		pathsep = "\\"
	} else {
		pathsep = "/"
	}

	cfpath := path + pathsep + "conf.yaml"

	viper.AddConfigPath(path)
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		logs.Emergency("Can not read the config file, will recreate it! ")
		// 初始化配置
		ProConf.OPQBot.Url = "http://127.0.0.1:8888"
		ProConf.TGBot.Proxy.Enable = false
		ProConf.TGBot.Proxy.Url = "sock5://127.0.0.1:1080"
		if err = initConfig(cfpath); err != nil {
			logs.Error("%s", err.Error())
		}
		logs.Trace("Please edit the " + cfpath + "，then restart app")
		os.Exit(1)
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(ProConf); err != nil {
		logs.Error(err.Error())
	}

}
