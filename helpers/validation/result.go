package validation

type Result struct {
	paramsIsValid bool
	hasError      bool
	firstError    string
	errors        map[string]string
}

func (r Result) HasError() bool {
	return r.hasError
}

func (r Result) ParamsIsValid() bool {
	return r.paramsIsValid
}

func (r Result) FirstError() string {
	return r.firstError
}

func (r Result) Errors() map[string]string {
	return r.errors
}
