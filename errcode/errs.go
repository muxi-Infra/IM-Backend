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
)
