export interface RequestReceivedMessage {
  request: string;
}

export interface RequestConfirmedMessage {
  request: string;
}

export interface ImposterMatchedMessage {
  imposter: string;
}

export interface ImposterConfirmedMessage {
  imposter: string;
}

export interface ResponsePreparedMessage {
  response: string;
}

export interface ResponseConfirmedMessage {
  response: string;
}