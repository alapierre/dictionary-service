//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

var Columns = struct {
	Child struct {
		Key, Type, Tenant, ParentKey, Content string

		ParentKeyTypeTenant string
	}
	Dictionary struct {
		Key, Type, GroupID, Tenant, Content, Children string
	}
}{
	Child: struct {
		Key, Type, Tenant, ParentKey, Content string

		ParentKeyTypeTenant string
	}{
		Key:       "key",
		Type:      "type",
		Tenant:    "tenant",
		ParentKey: "parent_key",
		Content:   "content",

		ParentKeyTypeTenant: "ParentKeyTypeTenant",
	},
	Dictionary: struct {
		Key, Type, GroupID, Tenant, Content, Children string
	}{
		Key:      "key",
		Type:     "type",
		GroupID:  "group_id",
		Tenant:   "tenant",
		Content:  "content",
		Children: "children",
	},
}

var Tables = struct {
	Child struct {
		Name, Alias string
	}
	Dictionary struct {
		Name, Alias string
	}
}{
	Child: struct {
		Name, Alias string
	}{
		Name:  "child",
		Alias: "t",
	},
	Dictionary: struct {
		Name, Alias string
	}{
		Name:  "dictionary",
		Alias: "t",
	},
}

//noinspection ALL
type Child struct {
	tableName struct{} `sql:"child,alias:t" pg:",discard_unknown_columns"`

	Key       string                 `sql:"key,pk"`
	Type      string                 `sql:"type,pk"`
	Tenant    string                 `sql:"tenant,pk"`
	ParentKey string                 `sql:"parent_key,notnull"`
	Content   map[string]interface{} `sql:"content"`

	Parent *Dictionary `pg:"fk:parent_key,type,tenant" sql:"-"` // unsupported
}

//noinspection ALL
type Dictionary struct {
	tableName struct{} `sql:"dictionary,alias:t" pg:",discard_unknown_columns"`

	Key     string                 `sql:"key,pk"`
	Type    string                 `sql:"type,pk"`
	GroupID *string                `sql:"group_id"`
	Tenant  string                 `sql:"tenant,pk"`
	Content map[string]interface{} `sql:"content"`
}
