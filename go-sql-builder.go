/**
 * Author: Vinicius Gazolla Boneto
 * File: go-sql-builder.go
 */

package lib

import (
	"fmt"
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

func (q *Query) InsertValue(field string, v any) *Query {
	if v != "" && v != 0 && v != nil {
		q.insertFieldSql = append(q.insertFieldSql, field)
		q.args = append(q.args, v)
	}

	return q
}

func (q *Query) InsertValueOnEmpty(field string, v any) *Query {
	q.insertFieldSql = append(q.insertFieldSql, field)
	q.args = append(q.args, v)

	return q
}

func (q *Query) UpdateSet(field string, v any) *Query {
	if v != "" && v != 0 && v != nil {
		q.updateFieldsSql = append(q.updateFieldsSql, field)
		q.args = append(q.args, v)
	}

	return q
}

func (q *Query) UpdateSetOnEmpty(field string, v any) *Query {
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

	if v != "" && v != 0 && v != nil && v != false {
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

	if v1 != "" && v1 != 0 && v1 != nil && v2 != "" && v2 != 0 && v2 != nil {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, v1, v2)
	}
	return q
}

func (q *Query) AndInInt(s string, v []int) *Query {

	if len(v) != 0 {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, v)
	}
	return q
}

func (q *Query) AndInStr(s string, v []string) *Query {

	if len(v) != 0 {
		q.whereFieldsSql = append(q.whereFieldsSql, s)
		q.args = append(q.args, v)
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

func (q *Query) OrderBy(f string) *Query {
	q.orderBy = "ORDER BY " + f
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
