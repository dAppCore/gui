import { CommonModule } from '@angular/common';
import { Component, DestroyRef, OnDestroy, computed, effect, inject, signal } from '@angular/core';
import { ApiConfigService } from '../services/api-config.service';
import { ProviderDiscoveryService, type ProviderInfo } from '../services/provider-discovery.service';
import { TranslationService } from '../services/translation.service';
import { WebSocketService } from '../services/websocket.service';

@Component({
  selector: 'core-display',
  imports: [CommonModule],
  templateUrl: './app.html',
  standalone: true,
})
export class App implements OnDestroy {
  private readonly discovery = inject(ProviderDiscoveryService);
  private readonly apiConfig = inject(ApiConfigService);
  private readonly translations = inject(TranslationService);
  private readonly websocket = inject(WebSocketService);
  private readonly destroyRef = inject(DestroyRef);

  protected readonly title = signal('Core GUI');
  protected readonly subtitle = signal('Desktop orchestration console');
  protected readonly clock = signal(new Date());

  protected readonly providers = this.discovery.providers;
  protected readonly providerCount = computed(() => this.providers().length);
  protected readonly connected = this.websocket.connected;
  protected readonly apiBase = computed(() => this.apiConfig.baseUrl || window.location.origin);

  protected readonly featuredProviders = computed<ProviderInfo[]>(() =>
    this.providers().filter((provider) => provider.element?.tag).slice(0, 6),
  );

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

  refreshProviders(): Promise<void> {
    return this.discovery.refresh();
  }

  trackByProvider(_: number, provider: ProviderInfo): string {
    return provider.name;
  }
}
