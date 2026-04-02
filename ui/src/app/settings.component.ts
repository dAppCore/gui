import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ApiConfigService } from '../services/api-config.service';
import { ProviderDiscoveryService } from '../services/provider-discovery.service';
import { WebSocketService } from '../services/websocket.service';

@Component({
  selector: 'settings-view',
  imports: [CommonModule, FormsModule],
  template: `
    <main class="display-shell">
      <section class="hero settings-hero">
        <div class="hero-copy">
          <p class="eyebrow">Settings</p>
          <h1>Connection surface</h1>
          <p class="subtitle">Adjust the API endpoint the shell uses for discovery and previews.</p>
          <p class="body">
            The frontend can target the embedded Wails origin or a remote Core API during
            development. Changes apply immediately to future discovery and provider-host requests.
          </p>
        </div>

        <div class="hero-meta">
          <div class="meta-card">
            <span class="meta-label">Providers</span>
            <strong>{{ providerCount() }}</strong>
          </div>
          <div class="meta-card">
            <span class="meta-label">Connection</span>
            <strong [class.good]="connected()">{{ connected() ? 'Live' : 'Reconnecting' }}</strong>
          </div>
        </div>
      </section>

      <section class="content-grid single-column">
        <article class="feature-panel">
          <div class="panel-heading">
            <div>
              <p class="eyebrow">API</p>
              <h2>Base URL</h2>
            </div>
          </div>

          <div class="settings-form">
            <label class="settings-field">
              <span>API base URL</span>
              <input
                type="url"
                [(ngModel)]="draftBaseUrl"
                placeholder="http://127.0.0.1:8080"
                autocomplete="off"
                spellcheck="false"
              />
            </label>

            <div class="settings-actions">
              <button type="button" class="primary-action" (click)="applyBaseUrl()">
                Apply
              </button>
              <button type="button" class="secondary-action" (click)="resetBaseUrl()">
                Use local origin
              </button>
            </div>
          </div>
        </article>
      </section>
    </main>
  `,
})
export class SettingsComponent {
  private readonly apiConfig = inject(ApiConfigService);
  private readonly discovery = inject(ProviderDiscoveryService);
  private readonly websocket = inject(WebSocketService);

  draftBaseUrl = this.apiConfig.baseUrl;

  readonly providerCount = () => this.discovery.providers().length;
  readonly connected = () => this.websocket.connected();

  applyBaseUrl(): void {
    this.apiConfig.baseUrl = this.draftBaseUrl.trim();
    this.discovery.refresh();
    this.websocket.disconnect();
    this.websocket.connect();
  }

  resetBaseUrl(): void {
    this.draftBaseUrl = '';
    this.applyBaseUrl();
  }
}
