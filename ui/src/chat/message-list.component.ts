import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { ChatMessage } from './chat.types';

@Component({
  selector: 'chat-message-list',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="thread">
      <article *ngFor="let message of messages" class="bubble" [class.bubble--user]="message.role === 'user'">
        <header>
          <strong>{{ message.role }}</strong>
          <span>{{ message.model || 'local' }}</span>
        </header>
        <p>{{ message.content }}</p>
        <section *ngIf="message.thinking?.content" class="thinking">
          <strong>Thinking</strong>
          <p>{{ message.thinking?.content }}</p>
        </section>
        <section *ngIf="message.tool_calls?.length" class="tool">
          <strong>Tool calls</strong>
          <pre>{{ message.tool_calls | json }}</pre>
        </section>
        <section *ngIf="message.tool_results?.length" class="tool">
          <strong>Tool results</strong>
          <pre>{{ message.tool_results | json }}</pre>
        </section>
      </article>
    </div>
  `,
  styles: [
    `
      .thread { display: grid; gap: 1rem; align-content: start; }
      .bubble { max-width: 54rem; padding: 1rem 1.1rem; border-radius: 1.25rem; background: rgba(11, 27, 44, 0.88); border: 1px solid rgba(125, 211, 252, 0.12); box-shadow: 0 12px 40px rgba(0, 0, 0, 0.18); }
      .bubble--user { margin-left: auto; background: linear-gradient(135deg, rgba(8, 47, 73, 0.95), rgba(14, 116, 144, 0.7)); }
      header { display: flex; justify-content: space-between; gap: 1rem; margin-bottom: 0.65rem; color: #7dd3fc; text-transform: uppercase; letter-spacing: 0.08em; font-size: 0.72rem; }
      p { margin: 0; white-space: pre-wrap; line-height: 1.6; color: #ecfeff; }
      .thinking, .tool { margin-top: 0.85rem; padding-top: 0.85rem; border-top: 1px solid rgba(125, 211, 252, 0.12); color: #cbd5e1; }
      pre { margin: 0.35rem 0 0; overflow: auto; background: rgba(2, 6, 23, 0.6); padding: 0.85rem; border-radius: 0.8rem; }
    `,
  ],
})
export class MessageListComponent {
  @Input() messages: ChatMessage[] = [];
}
