package structs

import "fmt"

var ErrNoMatch = fmt.Errorf("no matching record")

var ErrPostgres = fmt.Errorf("postgres error")

var ErrSql = fmt.Errorf("databse/sql error")

var ErrDublicate = fmt.Errorf("record dublicate")
