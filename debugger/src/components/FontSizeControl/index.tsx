import React from 'react';
import { IconGlassMinus, IconGlassPlus, } from '../../icons';
import { TooltipButton } from '../TooltipButton';
import { useEventBus } from '../../bus/hooks';
import { EventBusMessages } from '../../bus/messages';
import { FONT_SIZE_CHANGED_TOPIC, } from '../../bus/topics';
import { ClickIconButton } from '../ClickIconButton';

export const FontSizeControl: React.FC = () => {
  const eventBus = useEventBus<EventBusMessages>();

  const fontSizeIncreased = () => eventBus.publish(
    {topic: FONT_SIZE_CHANGED_TOPIC, payload: {increased: true}});

  const fontSizeDecreased = () => eventBus.publish(
    {topic: FONT_SIZE_CHANGED_TOPIC, payload: {increased: false}});

  return <>
    <TooltipButton
      tooltip={"Make font smaller"}
      button={<ClickIconButton icon={<IconGlassMinus/>}
                               onClick={fontSizeDecreased}/>
      }/>
    <TooltipButton
      tooltip={"Make font larger"}
      button={<ClickIconButton icon={<IconGlassPlus/>}
                               onClick={fontSizeIncreased}/>}
    />
  </>
}