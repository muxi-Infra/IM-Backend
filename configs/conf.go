package configs

import (
	"IM-Backend/pkg"
	"fmt"
)

// Notifyer 通知更新配置
// 需要配置的需要实现该方法
type Notifyer interface {
	Callback(*AppConf)
}

type AppConf struct {
	Svc        []ServiceConfig `yaml:"svc"`
	DB         DBConfig        `yaml:"db"`
	Cache      CacheConfig     `yaml:"cache"`
	notifyList []Notifyer
}

// 加载配置
func (ac *AppConf) load(content []byte) error {
	err := pkg.ReadYamlContent(content, ac)
	if err != nil {
		return err
	}
	for _, nf := range ac.notifyList {
		nf.Callback(ac)
	}
	return nil
}

// AddNotifyer 添加通知者
func (ac *AppConf) AddNotifyer(nf ...Notifyer) {
	ac.notifyList = append(ac.notifyList, nf...)
}

func (ac *AppConf) InitConfig(ncc *NacosClient) {
	content, err := ncc.GetConfig()
	if err != nil {
		panic(err)
	}
	//加载配置
	err = ac.load(content)
	if err != nil {
		panic(err)
	}
	f := func(data string) {
		err := ac.load([]byte(data))
		if err != nil {
			//TODO:记录日志
			fmt.Println(err)
		}
	}
	//监听配置
	//当配置改变时，通知各个Notifyer更新配置
	if err := ncc.ListenConfig(f); err != nil {
		panic(err)
	}
}

type ServiceConfig struct {
	Name   string `yaml:"name"`   //服务名
	Secret string `yaml:"secret"` //秘钥
}

type DBConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}
type CacheConfig struct {
	Addr          string `yaml:"addr"`
	Password      string `yaml:"password"`
	AppKeyExpire  uint64 `yaml:"appkeyexpire"`  //以秒为单位
	PostExpire    uint64 `yaml:"postexpire"`    //帖子缓存
	CommentExpire uint64 `yaml:"commentexpire"` //评论缓存
}
