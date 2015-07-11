package conf

type Privilege struct {
	Key  string `json:"key"`
	Desc string `json:"desc"`
}

var (
	Privileges = []Privilege{
		{"lesscms.editor", "Editor"},
	}
)
