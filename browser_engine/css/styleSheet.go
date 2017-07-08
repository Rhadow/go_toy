package css

// StyleSheet - A stylesheet
type StyleSheet struct {
	Rules []Rule
}

// Rule - A stylesheet rule
type Rule struct {
	Selectors    []Selector
	Declarations []StyleDeclaration
}
