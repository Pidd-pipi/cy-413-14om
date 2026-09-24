package constants

const (
	MessageOK             = "ok"
	MessageCreated        = "created"
	MessageUnauthorized   = "请先登录后继续"
	MessageForbidden      = "当前角色无权操作"
	MessageMoodSaved      = "情绪已种入你的花园"
	MessageJournalSaved   = "日记已安全保存"
	MessageAssessmentDone = "测评结果已生成"
	MessageInternal       = "服务暂时不可用"

	MessageMoodTagNameBlank   = "标签名不能为空"
	MessageMoodTagNameTooLong = "标签名最多 8 个字"
	MessageMoodTagNameBuiltIn = "不能与内置标签同名"
	MessageMoodTagNameExists  = "同名标签已经存在"
	MessageMoodTagLimit       = "最多只能创建 12 个自定义标签"
	MessageMoodTagNotOwner    = "标签不存在或不属于当前账号"
	MessageMoodTagUnsupported = "只能使用当前可选择的情绪标签"
)
