export {};

declare global {
  interface Window {
    __coreGUIBridge?: {
      call?: (action: string, payload?: unknown) => Promise<unknown> | unknown;
      generate?: (prompt: string) => Promise<string> | string;
      invoke?: (action: string, payload?: unknown) => Promise<unknown> | unknown;
      on?: (channel: string, handler: (...args: unknown[]) => void) => (() => void) | void;
      syncStorage?: (origin: string, state: unknown) => Promise<void> | void;
    };
    core?: {
      ml?: {
        generate?: (prompt: string) => Promise<string>;
      };
      storage?: {
        exportState?: () => unknown;
        origin?: string;
        reset?: () => unknown;
        sync?: () => void;
      };
    };
    electron?: unknown;
    require?: (name: string) => unknown;
  }
}
