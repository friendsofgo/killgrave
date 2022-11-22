import React, {
  Dispatch,
  MouseEventHandler,
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
import { Tooltip } from 'flowbite-react';
import {
  DebuggerInfo,
  DebuggerMessageType,
  DebuggerState,
  DebuggerStateMachine,
  DebuggerTransition,
  WsMessage, WsMessageType
} from './Debugger';
import { IconGlassMinus, IconGlassPlus, IconMoon, IconPlay, IconStop, IconSun } from './icons';

function App() {
  // Font size
  const defaultFontSize = 12;
  const [fontSize, setFontSize] = useState<number>(defaultFontSize);
  const defaultLineNumberSize = 7;
  const [lineNumberSize, setLineNumberSize] = useState<number>(defaultLineNumberSize);

  const increaseFontSize = () => {
    setFontSize(fontSize + 1);
    setLineNumberSize(lineNumberSize + 0.1);
  }

  const decreaseFontSize = () => {
    setFontSize(fontSize - 1);
    setLineNumberSize(lineNumberSize - 0.1);
  }

  // Debugger status
  const [running, setRunning] = useState<boolean>(false);

  // Register actions
  const editorRef = useRef<null | monaco.editor.IStandaloneCodeEditor>(null);
  const monacoRef = useRef<null | Monaco>(null);

  function handleEditorDidMount(editor: monaco.editor.IStandaloneCodeEditor, monaco: Monaco) {
    editorRef.current = editor;
    monacoRef.current = monaco;
  }

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

  // Dark mode
  const [dark, setDark] = useState<boolean>(false);
  const toggleDark = (dark: boolean) => {
    if (dark) {
      monacoRef.current?.editor.setTheme('vs-dark')
      document.getElementById('root')?.classList.add('dark')
    } else {
      monacoRef.current?.editor.setTheme('vs')
      document.getElementById('root')?.classList.remove('dark')
    }
    setDark(dark);
  }

  return (
    <div className="App bg-white dark:bg-dark h-full">
      <div className="container">
        <h1 className="text-3xl pt-10 pb-5 font-bold text-purplegrave">Killgrave Debugger</h1>
        <div className="flex items-center justify-center text-purplegrave"><img className="h-24" src={logo}
                                                                                alt="Killgrave Logo"/></div>
        <div className="p-1 flex items-center justify-center text-purplegrave font-bold text-s italic">
          Tip: Use CTRL + SHIFT + L to format code
        </div>
        <div className="p-3 flex items-center justify-center">
          {/* Light / dark toggle button*/}
          <TooltipButton
            tooltip={"Toggle dark mode"}
            button={<ToggleIconButton
              show={dark}
              setShow={toggleDark}
              primary={<IconSun/>} secondary={<IconMoon/>}
            />}
          />

          {/* Smaller font button*/}
          <TooltipButton
            tooltip={"Make font smaller"}
            button={<ClickIconButton
              icon={<IconGlassMinus/>}
              onClick={decreaseFontSize}
            />}
          />

          {/* Larger font button*/}
          <TooltipButton
            tooltip={"Make font larger"}
            button={<ClickIconButton
              icon={<IconGlassPlus/>}
              onClick={increaseFontSize}
            />}
          />

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

        </div>
        <Editor
          height="55vh"
          className="border-4 h-full border-purplegrave p-3"
          defaultLanguage="json"
          defaultValue="// some comment"
          options={{
            fontSize: fontSize,
            minimap: {enabled: false},
            scrollbar: {verticalScrollbarSize: 5},
            lineNumbersMinChars: 1,
          }}
          onMount={handleEditorDidMount}
        />
      </div>

    </div>
  );
}

interface TooltipButtonProps {
  tooltip: string;
  button: JSX.Element;
}

const TooltipButton: React.FC<TooltipButtonProps> = ({tooltip, button}) =>
  // eslint-disable-next-line
  <Tooltip style={'auto'} arrow={false} content={tooltip} animation="duration-100">
    {button}
  </Tooltip>;

interface ToggleIconButtonProps {
  show: boolean;
  setShow: Dispatch<boolean>;
  primary: JSX.Element;
  secondary: JSX.Element;
  tooltip?: string;
}

const ToggleIconButton: React.FC<ToggleIconButtonProps> = ({show, setShow, primary, secondary, tooltip}) =>
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

interface ClickIconButtonProps {
  icon: JSX.Element;
  onClick: MouseEventHandler;
}

const ClickIconButton: React.FC<ClickIconButtonProps> = ({icon, onClick}) =>
  <button
    className="flex text-4xl font-bold text-purplegrave items-center cursor-pointer"
    onClick={onClick}
  >
    {icon}
  </button>;

export default App;

//////

function establishConnection(
  wsRef: MutableRefObject<null | WebSocket>,
  debuggerInfo: MutableRefObject<Partial<DebuggerInfo>>,
  debuggerState: MutableRefObject<null | DebuggerStateMachine>,
) {
  const href = window.location.href;
  const root = href.split("://")[1]

  wsRef.current = new WebSocket(`ws://${root}ws`);

  wsRef.current.onopen = () => {
    debuggerState.current?.transition(DebuggerState.Connected);
  };

  wsRef.current.onclose = () => {
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