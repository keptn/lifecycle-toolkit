package slos

type SLOComparison struct {
	CompareWith               string `json:"compare_with" yaml:"compare_with"`                           // single_result|several_results
	IncludeResultWithScore    string `json:"include_result_with_score" yaml:"include_result_with_score"` // all|pass|pass_or_warn
	NumberOfComparisonResults int    `json:"number_of_comparison_results" yaml:"number_of_comparison_results"`
	AggregateFunction         string `json:"aggregate_function" yaml:"aggregate_function"`
}

type SLOCriteria struct {
	Criteria []string `json:"criteria" yaml:"criteria"`
}

type SLO struct {
	SLI         string         `json:"sli" yaml:"sli"`
	DisplayName string         `json:"displayName" yaml:"displayName"`
	Pass        []*SLOCriteria `json:"pass" yaml:"pass"`
	Warning     []*SLOCriteria `json:"warning" yaml:"warning"`
	Weight      int            `json:"weight" yaml:"weight"`
	KeySLI      bool           `json:"key_sli" yaml:"key_sli"`
}

type SLOScore struct {
	Pass    string `json:"pass" yaml:"pass"`
	Warning string `json:"warning" yaml:"warning"`
}

// ServiceLevelObjectives describes SLO requirements
type ServiceLevelObjectives struct {
	SpecVersion string            `json:"spec_version" yaml:"spec_version"`
	Filter      map[string]string `json:"filter" yaml:"filter"`
	Comparison  *SLOComparison    `json:"comparison" yaml:"comparison"`
	Objectives  []*SLO            `json:"objectives" yaml:"objectives"`
	TotalScore  *SLOScore         `json:"total_score" yaml:"total_score"`
}
