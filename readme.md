# github.com/vineboneto/go-sql-builder

---

## Usage

---

### Initialize your module

```bash
$ go mod init example.com/my-golib-demo
```

### Get the go-sql-builder module

```bash
$ go get github.com/vineboneto/go-sql-builder
```

```go
package main

import (
	"fmt"

	sqlbuilder "github.com/vineboneto/go-sql-builder"
)

func main() {
	type Input struct {
		ID        int
		FirstName string
		LastName  string
		GroupId   []int
	}

	input := Input{ID: 2, GroupId: []int{1, 2, 3}, LastName: "Boneto"}

	sql, args := sqlbuilder.BuildPG().
		Raw("SELECT * FROM tb").
		Where().
		And("id = ?", input.ID).
		AndInInt("group_id IN ? ", input.GroupId).
		And("first_name = ?", input.FirstName).
		AndLike("last_name LIKE ?", input.LastName).
		String()

	fmt.Println(sql)
	fmt.Println(args)
}

```
