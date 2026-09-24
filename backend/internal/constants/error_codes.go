package constants

const (
	CodeOK           = 0
	CodeValidation   = 1001
	CodeUnauthorized = 1002
	CodeForbidden    = 1003
	CodeNotFound     = 1004
	CodeConflict     = 1005
	CodeInternal     = 1500
)

var ErrorCodeHints = map[string]string{"Mood[mood_level]": "1-10", "Mood[mood_tags]": "happy/anxious/tired/angry/calm 或本账号启用中的自定义标签", "MoodTag[name]": "非空、最多8个字、同账号不可重复、最多12个", "Assessment[category]": "anxiety/depression/stress/sleep"}
