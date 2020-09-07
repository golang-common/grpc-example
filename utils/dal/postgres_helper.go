package dal

// 对于一些公用的sql查询操作，可以使用该结构体去渲染查询到的数据
type PostgresCommonResult struct {
	Max     int           `json:"max"`
	Count   int           `json:"count"` // 若查询符合某条件的数据量，可用该结构体进行返回
	Data    []interface{} `json:"data"`
	Row     string        `json:"row,omitempty"` // 若查询SELECT DISTINCT (xx, xx) 可将查询到的结果记录在该字段
	Error   error         `json:"error"`
	Id      string        `json:"id"`
	NextVal int           `json:"nextval"`
}
