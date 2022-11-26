import React from 'react';
import { IconGlassMinus, IconGlassPlus, } from '../icons';
import { TooltipButton } from '../TooltipButton';
import { useEventBus } from '../../bus/hooks';
import { EventBusMessages } from '../../bus/messages';
import { ClickIconButton } from '../ClickIconButton';
import { Topics } from '../../bus/types';

export const FontSizeControl: React.FC = () => {
  const eventBus = useEventBus<EventBusMessages>();

  const fontSizeIncreased = () => eventBus.publish(
    {topic: Topics.FONT_SIZE_CHANGED, payload: {increased: true}});

  const fontSizeDecreased = () => eventBus.publish(
    {topic: Topics.FONT_SIZE_CHANGED, payload: {increased: false}});

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