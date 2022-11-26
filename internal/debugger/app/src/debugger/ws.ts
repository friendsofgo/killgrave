export enum WsMessageType {
  Debugger = 1,
}

export interface WsMessage {
  type: WsMessageType;
  payload: DebuggerMessage;
}

export enum DebuggerMessageType {
  RequestReceived = 1,
  RequestConfirmed,
  ImposterMatched,
  ImposterConfirmed,
  ResponsePrepared,
  ResponseConfirmed,
}

export interface DebuggerInfo {
  imposter: string;
  request: string;
  response: string;
}

export interface DebuggerMessage extends Partial<DebuggerInfo> {
  type: DebuggerMessageType;
}
