package action

// Validation - hold messages related to validation of something
type Validation struct {
	Errors   []string
	Warnings []string
}

// AddWarning - add warning to validation
func (v *Validation) AddWarning(w string) {
	v.Warnings = append(v.Warnings, w)
}

// AddError - add error to validation
func (v *Validation) AddError(e string) {
	v.Errors = append(v.Errors, e)
}

// Valid - return true if no errors or warnings
func (v *Validation) Valid() bool {
	return len(v.Errors) == 0 && len(v.Warnings) == 0
}
