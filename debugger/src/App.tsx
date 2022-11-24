import React, {
  MutableRefObject,
  useCallback,
  useEffect,
  useRef,
  useState
} from 'react';
import './App.css';
import logo from './logo.svg';
import Editor, { Monaco } from '@monaco-editor/react';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
import {
  DebuggerInfo,
  DebuggerMessageType,
  DebuggerState,
  DebuggerStateMachine,
  DebuggerTransition,
  WsMessage, WsMessageType
} from './Debugger';
import { IconPlay, IconStop } from './icons';
import { Console, HelpButton, SwitchThemeButton, ToggleIconButton, TooltipButton } from './components';
import { useEventBus } from './bus/hooks';
import { EventBusMessages, FontSizeChangedMessage, ThemeChangedMessage } from './bus/messages';
import { FONT_SIZE_CHANGED_TOPIC, NOTIFICATION_TOPIC, THEME_CHANGED_TOPIC } from './bus/topics';
import { EventBus } from './bus/types';
import { FontSizeControl } from './components/FontSizeControl';

const DEFAULT_FONT_SIZE = 14;

function App() {
  // Monaco / Editor
  const editorRef = useRef<null | monaco.editor.IStandaloneCodeEditor>(null);
  const monacoRef = useRef<null | Monaco>(null);

  function handleEditorDidMount(editor: monaco.editor.IStandaloneCodeEditor, monaco: Monaco) {
    editorRef.current = editor;
    monacoRef.current = monaco;
  }

  // Event bus
  const eventBus = useEventBus<EventBusMessages>();

  useEffect(() => {
    const themeListener = eventBus.subscribe(THEME_CHANGED_TOPIC, (m: ThemeChangedMessage) => {
      if (m.dark) {
        monacoRef.current?.editor.setTheme('vs-dark')
        document.getElementById('root')?.classList.add('dark')
      } else {
        monacoRef.current?.editor.setTheme('vs')
        document.getElementById('root')?.classList.remove('dark')
      }
    })

    const fontSizeListener = eventBus.subscribe(FONT_SIZE_CHANGED_TOPIC, (m: FontSizeChangedMessage) => {
      const currFontSize = editorRef.current?.getRawOptions().fontSize || DEFAULT_FONT_SIZE;
      const newFontSize = m.increased ? currFontSize + 1 : currFontSize - 1
      editorRef.current?.updateOptions({fontSize: newFontSize});
    })

    return () => {
      themeListener.unsubscribe();
      fontSizeListener.unsubscribe();
    };
  }, [eventBus]);

  // Debugger status
  const [running, setRunning] = useState<boolean>(false);

  // Set up handle key press
  const handleKeyPress = useCallback((event: any) => {
    if (event.ctrlKey && document.activeElement instanceof HTMLElement) {
      document.activeElement?.blur()
    }

    if (event.shiftKey === true && event.key === 'L') {
      editorRef.current?.getAction('editor.action.formatDocument').run()
      editorRef.current?.focus();
    }
  }, [editorRef]);

  useEffect(() => {
    // attach the event listener
    document.addEventListener('keydown', handleKeyPress);

    // remove the event listener
    return () => {
      document.removeEventListener('keydown', handleKeyPress);
    };
  }, [handleKeyPress]);

  // Set up debugger state machine
  const wsRef = useRef<null | WebSocket>(null);

  const debuggerInfo = useRef<Partial<DebuggerInfo>>({});

  const debuggerState = useRef<DebuggerStateMachine>(new DebuggerStateMachine());
  debuggerState.current.register(
    DebuggerTransition.ConnectionRequested,
    () => establishConnection(wsRef, debuggerInfo, debuggerState)
  );

  debuggerState.current.register(
    DebuggerTransition.RequestReceived,
    () => {
      setRunning(false);
      displayReceivedRequest(editorRef, debuggerInfo);
    }
  );

  debuggerState.current.register(
    DebuggerTransition.ImposterMatched,
    () => {
      setRunning(false);
      displayReceivedImposter(editorRef, debuggerInfo);
    }
  )

  debuggerState.current.register(
    DebuggerTransition.ResponsePrepared,
    () => {
      setRunning(false);
      displayReceivedResponse(editorRef, debuggerInfo);
    }
  )

  useEffect(() => {
    if (!running) {
      return;
    }

    switch (debuggerState.current.state) {
      case DebuggerState.Unknown:
        debuggerState.current?.transition(DebuggerState.WaitingForConnection);
        break;
      case DebuggerState.WaitingForRequestConfirmation:
        sendConfirmedRequest(wsRef, editorRef);
        break;
      case DebuggerState.WaitingForImposterConfirmation:
        sendConfirmedImposter(wsRef, editorRef);
        break;
      case DebuggerState.WaitingForResponseConfirmation:
        sendConfirmedResponse(wsRef, editorRef);
        break;
    }
  }, [running])

  return (
    <div className="App bg-white dark:bg-dark h-full">
      <div className="container">
        <h1 className="text-3xl pt-8 pb-4 font-bold text-purplegrave">Killgrave Debugger</h1>
        <div className="flex items-center justify-center text-purplegrave"><img className="h-24" src={logo}
                                                                                alt="Killgrave Logo"/></div>
        <div className="p-3 flex items-center justify-center">
          <SwitchThemeButton/>
          <FontSizeControl/>

          {/* Play / pause button*/}
          <TooltipButton
            tooltip={"Start / stop the debugger"}
            button={<ToggleIconButton
              show={running}
              setShow={(b: boolean) => {
                setRunning(b)
              }}
              primary={<IconStop/>} secondary={<IconPlay/>}
            />}
          />
          <HelpButton/>
        </div>
        <Editor
          height="45vh"
          className="border-4 h-full border-purplegrave p-3"
          defaultLanguage="json"
          defaultValue="// some comment"
          options={{
            fontSize: DEFAULT_FONT_SIZE,
            minimap: {enabled: false},
            scrollbar: {verticalScrollbarSize: 5},
            lineNumbersMinChars: 1,
          }}
          onMount={handleEditorDidMount}
        />

        <Console/>
      </div>

    </div>
  );
}

