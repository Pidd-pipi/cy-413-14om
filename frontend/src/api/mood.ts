import {request} from '../utils/request';
import type {Mood, MoodTag, MoodTagDef} from '../types';

export const listMoods = (date?: string) =>
  request<Mood[]>(`/moods${date ? `?date=${date}` : ''}`);

export const saveMood = (payload: {
  mood_level: number;
  mood_tags: MoodTag[];
  note: string;
  record_date: string;
}) => request<Mood>('/moods', {method: 'POST', body: JSON.stringify(payload)});

export const updateMood = (
  id: number,
  payload: {mood_level: number; mood_tags: MoodTag[]; note: string; record_date: string},
) => request<Mood>(`/moods/${id}`, {method: 'PUT', body: JSON.stringify(payload)});

export const deleteMood = (id: number) => request(`/moods/${id}`, {method: 'DELETE'});

// 当前账号当前可选择的自定义标签（停用的不会返回）。
export const listMoodTags = () => request<MoodTagDef[]>('/mood-tags');

export const createMoodTag = (name: string) =>
  request<MoodTagDef>('/mood-tags', {method: 'POST', body: JSON.stringify({name})});

export const deactivateMoodTag = (id: number) =>
  request(`/mood-tags/${id}`, {method: 'DELETE'});
