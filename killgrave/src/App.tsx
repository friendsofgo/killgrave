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
import { Buffer } from 'buffer';

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
    DebuggerTransition.ImposterReceived,
    () => {
      setRunning(false);
      displayReceivedImposter(editorRef, debuggerInfo);
    }
  )

  debuggerState.current.register(
    DebuggerTransition.ResponseReceived,
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
      <Sidebar/>
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
              primary={<SunIcon/>} secondary={<MoonIcon/>}
            />}
          />

          {/* Smaller font button*/}
          <TooltipButton
            tooltip={"Make font smaller"}
            button={<ClickIconButton
              icon={<GlassMinusIcon/>}
              onClick={decreaseFontSize}
            />}
          />

          {/* Larger font button*/}
          <TooltipButton
            tooltip={"Make font larger"}
            button={<ClickIconButton
              icon={<GlassPlusIcon/>}
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
              primary={<StopIcon/>} secondary={<PlayIcon/>}
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

function Sidebar() {
  const [showSidebar, setShowSidebar] = useState<boolean>(false);

  return (
    <>
      <div className="fixed left-12 top-12 z-50">
        <ToggleIconButton show={showSidebar} setShow={setShowSidebar} primary={<CrossIcon/>} secondary={<BarsIcon/>}/>
      </div>

      <div
        className={`top-0 right-0 w-[35vw] bg-orangegrave p-10 pl-20 text-white fixed h-full z-40 ease-in-out duration-700 ${
          showSidebar ? "translate-x-0 " : "translate-x-full"
        }`}
      >
        <h3 className="mt-20 text-4xl font-semibold text-white">
          I am a sidebar
        </h3>
      </div>
    </>
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

const CrossIcon = () =>
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={3}
       stroke="currentColor" className="w-8 h-8">
    <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12"/>
  </svg>;

const BarsIcon = () =>
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={3}
       stroke="currentColor" className="w-8 h-8">
    <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"/>
  </svg>;

const MoonIcon = () =>
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor"
       className="w-6 h-6">
    <path strokeLinecap="round" strokeLinejoin="round"
          d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z"/>
  </svg>;

const SunIcon = () =>
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor"
       className="w-6 h-6">
    <path strokeLinecap="round" strokeLinejoin="round"
          d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"/>
  </svg>;

const GlassPlusIcon = () =>
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor"
       className="w-6 h-6">
    <path strokeLinecap="round" strokeLinejoin="round"
          d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607zM10.5 7.5v6m3-3h-6"/>
  </svg>


const GlassMinusIcon = () =>
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor"
       className="w-6 h-6">
    <path strokeLinecap="round" strokeLinejoin="round"
          d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607zM13.5 10.5h-6"/>
  </svg>

const PlayIcon = () =>
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor"
       className="w-6 h-6">
    <path strokeLinecap="round" strokeLinejoin="round" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
    <path strokeLinecap="round" strokeLinejoin="round"
          d="M15.91 11.672a.375.375 0 010 .656l-5.603 3.113a.375.375 0 01-.557-.328V8.887c0-.286.307-.466.557-.327l5.603 3.112z"/>
  </svg>


const StopIcon = () =>
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor"
       className="w-6 h-6">
    <path strokeLinecap="round" strokeLinejoin="round" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
    <path strokeLinecap="round" strokeLinejoin="round"
          d="M9 9.563C9 9.252 9.252 9 9.563 9h4.874c.311 0 .563.252.563.563v4.874c0 .311-.252.563-.563.563H9.564A.562.562 0 019 14.437V9.564z"/>
  </svg>


export default App;


//////

function establishConnection(
  wsRef: MutableRefObject<null | WebSocket>,
  debuggerInfo: MutableRefObject<Partial<DebuggerInfo>>,
  debuggerState: MutableRefObject<null | DebuggerStateMachine>,
) {
  wsRef.current = new WebSocket('ws://localhost:8080/ws'); // TODO: Parameterize and use ref

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
    const msg: WsMessage = parseAndDecode(evt);

    switch (msg.payload.type) {
      case DebuggerMessageType.RequestReceived:
        debuggerInfo.current = {...msg.payload};
        debuggerState.current?.transition(DebuggerState.WaitingForRequestConfirmation);
        break;
      case DebuggerMessageType.ImposterReceived:
        debuggerInfo.current = {...msg.payload};
        debuggerState.current?.transition(DebuggerState.WaitingForImposterConfirmation);
        break;
      case DebuggerMessageType.ResponseReceived:
        debuggerInfo.current = {...msg.payload};
        debuggerState.current?.transition(DebuggerState.WaitingForResponseConfirmation);
        break;
    }
  }
}

const parseAndDecode = (evt: MessageEvent): WsMessage => {
  // Parse message from event
  const msg: WsMessage = JSON.parse(evt.data);
  console.log(msg);

  // Decode
  msg.payload.imposter = decode(msg.payload.imposter);
  msg.payload.request = decode(msg.payload.request);
  msg.payload.response = decode(msg.payload.response);
  console.log(msg.payload);

  return msg;
}

const decode = (s?: string): string => {
  const enc: string = s || "";
  const buff: Buffer = Buffer.from(enc, 'base64');
  return buff.toString();
}

const encode = (s?: string): string => {
  const dec: string = s || "";
  const buff: Buffer = Buffer.from(dec);
  return buff.toString('base64');
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
      type: DebuggerMessageType.RequestContinued,
      request: encode(editorRef.current?.getValue())
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
      type: DebuggerMessageType.ImposterContinued,
      imposter: encode(editorRef.current?.getValue())
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
      type: DebuggerMessageType.ResponseContinued,
      response: encode(editorRef.current?.getValue())
    },
  }));
}