export default App;

//////

function establishConnection(
  wsRef: MutableRefObject<null | WebSocket>,
  debuggerInfo: MutableRefObject<Partial<DebuggerInfo>>,
  debuggerState: MutableRefObject<null | DebuggerStateMachine>,
) {
  const eventBus = EventBus<EventBusMessages>();

  wsRef.current = new WebSocket(`ws://localhost:3030/ws`);

  wsRef.current.onopen = () => {
    eventBus.publish({topic: NOTIFICATION_TOPIC, payload: {text: "Connection established"}});
    debuggerState.current?.transition(DebuggerState.Connected);
  };

  wsRef.current.onclose = () => {
    eventBus.publish({topic: NOTIFICATION_TOPIC, payload: {text: "Connection closed"}});
    debuggerState.current?.transition(DebuggerState.Unknown);
  };

  wsRef.current.onmessage = buildOnMessage(debuggerInfo, debuggerState);
}

function buildOnMessage(
  debuggerInfo: MutableRefObject<Partial<DebuggerInfo>>,
  debuggerState: MutableRefObject<null | DebuggerStateMachine>,
): (evt: MessageEvent) => void {
  return function (evt: MessageEvent) {
    const msg: WsMessage = JSON.parse(evt.data);

    switch (msg.payload.type) {
      case DebuggerMessageType.RequestReceived:
        debuggerInfo.current = {request: msg.payload.request};
        debuggerState.current?.transition(DebuggerState.WaitingForRequestConfirmation);
        break;
      case DebuggerMessageType.ImposterMatched:
        debuggerInfo.current = {imposter: msg.payload.imposter};
        debuggerState.current?.transition(DebuggerState.WaitingForImposterConfirmation);
        break;
      case DebuggerMessageType.ResponsePrepared:
        debuggerInfo.current = {response: msg.payload.response};
        debuggerState.current?.transition(DebuggerState.WaitingForResponseConfirmation);
        break;
    }
  }
}

function displayReceivedRequest(
  editorRef: MutableRefObject<null | monaco.editor.IStandaloneCodeEditor>,
  debuggerInfo: MutableRefObject<Partial<DebuggerInfo>>,
) {
  editorRef.current?.setValue(debuggerInfo.current.request || "");
  editorRef.current?.getAction('editor.action.formatDocument').run()
  editorRef.current?.focus();
}

function displayReceivedImposter(
  editorRef: MutableRefObject<null | monaco.editor.IStandaloneCodeEditor>,
  debuggerInfo: MutableRefObject<Partial<DebuggerInfo>>,
) {
  editorRef.current?.setValue(debuggerInfo.current.imposter || "");
  editorRef.current?.getAction('editor.action.formatDocument').run()
  editorRef.current?.focus();
}

function displayReceivedResponse(
  editorRef: MutableRefObject<null | monaco.editor.IStandaloneCodeEditor>,
  debuggerInfo: MutableRefObject<Partial<DebuggerInfo>>,
) {
  editorRef.current?.setValue(debuggerInfo.current.response || "");
  editorRef.current?.getAction('editor.action.formatDocument').run()
  editorRef.current?.focus();
}

function sendConfirmedRequest(
  wsRef: MutableRefObject<null | WebSocket>,
  editorRef: MutableRefObject<null | monaco.editor.IStandaloneCodeEditor>,
) {
  wsRef.current?.send(JSON.stringify({
    type: WsMessageType.Debugger,
    payload: {
      type: DebuggerMessageType.RequestConfirmed,
      request: editorRef.current?.getValue()
    },
  }));
}

function sendConfirmedImposter(
  wsRef: MutableRefObject<null | WebSocket>,
  editorRef: MutableRefObject<null | monaco.editor.IStandaloneCodeEditor>,
) {
  wsRef.current?.send(JSON.stringify({
    type: WsMessageType.Debugger,
    payload: {
      type: DebuggerMessageType.ImposterConfirmed,
      imposter: editorRef.current?.getValue()
    },
  }));
}

function sendConfirmedResponse(
  wsRef: MutableRefObject<null | WebSocket>,
  editorRef: MutableRefObject<null | monaco.editor.IStandaloneCodeEditor>,
) {
  wsRef.current?.send(JSON.stringify({
    type: WsMessageType.Debugger,
    payload: {
      type: DebuggerMessageType.ResponseConfirmed,
      response: editorRef.current?.getValue()
    },
  }));
}