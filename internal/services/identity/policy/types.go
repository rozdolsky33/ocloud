package policy

// IndexablePolicy represents a policy structure optimized for indexing and searching operations.
// It includes fields for policy name, description, statements, and flattened tag information.
type IndexablePolicy struct {
	Name        string
	Description string
	Statement   string
	Tags        string
	TagValues   string
}
