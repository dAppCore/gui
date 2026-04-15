import { CommonModule, DatePipe } from '@angular/common';
import { Component, Input } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { ChatMessage, ImageAttachment } from './chat.types';

@Component({
  selector: 'chat-message-list',
  standalone: true,
  imports: [CommonModule, DatePipe],
  template: `
    <div class="thread">
      <article *ngFor="let message of messages" class="bubble" [class.bubble--user]="message.role === 'user'">
        <header>
          <strong>{{ message.role }}</strong>
          <span>{{ message.model || 'local' }}</span>
        </header>
        <section class="content" [innerHTML]="renderMarkdown(message.content)"></section>
        <section *ngIf="message.attachments?.length" class="attachments">
          <figure *ngFor="let attachment of message.attachments" class="attachment">
            <img [src]="attachmentSource(attachment)" [alt]="attachment.filename" />
            <figcaption>{{ attachment.filename }} · {{ attachment.width }}×{{ attachment.height }}</figcaption>
          </figure>
        </section>
        <details *ngIf="message.thinking?.content" class="thinking">
          <summary>Thinking · {{ thinkingDuration(message) }}</summary>
          <p>{{ message.thinking?.content }}</p>
        </details>
        <section *ngIf="message.tool_calls?.length" class="tool">
          <strong>Tool calls</strong>
          <div *ngFor="let call of message.tool_calls" class="tool__block">
            <div class="tool__title">{{ call.name }}</div>
            <pre>{{ call.arguments | json }}</pre>
          </div>
        </section>
        <section *ngIf="message.tool_results?.length" class="tool">
          <strong>Tool results</strong>
          <div *ngFor="let result of message.tool_results" class="tool__block">
            <div class="tool__title">{{ result.tool_call_id }}</div>
            <pre>{{ result.content }}</pre>
          </div>
        </section>
        <footer>{{ message.created_at | date: 'MMM d, HH:mm:ss' }}</footer>
      </article>
    </div>
  `,
  styles: [
    `
      .thread { display: grid; gap: 1rem; align-content: start; }
      .bubble { max-width: 54rem; padding: 1rem 1.1rem; border-radius: 1.25rem; background: rgba(11, 27, 44, 0.88); border: 1px solid rgba(125, 211, 252, 0.12); box-shadow: 0 12px 40px rgba(0, 0, 0, 0.18); }
      .bubble--user { margin-left: auto; background: linear-gradient(135deg, rgba(8, 47, 73, 0.95), rgba(14, 116, 144, 0.7)); }
      header, footer { display: flex; justify-content: space-between; gap: 1rem; color: #7dd3fc; font-size: 0.72rem; }
      header { margin-bottom: 0.65rem; text-transform: uppercase; letter-spacing: 0.08em; }
      footer { margin-top: 0.9rem; color: #94a3b8; }
      .content { color: #ecfeff; line-height: 1.6; }
      .attachments { display: flex; gap: 0.8rem; flex-wrap: wrap; margin-top: 0.9rem; }
      .attachment { width: min(16rem, 100%); margin: 0; display: grid; gap: 0.4rem; }
      .attachment img { width: 100%; border-radius: 1rem; }
      .attachment figcaption { color: #cbd5e1; font-size: 0.76rem; }
      .thinking, .tool { margin-top: 0.85rem; padding-top: 0.85rem; border-top: 1px solid rgba(125, 211, 252, 0.12); color: #cbd5e1; }
      .tool { display: grid; gap: 0.65rem; }
      .tool__block { display: grid; gap: 0.35rem; }
      .tool__title { color: #f8fafc; font-weight: 700; }
      pre { margin: 0; overflow: auto; background: rgba(2, 6, 23, 0.6); padding: 0.85rem; border-radius: 0.8rem; white-space: pre-wrap; }
    `,
  ],
})
export class MessageListComponent {
  @Input() messages: ChatMessage[] = [];

  constructor(private readonly sanitizer: DomSanitizer) {}

  attachmentSource(attachment: ImageAttachment): string {
    return `data:${attachment.mime_type};base64,${attachment.data}`;
  }

  thinkingDuration(message: ChatMessage): string {
    const duration = message.thinking?.duration_ms ?? 0;
    return duration > 0 ? `${(duration / 1000).toFixed(1)}s` : 'in progress';
  }

  renderMarkdown(content: string): SafeHtml {
    const escaped = this.escapeHTML(content ?? '');
    const blocks = escaped.replace(/```([\w-]+)?\n([\s\S]*?)```/g, (_, language, code) => {
      const label = language ? `<div class="code-lang">${language}</div>` : '';
      return `${label}<pre><code>${code.trim()}</code></pre>`;
    });
    const inline = blocks
      .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.+?)\*/g, '<em>$1</em>')
      .replace(/`([^`]+)`/g, '<code>$1</code>')
      .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noreferrer">$1</a>')
      .replace(/\n/g, '<br>');
    return this.sanitizer.bypassSecurityTrustHtml(inline);
  }

  private escapeHTML(value: string): string {
    return value
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll("'", '&#39;');
  }
}
