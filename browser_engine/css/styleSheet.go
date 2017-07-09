package css

// StyleSheet - A stylesheet
type StyleSheet struct {
	Rules []Rule
}

func (s *StyleSheet) String() string {
	return "Rules: "
}

// Rule - A stylesheet rule
type Rule struct {
	Selectors    []Selector
	Declarations []StyleDeclaration
}
