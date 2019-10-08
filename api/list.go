package api

type List struct {
	List []Snapshots `json:"snapshots"`
}

type Snapshots struct {
	Name string `json:"name"`
	Times []int64 `json:"times"`
}
