import {useEffect, useState} from 'react';
import {customTagStore} from '../stores/moodStore';
import type {MoodTagDef} from '../types';

// 读取当前账号可选择的自定义标签；任一页面新建/停用后所有使用方即时刷新。
export function useCustomTags(): MoodTagDef[] {
  const [tags, setTags] = useState<MoodTagDef[]>(customTagStore.get());
  useEffect(() => {
    const unsub = customTagStore.subscribe(() => setTags(customTagStore.get()));
    customTagStore.load().catch(() => undefined);
    return () => {
      unsub();
    };
  }, []);
  return tags;
}
