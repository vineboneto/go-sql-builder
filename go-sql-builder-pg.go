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

type query struct {
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

func BuildPG() *query {
	return &query{}

}

func (q *query) InsertOnlyValue(field string, v any) *query {
	if checkHasValue(v) {
		q.insertFieldSql = append(q.insertFieldSql, field)
		q.args = append(q.args, v)
	}

	return q
}

func (q *query) Insert(field string, v any) *query {
	q.insertFieldSql = append(q.insertFieldSql, field)
	q.args = append(q.args, v)

	return q
}

func (q *query) UpdateOnlyValue(field string, v any) *query {
	if checkHasValue(v) {
		q.updateFieldsSql = append(q.updateFieldsSql, field)
		q.args = append(q.args, v)
	}

	return q
}

func (q *query) Update(field string, v any) *query {
	q.updateFieldsSql = append(q.updateFieldsSql, field)
	q.args = append(q.args, v)

	return q
}

func (q *query) InsertEnd(sql string) *query {
	q.insertEndSql = sql
	return q
}

func (q *query) Raw(raw string) *query {
	q.raw = raw
	return q
}

func (q *query) SubRaw(raw string) *query {
	q.subRaws = append(q.subRaws, raw)
	return q
}

func (q *query) Where() *query {
	q.whereFieldsSql = append(q.whereFieldsSql, "WHERE 1 = 1")
	return q
}

func (q *query) AndRaw(raw string) *query {
	q.whereFieldsSql = append(q.whereFieldsSql, raw)
	return q
}

func (q *query) And(s string, v any) *query {

	if checkHasValue(v) {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, v)
	}
	return q
}

func (q *query) AndRawCondition(s string, v bool) *query {

	if v {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
	}
	return q
}

func (q *query) AndLike(s string, v string) *query {

	if v != "" {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, "%"+v+"%")
	}
	return q
}

func (q *query) AndBetween(s string, v1 any, v2 any) *query {

	if checkHasValue(v1) && checkHasValue(v2) {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, v1, v2)
	}
	return q
}

func (q *query) AndIn(s string, v any) *query {

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

func (q *query) Offset(v int) *query {
	if v > 0 {
		q.offset = v
	}
	return q
}

func (q *query) Limit(v int) *query {
	if v > 0 {
		q.limit = v
	}
	return q
}

func (q *query) OrderBy(v string, order string) *query {
	if v != "" {
		q.orderBy = "ORDER BY " + v + " " + order
	}
	return q
}

func (q *query) String() (string, []any) {
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
