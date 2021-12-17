package gowikidata

import (
	"encoding/json"
	"strconv"
)

// WikiDataGetEntitiesRequest stores entities request url
type WikiDataGetEntitiesRequest struct {
	URL string
}

func (r *WikiDataGetEntitiesRequest) setParam(param string, values *[]string) {
	r.URL += createParam(param, *values)
}

// WikiDataGetClaimsRequest stores claims request url
type WikiDataGetClaimsRequest struct {
	URL string
}

func (r *WikiDataGetClaimsRequest) setParam(param string, values *[]string) {
	r.URL += createParam(param, *values)
}

// WikiDataSearchEntitiesRequest stores parameters for entities search
type WikiDataSearchEntitiesRequest struct {
	URL            string
	Limit          int
	Language       string
	Type           string
	Props          []string
	StrictLanguage bool
	Search         string
}

func (r *WikiDataSearchEntitiesRequest) setParam(param string, values *[]string) {
	r.URL += createParam(param, *values)
}

// Entity represents wikidata entities data
type Entity struct {
	ID           string                 `json:"id,omitempty"`
	PageID       int                    `json:"pageid,omitempty"`
	NS           int                    `json:"ns,omitempty"`
	Title        string                 `json:"title,omitempty"`
	LastRevID    int                    `json:"lastrevid,omitempty"`
	Modified     string                 `json:"modified,omitempty"`
	Type         string                 `json:"type,omitempty"`
	Label        string                 `json:"label,omitempty"`
	Labels       map[string]Label       `json:"labels,omitempty"`
	Descriptions map[string]Description `json:"descriptions,omitempty"`
	Aliases      map[string][]Alias     `json:"aliases,omitempty"`
	Claims       map[string][]Claim     `json:"claims,omitempty"`
	SiteLinks    map[string]SiteLink    `json:"sitelinks,omitempty"`
}

// GetDescription returns entity description in the given language code
func (e *Entity) GetDescription(languageCode string) string {
	return e.Descriptions[languageCode].Value
}

// GetLabel returns entity label in the given language code
func (e *Entity) GetLabel(languageCode string) string {
	return e.Labels[languageCode].Value
}

// Label represents wikidata labels data
type Label struct {
	Language    string `json:"language,omitempty"`
	Value       string `json:"value,omitempty"`
	ForLanguage string `json:"for-language,omitempty"`
}

// Description represents wikidata descriptions data
type Description struct {
	Language    string `json:"language,omitempty"`
	Value       string `json:"value,omitempty"`
	ForLanguage string `json:"for-language,omitempty"`
}

// Alias represents wikidata aliases data
type Alias struct {
	Language string `json:"language,omitempty"`
	Value    string `json:"value,omitempty"`
}

// SiteLink represents wikidata site links data
type SiteLink struct {
	Site   string   `json:"site,omitempty"`
	Title  string   `json:"title,omitempty"`
	Badges []string `json:"badges,omitempty"`
}

// Claim represents wikidata claims data
type Claim struct {
	ID              string            `json:"id,omitempty"`
	Label           string            `json:"label,omitempty"`
	Rank            string            `json:"rank,omitempty"`
	Type            string            `json:"type,omitempty"`
	MainSnak        Snak              `json:"mainsnak,omitempty"`
	Qualifiers      map[string][]Snak `json:"qualifiers,omitempty"`
	QualifiersOrder []string          `json:"qualifiers-order,omitempty"`
}

// Snak represents wikidata snak values
type Snak struct {
	SnakType  string    `json:"snaktype,omitempty"`
	Property  string    `json:"property,omitempty"`
	Hash      string    `json:"hash,omitempty"`
	DataType  string    `json:"datatype,omitempty"`
	DataValue DataValue `json:"datavalue,omitempty"`
}

// DataValue represents wikidata values
// Wikidata values can be either string or number
// It will store the data type so you can work with it accordingly
type DataValue struct {
	Type  string           `json:"type,omitempty"`
	Value DynamicDataValue `json:"value,omitempty"`
}

// DynamicDataValue represents wikidata values for DataValue struct
type DynamicDataValue struct {
	Data        interface{}
	S           string
	I           int
	ValueFields DataValueFields
	Type        string
}

