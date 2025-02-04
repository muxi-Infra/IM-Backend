package global

import "sync"

//此包用来管理一些项目中的全局变量的

// AppLock 保证处理请求和配置热更新不会同时发生
var AppLock sync.RWMutex
