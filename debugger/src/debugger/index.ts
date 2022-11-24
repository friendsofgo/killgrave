import { EventBus, Topics } from '../bus/types';
import {
  EventBusMessages,
  ImposterConfirmedMessage,
  RequestConfirmedMessage, ResponseConfirmedMessage
} from '../bus/messages';
import { DebuggerMessageType, WsMessage, WsMessageType } from './ws';
import { nanoid } from 'nanoid';
import React from 'react';
import { isDevelopment } from '../App';

export enum DebuggerState {
  Default,
  RequestReceived,
  ImposterMatched,
  ResponsePrepared,
}

export class Debugger {
  bus?: any;
  state: DebuggerState;
  wsRef?: WebSocket;

  constructor() {
    console.log(nanoid());
    this.state = DebuggerState.Default;
    this.bus = EventBus<EventBusMessages>();
    this.bus.subscribe(Topics.DEBUGGER_STARTED, this.onDebuggerStarted);
    this.bus.subscribe(Topics.DEBUGGER_STOPPED, this.onDebuggerStopped);
    this.bus.subscribe(Topics.REQUEST_RECEIVED, this.onRequestReceived);
    this.bus.subscribe(Topics.REQUEST_CONFIRMED, this.onRequestConfirmed);
    this.bus.subscribe(Topics.IMPOSTER_MATCHED, this.onImposterMatched);
    this.bus.subscribe(Topics.IMPOSTER_CONFIRMED, this.onImposterConfirmed);
    this.bus.subscribe(Topics.RESPONSE_PREPARED, this.onResponsePrepared);
    this.bus.subscribe(Topics.RESPONSE_CONFIRMED, this.onResponseConfirmed);
  }

  onDebuggerStarted = () => {
    if (this.wsRef) {
      return;
    }

    const remote = isDevelopment
      ? `ws://localhost:3030/ws`
      : `ws://${window.location.href.split("://")[1]}ws`;

    this.wsRef = new WebSocket(remote);
    console.log(`Trying to connect to... ${remote}`)

    this.wsRef.onopen = () => {
      this.bus.publish({topic: Topics.CONNECTION_ESTABLISHED, payload: {}});
      this.bus.publish({topic: Topics.NOTIFICATION, payload: {text: "Connection established"}});
    };

    this.wsRef.onerror = () => {
      this.bus.publish({topic: Topics.NOTIFICATION, payload: {text: "Connection error"}});
    }

    this.wsRef.onclose = () => {
      this.bus.publish({topic: Topics.CONNECTION_CLOSED, payload: {}});
    };

    this.wsRef.onmessage = (evt: MessageEvent) => {
      const msg: WsMessage = JSON.parse(evt.data);

      switch (msg.payload.type) {
        case DebuggerMessageType.RequestReceived:
          const request = msg.payload.request;
          this.bus.publish({topic: Topics.REQUEST_RECEIVED, payload: {request}});
          break;
        case DebuggerMessageType.ImposterMatched:
          const imposter = msg.payload.imposter;
          this.bus.publish({topic: Topics.IMPOSTER_MATCHED, payload: {imposter}});
          break;
        case DebuggerMessageType.ResponsePrepared:
          const response = msg.payload.response;
          this.bus.publish({topic: Topics.RESPONSE_PREPARED, payload: {response}});
          break;
      }
    }
  }

  onDebuggerStopped = () => {
    this.wsRef?.close(1000)
    this.wsRef = undefined;
    this.bus.publish({topic: Topics.NOTIFICATION, payload: {text: "Connection closed"}});
  }

  onRequestReceived = () => {
    this.state = DebuggerState.RequestReceived;
  }

  onRequestConfirmed = (m: RequestConfirmedMessage) => {
    this.wsRef?.send(JSON.stringify({
      type: WsMessageType.Debugger,
      payload: {
        type: DebuggerMessageType.RequestConfirmed,
        request: m.request
      },
    }));
  }

  onImposterMatched = () => {
    this.state = DebuggerState.ImposterMatched;
  }

  onImposterConfirmed = (m: ImposterConfirmedMessage) => {
    this.wsRef?.send(JSON.stringify({
      type: WsMessageType.Debugger,
      payload: {
        type: DebuggerMessageType.ImposterConfirmed,
        imposter: m.imposter
      },
    }));
  }

  onResponsePrepared = () => {
    this.state = DebuggerState.ResponsePrepared;
  }

  onResponseConfirmed = (m: ResponseConfirmedMessage) => {
    this.wsRef?.send(JSON.stringify({
      type: WsMessageType.Debugger,
      payload: {
        type: DebuggerMessageType.ResponseConfirmed,
        response: m.response
      },
    }));
  }
}

export const DebuggerContext = React.createContext(new Debugger());