// UnmarshalJSON unmarshales given json result to DynamicDataValue
// It's main job is to find the data type and set the fields accordingly
func (d *DynamicDataValue) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// If value starts with " and also ends with "
	// Then its string
	if string(s[0]) == "\"" && string(s[len(s)-1]) == "\"" {
		// Remove extra " from both sides of the string.
		cleaned := s[1 : len(s)-1]
		d.Data = cleaned
		d.S = cleaned
		d.Type = "String"
	} else {
		// If its int
		i, err := strconv.Atoi(s)
		if err != nil {
			// If its not int or string
			// Use DataValueFields
			values := DataValueFields{}
			err := json.Unmarshal(b, &values)
			if err != nil {
				return err
			}
			d.Type = "DataValueFields"
			d.ValueFields = values
			d.Data = values
		} else {
			// set value
			d.Type = "Int"
			d.I = i
			d.Data = i
		}
	}
	return
}

// DataValueFields represents wikidata value fields
type DataValueFields struct {
	EntityType    string  `json:"entity-type,omitempty"`
	NumericID     int     `json:"numeric-id,omitempty"`
	ID            string  `json:"id,omitempty"`
	Type          string  `json:"type,omitempty"`
	Value         string  `json:"value,omitempty"`
	Time          string  `json:"time,omitempty"`
	Precision     float64 `json:"precision,omitempty"`
	Before        int     `json:"before,omitempty"`
	After         int     `json:"after,omitempty"`
	TimeZone      int     `json:"timezone,omitempty"`
	CalendarModel string  `json:"calendarmodel,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	Globe         string  `json:"globe,omitempty"`
	Amount        string  `json:"amount,omitempty"`
	LowerBound    string  `json:"lowerbound,omitempty"`
	UpperBound    string  `json:"upperbound,omitempty"`
	Unit          string  `json:"unit,omitempty"`
	Text          string  `json:"text,omitempty"`
	Language      string  `json:"language,omitempty"`
}

// Reference represents wikidata references
type Reference struct {
	Hash       string            `json:"hash,omitempty"`
	Snaks      map[string][]Snak `json:"snaks,omitempty"`
	SnaksOrder []string          `json:"snaks-order,omitempty"`
}

// GetEntitiesResponse represents wikidata entities response
type GetEntitiesResponse struct {
	Entities map[string]Entity `json:"entities,omitempty"`
	Success  uint              `json:"success,omitempty"`
}

// GetClaimsResponse represents wikidata claims response
type GetClaimsResponse struct {
	Claims map[string][]Claim `json:"claims,omitempty"`
}

// SearchEntity represents wikidata entities search
type SearchEntity struct {
	Repository  string      `json:"repository,omitempty"`
	ID          string      `json:"id,omitempty"`
	ConceptURI  string      `json:"concepturi,omitempty"`
	Title       string      `json:"title,omitempty"`
	PageID      int         `json:"pageid,omitempty"`
	URL         string      `json:"url,omitempty"`
	Label       string      `json:"label,omitempty"`
	Description string      `json:"description,omitempty"`
	Match       SearchMatch `json:"match,omitempty"`
	DataType    string      `json:"datatype,omitempty"`
}

// SearchMatch represents wikidata search match value
type SearchMatch struct {
	Type     string `json:"type,omitempty"`
	Language string `json:"language,omitempty"`
	Text     string `json:"text,omitempty"`
}

// SearchInfo represents wikidata search info
type SearchInfo struct {
	Search string `json:"search,omitempty"`
}

// SearchEntitiesResponse represents wikidata entities search response
type SearchEntitiesResponse struct {
	SearchInfo      SearchInfo     `json:"searchinfo,omitempty"`
	SearchResult    []SearchEntity `json:"search,omitempty"`
	SearchContinue  int            `json:"search-continue,omitempty"`
	Success         uint           `json:"success,omitempty"`
	CurrentContinue int
	SearchRequest   WikiDataSearchEntitiesRequest
}

// WikiPediaQuery represents wikipedia query
type WikiPediaQuery struct {
	BatchComplete string `json:"batchcomplete,omitempty"`
	Query         struct {
		Pages map[string]struct {
			PageProps struct {
				WikiBaseItem string `json:"wikibase_item,omitempty"`
			} `json:"pageprops,omitempty"`
		} `json:"pages,omitempty"`
	} `json:"query,omitempty"`
}
