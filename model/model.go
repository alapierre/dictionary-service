//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

//noinspection ALL
type Dictionary struct {
	tableName struct{} `pg:"dictionary,alias:p,discard_unknown_columns"`

	Key       string                 `pg:"key,pk"`
	Type      string                 `pg:"type,pk"`
	Name      string                 `pg:"name"`
	GroupID   *string                `pg:"group_id"`
	Tenant    string                 `pg:"tenant,pk,use_zero"`
	Content   map[string]interface{} `pg:"content"`
	ParentKey *string                `pg:"parent_key"`
}

// helper indirect model

type ParentDictionary struct {
	Key      string
	Type     string
	Name     string
	GroupId  *string
	Tenant   string
	Content  map[string]interface{}
	Children []ChildDictionary
}

type ChildDictionary struct {
	Key       string
	Name      string
	ParentKey string
	Content   map[string]interface{}
}
