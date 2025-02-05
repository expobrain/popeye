package internal

// ContextKey represents context key.
type ContextKey string

// A collection of context keys.
const (
	KeyFactory    ContextKey = "factory"
	KeyLabels     ContextKey = "labels"
	KeyFields     ContextKey = "fields"
	KeyOverAllocs ContextKey = "overAllocs"
	KeyRunInfo    ContextKey = "runInfo"
	KeyConfig     ContextKey = "config"
)
