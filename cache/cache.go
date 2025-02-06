package cache

import (
	"IM-Backend/configs"
	"log"
)

type KV interface {
	ReadFromStrVal(jsonStr string) error

	GetStrKey() string
	GetStrVal() string

	GetSetKey() string
	GetSetVal() string
}

// 使用map来存储service的信息
type SvcManager struct {
	mp map[string]string // key:name, value:secret
}

func NewSvcManager(conf configs.AppConf) *SvcManager {
	//分配一个map
	mp := make(map[string]string, len(conf.Svc))
	for _, svc := range conf.Svc {
		mp[svc.Name] = svc.Secret
	}
	log.Printf("init svc: %+v\n", mp)
	return &SvcManager{mp: mp}
}

func (m *SvcManager) Callback(conf configs.AppConf) {
	//重新分配一个map
	m.mp = make(map[string]string, len(conf.Svc))
	for _, svc := range conf.Svc {
		m.mp[svc.Name] = svc.Secret
	}
	log.Printf("init svc: %+v\n", m.mp)
}

func (m *SvcManager) GetAllServices() []string {
	services := make([]string, 0, len(m.mp))
	for name := range m.mp {
		services = append(services, name)
	}
	return services
}

func (m *SvcManager) GetSecretByName(name string) string {
	return m.mp[name]
}
