import { TooltipButton } from '../TooltipButton';
import { ToggleIconButton } from '../ToggleIconButton';
import { IconPlay, IconStop } from '../icons';
import React, { useEffect, useState } from 'react';
import { useEventBus } from '../../bus/hooks';
import { EventBusMessages } from '../../bus/messages';
import { Listener, Topics } from '../../bus/types';

export const PlayButton: React.FC = () => {
  const eventBus = useEventBus<EventBusMessages>();
  const [running, setRunning] = useState<boolean>(false);

  useEffect(() => {
    const listeners: Listener[] = [
      eventBus.subscribe(Topics.REQUEST_RECEIVED, () => setRunning(false)),
      eventBus.subscribe(Topics.IMPOSTER_MATCHED, () => setRunning(false)),
      eventBus.subscribe(Topics.RESPONSE_PREPARED, () => setRunning(false)),
    ];

    return () => {
      listeners.forEach((l) => l.unsubscribe());
    };
  }, [eventBus]);

  const onDebugger = (running: boolean) => {
    setRunning(running);
    const topic = running ? Topics.DEBUGGER_STARTED : Topics.DEBUGGER_STOPPED;
    eventBus.publish({topic: topic, payload: {}});
  }

  return <TooltipButton
    tooltip={"Start / stop the debugger"}
    button={<ToggleIconButton
      show={running}
      setShow={onDebugger}
      primary={<IconStop/>} secondary={<IconPlay/>}
    />}
  />;
}