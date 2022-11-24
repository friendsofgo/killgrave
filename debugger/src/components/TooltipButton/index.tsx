import React from 'react';
import { Tooltip } from 'flowbite-react';

export interface TooltipButtonProps {
  tooltip: string;
  button: JSX.Element;
}

export const TooltipButton: React.FC<TooltipButtonProps> = ({tooltip, button}) =>
  // eslint-disable-next-line
  <Tooltip style={'auto'} arrow={false} content={tooltip} animation="duration-100">
    {button}
  </Tooltip>;