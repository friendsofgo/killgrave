import React, {
  useEffect,
} from 'react';
import './App.css';
import { Console, CodeEditor, Toolbar, AppLogo } from './components';
import { useEventBus } from './bus/hooks';
import { EventBusMessages, ThemeChangedMessage } from './bus/messages';
import { Topics } from './bus/types';
import process from 'process';

export const isDevelopment: boolean = !process.env.NODE_ENV || process.env.NODE_ENV === 'development';

function App() {
  const eventBus = useEventBus<EventBusMessages>();

  useEffect(() => {
    const themeListener = eventBus.subscribe(Topics.THEME_CHANGED, (m: ThemeChangedMessage) => {
      const root = document.getElementById('root');
      m.dark ? root?.classList.add('dark') : root?.classList.remove('dark')
    })

    return () => {
      themeListener.unsubscribe();
    };
  }, [eventBus]);

  return (
    <div className="App bg-white dark:bg-dark h-full">
      <div className="container">
        <h1 className="text-3xl pt-8 pb-4 font-bold text-purplegrave">
          Killgrave Debugger
        </h1>
        <div className="flex items-center justify-center text-purplegrave">
          <AppLogo/>
        </div>
        <div className="p-3 flex justify-center w-full">
          <div className="flex flex-col items-center justify-center z-10">
            <Toolbar/>
          </div>
          <div className="p-3 flex justify-center w-full">
            <CodeEditor/>
          </div>
        </div>
        <div className="pl-12 pr-5"><Console/></div>
      </div>
    </div>
  );
}

export default App;