package errcode

var (
	ERRNoTable          = NewErr(10001, "the table of the service has not been created")
	ERRCreateTable      = NewErr(10002, "can not create table of the service")
	ERRCreateData       = NewErr(10003, "create data failed")
	ERRUpdateData       = NewErr(10004, "update data failed")
	ERRDeleteData       = NewErr(10005, "delete data failed")
	ERRFindData         = NewErr(10006, "find data failed")
	ERRCount            = NewErr(10007, "count data failed")
	ERRFindQueryIsEmpty = NewErr(10008, "find query is empty")
	ERRSetKV            = NewErr(10009, "set k-v in cache failed")
	ERRAddSet           = NewErr(10010, "sadd val in cache failed")
	ERRGetKV            = NewErr(10011, "get k-v from cache failed")
	ERRConvertJson      = NewErr(10012, "convert json failed")
	ERRGetSet           = NewErr(10013, "get set from cache failed")
	ERRNoRightRecord    = NewErr(10014, "no right  record found")
	ERRCacheMiss        = NewErr(10015, "cache miss")
	ERRUpdateQueryEmpty = NewErr(10016, "update query is empty")
	ERRGenerateID       = NewErr(10017, "generate id failed")
	ERRDelKV            = NewErr(10018, "del k-v in cache failed")
)
var ClientMsgMapping = map[int]string{
	10001: "服务初始化未完成，请稍后重试",
	10002: "服务初始化失败，请联系管理员",
	10003: "数据创建失败，请检查输入后重试",
	10004: "数据更新失败，请检查输入或联系支持",
	10005: "数据删除失败，请确认状态后重试",
	10006: "数据查询异常，请稍后重试",
	10007: "数据统计失败，请联系技术支持",
	10008: "查询条件不能为空",
	10009: "系统缓存繁忙，请稍后重试",
	10010: "系统缓存繁忙，请稍后重试",
	10011: "系统缓存繁忙，请稍后重试",
	10012: "数据格式解析错误，请检查输入",
	10013: "系统缓存繁忙，请稍后重试",
	10014: "未找到有效记录，请确认查询条件",
	10015: "正在重新加载数据，请稍后重试",
	10016: "更新条件不能为空",
	10017: "系统资源分配失败，请重试",
	10018: "系统缓存繁忙，请稍后重试",
}
