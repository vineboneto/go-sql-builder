package lib

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type offset struct {
	value int
	use   bool
}

type querySqlServer struct {
	raw             string
	subRaws         []any
	insertFieldSql  []string
	updateFieldsSql []string
	insertEndSql    string
	whereFieldsSql  []string

	args    []any
	limit   int
	offset  offset
	orderBy string
}

func BuildSQLServer() *querySqlServer {
	return &querySqlServer{}

}

func (q *querySqlServer) InsertOnlyValue(field string, v any) *querySqlServer {
	if q.checkHashValue(v) {
		q.insertFieldSql = append(q.insertFieldSql, field)
		q.args = append(q.args, v)
	}

	return q
}

func (q *querySqlServer) Insert(field string, v any) *querySqlServer {
	q.insertFieldSql = append(q.insertFieldSql, field)
	q.args = append(q.args, v)

	return q
}

func (q *querySqlServer) UpdateOnlyValue(field string, v any) *querySqlServer {
	if q.checkHashValue(v) {
		q.updateFieldsSql = append(q.updateFieldsSql, field)
		q.args = append(q.args, v)
	}

	return q
}

func (q *querySqlServer) Update(field string, v any) *querySqlServer {
	q.updateFieldsSql = append(q.updateFieldsSql, field)
	q.args = append(q.args, v)

	return q
}

func (q *querySqlServer) InsertEnd(sql string) *querySqlServer {
	q.insertEndSql = sql
	return q
}

func (q *querySqlServer) Raw(raw string) *querySqlServer {
	q.raw = raw
	return q
}

func (q *querySqlServer) SubRaw(raw string) *querySqlServer {
	q.subRaws = append(q.subRaws, raw)
	return q
}

func (q *querySqlServer) Where() *querySqlServer {
	q.whereFieldsSql = append(q.whereFieldsSql, "WHERE 1 = 1")
	return q
}

func (q *querySqlServer) AndRaw(raw string) *querySqlServer {
	q.whereFieldsSql = append(q.whereFieldsSql, raw)
	return q
}

func (q *querySqlServer) And(s string, v any) *querySqlServer {

	if q.checkHashValue(v) {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, v)
	}
	return q
}

func (q *querySqlServer) AndRawCondition(s string, v bool) *querySqlServer {

	if v {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
	}
	return q
}

func (q *querySqlServer) AndLike(s string, v string) *querySqlServer {

	if v != "" {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, "%"+v+"%")
	}
	return q
}

func (q *querySqlServer) AndBetween(s string, v1 any, v2 any) *querySqlServer {

	if q.checkHashValue(v1) && q.checkHashValue(v2) {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, v1, v2)
	}
	return q
}

func (q *querySqlServer) AndIn(s string, v any) *querySqlServer {

	isArray := func(v any) bool {
		return reflect.ValueOf(v).Kind() == reflect.Slice
	}

	if isArray(v) {

		array := reflect.ValueOf(v)

		if array.Len() != 0 {
			q.whereFieldsSql = append(q.whereFieldsSql, s)
			q.args = append(q.args, v)
		}

	}

	return q
}

func (q *querySqlServer) Offset(v int) *querySqlServer {
	if v >= 0 {
		q.offset = offset{value: v, use: true}
	}
	return q
}

func (q *querySqlServer) Limit(v int) *querySqlServer {
	if v > 0 {
		q.limit = v
	}
	return q
}

func (q *querySqlServer) OrderBy(v string, order string) *querySqlServer {
	if v != "" {
		q.orderBy = "ORDER BY " + v + " " + order
	}
	return q
}

func (q *querySqlServer) String() (string, []any) {
	limitSql := ""
	offsetSql := ""
	insertFieldsSql := ""
	insertValuesSql := ""
	updateSql := ""
	initialSql := q.raw

	if q.offset.use {
		q.args = append(q.args, q.offset.value)
		offsetSql = " OFFSET ? ROWS"
	}

	if q.limit > 0 {
		q.args = append(q.args, q.limit)
		limitSql = "FETCH NEXT ? ROWS ONLY"
	}

	if len(q.subRaws) > 0 {
		initialSql = fmt.Sprintf(initialSql, q.subRaws...)
	}

	if len(q.insertFieldSql) > 0 {
		insertParams := []string{}
		for i := 0; i < len(q.args); i++ {
			insertParams = append(insertParams, "?")
		}
		insertFieldsSql = fmt.Sprintf("(%s)", strings.Join(q.insertFieldSql, ", "))
		insertValuesSql = fmt.Sprintf("VALUES (%s)", strings.Join(insertParams, ", "))
	}

	if len(q.updateFieldsSql) > 0 {
		updateSql = fmt.Sprintf("SET %s", strings.Join(q.updateFieldsSql, ", "))
	}

	str := []string{
		initialSql,
		updateSql,
		insertFieldsSql,
		insertValuesSql,
		q.insertEndSql,
		strings.Join(q.whereFieldsSql, " AND "),
		q.orderBy,
		offsetSql,
		limitSql,
	}

	pattern := regexp.MustCompile(`\s+`)

	trimJoin := strings.TrimSpace(strings.Join(str, " "))

	return pattern.ReplaceAllString(trimJoin, " "), q.args
}

func (q *querySqlServer) checkHashValue(v any) bool {
	switch v := v.(type) {
	case int:
		if v != 0 {
			return true
		}
	case string:
		if v != "" {
			return true
		}
	case float64:
		if v != float64(0) {
			return true
		}
	case bool:
		return v
	case nil:
		return false
	default:
		return false
	}
	return false
}
