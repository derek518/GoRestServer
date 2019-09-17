package base

// 统一查询格式
type QueryCondition struct {
	/** 当前页  */
	Page int `json:"page"`
	/** 每页显示的最大行数 */
	Size int `json:"size"`
	/** 排序方式 */
	Sort string `json:"sort,omitempty"`
	/** 指定返回的字段 */
	Selection string `json:"selection,omitempty"`
	/** 指定And查询条件 */
	AndCons interface{} `json:"and_cons,omitempty"`
	/** 指定Or查询条件 */
	OrCons interface{} `json:"or_cons,omitempty"`
}

// 统一 json 结构体
type JsonObject struct {
	/** 状态码: 0: 成功; 1: 异常 */
	Code int `json:"code"`
	/** 内容体 */
	Content interface{} `json:"content,omitempty"`
	/** 消息 */
	Message string `json:"message,omitempty"`
}

// 全局分页对象
type PageBean struct {
	/** 当前页  */
	Page int `json:"page"`
	/** 每页显示的最大行数 */
	Size int `json:"size"`
	/** 总记录数 */
	Total int `json:"total"`
	/** 每行的数据 */
	Rows interface{} `json:"rows"`
}
