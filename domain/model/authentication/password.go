package authentication

type Password string

// ToString
// Summary: This is the function which converts the type Password to string.
// output: (string) string
func (p Password) ToString() string {
	return string(p)
}
