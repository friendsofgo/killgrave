export enum WsMessageType {
  Debugger = 1,
}

export interface WsMessage {
  type: WsMessageType;
  payload: DebuggerMessage;
}

export enum DebuggerMessageType {
  RequestReceived = 1,
  RequestContinued,
  ImposterReceived,
  ImposterContinued,
  ResponseReceived,
  ResponseContinued,
}

export interface DebuggerInfo {
  imposter: string;
  request: string;
  response: string;
}

export interface DebuggerMessage extends Partial<DebuggerInfo> {
  type: DebuggerMessageType;
}

export enum DebuggerState {
  Unknown = 'Unknown',
  WaitingForConnection = 'WaitingForConnection',
  Connected = 'Connected',
  WaitingForRequestConfirmation = 'WaitingForRequestConfirmation',
  WaitingForImposterConfirmation = 'WaitingForImposterConfirmation',
  WaitingForResponseConfirmation = 'WaitingForResponseConfirmation',
}

export enum DebuggerTransition {
  Unknown = 'Unknown',
  ConnectionRequested = 'ConnectionRequested',
  ConnectionEstablished = 'ConnectionEstablished',
  ConnectionClosed = 'ConnectionClosed',
  RequestReceived = 'RequestReceived',
  ImposterReceived = 'ImposterReceived',
  ResponseReceived = 'ResponseReceived',
}

export type TransitionCallback = () => void;

export class DebuggerStateMachine {
  state: DebuggerState;
  callbacks: Partial<TransitionCallbacksMap>;

  constructor() {
    this.state = DebuggerState.Unknown;
    this.callbacks = {};
  }

  register(transition: DebuggerTransition, callback: TransitionCallback) {
    this.callbacks[transition] = callback;
  }

  transition(to: DebuggerState) {
    const transition = getTransition(this.state, to)
    this.callbacks[transition]?.();
    this.state = to;

    console.log(`Transitioned from ${this.state} to ${to}`)
    console.log(`The transition name is ${transition}`)
  }
}

export type TransitionCallbacksMap = {
  [key in DebuggerTransition]: TransitionCallback;
};

// TODO: Review unknowns (although they could make sense)
// TODO: Evaluate the usage of Partial (plus fallback)
const getTransition = (from: DebuggerState, to: DebuggerState): DebuggerTransition => {
  const transitions: { [key in DebuggerState]: { [key in DebuggerState]: DebuggerTransition } } = {
    [DebuggerState.Unknown]: {
      [DebuggerState.Unknown]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForConnection]: DebuggerTransition.ConnectionRequested,
      [DebuggerState.Connected]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForRequestConfirmation]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForImposterConfirmation]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForResponseConfirmation]: DebuggerTransition.Unknown,
    },
    [DebuggerState.WaitingForConnection]: {
      [DebuggerState.Unknown]: DebuggerTransition.ConnectionClosed,
      [DebuggerState.WaitingForConnection]: DebuggerTransition.ConnectionRequested,
      [DebuggerState.Connected]: DebuggerTransition.ConnectionEstablished,
      [DebuggerState.WaitingForRequestConfirmation]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForImposterConfirmation]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForResponseConfirmation]: DebuggerTransition.Unknown,
    },
    [DebuggerState.Connected]: {
      [DebuggerState.Unknown]: DebuggerTransition.ConnectionClosed,
      [DebuggerState.WaitingForConnection]: DebuggerTransition.Unknown,
      [DebuggerState.Connected]: DebuggerTransition.ConnectionEstablished,
      [DebuggerState.WaitingForRequestConfirmation]: DebuggerTransition.RequestReceived,
      [DebuggerState.WaitingForImposterConfirmation]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForResponseConfirmation]: DebuggerTransition.Unknown,
    },
    [DebuggerState.WaitingForRequestConfirmation]: {
      [DebuggerState.Unknown]: DebuggerTransition.ConnectionClosed,
      [DebuggerState.WaitingForConnection]: DebuggerTransition.Unknown,
      [DebuggerState.Connected]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForRequestConfirmation]: DebuggerTransition.RequestReceived,
      [DebuggerState.WaitingForImposterConfirmation]: DebuggerTransition.ImposterReceived,
      [DebuggerState.WaitingForResponseConfirmation]: DebuggerTransition.Unknown,
    },
    [DebuggerState.WaitingForImposterConfirmation]: {
      [DebuggerState.Unknown]: DebuggerTransition.ConnectionClosed,
      [DebuggerState.WaitingForConnection]: DebuggerTransition.Unknown,
      [DebuggerState.Connected]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForRequestConfirmation]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForImposterConfirmation]: DebuggerTransition.ImposterReceived,
      [DebuggerState.WaitingForResponseConfirmation]: DebuggerTransition.ResponseReceived,
    },
    [DebuggerState.WaitingForResponseConfirmation]: {
      [DebuggerState.Unknown]: DebuggerTransition.ConnectionClosed,
      [DebuggerState.WaitingForConnection]: DebuggerTransition.Unknown,
      [DebuggerState.Connected]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForRequestConfirmation]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForImposterConfirmation]: DebuggerTransition.Unknown,
      [DebuggerState.WaitingForResponseConfirmation]: DebuggerTransition.ResponseReceived,
    }
  }

  return transitions[from][to];
}

