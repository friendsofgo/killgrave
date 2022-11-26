import React, { useCallback, useContext, useEffect, useRef } from 'react';
import Editor, { Monaco } from '@monaco-editor/react';
import { useEventBus } from '../../bus/hooks';
import {
  EventBusMessages,
  FontSizeChangedMessage, ImposterMatchedMessage,
  RequestReceivedMessage, ResponsePreparedMessage,
  ThemeChangedMessage
} from '../../bus/messages';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
import { Listener, Topics } from '../../bus/types';
import { Debugger, DebuggerContext, DebuggerState } from '../../debugger';

const DEFAULT_FONT_SIZE = 15;
const DEFAULT_EDITOR_TEXT = "// Waiting...";

export const CodeEditor = () => {
  // Monaco / Editor
  const editorRef = useRef<null | monaco.editor.IStandaloneCodeEditor>(null);
  const monacoRef = useRef<null | Monaco>(null);

  function handleEditorDidMount(editor: monaco.editor.IStandaloneCodeEditor, monaco: Monaco) {
    editorRef.current = editor;
    monacoRef.current = monaco;
  }

  const debug = useContext<Debugger>(DebuggerContext);
  const eventBus = useEventBus<EventBusMessages>();

  {
    useEffect(() => {
      const listeners: Listener[] = [
        eventBus.subscribe(Topics.THEME_CHANGED, onThemeChanged),
        eventBus.subscribe(Topics.FONT_SIZE_CHANGED, onFontSizeChanged),
        eventBus.subscribe(Topics.REQUEST_RECEIVED, onRequestReceived),
        eventBus.subscribe(Topics.IMPOSTER_MATCHED, onImposterMatched),
        eventBus.subscribe(Topics.RESPONSE_PREPARED, onResponsePrepared),
        eventBus.subscribe(Topics.RESPONSE_CONFIRMED, onResponseConfirmed),
        eventBus.subscribe(Topics.DEBUGGER_STARTED, onDebuggerStarted),
      ];

      return () => {
        listeners.forEach((l) => l.unsubscribe());
      };
    });

    const onThemeChanged = (m: ThemeChangedMessage) => {
      const themeName = m.dark ? 'vs-dark' : 'vs';
      monacoRef.current?.editor.setTheme(themeName);
    }

    const onFontSizeChanged = (m: FontSizeChangedMessage) => {
      const currFontSize = editorRef.current?.getRawOptions().fontSize || DEFAULT_FONT_SIZE;
      const newFontSize = m.increased ? currFontSize + 1 : currFontSize - 1
      editorRef.current?.updateOptions({fontSize: newFontSize});
    }

    const onRequestReceived = (m: RequestReceivedMessage) => {
      updateEditorCode(m.request);
    }

    const onImposterMatched = (m: ImposterMatchedMessage) => {
      updateEditorCode(m.imposter);
    }

    const onResponsePrepared = (m: ResponsePreparedMessage) => {
      updateEditorCode(m.response);
    }

    const onResponseConfirmed = () => {
      updateEditorCode(DEFAULT_EDITOR_TEXT);
    }

    const updateEditorCode = (content?: string) => {
      editorRef.current?.setValue(content || "");
      editorRef.current?.getAction('editor.action.formatDocument').run()
      editorRef.current?.focus();
    }

    const onDebuggerStarted = () => {
      switch (debug.state) {
        case DebuggerState.RequestReceived:
          const request = editorRef.current?.getValue() || "";
          eventBus.publish({topic: Topics.REQUEST_CONFIRMED, payload: {request}});
          break;
        case DebuggerState.ImposterMatched:
          const imposter = editorRef.current?.getValue() || "";
          eventBus.publish({topic: Topics.IMPOSTER_CONFIRMED, payload: {imposter}});
          break;
        case DebuggerState.ResponsePrepared:
          const response = editorRef.current?.getValue() || "";
          eventBus.publish({topic: Topics.RESPONSE_CONFIRMED, payload: {response}});
          break;
      }
    }
  }

  // Set up handle key press
  {
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
      document.addEventListener('keydown', handleKeyPress);
      return () => {
        document.removeEventListener('keydown', handleKeyPress);
      };
    }, [handleKeyPress]);
  }

  return <Editor
    height="50vh"
    className="border-4 h-full border-purplegrave p-3"
    defaultLanguage="json"
    defaultValue={DEFAULT_EDITOR_TEXT}
    options={{
      renderValidationDecorations: 'off',
      fontSize: DEFAULT_FONT_SIZE,
      language: 'json',
      minimap: {enabled: false},
      scrollbar: {verticalScrollbarSize: 5},
      lineNumbersMinChars: 1,
    }}
    onMount={handleEditorDidMount}
  />;
}