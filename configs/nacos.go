package configs

import (
	"IM-Backend/global"
	"IM-Backend/pkg"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosConfig struct {
	IpAddr      string `yaml:"ipaddr"`
	Port        uint64 `yaml:"port"`
	Namespaceid string `yaml:"namespaceid"`
	Group       string `yaml:"group"`
	DataId      string `yaml:"dataid"`
}

func NewNacosConfig(path string) NacosConfig {
	//获取NacosConfig
	var nc NacosConfig
	err := pkg.ReadYaml(path, &nc)
	if err != nil {
		panic(err)
	}
	return nc
}

type NacosClient struct {
	nc     NacosConfig
	client config_client.IConfigClient
}

func NewNacosClient(nc NacosConfig) *NacosClient {
	sc := []constant.ServerConfig{
		{
			IpAddr: nc.IpAddr,
			Port:   nc.Port,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         nc.Namespaceid, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/log",
		CacheDir:            "tmp/cache",
		LogLevel:            "info",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(fmt.Sprintf("connect to nocos failed:%v", err))
	}
	global.Log.Info("connect to nocos successfully")
	return &NacosClient{
		nc:     nc,
		client: configClient,
	}
}
func (ncli *NacosClient) GetConfig() ([]byte, error) {
	content, err := ncli.client.GetConfig(vo.ConfigParam{
		DataId: ncli.nc.DataId,
		Group:  ncli.nc.Group,
	})
	if err != nil {
		return nil, err
	}
	return []byte(content), err
}
func (ncli *NacosClient) ListenConfig(f func(data string)) error {
	return ncli.client.ListenConfig(vo.ConfigParam{
		DataId: ncli.nc.DataId,
		Group:  ncli.nc.Group,
		OnChange: func(namespace, group, dataId, data string) {
			//fmt.Println("配置文件发生了变化...")
			//fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
			f(data)
		},
	})

}
