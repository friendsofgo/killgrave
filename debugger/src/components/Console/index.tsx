import React, { useEffect, useState } from 'react';
import Terminal, { ColorMode, TerminalInput } from 'react-terminal-ui';
import { useEventBus } from '../../bus/hooks';
import { EventBusMessages, NotificationMessage, ThemeChangedMessage } from '../../bus/messages';
import { NOTIFICATION_TOPIC, THEME_CHANGED_TOPIC } from '../../bus/topics';
import { nanoid } from 'nanoid';

const DEFAULT_COLOR = ColorMode.Light;

export const Console = (props = {}) => {
  const eventBus = useEventBus<EventBusMessages>();
  const [color, setColor] = useState<ColorMode>(DEFAULT_COLOR);
  const [input, setInput] = useState<TerminalInput>(<TerminalInput>{"..."}</TerminalInput>);

  useEffect(() => {
    const themeListener = eventBus.subscribe(THEME_CHANGED_TOPIC, (m: ThemeChangedMessage) =>
      m.dark ? setColor(ColorMode.Dark) : setColor(ColorMode.Light));

    const notificationListener = eventBus.subscribe(NOTIFICATION_TOPIC, (m: NotificationMessage) =>
      setInput(<TerminalInput key={nanoid()}>{`${m.text}...`}</TerminalInput>));

    return () => {
      themeListener.unsubscribe();
      notificationListener.unsubscribe();
    };
  }, [eventBus]);

  return (
    <div className="pt-3 pb-3">
      <Terminal name='Console' colorMode={color}>{input}</Terminal>
    </div>
  )
};