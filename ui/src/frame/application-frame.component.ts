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
        position: relative;
      }

      .frame-main {
        min-height: calc(100vh - 6.5rem);
        position: relative;
        z-index: 0;
      }

      .application-frame .frame-header {
        backdrop-filter: blur(18px);
        background: linear-gradient(180deg, rgba(8, 12, 22, 0.94), rgba(8, 12, 22, 0.82));
        border-bottom-color: rgba(255, 255, 255, 0.06);
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.18);
      }

      .application-frame .frame-nav {
        position: relative;
        z-index: 30;
      }

      .application-frame .frame-nav .lg\\:fixed {
        background: linear-gradient(180deg, rgba(5, 9, 18, 0.96), rgba(7, 12, 22, 0.84));
        backdrop-filter: blur(18px);
        border-right: 1px solid rgba(255, 255, 255, 0.06);
        box-shadow: 12px 0 40px rgba(0, 0, 0, 0.16);
      }

      .application-frame .frame-nav a {
        transition:
          transform 140ms ease,
          background 140ms ease,
          color 140ms ease;
      }

      .application-frame .frame-nav a:hover {
        transform: translateX(1px);
      }

      .application-frame .frame-main {
        background:
          radial-gradient(circle at 0% 0%, rgba(20, 184, 166, 0.08), transparent 22%),
          linear-gradient(180deg, rgba(255, 255, 255, 0.01), transparent 24%);
      }

      .application-frame .frame-main .px-0 {
        padding-left: clamp(1rem, 2vw, 1.5rem);
        padding-right: clamp(1rem, 2vw, 1.5rem);
      }

      .application-frame .frame-main router-outlet {
        display: block;
      }

      .application-frame .frame-header input {
        color: var(--text);
        caret-color: var(--accent-strong);
      }

      .application-frame .frame-header button,
      .application-frame .frame-header a {
        transition:
          transform 140ms ease,
          color 140ms ease,
          background 140ms ease;
      }

      .application-frame .frame-header button:hover,
      .application-frame .frame-header a:hover {
        transform: translateY(-1px);
      }

      .application-frame .frame-header .fa-bell,
      .application-frame .frame-header .fa-bars {
        color: var(--accent-strong);
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
