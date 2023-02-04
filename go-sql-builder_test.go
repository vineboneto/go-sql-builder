package lib

import (
	"log"
	"reflect"
	"testing"
)

func TestSelectWhere(t *testing.T) {

	type Input struct {
		ID        int
		FirstName string
		LastName  string
		GroupId   []int
		Active    bool
	}

	input := Input{ID: 2, GroupId: []int{1, 2, 3}, LastName: "Boneto", Active: true}

	expected_sql := "SELECT *.tb, (SELECT json_agg(g.id) FROM tbl_group g WHERE g.id_user = tb.id AND g.active = ?) AS groups FROM tbl tb WHERE 1 = 1 AND tb.id = ? AND tb.group_id IN ? AND tb.active = 1 AND tb.last_name LIKE ? ORDER BY tb.id OFFSET ? LIMIT ?"
	expected_args := []any{1, 2, []int{1, 2, 3}, "%Boneto%", 10, 20}

	subRaw, args1 := BuildPG().
		Raw("SELECT json_agg(g.id) FROM tbl_group g").
		AndRaw("WHERE g.id_user = tb.id").
		And("g.active = ?", 1).
		String()

	sql, args2 := BuildPG().Raw(`
		SELECT *.tb, (%s) AS groups FROM tbl tb
	`).
		SubRaw(subRaw).
		Where().
		And("tb.id = ?", input.ID).
		AndInInt("tb.group_id IN ? ", input.GroupId).
		AndRawCondition("tb.active = 1", input.Active == true).
		And("tb.first_name = ?", input.FirstName).
		AndLike("tb.last_name LIKE ?", input.LastName).
		Offset(10).
		Limit(20).
		OrderBy("tb.id").
		String()

	var args []any

	args = append(args, args1...)
	args = append(args, args2...)

	log.Println(len(sql))
	log.Println(len(expected_sql))

	if expected_sql != sql || !reflect.DeepEqual(expected_args, args) {
		t.Errorf("Invalid Where, expected %s, receive %s", expected_sql, sql)
		t.Error("expected args:", expected_args)
		t.Error("receive args:", args)
	}
}

func TestInsertWhere(t *testing.T) {

	type Input struct {
		ID        int
		FirstName string
		LastName  string
		Phone     string
	}

	input := Input{FirstName: "Vinicius", Phone: ""}

	expected_sql := "INSERT INTO tbl (first_name, phone) VALUES (?, ?)"
	expected_args := []any{"Vinicius", ""}

	sql, args := BuildPG().Raw(`
		INSERT INTO tbl
	`).
		InsertValue("first_name", input.FirstName).
		InsertValue("last_name", input.LastName).
		InsertValueOnEmpty("phone", input.Phone).
		String()

	if expected_sql != sql || !reflect.DeepEqual(expected_args, args) {
		t.Errorf("Invalid Insert, expected %s, receive %s", expected_sql, sql)
		t.Error("expected args:", expected_args)
		t.Error("receive args:", args)
	}
}

func TestUpdateWhere(t *testing.T) {

	type Input struct {
		ID        int
		FirstName string
		LastName  string
		Phone     string
	}

	input := Input{ID: 2, FirstName: "Vinicius", Phone: ""}

	expected_sql := "UPDATE tbl SET first_name = ?, phone = ? WHERE 1 = 1 AND id = ?"
	expected_args := []any{"Vinicius", "", 2}

	sql, args := BuildPG().Raw(`
		UPDATE tbl
	`).
		UpdateSet("first_name = ?", input.FirstName).
		UpdateSet("last_name = ?", input.LastName).
		UpdateSetOnEmpty("phone = ?", input.Phone).
		Where().
		And("id = ?", input.ID).
		String()

	if expected_sql != sql || !reflect.DeepEqual(expected_args, args) {
		t.Errorf("Invalid Update, expected %s, receive %s", expected_sql, sql)
		t.Error("expected args:", expected_args)
		t.Error("receive args:", args)
	}
}
