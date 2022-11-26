import React, { MouseEventHandler } from 'react';

export interface ClickIconButtonProps {
  icon: JSX.Element;
  onClick: MouseEventHandler;
}

export const ClickIconButton: React.FC<ClickIconButtonProps> = ({icon, onClick}) =>
  <button
    className="flex text-4xl font-bold text-purplegrave items-center cursor-pointer"
    onClick={onClick}
  >
    {icon}
  </button>;