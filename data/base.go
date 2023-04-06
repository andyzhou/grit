package data

/*
 * base data face
 */

//data face
type BaseData struct {
	Kind int
}

type CountData struct {
	Fields map[string]int64 `json:"fields"`
}