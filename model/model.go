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
type Dictionary struct {
	tableName struct{} `pg:"dictionary,alias:p" pg:",discard_unknown_columns"`

	Key       string                 `pg:"key,pk"`
	Type      string                 `pg:"type,pk"`
	Name      string                 `pg:"name"`
	GroupID   *string                `pg:"group_id"`
	Tenant    string                 `pg:"tenant,pk"`
	Content   map[string]interface{} `pg:"content"`
	Parent    bool                   `pg:"parent"`
	ParentKey *string                `pg:"parent_key"`
}
