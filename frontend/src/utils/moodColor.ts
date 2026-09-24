export function moodColor(level: number) {
  return level >= 8 ? '#52c41a' : level >= 6 ? '#95de64' : level >= 4 ? '#faad14' : '#ff7875';
}

const BUILTIN_TAG_COLORS: Record<string, string> = {
  happy: 'gold',
  anxious: 'orange',
  tired: 'blue',
  angry: 'red',
  calm: 'green',
};

const CUSTOM_TAG_PALETTE = ['magenta', 'geekblue', 'purple', 'volcano', 'cyan', 'lime'];

// 内置标签沿用固定颜色；自定义标签按名字稳定取色，
// 停用后历史记录里的标签仍能渲染出同一颜色。
export function tagColor(tag: string): string {
  if (BUILTIN_TAG_COLORS[tag]) return BUILTIN_TAG_COLORS[tag];
  let hash = 0;
  for (let i = 0; i < tag.length; i++) {
    hash = (hash * 31 + tag.charCodeAt(i)) >>> 0;
  }
  return CUSTOM_TAG_PALETTE[hash % CUSTOM_TAG_PALETTE.length];
}
