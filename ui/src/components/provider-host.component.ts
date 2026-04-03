// SPDX-Licence-Identifier: EUPL-1.2

import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  ElementRef,
  DestroyRef,
  Input,
  OnChanges,
  OnInit,
  inject,
  Renderer2,
  ViewChild,
} from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ApiConfigService } from '../services/api-config.service';
import { ProviderDiscoveryService } from '../services/provider-discovery.service';

/**
 * ProviderHostComponent renders any custom element by tag name using
 * Angular's Renderer2 for safe DOM manipulation. It reads the :provider
 * route parameter to look up the element tag from the discovery service.
 */
@Component({
  selector: 'provider-host',
  standalone: true,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  template: '<div #container class="provider-host"></div>',
  styles: [
    `
      :host {
        display: block;
        width: 100%;
        height: 100%;
      }
      .provider-host {
        width: 100%;
        height: 100%;
      }

      .provider-host-empty {
        display: grid;
        place-items: center;
        min-height: 100%;
        padding: 1.5rem;
        color: rgb(156 163 175);
        background: rgba(255, 255, 255, 0.02);
        text-align: center;
      }
    `,
  ],
})
export class ProviderHostComponent implements OnInit, OnChanges {
  /** The custom element tag to render. Can be set via input or route param. */
  @Input() tag = '';

  /** API URL attribute passed to the custom element. */
  @Input() apiUrl = '';

  @ViewChild('container', { static: true }) container!: ElementRef;
  private readonly destroyRef = inject(DestroyRef);

  constructor(
    private renderer: Renderer2,
    private route: ActivatedRoute,
    private apiConfig: ApiConfigService,
    private providerService: ProviderDiscoveryService,
  ) {}

  ngOnInit(): void {
    this.route.params.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((params) => {
      const providerName = this.normalizeProviderName(params['provider']);
      if (providerName) {
        const provider = this.providerService
          .providers()
          .find((p) => p.name.toLowerCase() === providerName.toLowerCase());
        if (provider?.element?.tag) {
          this.tag = provider.element.tag;
          this.renderElement();
          return;
        }
      }

      this.tag = '';
      this.renderEmptyState();
    });
  }

  ngOnChanges(): void {
    this.renderElement();
  }

  private renderElement(): void {
    if (!this.container) {
      return;
    }

    if (!this.tag) {
      this.renderEmptyState();
      return;
    }

    const native = this.container.nativeElement;

    // Clear previous element safely
    while (native.firstChild) {
      this.renderer.removeChild(native, native.firstChild);
    }

    // Create and append the custom element
    const el = this.renderer.createElement(this.tag);
    const url = this.apiUrl || this.apiConfig.effectiveBaseUrl;
    if (url) {
      this.renderer.setAttribute(el, 'api-url', url);
    }
    this.renderer.appendChild(native, el);
  }

  private renderEmptyState(): void {
    if (!this.container) {
      return;
    }

    const native = this.container.nativeElement;
    while (native.firstChild) {
      this.renderer.removeChild(native, native.firstChild);
    }

    const empty = this.renderer.createElement('div');
    this.renderer.addClass(empty, 'provider-host-empty');
    this.renderer.appendChild(
      empty,
      this.renderer.createText('Select a renderable provider to preview its custom element.'),
    );
    this.renderer.appendChild(native, empty);
  }

  private normalizeProviderName(value: unknown): string {
    if (typeof value !== 'string') {
      return '';
    }

    try {
      return decodeURIComponent(value).trim();
    } catch {
      return value.trim();
    }
  }
}
