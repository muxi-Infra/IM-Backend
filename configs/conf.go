package configs

import (
	"IM-Backend/global"
	"IM-Backend/pkg"
	"fmt"
)

// Notifyer 通知更新配置
// 需要配置的需要实现该方法
type Notifyer interface {
	Callback(AppConf)
}

type AppConf struct {
	Svc        []ServiceConfig `yaml:"svc"`
	DB         DBConfig        `yaml:"db"`
	Cache      CacheConfig     `yaml:"cache"`
	Clean      CleanConfig     `yaml:"clean"`
	notifyList []Notifyer
}

// 加载配置
func (ac *AppConf) load(content []byte) error {
	err := pkg.ReadYamlContent(content, ac)
	if err != nil {
		return err
	}
	return nil
}
func (ac *AppConf) notify() {
	for _, nf := range ac.notifyList {
		nf.Callback(*ac)
	}
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
}

func (ac *AppConf) StartListen(ncc *NacosClient) {
	f := func(data string) {
		global.AppLock.Lock()
		defer global.AppLock.Unlock()
		//先读取配置，再通知各个Notifyer更新配置
		err := ac.load([]byte(data))
		if err != nil {
			//TODO:记录日志
			fmt.Println(err)
			return
		}
		ac.notify()
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
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	PassWord string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Port     int    `yaml:"port"`
}
type CacheConfig struct {
	Addr              string `yaml:"addr"`
	Password          string `yaml:"password"`
	PostExpire        uint64 `yaml:"post_expire"`        //帖子缓存
	CommentExpire     uint64 `yaml:"comment_expire"`     //评论缓存
	PostLikeExpire    uint64 `yaml:"postlike_expire"`    //帖子点赞信息缓存过期时间
	CommentLikeExpire uint64 `yaml:"commentlike_expire"` //帖子点赞信息缓存过期时间
}

// CleanConfig 清理数据库中垃圾数据的配置
type CleanConfig struct {
	CommentBatch     int `yaml:"commentbatch"`     //删除comment的一批的数量
	PostLikeBatch    int `yaml:"postlikebatch"`    // 删除post like的一批的数量
	CommentLikeBatch int `yaml:"commentlikebatch"` // 删除comment like的一批的数量
}
