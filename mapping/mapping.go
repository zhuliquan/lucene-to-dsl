package mapping

type Mapping struct {
	Boost        float32 `json:"boost"`
	DocValues    bool    `json:"doc_values"`
	Index        bool    `json:"index"`      // Should the field be searchable? Accepts true (default) and false.
	NullValue    string  `json:"null_value"` //
	IncludeInAll bool    `json:"include_in_all"`
}

type KeywordMapping struct {
	Boost     float32 `json:"boost"`      //
	DocValues bool    `json:"doc_values"` // 字段是否需要存储
	Index     bool    `json:"index"`      // 是否需要构建索引
}

type AliasType struct {
	Path string `json:"path"`
}
