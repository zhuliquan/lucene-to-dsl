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
	Index bool `json:"index,omitempty"`

	// A null value cannot be indexed or searched. When a field is set to null, (or an empty array or an array of null values) it is treated as though that field has no values.
	// The null_value parameter allows you to replace explicit null values with the specified value so that it can be indexed and searched. For instance:
	// IMPORTANT: The null_value needs to be the same data type as the field. For instance, a long field cannot have a string null_value.
	// NOTE: The null_value only influences how data is indexed, it doesn’t modify the _source document.
	NullValue string `json:"null_value,omitempty"`

	// Metadata attached to the field. This metadata is opaque to Elasticsearch, it is only useful for multiple applications that work on the same indices to share meta information about fields such as units
	// NOTE: Field metadata enforces at most 5 entries, that keys have a length that is less than or equal to 20, and that values are strings whose length is less than or equal to 50.
	// NOTE: Field metadata is updatable by submitting a mapping update. The metadata of the update will override the metadata of the existing field.
	// Elastic products use the following standard metadata entries for fields. You can follow these same metadata conventions to get a better out-of-the-box experience with your data.
	// unit
	// The unit associated with a numeric field: "percent", "byte" or a time unit.
	// By default, a field does not have a unit. Only valid for numeric fields.
	// The convention for percents is to use value 1 to mean 100%.
	// metric_type
	// The metric type of a numeric field: "gauge" or "counter".
	// A gauge is a single-value measurement that can go up or down over time, such as a temperature.
	// A counter is a single-value cumulative counter that only goes up, such as the number of requests processed by a web server.
	// By default, no metric type is associated with a field. Only valid for numeric fields.
	Meta *Meta `json:"meta,omitempty"`

	// It is often useful to index the same field in different ways for different purposes.
	// This is the purpose of multi-fields. For instance, a string field could be mapped as a text field for full-text search, and as a keyword field for sorting or aggregations:
	// You can add multi-fields to an existing field using the update mapping API.
	// NOTE: A multi-field mapping is completely separate from the parent field’s mapping.
	// NOTE: A multi-field doesn’t inherit any mapping options from its parent field. Multi-fields don’t change the original _source field.
	Fields map[string]*Property `json:"fields,omitempty"`

	// path is parameter for alias type
	// The path to the target field. Note that this must be the full path, including any parent objects (e.g. object1.object2.field).
	// There are a few restrictions on the target of an alias:
	//   1、The target must be a concrete field, and not an object or another field alias.
	//   2、The target field must exist at the time the alias is created.
	//   3、If nested objects are defined, a field alias must have the same nested scope as its target.
	// Additionally, a field alias can only have one target. This means that it is not possible to use a field alias to query over multiple target fields in a single clause.
	// An alias can be changed to refer to a new target through a mappings update. A known limitation is that if any stored percolator queries contain the field alias, they will still refer to its original target. More information can be found in the percolator documentation.
	Path string `json:"path,omitempty"`

	// relations is parameter for join type
	// The join data type is a special field that creates parent/child relation within documents of the same index. The relations section defines a set of possible relations within the documents, each relation being a parent name and a child name.
	// relations map you can define map[string]string / map[string][]string
	Relations map[string]interface{} `json:"relations,omitempty"`

	// If enabled, two-term word combinations (shingles) are indexed into a separate field.
	// This allows exact phrase queries (no slop) to run more efficiently, at the expense of a larger index.
	// Note that this works best when stopwords are not removed, as phrases containing stopwords will not use the subsidiary field and will fall back to a standard phrase query. Accepts true or false (default).
	IndexPhrases bool `json:"index_phrases,omitempty"`

	// The index_prefixes parameter enables the indexing of term prefixes to speed up prefix searches.
	// It accepts the following optional settings:
	// 		1. min_chars: The minimum prefix length to index. Must be greater than 0, and defaults to 2. The value is inclusive.
	// 		2. max_chars: The maximum prefix length to index. Must be less than 20, and defaults to 5. The value is inclusive.
	// This example creates a text field using the default prefix length settings:
	// 	PUT my-index-000001
	// {
	//   "mappings": {
	//     "properties": {
	//       "body_text": {
	//         "type": "text",
	//         "index_prefixes": { }
	//       }
	//     }
	//   }
	// }
	IndexPrefixes bool `json:"index_prefixes,omitempty"`

	// The normalizer property of keyword fields is similar to analyzer except that it guarantees that the analysis chain produces a single token.
	// The normalizer is applied prior to indexing the keyword, as well as at search-time when the keyword field is searched via a query parser such as the match query or via a term-level query such as the term query.
	// A simple normalizer called lowercase ships with elasticsearch and can be used.
	// Custom normalizers can be defined as part of analysis settings as follows.
	Normalizer interface{} `json:"normalizer,omitempty"`

	// customized properties
	ExtProperties map[string]interface{} `json:"ext_properties,omitempty"`

	// IMPORTANT: below parameters not used
	// include in _all
	IncludeInAll bool `json:"include_in_all,omitempty"`

	// WARNING: Only text fields support the analyzer mapping parameter.
	//
	// The analyzer parameter specifies the analyzer used for text analysis when indexing or searching a text field.
	// Unless overridden with the search_analyzer mapping parameter, this analyzer is used for both index and search analysis. See Specify an analyzer.
	//
	// Tip: We recommend testing analyzers before using them in production. See Test an analyzer.
	// Tip: The analyzer setting can not be updated on existing fields using the update mapping API.
	//
	// search_quote_analyzeredit
	// The search_quote_analyzer setting allows you to specify an analyzer for phrases, this is particularly useful when dealing with disabling stop words for phrase queries.
	// To disable stop words for phrases a field utilising three analyzer settings will be required:
	// 1. An analyzer setting for indexing all terms including stop words
	// 2. A search_analyzer setting for non-phrase queries that will remove stop words
	// 3. A search_quote_analyzer setting for phrase queries that will not remove stop words
	Analyzer            string `json:"analyzer,omitempty"`
	SearchAnalyzer      string `json:"search_analyzer,omitempty"`
	SearchQuoteAnalyzer string `json:"search_quote_analyzer,omitempty"`

	// Norms store various normalization factors that are later used at query time in order to compute the score of a document relatively to a query.
	// Although useful for scoring, norms also require quite a lot of disk (typically in the order of one byte per document per field in your index, even for documents that don’t have this specific field).
	// As a consequence, if you don’t need scoring on a specific field, you should disable norms on that field. In particular, this is the case for fields that are used solely for filtering or aggregations.
	// Norms can be disabled on existing fields using the update mapping API.
	// Norms can be disabled (but not reenabled after the fact), using the update mapping API like so:
	// NOTE: Norms will not be removed instantly, but will be removed as old segments are merged into new segments as you continue indexing new documents.
	// Any score computation on a field that has had norms removed might return inconsistent results since some documents won’t have norms anymore while other documents might still have norms.
	Norms bool `json:"norms,omitempty"`

	DepthLimit int `json:"depth_limit,omitempty"`

	// The index_options parameter controls what information is added to the inverted index for search and highlighting purposes.
	// WARNING: The index_options parameter is intended for use with text fields only. Avoid using index_options with other field data types.
	// The parameter accepts one of the following values. Each value retrieves information from the previous listed values. For example, freqs contains docs; positions contains both freqs and docs.
	// docs: Only the doc number is indexed. Can answer the question Does this term exist in this field?
	// freqs: Doc number and term frequencies are indexed. Term frequencies are used to score repeated terms higher than single terms.
	// positions (default): Doc number, term frequencies, and term positions (or order) are indexed. Positions can be used for proximity or phrase queries.
	// offsets: Doc number, term frequencies, positions, and start and end character offsets (which map the term back to the original string) are indexed. Offsets are used by the unified highlighter to speed up highlighting.
	IndexOptions IndexOptions `json:"index_options,omitempty"`

	// Individual fields can be boosted automatically — count more towards the relevance score — at query time, with the boost parameter as follows:
	// NOTE: The boost is applied only for term queries (prefix, range and fuzzy queries are not boosted).
	// WARNING: Deprecated in 5.0.0. Index time boost is deprecated. Instead, the field mapping boost is applied at query time. For indices created before 5.0.0, the boost will still be applied at index time.
	// WARNING: Why index time boosting is a bad idea
	// We advise against using index time boosting for the following reasons:
	// 		1. You cannot change index-time boost values without reindexing all of your documents.
	// 		2. Every query supports query-time boosting which achieves the same effect.The difference is that you can tweak the boost value without having to reindex.
	// 		3. Index-time boosts are stored as part of the norm, which is only one byte. This reduces the resolution of the field length normalization factor which can lead to lower quality relevance calculations.
	Boost float64 `json:"boost,omitempty"`

	// By default, field values are indexed to make them searchable, but they are not stored.
	// This means that the field can be queried, but the original field value cannot be retrieved.
	// Usually this doesn’t matter. The field value is already part of the _source field, which is stored by default.
	// If you only want to retrieve the value of a single field or of a few fields, instead of the whole _source, then this can be achieved with source filtering.
	// In certain situations it can make sense to store a field. For instance, if you have a document with a title, a date, and a very large content field, you may want to retrieve just the title and the date without having to extract those fields from a large _source field:
	// You can get stored fields by:
	// 	GET my-index-000001/_search
	//  {
	//     "stored_fields": [ "title", "date" ]
	//  }
	// NOTE: Stored fields returned as arrays
	// For consistency, stored fields are always returned as an array because there is no way of knowing if the original field value was a single value, multiple values, or an empty array.
	// If you need the original value, you should retrieve it from the _source field instead.
	Store bool `json:"store,omitempty"`

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
	Format string `json:"format,omitempty"`

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

	// similarity used for compute score
	// Elasticsearch allows you to configure a scoring algorithm or similarity per field.
	// The similarity setting provides a simple way of choosing a similarity algorithm other than the default BM25, such as TF/IDF.
	// Similarities are mostly useful for text fields, but can also apply to other field types.
	// Custom similarities can be configured by tuning the parameters of the built-in similarities.
	// For more details about this expert options, see the similarity module.
	// The only similarities which can be used out of the box, without any further configuration are:
	// 		1. BM25:
	// 		The Okapi BM25 algorithm. The algorithm used by default in Elasticsearch and Lucene.
	// 		2. classic:
	// 		[7.0.0] Deprecated in 7.0.0.The TF/IDF algorithm, the former default in Elasticsearch and Lucene.
	// 		boolean:
	// 		A simple boolean similarity, which is used when full-text ranking is not needed and the score should only be based on whether the query terms match or not. Boolean similarity gives terms a score equal to their query boost.
	// The similarity can be set on the field level when a field is first created, as follows:
	Similarity Similarity `json:"similarity,omitempty"`

	// used for search text array. you can treat text array as text separated by space when match_phrase' slop greater than position_increment_gap.
	// Analyzed text fields take term positions into account, in order to be able to support proximity or phrase queries. When indexing text fields with multiple values a "fake" gap is added between the values to prevent most phrase queries from matching across the values. The size of this gap is configured using position_increment_gap and defaults to 100.
	PositionIncrementGap int `json:"position_increment_gap,omitempty"`

	// term_vector is used in hightlighter, The fast vector highlighter will be used by default for the text field because term vectors are enabled.
	// Term vectors contain information about the terms produced by the analysis process, including: a list of terms.
	// the position (or order) of each term. the start and end character offsets mapping the term to its origin in the original string.
	// payloads (if they are available) — user-defined binary data associated with each term position.
	// These term vectors can be stored so that they can be retrieved for a particular document.
	// The term_vector setting accepts:
	// 		1. no: No term vectors are stored. (default)
	// 		2. yes: Just the terms in the field are stored.
	// 		3. with_positions: Terms and positions are stored.
	// 		4. with_offsets: Terms and character offsets are stored.
	// 		5. with_positions_offsets: Terms, positions, and character offsets are stored.
	// 		6. with_positions_payloads: Terms, positions, and payloads are stored.
	// 		7. with_positions_offsets_payloads: Terms, positions, offsets and payloads are stored.
	// NOTE: The fast vector highlighter requires with_positions_offsets. The term vectors API can retrieve whatever is stored.
	// WARNING: Setting with_positions_offsets will double the size of a field’s index.
	TermVector TermVector `json:"term_vector,omitempty"`
}

