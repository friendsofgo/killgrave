import { SwitchThemeButton } from '../SwitchThemeButton';
import { FontSizeControl } from '../FontSizeControl';
import { HelpButton } from '../HelpButton';
import React from 'react';
import { PlayButton } from '../PlayButton';

export const Toolbar = () =>
  <>
    <SwitchThemeButton/>
    <FontSizeControl/>
    <PlayButton/>
    <HelpButton/>
  </>;