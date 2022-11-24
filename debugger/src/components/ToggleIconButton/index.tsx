import React, { Dispatch } from 'react';

export interface ToggleIconButtonProps {
  show: boolean;
  setShow: Dispatch<boolean>;
  primary: JSX.Element;
  secondary: JSX.Element;
  tooltip?: string;
}

export const ToggleIconButton: React.FC<ToggleIconButtonProps> = ({show, setShow, primary, secondary, tooltip}) =>
  (show ? (
    <button
      className="flex text-4xl font-bold text-purplegrave items-center cursor-pointer"
      onClick={() => setShow(!show)}
    >
      {primary}
    </button>
  ) : (
    <button
      className="flex text-4xl font-bold text-purplegrave items-center cursor-pointer"
      onClick={() => setShow(!show)}
    >
      {secondary}
    </button>
  ));