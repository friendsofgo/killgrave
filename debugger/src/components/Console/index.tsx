import React, { useEffect, useState } from 'react';
import Terminal, { ColorMode, TerminalInput } from 'react-terminal-ui';
import { useEventBus } from '../../bus/hooks';
import { EventBusMessages, NotificationMessage, ThemeChangedMessage } from '../../bus/messages';
import { nanoid } from 'nanoid';
import { Topics } from '../../bus/types';

const DEFAULT_COLOR = ColorMode.Light;
const DEFAULT_TEXT = "Click play button to start";

export const Console = (props = {}) => {
  const eventBus = useEventBus<EventBusMessages>();
  const [color, setColor] = useState<ColorMode>(DEFAULT_COLOR);
  const [input, setInput] = useState<TerminalInput>(<TerminalInput>{`${DEFAULT_TEXT}...`}</TerminalInput>);

  useEffect(() => {
    const themeListener = eventBus.subscribe(Topics.THEME_CHANGED, (m: ThemeChangedMessage) =>
      m.dark ? setColor(ColorMode.Dark) : setColor(ColorMode.Light));

    const notificationListener = eventBus.subscribe(Topics.NOTIFICATION, (m: NotificationMessage) =>
      setInput(<TerminalInput key={nanoid()}>{`${m.text}...`}</TerminalInput>));

    return () => {
      themeListener.unsubscribe();
      notificationListener.unsubscribe();
    };
  }, [eventBus]);

  return <Terminal name='Console' colorMode={color}>{input}</Terminal>;
};