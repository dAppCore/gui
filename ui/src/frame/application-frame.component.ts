// SPDX-Licence-Identifier: EUPL-1.2

import { Component, CUSTOM_ELEMENTS_SCHEMA, Input, OnInit, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { TranslationService } from '../services/translation.service';
import { ProviderDiscoveryService } from '../services/provider-discovery.service';
import { WebSocketService } from '../services/websocket.service';
import { StatusBarComponent } from '../components/status-bar.component';

interface NavItem {
  name: string;
  href: string;
  icon: string;
}

/**
 * ApplicationFrameComponent is the HLCRF (Header, Left nav, Content, Right, Footer)
 * shell for all Core Wails applications. It provides:
 *
 * - Dynamic sidebar navigation populated from ProviderDiscoveryService
 * - Content area rendered via router-outlet for child routes
 * - Footer status bar with time, version, and provider status
 * - Mobile-responsive sidebar with expand/collapse
 * - Dark mode support
 *
 * Ported from core-gui/cmd/lthn-desktop/frontend/src/frame/application.frame.ts
 * with navigation made dynamic via provider discovery.
 */
@Component({
  selector: 'application-frame',
  standalone: true,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive, StatusBarComponent],
  templateUrl: './application-frame.component.html',
  styles: [
    `
      .application-frame {
        min-height: 100vh;
      }

      .frame-main {
        min-height: calc(100vh - 6.5rem);
      }

    `,
  ],
})
export class ApplicationFrameComponent implements OnInit {
  @Input() version = 'v0.1.0';

  sidebarOpen = false;
  userMenuOpen = false;

  /** Static navigation items set by the host application. */
  @Input() staticNavigation: NavItem[] = [];

  /** Combined navigation: static + dynamic from providers. */
  readonly navigation = computed<NavItem[]>(() => {
    const dynamicItems = this.providerService
      .providers()
      .filter((p) => p.element)
      .map((p) => ({
        name: p.name,
        href: `/provider/${encodeURIComponent(p.name.toLowerCase())}`,
        icon: 'fa-regular fa-puzzle-piece fa-2xl shrink-0',
      }));

    return [...this.staticNavigation, ...dynamicItems];
  });

  userNavigation: NavItem[] = [];

  constructor(
    public t: TranslationService,
    private providerService: ProviderDiscoveryService,
    private wsService: WebSocketService,
  ) {}

  async ngOnInit(): Promise<void> {
    await this.t.onReady();
    this.initUserNavigation();

    // Discover providers and build navigation
    await this.providerService.discover();

    // Connect WebSocket for real-time updates
    this.wsService.connect();
  }

  private initUserNavigation(): void {
    this.userNavigation = [
      {
        name: this.t._('menu.settings'),
        href: '/settings',
        icon: 'fa-regular fa-gear',
      },
    ];
  }

}
