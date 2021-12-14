package mapping

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Meta struct {
	Unit        MetaUnitType    `json:"unit,omitempty"`
	MetricsType MetaMetricsType `json:"metrics_type,omitempty"`
}

type Source struct {
	Enabled  bool     `json:"enabled,omitempty"`
	Includes []string `json:"includes,omitempty"`
	Excludes []string `json:"excludes,omitempty"`
}

type All struct {
	Enable bool `json:"anabled,omitempty"`
	Store  bool `json:"store,omitempty"`
}

type Property struct {
	Mapping
	Type FieldType `json:"type,omitempty"`

	// The index option controls whether field values are indexed. It accepts true or false and defaults to true.
	// Fields that are not indexed are not queryable.
	Index     bool   `json:"index,omitempty"`
	Store     bool   `json:"store,omitempty"`
	NullValue string `json:"null_value,omitempty"`
	Path      string `json:"path,omitempty"`
	Meta *Meta `json:"meta,omitempty"`
	// The index_options parameter controls what information is added to the inverted index for search and highlighting purposes.
	// WARNING: The index_options parameter is intended for use with text fields only. Avoid using index_options with other field data types.
	IndexOptions IndexOptions `json:"index_options,omitempty"`
	IncludeInAll bool         `json:"include_in_all,omitempty"`

	DepthLimit int `json:"depth_limit,omitempty"`

	// customized properties
	ExtProperties map[string]interface{} `json:"ext_properties,omitempty"`

	// NOTE: below parameters not used
	// Individual fields can be boosted automatically — count more towards the relevance score — at query time, with the boost parameter as follows:
	// NOTE: The boost is applied only for term queries (prefix, range and fuzzy queries are not boosted).
	// WARNING: Deprecated in 5.0.0. Index time boost is deprecated. Instead, the field mapping boost is applied at query time. For indices created before 5.0.0, the boost will still be applied at index time.
	// WARNING: Why index time boosting is a bad idea
	// We advise against using index time boosting for the following reasons:
	// 		1. You cannot change index-time boost values without reindexing all of your documents.
	// 		2. Every query supports query-time boosting which achieves the same effect.The difference is that you can tweak the boost value without having to reindex.
	// 		3. Index-time boosts are stored as part of the norm, which is only one byte. This reduces the resolution of the field length normalization factor which can lead to lower quality relevance calculations.
	Boost float64 `json:"boost,omitempty"`

	// The copy_to parameter allows you to copy the values of multiple fields into a group field, which can then be queried as a single field.
	// TIP: If you often search multiple fields, you can improve search speeds by using copy_to to search fewer fields. See Search as few fields as possible.
	// Some important points:
	// 		1. It is the field value which is copied, not the terms (which result from the analysis process).
	// 		2. The original _source field will not be modified to show the copied values.
	// 		3. The same value can be copied to multiple fields, with "copy_to": [ "field_1", "field_2" ]
	// 		4. You cannot copy recursively via intermediary fields such as a copy_to on field_1 to field_2 and copy_to on field_2 to field_3 expecting indexing into field_1 will eventuate in field_3, instead use copy_to directly to multiple fields from the originating field.
	// NOTE: copy-to is not supported for field types where values take the form of objects, e.g. date_range
	CopyTo interface{} `json:"copy_to,omitempty"`

	// Most fields are indexed by default, which makes them searchable.
	// The inverted index allows queries to look up the search term in unique sorted list of terms, and from that immediately have access to the list of documents that contain the term.
	// Sorting, aggregations, and access to field values in scripts requires a different data access pattern.
	// Instead of looking up the term and finding documents, we need to be able to look up the document and find the terms that it has in a field.
	// Doc values are the on-disk data structure, built at document index time, which makes this data access pattern possible.
	// They store the same values as the _source but in a column-oriented fashion that is way more efficient for sorting and aggregations.
	// Doc values are supported on almost all field types, with the notable exception of text and annotated_text fields.
	// All fields which support doc values have them enabled by default.
	// If you are sure that you don’t need to sort or aggregate on a field, or access the field value from a script, you can disable doc values in order to save disk space:
	DocValues bool `json:"doc_values,omitempty"`

	// Most fields are indexed by default, which makes them searchable.
	// The inverted index allows queries to look up the search term in unique sorted list of terms, and from that immediately have access to the list of documents that contain the term.
	// Sorting, aggregations, and access to field values in scripts requires a different data access pattern.
	// Instead of looking up the term and finding documents, we need to be able to look up the document and find the terms that it has in a field.
	// Doc values are the on-disk data structure, built at document index time, which makes this data access pattern possible.
	// They store the same values as the _source but in a column-oriented fashion that is way more efficient for sorting and aggregations.
	// Doc values are supported on almost all field types, with the notable exception of text and annotated_text fields.
	// All fields which support doc values have them enabled by default.
	// If you are sure that you don’t need to sort or aggregate on a field, or access the field value from a script, you can disable doc values in order to save disk space:
	EagerGlobalOrdinals bool `json:"eager_global_ordinals,omitempty"`

	// Strings longer than the ignore_above setting will not be indexed or stored.
	// For arrays of strings, ignore_above will be applied for each array element separately and string elements longer than ignore_above will not be indexed or stored.
	// NOTE: All strings/array elements will still be present in the _source field, if the latter is enabled which is the default in Elasticsearch.
	IgnoreAbove string `json:"ignore_above,omitempty"`

	// Sometimes you don’t have much control over the data that you receive.
	// One user may send a login field that is a date, and another sends a login field that is an email address.
	// Trying to index the wrong data type into a field throws an exception by default, and rejects the whole document.
	// The ignore_malformed parameter, if set to true, allows the exception to be ignored.
	// The malformed field is not indexed, but other fields in the document are processed normally.
	IgnoreMalformed bool `json:"ignore_malformed,omitempty"`

	// In JSON documents, dates are represented as strings.
	// Elasticsearch uses a set of preconfigured formats to recognize and parse these strings into a long value representing milliseconds-since-the-epoch in UTC.
	// Besides the built-in formats, your own custom formats can be specified using the familiar yyyy/MM/dd syntax:
	Format TimeFormat `json:"format,omitempty"`

	// 	Elasticsearch tries to index all of the fields you give it, but sometimes you want to just store the field without indexing it.
	// For instance, imagine that you are using Elasticsearch as a web session store.
	// You may want to index the session ID and last update time, but you don’t need to query or run aggregations on the session data itself.
	// The enabled setting, which can be applied only to the top-level mapping definition and to object fields, causes Elasticsearch to skip parsing of the contents of the field entirely.
	// The JSON can still be retrieved from the _source field, but it is not searchable or stored in any other way:
	Enabled bool `json:"enabled"`

	// Data is not always clean.
	// Depending on how it is produced a number might be rendered in the JSON body as a true JSON number, e.g. 5, but it might also be rendered as a string, e.g. "5".
	// Alternatively, a number that should be an integer might instead be rendered as a floating point, e.g. 5.0, or even "5.0".
	// Coercion attempts to clean up dirty values to fit the data type of a field. For instance:
	//     1、Strings will be coerced to numbers.
	//     2、Floating points will be truncated for integer values.
	Coerce bool `json:"coerce,omitempty"`
}

type Mapping struct {
	MappingType  MappingType          `json:"dynamic,omitempty"`
	IncludeInAll bool                 `json:"include_in_all,omitempty"`
	Source       *Source              `json:"_source,omitempty"`
	All          *All                 `json:"_all,omitempty"`
	Properties   map[string]*Property `json:"properties,omitempty"`
}

func (m *Mapping) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func LoadMapping(mappingPath string) (*Mapping, error) {
	if data, err := ioutil.ReadFile(mappingPath); err != nil {
		return nil, fmt.Errorf("failed to read mapping file: %s, err: %+v", mappingPath, err)
	} else {
		var mappingData = &Mapping{}
		if err := json.Unmarshal(data, mappingData); err != nil {
			return nil, fmt.Errorf("failed to parser mapping data: %s, err: %+v", data, err)
		} else {
			return mappingData, nil
		}
	}
}
