package common_error

func ErrorsToString(errs []error) []string {
	var strErrs []string
	for _, err := range errs {
		strErrs = append(strErrs, err.Error())
	}
	return strErrs
}
