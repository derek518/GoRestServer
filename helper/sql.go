package helper

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"strings"
)

const (
	SymbolGT      = "gt"
	SymbolLT      = "lt"
	SymbolGTE     = "gte"
	SymbolLTE     = "lte"
	SymbolNot     = "not"
	SymbolBetween = "between"
	SymbolIn      = "in"
	SymbolLike    = "like"
)

var SymbolMap = map[string]string{
	"gt":      ">",
	"lt":      "<",
	"gte":     ">=",
	"lte":     "<=",
	"not":     "<>",
	"between": "BETWEEN",
	"like":    "LIKE",
	"in":      "IN",
}

// json格式查询条件转换成SQL Where格式查询条件
func BuildSqlWhere(query map[string]interface{}) map[string]interface{} {
	where := make(map[string]interface{})

	if query != nil && len(query) > 0 {
		for col, value := range query {
			condMap, ok := value.(map[string]interface{})
			if ok {
				for symbol, v := range condMap {
					sqlSym, ok := SymbolMap[symbol]
					if ok {
						if sqlSym == "IN" {
							where[fmt.Sprintf("%s %s (?)", col, sqlSym)] = v
						} else if sqlSym == "BETWEEN" {
							where[fmt.Sprintf("%s BETWEEN ? AND ?", col)] = v
						} else {
							where[fmt.Sprintf("%s %s ?", col, sqlSym)] = v
						}
					}
				}
			} else {
				where[fmt.Sprintf("%s = ?", col)] = value
			}
		}
	}

	log.Debug().Caller().Fields(where).Msg("sql where")

	return where
}

func BuildWhereIntoDB(db *gorm.DB, andCons, orCons map[string]interface{}) *gorm.DB {
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			if strings.Contains(k, "BETWEEN") {
				if arr, ok := v.([]interface{}); ok {
					db = db.Where(k, arr...)
				}
			} else {
				db = db.Where(k, v)
			}
		}
	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			if strings.Contains(k, "BETWEEN") {
				if arr, ok := v.([]interface{}); ok {
					db = db.Or(k, arr...)
				}
			} else {
				db = db.Or(k, v)
			}
		}
	}

	return db
}

func AddTablePrefixForWhere(cons map[string]interface{}, prefix, tblName string) {
	if strings.TrimSpace(prefix) == "" {
		return
	}

	if cons != nil && len(cons) > 0 {
		for k, v := range cons {
			if strings.HasPrefix(k, tblName) {
				delete(cons, k)
				cons[prefix+k] = v
				break
			}
		}
	}
}
