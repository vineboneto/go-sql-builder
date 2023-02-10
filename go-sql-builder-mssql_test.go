package lib

import (
	"reflect"
	"testing"
)

func TestSelectWhereMSSQL(t *testing.T) {

	type Input struct {
		ID          int
		FirstName   string
		LastName    string
		GroupId     []int
		Permissions []string
		Values      []float64
		Active      bool
	}

	input := Input{ID: 2, GroupId: []int{1, 2, 3}, Permissions: []string{"USER", "ADMIN"}, LastName: "Boneto", Active: true}

	expected_sql := "SELECT *.tb, (SELECT g.id FROM tbl_group g WHERE g.id_user = tb.id AND g.active = ? FOR JSON PATH) AS groups FROM tbl tb WHERE 1 = 1 AND tb.id = ? AND tb.group_id IN ? AND tb.permission_id IN ? AND tb.active = 1 AND tb.last_name LIKE ? ORDER BY tb.id asc OFFSET ? ROWS FETCH NEXT ? ROWS ONLY"
	expected_args := []any{1, 2, []int{1, 2, 3}, []string{"USER", "ADMIN"}, "%Boneto%", 10, 20}

	subRaw, args1 := BuildPG().
		Raw("SELECT g.id FROM tbl_group g").
		AndRaw("WHERE g.id_user = tb.id").
		And("g.active = ?", 1).
		String()

	sql, args2 := BuildSQLServer().Raw("SELECT *.tb, (%s FOR JSON PATH) AS groups FROM tbl tb").
		SubRaw(subRaw).
		Where().
		And("tb.id = ?", input.ID).
		AndIn("tb.group_id IN ? ", input.GroupId).
		AndIn("tb.permission_id IN ? ", input.Permissions).
		AndIn("tb.values IN ? ", input.Values).
		AndRawCondition("tb.active = 1", input.Active == true).
		And("tb.first_name = ?", input.FirstName).
		AndLike("tb.last_name LIKE ?", input.LastName).
		Offset(10).
		Limit(20).
		OrderBy("tb.id", "asc").
		String()

	var args []any

	args = append(args, args1...)
	args = append(args, args2...)

	if expected_sql != sql || !reflect.DeepEqual(expected_args, args) {
		t.Errorf("Invalid Where, expected %s, receive %s", expected_sql, sql)
		t.Error("expected args:", expected_args)
		t.Error("receive args:", args)
	}
}

func TestInsertWhereMSSQL(t *testing.T) {

	type Input struct {
		ID        int
		FirstName string
		LastName  string
		Phone     string
	}

	input := Input{FirstName: "Vinicius", Phone: ""}

	expected_sql := "INSERT INTO tbl (first_name, phone) VALUES (?, ?)"
	expected_args := []any{"Vinicius", ""}

	sql, args := BuildSQLServer().Raw(`
		INSERT INTO tbl
	`).
		InsertOnlyValue("first_name", input.FirstName).
		InsertOnlyValue("last_name", input.LastName).
		Insert("phone", input.Phone).
		String()

	if expected_sql != sql || !reflect.DeepEqual(expected_args, args) {
		t.Errorf("Invalid Insert, expected %s, receive %s", expected_sql, sql)
		t.Error("expected args:", expected_args)
		t.Error("receive args:", args)
	}
}

func TestUpdateWhereMSSQL(t *testing.T) {

	type Input struct {
		ID        int
		FirstName string
		LastName  string
		Phone     string
	}

	input := Input{ID: 2, FirstName: "Vinicius", Phone: ""}

	expected_sql := "UPDATE tbl SET first_name = ?, phone = ? WHERE 1 = 1 AND id = ?"
	expected_args := []any{"Vinicius", "", 2}

	sql, args := BuildSQLServer().Raw("UPDATE tbl").
		UpdateOnlyValue("first_name = ?", input.FirstName).
		UpdateOnlyValue("last_name = ?", input.LastName).
		Update("phone = ?", input.Phone).
		Where().
		And("id = ?", input.ID).
		String()

	if expected_sql != sql || !reflect.DeepEqual(expected_args, args) {
		t.Errorf("Invalid Update, expected %s, receive %s", expected_sql, sql)
		t.Error("expected args:", expected_args)
		t.Error("receive args:", args)
	}
}
