/**
 * Author: Vinicius Gazolla Boneto
 * File: go-sql-builder.go
 */

package lib

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Query struct {
	raw             string
	subRaws         []any
	insertFieldSql  []string
	updateFieldsSql []string
	insertEndSql    string
	whereFieldsSql  []string

	args    []any
	limit   int
	offset  int
	orderBy string
}

func BuildPG() *Query {
	return &Query{}

}

func (q *Query) InsertOnlyValue(field string, v any) *Query {
	if q.checkHashValue(v) {
		q.insertFieldSql = append(q.insertFieldSql, field)
		q.args = append(q.args, v)
	}

	return q
}

func (q *Query) Insert(field string, v any) *Query {
	q.insertFieldSql = append(q.insertFieldSql, field)
	q.args = append(q.args, v)

	return q
}

func (q *Query) UpdateOnlyValue(field string, v any) *Query {
	if q.checkHashValue(v) {
		q.updateFieldsSql = append(q.updateFieldsSql, field)
		q.args = append(q.args, v)
	}

	return q
}

func (q *Query) Update(field string, v any) *Query {
	q.updateFieldsSql = append(q.updateFieldsSql, field)
	q.args = append(q.args, v)

	return q
}

func (q *Query) InsertEnd(sql string) *Query {
	q.insertEndSql = sql
	return q
}

func (q *Query) Raw(raw string) *Query {
	q.raw = raw
	return q
}

func (q *Query) SubRaw(raw string) *Query {
	q.subRaws = append(q.subRaws, raw)
	return q
}

func (q *Query) Where() *Query {
	q.whereFieldsSql = append(q.whereFieldsSql, "WHERE 1 = 1")
	return q
}

func (q *Query) AndRaw(raw string) *Query {
	q.whereFieldsSql = append(q.whereFieldsSql, raw)
	return q
}

func (q *Query) And(s string, v any) *Query {

	if q.checkHashValue(v) {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, v)
	}
	return q
}

func (q *Query) AndRawCondition(s string, v bool) *Query {

	if v {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
	}
	return q
}

func (q *Query) AndLike(s string, v string) *Query {

	if v != "" {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, "%"+v+"%")
	}
	return q
}

func (q *Query) AndBetween(s string, v1 any, v2 any) *Query {

	if q.checkHashValue(v1) && q.checkHashValue(v2) {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, v1, v2)
	}
	return q
}

func (q *Query) AndIn(s string, v any) *Query {

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

func (q *Query) Offset(v int) *Query {
	if v > 0 {
		q.offset = v
	}
	return q
}

func (q *Query) Limit(v int) *Query {
	if v > 0 {
		q.limit = v
	}
	return q
}

func (q *Query) OrderBy(v string, order string) *Query {
	if v != "" {
		q.orderBy = "ORDER BY " + v + " " + order
	}
	return q
}

func (q *Query) String() (string, []any) {
	limitSql := ""
	offsetSql := ""
	insertFieldsSql := ""
	insertValuesSql := ""
	updateSql := ""
	initialSql := q.raw

	if q.offset > 0 {
		q.args = append(q.args, q.offset)
		offsetSql = " OFFSET ? "
	}

	if q.limit > 0 {
		q.args = append(q.args, q.limit)
		limitSql = " LIMIT ? "
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

func (q *Query) checkHashValue(v any) bool {
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
