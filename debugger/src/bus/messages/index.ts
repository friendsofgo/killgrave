export interface ThemeChangedMessage {
  dark: boolean;
}

export interface FontSizeChangedMessage {
  increased: boolean;
}

export interface NotificationMessage {
  text: string;
}

export interface EventBusMessages {
  ThemeChanged: ThemeChangedMessage;
  FontSizeChanged: FontSizeChangedMessage;
  Notification: NotificationMessage;
}
