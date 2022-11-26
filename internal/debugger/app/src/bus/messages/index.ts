import { NotificationMessage } from './notification';
import { FontSizeChangedMessage, ThemeChangedMessage } from './ui';
import {
  RequestReceivedMessage,
  RequestConfirmedMessage,
  ImposterMatchedMessage,
  ImposterConfirmedMessage,
  ResponsePreparedMessage,
  ResponseConfirmedMessage
} from './debugger';

export type { NotificationMessage } from './notification';
export type { FontSizeChangedMessage, ThemeChangedMessage } from './ui';
export type {
  RequestReceivedMessage,
  RequestConfirmedMessage,
  ImposterMatchedMessage,
  ImposterConfirmedMessage,
  ResponsePreparedMessage,
  ResponseConfirmedMessage
} from './debugger';


export interface EventBusMessages {
  // notification
  Notification: NotificationMessage;

  // ui
  ThemeChanged: ThemeChangedMessage;
  FontSizeChanged: FontSizeChangedMessage;

  // websocket
  ConnectionEstablished: {};
  ConnectionClosed: {};

  // debugger
  DebuggerStarted: {};
  DebuggerStopped: {};
  RequestReceived: RequestReceivedMessage;
  RequestConfirmed: RequestConfirmedMessage;
  ImposterMatched: ImposterMatchedMessage;
  ImposterConfirmed: ImposterConfirmedMessage;
  ResponsePrepared: ResponsePreparedMessage;
  ResponseConfirmed: ResponseConfirmedMessage;
}
