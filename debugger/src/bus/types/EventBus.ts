import { isDevelopment } from '../../App';

const defaultPublishOptions: PublishOptions = {
  targetOrigin: '*',
  targetWindow: window
};

const defaultSubscribeOptions: SubscribeOptions = {
  targetWindow: window
};

export const EventBus = <MessagesMap extends Record<string, any>>() => {
  const publish = <Topic extends keyof MessagesMap>(
    message: Message<MessagesMap, Topic>,
    options: PublishOptions = defaultPublishOptions
  ) => {
    if (isDevelopment) console.log(message);
    options.targetWindow.postMessage(message, options.targetOrigin);
  };

  const subscribe = <Topic extends keyof MessagesMap>(
    topic: Topic,
    handler: Handler<MessagesMap[Topic]>,
    options: SubscribeOptions = defaultSubscribeOptions
  ) => {
    const messageEventHandler = (
      event: MessageEvent<Message<MessagesMap, Topic>>
    ) => subscriptionHandler(event, topic, handler);

    const attachEventListener = () => {
      options.targetWindow.addEventListener(
        'message',
        messageEventHandler
      );
    };
    const detachEventListener = () => {
      options.targetWindow.removeEventListener(
        'message',
        messageEventHandler
      );
    };

    attachEventListener();
    return {unsubscribe: detachEventListener};
  };

  const subscriptionHandler = <Topic extends keyof MessagesMap>(
    event: MessageEvent<Message<MessagesMap, Topic>>,
    topic: Topic,
    handler: Handler<MessagesMap[Topic]>
  ) => {
    if (event.data.topic === topic) {
      handler(event.data.payload);
    }
  };

  return {publish, subscribe};
};

export interface Message<M extends Record<string, any>, T extends keyof M> {
  topic: T;
  payload: M[T];
}

export interface PublishOptions {
  targetOrigin: string;
  targetWindow: Window;
}

export interface SubscribeOptions {
  targetWindow: Window;
}

export type Handler<Payload extends any> = (payload: Payload) => void;