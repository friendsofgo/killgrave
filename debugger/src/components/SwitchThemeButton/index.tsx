import React, { useState } from 'react';
import { IconMoon, IconSun } from '../../icons';
import { ToggleIconButton } from '../ToggleIconButton';
import { TooltipButton } from '../TooltipButton';
import { useEventBus } from '../../bus/hooks';
import { EventBusMessages } from '../../bus/messages';
import { THEME_CHANGED_TOPIC } from '../../bus/topics';

export const SwitchThemeButton: React.FC = () => {
  const eventBus = useEventBus<EventBusMessages>();
  const [dark, setDark] = useState<boolean>(false);

  const toggleDark = (dark: boolean) => {
    eventBus.publish({topic: THEME_CHANGED_TOPIC, payload: {dark}});
    setDark(dark);
  }

  return <TooltipButton
    tooltip={"Switch light/dark mode"}
    button={<ToggleIconButton
      show={dark}
      setShow={toggleDark}
      primary={<IconSun/>} secondary={<IconMoon/>}
    />}
  />;
}