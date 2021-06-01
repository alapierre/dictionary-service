//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

//noinspection ALL
type Dictionary struct {
	tableName struct{} `pg:"dictionary,alias:p,discard_unknown_columns"`

	Key       string                 `pg:"key,pk" json:"key"`
	Type      string                 `pg:"type,pk" json:"type"`
	Name      string                 `pg:"name" json:"name"`
	GroupId   *string                `pg:"group_id" json:"group_id"`
	Tenant    string                 `pg:"tenant,pk,use_zero" json:"-"`
	Content   map[string]interface{} `pg:"content" json:"content"`
	ParentKey *string                `pg:"parent_key" json:"parent_key"`
}

//noinspection ALL
type Translation struct {
	tableName struct{} `pg:"translation,alias:t"`

	Key      string `pg:"key,pk"  json:"key"`
	Type     string `pg:"type,pk" json:"type"`
	Tenant   string `pg:"tenant,pk,use_zero" json:"tenant"`
	Language string `pg:"language,pk" json:"language"`
	Name     string `pg:"name" json:"name"`
}

//noinspection ALL
type DictionaryMetadata struct {
	tableName struct{} `pg:"dictionary_metadata,alias:m"`

	Type    string `pg:"type,pk" json:"type"`
	Tenant  string `pg:"tenant,pk,use_zero" json:"tenant"`
	Content string `pg:"content" json:"content"`
}

// helper indirect model

type ParentDictionary struct {
	Key      string                 `json:"key"`
	Type     string                 `json:"type"`
	Name     string                 `json:"name"`
	GroupId  *string                `json:"group_id"`
	Tenant   string                 `json:"-"`
	Content  map[string]interface{} `json:"content"`
	Children []ChildDictionary      `json:"children"`
}

type ChildDictionary struct {
	Key       string                 `json:"key"`
	Name      string                 `json:"name"`
	ParentKey string                 `json:"parent_key"`
	Content   map[string]interface{} `json:"content"`
}
