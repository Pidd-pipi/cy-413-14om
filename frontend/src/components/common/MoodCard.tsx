import {Card, Tag, Typography} from 'antd';
import {moodTagLabel} from '../../constants/mood';
import {moodColor, tagColor} from '../../utils/moodColor';
import type {Mood} from '../../types';

export function MoodCard({mood}: {mood: Mood}) {
  let tags: string[] = [];
  try {
    const parsed: unknown = JSON.parse(mood.mood_tags);
    if (Array.isArray(parsed)) tags = parsed.filter((t): t is string => typeof t === 'string');
  } catch {
    // 历史脏数据保持为空标签展示，不影响其他字段。
  }
  return (
    <Card size="small" className="mood-card">
      <div className="mood-card__top">
        <b style={{color: moodColor(mood.mood_level)}}>心情 {mood.mood_level}/10</b>
        <span>{mood.record_date.slice(0, 10)}</span>
      </div>
      <div>
        {tags.map(tag => (
          <Tag color={tagColor(tag)} key={tag}>
            {moodTagLabel(tag)}
          </Tag>
        ))}
      </div>
      {mood.note && <Typography.Paragraph className="muted">{mood.note}</Typography.Paragraph>}
    </Card>
  );
}