type Mapping struct {
	All         *All                 `json:"_all,omitempty"`
	Source      *Source              `json:"_source,omitempty"`
	MappingType MappingType          `json:"dynamic,omitempty"`
	Properties  map[string]*Property `json:"properties,omitempty"`
}

func (m *Mapping) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

type PropertyMapping struct {
	_mapping  *Mapping
	_cacheMap map[string]*Property
	_aliasMap map[string]string
	_extFuncs map[string]func(interface{}, map[string]interface{}) (interface{}, error)
}

func Init(mappingPath string,
	extFuncs map[string]func(interface{}, map[string]interface{}) (interface{}, error)) (*PropertyMapping, error) {
	if data, err := ioutil.ReadFile(mappingPath); err != nil {
		return nil, fmt.Errorf("failed to read mapping file: %s, err: %+v", mappingPath, err)
	} else {
		var mappingData = &Mapping{}
		if err := json.Unmarshal(data, mappingData); err != nil {
			return nil, fmt.Errorf("failed to parser mapping data: %s, err: %s", data, err)
		} else {
			var pm = &PropertyMapping{
				_mapping:  mappingData,
				_cacheMap: map[string]*Property{},
				_aliasMap: map[string]string{},
				_extFuncs: extFuncs,
			}
			if aliasMap, err := getAliasMap(pm); err != nil {
				return nil, err
			} else {
				pm._aliasMap = aliasMap
			}
			return pm, nil
		}
	}
}

func (m *PropertyMapping) GetProperty(field string) (*Property, error) {
	if target, have := m._aliasMap[field]; have {
		field = target
	}
	if property, have := m._cacheMap[field]; have {
		if checkTypeSupportLucene(property.Type) {
			return property, nil
		} else {
			return nil, fmt.Errorf("filed: %s type: %s don't support lucene query", field, property.Type)
		}
	} else {
		// 从 mapping 中去获取
		if property, err := fetchProperty(m, field); err != nil {
			return nil, err
		} else {
			m._cacheMap[field] = property
			return property, nil
		}
	}
}

func (m *PropertyMapping) GetExtFuncs(field string) func(interface{}, map[string]interface{}) (interface{}, error) {
	return m._extFuncs[field]
}
