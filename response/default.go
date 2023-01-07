package response

type ResponseDict struct {
	ReturnRecord bool /* 是否返回新增或者删除的记录 */
	ReturnCount  int /* 是否返回所有条数 */
	ReturnPage   int /* 是否返回当前页数 */
	ReturnSize   int /* 是否返回当前页码数 */
	ReturnId     string /* 是否返回当前操作的主键 */
}

type ResponseDicts interface {
	apply(*ResponseDict)
}

type ResponseDictFun func(*ResponseDict)

type FuncResponseMap struct {
	f ResponseDictFun
}

func (funcs *FuncResponseMap) apply(r *ResponseDict) {
	funcs.f(r)
}

func newDict(f ResponseDictFun) FuncResponseMap {
	return FuncResponseMap{
		f: f,
	}
}

var defaultValue = ResponseDict{
	ReturnRecord: false,
	ReturnCount:  0,
	ReturnPage:   0,
	ReturnSize:   0,
	ReturnId:     "",
}

/* 注入值 */
func WithInstallRecord(val bool) FuncResponseMap {
	return newDict(func(rd *ResponseDict) {
		rd.ReturnRecord = val
	})
}

func WithInstallCount(val int) FuncResponseMap {
	return newDict(func(rd *ResponseDict) {
		rd.ReturnCount = val
	})
}

func WithInstallPage(val int) FuncResponseMap {
	return newDict(func(rd *ResponseDict) {
		rd.ReturnPage = val
	})
}

func WithInstallSize(val int) FuncResponseMap {
	return newDict(func(rd *ResponseDict) {
		rd.ReturnSize = val
	})
}

func WithInstallId(val string) FuncResponseMap {
	return newDict(func(rd *ResponseDict) {
		rd.ReturnId = val
	})
}

func NewResponseDict(opts ...FuncResponseMap) ResponseDict {
	var res ResponseDict = defaultValue
	for _, opt := range opts {
		opt.apply(&res)
	}
	return res
}
