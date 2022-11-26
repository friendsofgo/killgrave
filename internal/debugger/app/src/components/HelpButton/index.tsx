import React from 'react';
import { IconQuestionMark } from '../icons';
import { TooltipButton } from '../TooltipButton';

export const HelpButton: React.FC = () => {
  return <TooltipButton
    tooltip={"Use CTRL + SHIFT + L to format code"}
    button={<button className="flex text-4xl font-bold text-purplegrave items-center cursor-pointer">
      {<IconQuestionMark/>}
    </button>}
  />;
}