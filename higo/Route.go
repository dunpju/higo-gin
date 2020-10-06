package higo

type Route struct {
	Method       string      // 请求方法 GET/POST/DELETE/PATCH/OPTIONS/HEAD
	RelativePath string      // 后端 api relativePath
	Handle       interface{} // 后端控制器函数
	Flag         string      // 后端控制器函数标记
	FrontPath    string      // 前端 path(前端菜单路由)
	Desc         string      // 描述
}

