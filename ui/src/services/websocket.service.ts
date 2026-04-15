// SPDX-Licence-Identifier: EUPL-1.2

import { Injectable, OnDestroy, signal } from '@angular/core';
import { ApiConfigService } from './api-config.service';

export interface WSMessage {
  channel: string;
  data: unknown;
}

/**
 * WebSocketService manages a persistent WebSocket connection with automatic
 * reconnection. Follows the same pattern used by Mining's websocket.service.ts.
 */
@Injectable({ providedIn: 'root' })
export class WebSocketService implements OnDestroy {
  private ws: WebSocket | null = null;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private listeners = new Map<string, Set<(data: unknown) => void>>();
  private reconnectDelay = 1000;
  private maxReconnectDelay = 30000;
  private shouldReconnect = true;

  readonly connected = signal(false);

  constructor(private apiConfig: ApiConfigService) {}

  /** Open the WebSocket connection. */
  connect(path = '/ws'): void {
    if (this.ws) {
      return;
    }

    this.shouldReconnect = true;
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const base = this.apiConfig.baseUrl || window.location.origin;
    const wsBase = base.replace(/^http/, 'ws');
    const url = `${wsBase.length > 0 ? wsBase : `${protocol}//${window.location.host}`}${path}`;

    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      this.connected.set(true);
      this.reconnectDelay = 1000;
    };

    this.ws.onclose = () => {
      this.connected.set(false);
      this.ws = null;
      this.scheduleReconnect(path);
    };

    this.ws.onerror = () => {
      this.ws?.close();
    };

    this.ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data as string);
        this.dispatch(msg.channel, msg.data);
      } catch {
        // Silently ignore malformed messages
      }
    };
  }

  /** Subscribe to a channel. Returns an unsubscribe function. */
  on(channel: string, callback: (data: unknown) => void): () => void {
    if (!this.listeners.has(channel)) {
      this.listeners.set(channel, new Set());
    }
    this.listeners.get(channel)!.add(callback);

    return () => {
      const set = this.listeners.get(channel);
      if (set) {
        set.delete(callback);
        if (set.size === 0) {
          this.listeners.delete(channel);
        }
      }
    };
  }

  /** Send a message over the WebSocket. */
  send(channel: string, data: unknown): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ channel, data }));
    }
  }

  /** Disconnect and stop reconnecting. */
  disconnect(): void {
    this.shouldReconnect = false;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    this.ws?.close();
    this.ws = null;
  }

  ngOnDestroy(): void {
    this.disconnect();
  }

  private dispatch(channel: string, data: unknown): void {
    // Exact match
    this.listeners.get(channel)?.forEach((cb) => cb(data));
    // Wildcard match
    this.listeners.get('*')?.forEach((cb) => cb({ channel, data }));
  }

  private scheduleReconnect(path: string): void {
    if (!this.shouldReconnect) {
      return;
    }

    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.connect(path);
    }, this.reconnectDelay);

    this.reconnectDelay = Math.min(this.reconnectDelay * 2, this.maxReconnectDelay);
  }
}
