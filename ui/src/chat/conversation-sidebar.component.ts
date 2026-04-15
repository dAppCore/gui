import { CommonModule, DatePipe } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ConversationSummary } from './chat.types';

type ConversationGroup = { label: string; items: ConversationSummary[] };

@Component({
  selector: 'chat-conversation-sidebar',
  standalone: true,
  imports: [CommonModule, FormsModule, DatePipe],
  template: `
    <aside class="sidebar">
      <div class="sidebar__head">
        <button type="button" class="ghost" (click)="create.emit()">New chat</button>
        <input [ngModel]="query" (ngModelChange)="queryChange.emit($event)" placeholder="Search history" />
      </div>

      <div class="sidebar__list">
        <section *ngFor="let group of groupedConversations" class="group">
          <p class="group__label">{{ group.label }}</p>

          <article
            *ngFor="let item of group.items"
            class="conversation"
            [class.conversation--active]="item.id === activeId"
          >
            <button type="button" class="conversation__select" (click)="select.emit(item.id)">
              <span class="conversation__title">{{ item.title }}</span>
              <span class="conversation__meta">{{ item.model }} · {{ item.updated_at | date: 'MMM d, HH:mm' }}</span>
            </button>

            <div class="conversation__actions">
              <button type="button" (click)="beginRename(item)">Rename</button>
              <button type="button" (click)="export.emit(item.id)">Export</button>
              <button type="button" class="danger" (click)="removeConversation(item)">Delete</button>
            </div>

            <form *ngIf="editingId === item.id" class="conversation__rename" (ngSubmit)="commitRename(item.id)">
              <input
                [ngModel]="draftTitle"
                (ngModelChange)="draftTitle = $event"
                name="rename-{{ item.id }}"
                autocomplete="off"
              />
              <div class="conversation__rename-actions">
                <button type="submit">Save</button>
                <button type="button" class="ghost" (click)="cancelRename()">Cancel</button>
              </div>
            </form>
          </article>
        </section>
      </div>
    </aside>
  `,
  styles: [
    `
      .sidebar { display: grid; gap: 1rem; padding: 1rem; background: rgba(9, 20, 34, 0.72); border-right: 1px solid rgba(124, 156, 191, 0.16); }
      .sidebar__head { display: grid; gap: 0.75rem; }
      .sidebar__list { display: grid; gap: 1rem; align-content: start; overflow: auto; }
      .group { display: grid; gap: 0.5rem; }
      .group__label { margin: 0; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.14em; font-size: 0.7rem; }
      .ghost, input, .conversation__actions button, .conversation__rename-actions button {
        border-radius: 0.9rem;
        border: 1px solid rgba(124, 156, 191, 0.2);
        background: rgba(11, 27, 44, 0.7);
        color: #eaf4ff;
        padding: 0.8rem 0.9rem;
      }
      .ghost { cursor: pointer; text-align: left; font-weight: 700; }
      .conversation { display: grid; gap: 0.7rem; padding: 0.9rem; border: 0; border-radius: 1rem; background: rgba(8, 21, 35, 0.55); color: #dbeafe; text-align: left; }
      .conversation--active { background: linear-gradient(135deg, rgba(14, 116, 144, 0.55), rgba(8, 47, 73, 0.82)); box-shadow: inset 0 0 0 1px rgba(125, 211, 252, 0.28); }
      .conversation__select { display: grid; gap: 0.35rem; border: 0; padding: 0; background: transparent; color: inherit; text-align: left; cursor: pointer; }
      .conversation__title { font-weight: 700; }
      .conversation__meta { color: #94a3b8; font-size: 0.8rem; }
      .conversation__actions { display: flex; gap: 0.45rem; flex-wrap: wrap; }
      .conversation__actions button, .conversation__rename-actions button { padding: 0.35rem 0.6rem; font-size: 0.72rem; cursor: pointer; }
      .conversation__rename { display: grid; gap: 0.5rem; }
      .conversation__rename-actions { display: flex; gap: 0.45rem; flex-wrap: wrap; }
      .danger { color: #fecaca; border-color: rgba(248, 113, 113, 0.28); }
    `,
  ],
})
export class ConversationSidebarComponent {
  @Input() conversations: ConversationSummary[] = [];
  @Input() activeId = '';
  @Input() query = '';
  @Output() queryChange = new EventEmitter<string>();
  @Output() select = new EventEmitter<string>();
  @Output() create = new EventEmitter<void>();
  @Output() rename = new EventEmitter<{ id: string; title: string }>();
  @Output() delete = new EventEmitter<string>();
  @Output() export = new EventEmitter<string>();

  editingId: string | null = null;
  draftTitle = '';

  get groupedConversations(): ConversationGroup[] {
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime();
    const yesterday = today - 24 * 60 * 60 * 1000;
    const lastWeek = today - 7 * 24 * 60 * 60 * 1000;
    const groups = new Map<string, ConversationSummary[]>();

    for (const conversation of this.conversations) {
      const updatedAt = Date.parse(conversation.updated_at);
      let label = 'Older';
      if (updatedAt >= today) {
        label = 'Today';
      } else if (updatedAt >= yesterday) {
        label = 'Yesterday';
      } else if (updatedAt >= lastWeek) {
        label = 'Previous 7 Days';
      }
      groups.set(label, [...(groups.get(label) ?? []), conversation]);
    }

    return ['Today', 'Yesterday', 'Previous 7 Days', 'Older']
      .filter((label) => groups.has(label))
      .map((label) => ({ label, items: groups.get(label) ?? [] }));
  }

  beginRename(item: ConversationSummary): void {
    this.editingId = item.id;
    this.draftTitle = item.title;
  }

  commitRename(id: string): void {
    const title = this.draftTitle.trim();
    if (title) {
      this.rename.emit({ id, title });
    }
    this.cancelRename();
  }

  cancelRename(): void {
    this.editingId = null;
    this.draftTitle = '';
  }

  removeConversation(item: ConversationSummary): void {
    if (window.confirm(`Delete "${item.title}"?`)) {
      this.delete.emit(item.id);
    }
  }
}
