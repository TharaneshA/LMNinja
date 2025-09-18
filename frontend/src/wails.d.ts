interface Window {
  runtime: {
    EventsOn(eventName: string, callback: (...args: any[]) => void): void;
  };
}