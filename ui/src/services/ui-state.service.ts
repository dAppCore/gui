// SPDX-Licence-Identifier: EUPL-1.2

import { Injectable, computed, signal } from '@angular/core';

/**
 * UiStateService stores shell-wide client state shared across the frame and
 * dashboard. Keeping it in one place avoids threading search state through
 * unrelated routes.
 */
@Injectable({ providedIn: 'root' })
export class UiStateService {
  private readonly _searchQuery = signal('');

  readonly searchQuery = this._searchQuery.asReadonly();
  readonly hasSearch = computed(() => this.searchQuery().length > 0);

  setSearchQuery(value: string): void {
    this._searchQuery.set(value.trim());
  }

  clearSearchQuery(): void {
    this._searchQuery.set('');
  }
}
