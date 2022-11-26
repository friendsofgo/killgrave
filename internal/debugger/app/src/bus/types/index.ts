export { EventBus } from './EventBus'
export { Topics } from './topics'

export interface Listener {
  unsubscribe(): void;
}