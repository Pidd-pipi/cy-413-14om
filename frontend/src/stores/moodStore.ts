import {listMoodTags, createMoodTag, deactivateMoodTag} from '../api/mood';
import type {MoodTagDef} from '../types';

let moods = [] as import('../types').Mood[];
export const moodStore = {
  get: () => moods,
  set: (values: import('../types').Mood[]) => {
    moods = values;
  },
};

// 自定义标签在 Dashboard 快速记录与 Moods 新增记录间共享，
// 用极简订阅让两个页面在新建/停用后立即同步选择区。
let customTags: MoodTagDef[] = [];
const listeners = new Set<() => void>();
let loaded = false;
let inflight: Promise<MoodTagDef[]> | null = null;

function emit() {
  listeners.forEach(fn => fn());
}

export const customTagStore = {
  get: () => customTags,
  subscribe(fn: () => void): () => void {
    listeners.add(fn);
    return () => {
      listeners.delete(fn);
    };
  },
  async load(force = false): Promise<MoodTagDef[]> {
    if (loaded && !force) return customTags;
    if (inflight) return inflight;
    inflight = listMoodTags()
      .then(vs => {
        customTags = vs;
        loaded = true;
        emit();
        return vs;
      })
      .finally(() => {
        inflight = null;
      });
    return inflight;
  },
  async create(name: string): Promise<MoodTagDef> {
    const v = await createMoodTag(name);
    // 服务端对“停用后同名重建”会原地恢复；以返回的 id 去重后再放入选择区。
    customTags = [...customTags.filter(t => t.id !== v.id), v];
    emit();
    return v;
  },
  async deactivate(id: number) {
    await deactivateMoodTag(id);
    customTags = customTags.filter(t => t.id !== id);
    emit();
  },
};
