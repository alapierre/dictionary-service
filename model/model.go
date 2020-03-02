//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

//noinspection ALL
type Dictionary struct {
	tableName struct{} `pg:"dictionary,alias:p,discard_unknown_columns"`

	Key       string                 `pg:"key,pk"`
	Type      string                 `pg:"type,pk"`
	Name      string                 `pg:"name"`
	GroupId   *string                `pg:"group_id"`
	Tenant    string                 `pg:"tenant,pk,use_zero"`
	Content   map[string]interface{} `pg:"content"`
	ParentKey *string                `pg:"parent_key"`
}

//noinspection ALL
type Translation struct {
	tableName struct{} `pg:"translation,alias:t"`

	Key      string `pg:"key,pk"`
	Type     string `pg:"type,pk"`
	Tenant   string `pg:"tenant,pk,use_zero"`
	Language string `pg:"language,pk"`
	Name     string `pg:"name"`
}

// helper indirect model

type ParentDictionary struct {
	Key      string                 `json:"key"`
	Type     string                 `json:"type"`
	Name     string                 `json:"name"`
	GroupId  *string                `json:"group_id"`
	Tenant   string                 `json:"tenant"`
	Content  map[string]interface{} `json:"content"`
	Children []ChildDictionary      `json:"children"`
}

type ChildDictionary struct {
	Key       string                 `json:"key"`
	Name      string                 `json:"name"`
	ParentKey string                 `json:"parent_key"`
	Content   map[string]interface{} `json:"content"`
}
