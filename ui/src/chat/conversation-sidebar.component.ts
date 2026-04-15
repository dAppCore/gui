import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ConversationSummary } from './chat.types';

@Component({
  selector: 'chat-conversation-sidebar',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <aside class="sidebar">
      <div class="sidebar__head">
        <button type="button" class="ghost" (click)="create.emit()">New chat</button>
        <input [(ngModel)]="query" (ngModelChange)="queryChange.emit($event)" placeholder="Search history" />
      </div>
      <div class="sidebar__list">
        <button
          *ngFor="let item of conversations"
          type="button"
          class="conversation"
          [class.conversation--active]="item.id === activeId"
          (click)="select.emit(item.id)"
        >
          <span class="conversation__title">{{ item.title }}</span>
          <span class="conversation__meta">{{ item.model }} · {{ item.message_count }} msgs</span>
        </button>
      </div>
    </aside>
  `,
  styles: [
    `
      .sidebar { display: grid; gap: 1rem; padding: 1rem; background: rgba(9, 20, 34, 0.72); border-right: 1px solid rgba(124, 156, 191, 0.16); }
      .sidebar__head { display: grid; gap: 0.75rem; }
      .sidebar__list { display: grid; gap: 0.5rem; align-content: start; overflow: auto; }
      .ghost, input { width: 100%; border-radius: 0.9rem; border: 1px solid rgba(124, 156, 191, 0.2); background: rgba(11, 27, 44, 0.7); color: #eaf4ff; padding: 0.8rem 0.9rem; }
      .ghost { cursor: pointer; text-align: left; font-weight: 700; }
      .conversation { display: grid; gap: 0.2rem; padding: 0.9rem; border: 0; border-radius: 1rem; background: rgba(8, 21, 35, 0.55); color: #dbeafe; text-align: left; cursor: pointer; }
      .conversation--active { background: linear-gradient(135deg, rgba(14, 116, 144, 0.55), rgba(8, 47, 73, 0.82)); box-shadow: inset 0 0 0 1px rgba(125, 211, 252, 0.28); }
      .conversation__title { font-weight: 700; }
      .conversation__meta { color: #94a3b8; font-size: 0.8rem; }
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
}
