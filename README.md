# GoRestServer-简单的后端基本框架代码
基于go web框架[gin]()和数据库orm框架[gorm]()实现一个简单的支持权限管理功能的基本RestAPI服务器代码，
根据项目需要可直接添加相应的业务数据表代码即可。

在model下添加数据表和json原型结构，在service下添加相应的服务接口，在api下添加相应的gin路由接口即可。

## 实现功能
- 支持YAML/JSON等多种格式配置文件
- 采用Controller(API)-Service-Model架构
- 支持灵活的不定条件复杂检索并翻页查询功能,比如gt,lt,gte,lte,not,between,like,in等复杂条件
- 支持swagger API文档自动生成并在线浏览功能
- 支持JWT中间件进行身份认证
- 支持集成redis对数据表进行可配置缓存功能

## 不定条件复杂检索功能
```go
// 统一查询格式
type QueryCondition struct {
	/** 当前页  */
	Page   int	`json:"page"`
	/** 每页显示的最大行数 */
	Size  int	`json:"size"`
	/** 排序方式 */
	Sort string	`json:"sort,omitempty"`
	/** 指定返回的字段 */
	Selection string	`json:"selection,omitempty"`
	/** 指定And查询条件 */
	AndCons interface{}	`json:"and_cons,omitempty"`
	/** 指定Or查询条件 */
	OrCons interface{}	`json:"or_cons,omitempty"`
}
```
 - 灵活 Page和Size组合支持翻页功能
 - Sort支持排序功能，语法同SQL中的Order By，支持跨表排序
 ```
 "name desc, role.seq asc"
```
 - Selection支持按需返回数据字段
 - AndCons和OrCons查询条件，支持复杂条件查询，支持跨表关联查询
 ```json
{
  "name": "hello",
  "desc": {"like": "%hello%"},
  "role.key": {"between": ["value1", "value2"]},  // 跨表
  "created_at": {"gte": "time1"},
  "key1": {"in": ["value1", "value2", "value3"]}
}
```

## 编译运行
下载代码到$GOPATH/src下面，修改conf配置文件中数据库和Redis配置
```
go get
go build
```
运行编译后的二进制程序即可

API文档路径为：
http://localhost:8080/swagger/index.html
 