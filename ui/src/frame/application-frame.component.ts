import {
  Component,
  HostListener,
  computed,
  inject,
  signal,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { StatusBarComponent } from '../components/status-bar.component';
import { TranslationService } from '../services/translation.service';
import { UiStateService } from '../services/ui-state.service';

interface NavItem {
  name: string;
  href: string;
  icon: string;
}

@Component({
  selector: 'application-frame',
  standalone: true,
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
      }

      .application-frame .frame-header {
        backdrop-filter: blur(18px);
        background: linear-gradient(180deg, rgba(8, 12, 22, 0.94), rgba(8, 12, 22, 0.82));
        border-bottom-color: rgba(255, 255, 255, 0.06);
      }

      .application-frame .frame-nav .lg\\:fixed {
        background: linear-gradient(180deg, rgba(7, 12, 22, 0.98), rgba(9, 15, 28, 0.92));
        border-right: 1px solid rgba(255, 255, 255, 0.08);
      }

      .search-shell {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        min-height: 2.75rem;
        padding: 0 0.875rem;
        border-radius: 999px;
        border: 1px solid rgba(255, 255, 255, 0.08);
        background: rgba(255, 255, 255, 0.04);
      }

      .search-shell input {
        min-width: 0;
        flex: 1;
        border: 0;
        background: transparent;
        outline: none;
        color: #f8fafc;
      }

      .search-clear,
      .search-count {
        color: rgba(226, 232, 240, 0.72);
      }

      .search-clear {
        border: 0;
        background: transparent;
        cursor: pointer;
      }

      .search-count {
        white-space: nowrap;
      }
    `,
  ],
})
export class ApplicationFrameComponent {
  readonly sidebarOpen = signal(false);
  readonly userMenuOpen = signal(false);
  protected readonly version = 'v0.1.0';

  private readonly uiState = inject(UiStateService);
  protected readonly t = inject(TranslationService);

  private readonly navigationItems: NavItem[] = [
    { name: 'Chat', href: '/', icon: 'fa-regular fa-comments fa-2xl shrink-0' },
    { name: 'Settings', href: '/settings', icon: 'fa-regular fa-sliders fa-2xl shrink-0' },
  ];

  readonly visibleNavigation = computed(() => {
    const query = this.uiState.searchQuery().toLowerCase();
    if (!query) {
      return this.navigationItems;
    }
    return this.navigationItems.filter((item) =>
      `${item.name} ${item.href}`.toLowerCase().includes(query),
    );
  });

  readonly searchQuery = this.uiState.searchQuery;

  readonly userNavigation: NavItem[] = [
    { name: 'Chat', href: '/', icon: 'fa-regular fa-comments' },
    { name: 'Settings', href: '/settings', icon: 'fa-regular fa-sliders' },
  ];

  @HostListener('document:keydown.escape')
  onEscape(): void {
    if (this.sidebarOpen()) {
      this.sidebarOpen.set(false);
      return;
    }
    if (this.userMenuOpen()) {
      this.userMenuOpen.set(false);
      return;
    }
    this.uiState.clearSearchQuery();
  }

  protected onSearchInput(value: string): void {
    this.uiState.setSearchQuery(value);
  }

  protected clearSearch(): void {
    this.uiState.clearSearchQuery();
  }
}
