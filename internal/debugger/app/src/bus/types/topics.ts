export enum Topics {
  // notification
  NOTIFICATION = 'Notification',

  // ui
  THEME_CHANGED = 'ThemeChanged',
  FONT_SIZE_CHANGED = 'FontSizeChanged',

  // websocket
  CONNECTION_ESTABLISHED = 'ConnectionEstablished',
  CONNECTION_CLOSED = 'ConnectionClosed',

  // debugger
  DEBUGGER_STARTED = 'DebuggerStarted',
  DEBUGGER_STOPPED = 'DebuggerStopped',
  REQUEST_RECEIVED = 'RequestReceived',
  REQUEST_CONFIRMED = 'RequestConfirmed',
  IMPOSTER_MATCHED = 'ImposterMatched',
  IMPOSTER_CONFIRMED = 'ImposterConfirmed',
  RESPONSE_PREPARED = 'ResponsePrepared',
  RESPONSE_CONFIRMED = 'ResponseConfirmed',
}

