import {useState} from 'react';
import {Button, Input, Popconfirm, Popover, Space, Tooltip, message} from 'antd';
import {PlusOutlined} from './icons';
import {
  BUILT_IN_MOOD_TAGS,
  CUSTOM_TAG_MAX_COUNT,
  CUSTOM_TAG_MAX_LENGTH,
  isBuiltInMoodTag,
  moodTagEmoji,
  moodTagLabel,
} from '../../constants/mood';
import {useCustomTags} from '../../hooks/useCustomTags';
import {customTagStore} from '../../stores/moodStore';

interface MoodSelectorProps {
  value: string[];
  onChange: (v: string[]) => void;
  // 停用某标签时，把当前已勾选的同名标签一并从 value 中移除。
  onTagDeactivated?: (name: string) => void;
}

export function MoodSelector({value, onChange, onTagDeactivated}: MoodSelectorProps) {
  const customTags = useCustomTags();
  const [adding, setAdding] = useState(false);
  const [draft, setDraft] = useState('');
  const [saving, setSaving] = useState(false);

  const toggle = (tag: string) =>
    onChange(value.includes(tag) ? value.filter(v => v !== tag) : [...value, tag]);

  const submit = async () => {
    const name = draft.trim();
    if (!name) {
      message.warning('标签名不能为空');
      return;
    }
    if ([...name].length > CUSTOM_TAG_MAX_LENGTH) {
      message.warning(`标签名最多 ${CUSTOM_TAG_MAX_LENGTH} 个字`);
      return;
    }
    if (isBuiltInMoodTag(name) || BUILT_IN_MOOD_TAGS.includes(name)) {
      message.warning('该名字与内置标签重复了');
      return;
    }
    if (customTags.some(t => t.name === name)) {
      message.warning('同名标签已经存在');
      return;
    }
    if (customTags.length >= CUSTOM_TAG_MAX_COUNT) {
      message.warning(`最多只能创建 ${CUSTOM_TAG_MAX_COUNT} 个自定义标签`);
      return;
    }
    setSaving(true);
    try {
      const created = await customTagStore.create(name);
      // 新建成功后马上能勾选：自动加入当前选择。
      if (!value.includes(created.name)) onChange([...value, created.name]);
      message.success('标签已添加');
      setDraft('');
      setAdding(false);
    } catch (e) {
      message.error((e as Error).message);
    } finally {
      setSaving(false);
    }
  };

  const remove = async (id: number, name: string) => {
    try {
      await customTagStore.deactivate(id);
      onChange(value.filter(v => v !== name));
      onTagDeactivated?.(name);
      message.success('标签已停用，历史记录仍会保留');
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const addPanel = (
    <div style={{width: 220}}>
      <Input
        autoFocus
        maxLength={32}
        placeholder={`输入标签名（最多 ${CUSTOM_TAG_MAX_LENGTH} 个字）`}
        value={draft}
        onChange={e => setDraft(e.target.value)}
        onPressEnter={submit}
        showCount
      />
      <div style={{marginTop: 8, textAlign: 'right'}}>
        <Space>
          <Button size="small" onClick={() => setAdding(false)}>
            取消
          </Button>
          <Button size="small" type="primary" loading={saving} onClick={submit}>
            添加
          </Button>
        </Space>
      </div>
    </div>
  );

  return (
    <Space wrap>
      {BUILT_IN_MOOD_TAGS.map(tag => (
        <Tooltip title={moodTagLabel(tag)} key={tag}>
          <Button
            aria-label={`选择${moodTagLabel(tag)}`}
            type={value.includes(tag) ? 'primary' : 'default'}
            onClick={() => toggle(tag)}
          >
            {moodTagEmoji(tag)} {moodTagLabel(tag)}
          </Button>
        </Tooltip>
      ))}
      {customTags.map(t => (
        <Tooltip title={`${t.name}（点击勾选，× 停用）`} key={t.id}>
          <Button
            type={value.includes(t.name) ? 'primary' : 'default'}
            onClick={() => toggle(t.name)}
          >
            {moodTagEmoji(t.name)} {t.name}
            <Popconfirm
              title={`停用标签「${t.name}」？`}
              description="停用后不再出现在选择区，过去记录的文字与统计仍会保留。"
              okText="停用"
              cancelText="再想想"
              okButtonProps={{danger: true}}
              onConfirm={() => remove(t.id, t.name)}
            >
              <span
                aria-label={`停用${t.name}`}
                style={{marginLeft: 6, opacity: 0.6}}
                onClick={e => e.stopPropagation()}
              >
                ×
              </span>
            </Popconfirm>
          </Button>
        </Tooltip>
      ))}
      <Popover
        title="新增自定义标签"
        content={addPanel}
        trigger="click"
        open={adding}
        onOpenChange={setAdding}
      >
        <Button icon={<PlusOutlined />} disabled={customTags.length >= CUSTOM_TAG_MAX_COUNT}>
          自定义
        </Button>
      </Popover>
      <span className="muted" style={{fontSize: 12}}>
        {customTags.length}/{CUSTOM_TAG_MAX_COUNT}
      </span>
    </Space>
  );
}
