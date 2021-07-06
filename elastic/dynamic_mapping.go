package elastic

type DynamicTemplate struct {
	Settings map[string]interface{} `json:"settings,omitempty"`
	Mappings *Mappings              `json:"mappings,omitempty"`
}

type Mappings struct {
	DynamicTemplates []map[string]MatchConditions `json:"dynamic_templates"`
	Properties       map[string]Property          `json:"properties"`
}

type Property struct {
	Type string `json:"type"`
}

// NewESMappings
func NewMappings() *Mappings {
	return &Mappings{
		Properties: make(map[string]Property),
	}
}

// AddDynamicTemplate
func (es *Mappings) AddDynamicTemplate(name string, matchCondition MatchConditions) *Mappings {
	m := make(map[string]MatchConditions)
	m[name] = matchCondition

	es.DynamicTemplates = append(es.DynamicTemplates, m)
	return es
}

// AddPropertyWithType
func (es *Mappings) AddPropertyWithType(name, typ string) *Mappings {
	es.Properties[name] = Property{Type: typ}
	return es
}
