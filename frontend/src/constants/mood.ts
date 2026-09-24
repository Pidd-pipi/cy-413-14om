export const MoodTagValues = {
  HAPPY: 'happy',
  ANXIOUS: 'anxious',
  TIRED: 'tired',
  ANGRY: 'angry',
  CALM: 'calm',
} as const;

export const BUILT_IN_MOOD_TAGS: string[] = Object.values(MoodTagValues);

export const MOOD_LABELS: Record<string, string> = {
  happy: '开心',
  anxious: '焦虑',
  tired: '疲惫',
  angry: '愤怒',
  calm: '平静',
};

export const MOOD_EMOJI: Record<string, string> = {
  happy: '😊',
  anxious: '😟',
  tired: '😴',
  angry: '😤',
  calm: '😌',
};

// 自定义标签的展示规则：文字本身就是名字，无内置 emoji。
export const moodTagLabel = (tag: string): string => MOOD_LABELS[tag] ?? tag;
export const moodTagEmoji = (tag: string): string => MOOD_EMOJI[tag] ?? '🏷️';
export const isBuiltInMoodTag = (tag: string): boolean => tag in MOOD_LABELS;

// 自定义标签的创建约束，与后端 constants/mood_tag.go 保持一致。
export const CUSTOM_TAG_MAX_LENGTH = 8;
export const CUSTOM_TAG_MAX_COUNT = 12;
