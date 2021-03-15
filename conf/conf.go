package conf

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/spf13/viper"
)

// OPQBotConfig OPQ Bot 配置
type OPQBotConfig struct {
	Url   string
	QQ    int
	Group int
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

// Init 初始化函数
func Init() {

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

	viper.AddConfigPath(path)   // 设置读取的文件路径
	viper.SetConfigName("conf") // 设置读取的文件名
	viper.SetConfigType("yaml") // 设置文件的类型

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {            // 读取配置信息失败
		log.Warningln("Can not read the config file, will recreate it! ")
		// 初始化配置
		ProConf.OPQBot.Url = "http://127.0.0.1:8888"
		ProConf.TGBot.Proxy.Enable = false
		ProConf.TGBot.Proxy.Url = "sock5://127.0.0.1:1080"
		if err = initConfig(cfpath); err != nil { // 重新初始化配置文件
			log.Fatalln(err)
		}
		log.Println(errors.New("Please edit the " + cfpath + "，then restart app"))
		os.Exit(1)
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(ProConf); err != nil {
		log.Fatalln(err)
	}

}
