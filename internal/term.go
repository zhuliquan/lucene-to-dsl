package internal

import "strings"

type Term struct {
	PhraseTerm *PhraseTerm `parser:"  @@" json:"phrase_term"`
	RegexpTerm *RegexpTerm `parser:"| @@" json:"regexp_term"`
	SimpleTerm *SimpleTerm `parser:"| @@" json:"simple_term"`
}

// func (t *Term) IsWillcard() bool {

// }

// func (t *Term) Boost()

func (t *Term) String() string {
	if t == nil {
		return ""
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.String()
	} else if t.RegexpTerm != nil {
		return t.RegexpTerm.String()
	} else if t.SimpleTerm != nil {
		return t.SimpleTerm.String()
	} else {
		return ""
	}
}

type RegexpTerm struct {
	Value string `parser:"@REGEXP" json:"value"`
}

func (t *RegexpTerm) String() string {
	if t == nil {
		return ""
	} else if t.Value != "" {
		return t.Value[1 : len(t.Value)-1]
	} else {
		return ""
	}
}

type PhraseTerm struct {
	Value string `parser:"@STRING" json:"value"`
	Fuzzy string `parser:"@FUZZY?" json:"fuzzy"`
	Boost string `parser:"@BOOST?" json:"boost"`
}

func (t *PhraseTerm) String() string {
	if t == nil {
		return ""
	} else if t.Value != "" {
		var res = "\" " + t.Value[1:len(t.Value)-1] + " \""
		if t.Fuzzy != "" {
			res += " " + t.Fuzzy
		}
		if t.Boost != "" {
			res += " " + t.Boost
		}
		return res
	} else {
		return ""
	}
}

func (t *PhraseTerm) isWildCard() bool {
	for i := 1; i < len(t.Value)-1; i++ {
		if i > 1 && (t.Value[i] == '?' || t.Value[i] == '*' && t.Value[i-1] != '\\') {
			return true
		}
		if i == 1 && (t.Value[i] == '?' || t.Value[i] == '*') {
			return true
		}
	}
	return false
}

type SimpleTerm struct {
	Value []string `parser:"@(IDENT|WILDCARD)+" json:"value"`
	Fuzzy string   `parser:"@FUZZY?" json:"fuzzy"`
	Boost string   `parser:"@BOOST?" json:"boost"`
}

func (t *SimpleTerm) String() string {
	if t == nil {
		return ""
	} else if len(t.Value) != 0 {
		var res = strings.Join(t.Value, "#")
		if t.Fuzzy != "" {
			res += " " + t.Fuzzy
		}
		if t.Boost != "" {
			res += " " + t.Boost
		}
		return res
	} else {
		return ""
	}
}

func (t *SimpleTerm) isWildcard() bool {
	for i := 0; i < len(t.Value); i++ {
		if t.Value[i] == "?" || t.Value[i] == "*" {
			return true
		}
	}
	return false
}

// func (t *SimpleTerm) getFuzzy() (float64, bool) {
// 	if t.Fuzzy == "" {
// 		return 0.0, false
// 	} else {
// 		return strconv.ParseFloat()
// 	}
// }
