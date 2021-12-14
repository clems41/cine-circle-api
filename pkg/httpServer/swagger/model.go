package swagger

type TagName string
type TagDescription string

type Info struct {
	Title        string
	Description  string
	ContactName  string
	ContactEmail string
	Version      string
	Tags         map[TagName]TagDescription
}
