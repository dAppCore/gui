import { CommonModule } from '@angular/common';
import { Component, DestroyRef, computed, effect, inject, signal } from '@angular/core';
import { ApiConfigService } from '../services/api-config.service';
import { ProviderDiscoveryService, type ProviderInfo } from '../services/provider-discovery.service';
import { TranslationService } from '../services/translation.service';
import { WebSocketService } from '../services/websocket.service';
import { ProviderHostComponent } from '../components/provider-host.component';

@Component({
  selector: 'dashboard-view',
  imports: [CommonModule, ProviderHostComponent],
  template: `
    <main class="display-shell">
      <section class="hero">
        <div class="hero-copy">
          <p class="eyebrow">Core GUI</p>
          <h1>{{ title() }}</h1>
          <p class="subtitle">{{ subtitle() }}</p>
          <p class="body">
            A compact operator surface for desktop workflows, provider discovery, and realtime
            backend status.
          </p>

          <div class="hero-actions">
            <button type="button" class="primary-action" (click)="refreshProviders()">
              Refresh providers
            </button>
            <a class="secondary-action" [href]="apiBase()" target="_blank" rel="noreferrer">
              Open API endpoint
            </a>
          </div>
        </div>

        <div class="hero-meta">
          <div class="meta-card">
            <span class="meta-label">Connection</span>
            <strong [class.good]="connected()">{{ connected() ? 'Live' : 'Reconnecting' }}</strong>
          </div>
          <div class="meta-card">
            <span class="meta-label">Providers</span>
            <strong>{{ providerCount() }}</strong>
          </div>
          <div class="meta-card">
            <span class="meta-label">Local time</span>
            <strong>{{ clock() | date: 'mediumTime' }}</strong>
          </div>
          <div class="meta-card">
            <span class="meta-label">API base</span>
            <strong class="mono">{{ apiBase() }}</strong>
          </div>
        </div>
      </section>

      <section class="content-grid">
        <article class="feature-panel">
          <div class="panel-heading">
            <div>
              <p class="eyebrow">Discovered providers</p>
              <h2>Renderable capabilities</h2>
            </div>
            <span class="pill">{{ providerCount() }} total</span>
          </div>

          <div class="provider-list" *ngIf="featuredProviders().length; else emptyState">
            <button
              type="button"
              class="provider-row"
              *ngFor="let provider of featuredProviders(); trackBy: trackByProvider"
              [class.selected]="selectedProviderName() === provider.name"
              (click)="selectProvider(provider)"
            >
              <div class="provider-icon">
                <span>{{ provider.name.slice(0, 1).toUpperCase() }}</span>
              </div>
              <div class="provider-copy">
                <strong>{{ provider.name }}</strong>
                <span>{{ provider.basePath }}</span>
                <small *ngIf="provider.element?.tag" class="mono">
                  {{ provider.element?.tag }} · {{ provider.element?.source }}
                </small>
              </div>
            </button>
          </div>

          <ng-template #emptyState>
            <div class="empty-state">
              <strong>No providers discovered yet.</strong>
              <span>The shell will populate this view once the backend exposes provider metadata.</span>
            </div>
          </ng-template>
        </article>

        <article class="feature-panel accent">
          <div class="panel-heading">
            <div>
              <p class="eyebrow">Live wiring</p>
              <h2>What this shell keeps online</h2>
            </div>
          </div>

          <ul class="feature-list">
            <li>
              <strong>Provider discovery</strong>
              <span>Loads provider metadata and registers custom element scripts automatically.</span>
            </li>
            <li>
              <strong>Realtime status</strong>
              <span>Tracks the websocket connection used for backend events.</span>
            </li>
            <li>
              <strong>Desktop bridge</strong>
              <span>Renders in the Wails webview and stays responsive to the local runtime.</span>
            </li>
          </ul>
        </article>

        <article class="feature-panel preview-panel">
          <div class="panel-heading">
            <div>
              <p class="eyebrow">Provider preview</p>
              <h2>{{ selectedProviderTitle() }}</h2>
            </div>
            <span class="pill" [class.good]="hasRenderableSelection()">
              {{ hasRenderableSelection() ? 'Renderable' : 'Select one' }}
            </span>
          </div>

          @if (selectedRenderableProvider(); as provider) {
            <div class="preview-meta">
              <span class="preview-field">
                <label>Base path</label>
                <strong>{{ provider.basePath }}</strong>
              </span>
              <span class="preview-field">
                <label>Tag</label>
                <strong class="mono">{{ provider.element?.tag }}</strong>
              </span>
            </div>
            <div class="preview-host">
              <provider-host [tag]="provider.element?.tag || ''" [apiUrl]="apiBase()"></provider-host>
            </div>
          } @else {
            <div class="empty-state">
              <strong>No renderable provider selected.</strong>
              <span>Pick a provider with a custom element to load its live preview here.</span>
            </div>
          }
        </article>
      </section>
    </main>
  `,
})
export class DashboardComponent {
  private readonly discovery = inject(ProviderDiscoveryService);
  private readonly apiConfig = inject(ApiConfigService);
  private readonly translations = inject(TranslationService);
  private readonly websocket = inject(WebSocketService);
  private readonly destroyRef = inject(DestroyRef);

  protected readonly title = signal('Core GUI');
  protected readonly subtitle = signal('Desktop orchestration console');
  protected readonly clock = signal(new Date());
  protected readonly selectedProviderName = signal('');

  protected readonly providers = this.discovery.providers;
  protected readonly providerCount = computed(() => this.providers().length);
  protected readonly connected = this.websocket.connected;
  protected readonly apiBase = computed(() => this.apiConfig.effectiveBaseUrl);

  protected readonly featuredProviders = computed<ProviderInfo[]>(() =>
    this.providers().filter((provider) => provider.element?.tag).slice(0, 6),
  );

  protected readonly selectedRenderableProvider = computed<ProviderInfo | null>(() => {
    const selection = this.selectedProviderName();
    if (!selection) {
      return this.featuredProviders()[0] ?? null;
    }

    return (
      this.providers().find((provider) => provider.name === selection && provider.element?.tag) ??
      this.featuredProviders()[0] ??
      null
    );
  });

  protected readonly selectedProviderTitle = computed(() => {
    const provider = this.selectedRenderableProvider();
    return provider?.name ?? 'Preview';
  });

  constructor() {
    const tick = setInterval(() => this.clock.set(new Date()), 1000);
    this.destroyRef.onDestroy(() => clearInterval(tick));

    effect(() => {
      if (this.connected()) {
        document.documentElement.setAttribute('data-connected', 'true');
      } else {
        document.documentElement.removeAttribute('data-connected');
      }
    });
  }

  async ngOnInit(): Promise<void> {
    await this.translations.onReady();
    await this.discovery.discover();
    this.websocket.connect();
  }

  ngOnDestroy(): void {
    this.websocket.disconnect();
  }

  async refreshProviders(): Promise<void> {
    await this.discovery.refresh();
    if (!this.selectedRenderableProvider()) {
      this.selectedProviderName.set('');
    }
  }

  selectProvider(provider: ProviderInfo): void {
    if (provider.element?.tag) {
      this.selectedProviderName.set(provider.name);
    }
  }

  hasRenderableSelection(): boolean {
    return !!this.selectedRenderableProvider();
  }

  trackByProvider(_: number, provider: ProviderInfo): string {
    return provider.name;
  }
}
