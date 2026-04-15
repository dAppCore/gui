import { CommonModule, DatePipe } from '@angular/common';
import {
  AfterViewChecked,
  Component,
  ElementRef,
  Input,
  OnChanges,
  SimpleChanges,
  ViewChild,
} from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { ChatMessage, ImageAttachment } from './chat.types';

type RenderBlock =
  | { kind: 'text'; html: SafeHtml }
  | { kind: 'code'; code: string; language: string };

@Component({
  selector: 'chat-message-list',
  standalone: true,
  imports: [CommonModule, DatePipe],
  template: `
    <div class="thread-shell">
      <div #viewport class="thread" (scroll)="onScroll()">
        <article *ngFor="let message of messages; trackBy: trackByMessage" class="bubble" [class.bubble--user]="message.role === 'user'">
          <header>
            <strong>{{ message.role }}</strong>
            <span>{{ message.model || 'local' }}</span>
          </header>

          <ng-container *ngFor="let block of renderBlocks(message.content)">
            <section *ngIf="block.kind === 'text'" class="content" [innerHTML]="block.html"></section>
            <section *ngIf="block.kind === 'code'" class="code-block">
              <div class="code-block__bar">
                <span>{{ block.language || 'code' }}</span>
                <button type="button" class="code-block__copy" (click)="copyCode(block.code)">Copy</button>
              </div>
              <pre><code>{{ block.code }}</code></pre>
            </section>
          </ng-container>

          <section *ngIf="message.attachments?.length" class="attachments">
            <figure *ngFor="let attachment of message.attachments" class="attachment">
              <img [src]="attachmentSource(attachment)" [alt]="attachment.filename" />
              <figcaption>{{ attachment.filename }} · {{ attachment.width }}×{{ attachment.height }}</figcaption>
            </figure>
          </section>

          <section *ngIf="message.thinking" class="thinking">
            <div class="thinking__badge" [class.thinking__badge--active]="message.thinking.active">
              Thinking...
            </div>
            <details>
              <summary>Thought for {{ thinkingDuration(message) }}</summary>
              <p>{{ message.thinking.content || 'Thinking in progress' }}</p>
            </details>
          </section>

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

          <footer [attr.title]="message.created_at | date : 'full'">
            {{ message.created_at | date: 'MMM d, HH:mm:ss' }}
          </footer>
        </article>
      </div>

      <button *ngIf="streaming && !pinnedToBottom" type="button" class="scroll-pill" (click)="scrollToBottom()">
        Scroll to bottom
      </button>
    </div>
  `,
  styles: [
    `
      .thread-shell { position: relative; min-height: 100%; height: 100%; }
      .thread { display: grid; gap: 1rem; align-content: start; min-height: 0; max-height: 100%; overflow: auto; padding-right: 0.25rem; }
      .bubble { max-width: 54rem; padding: 1rem 1.1rem; border-radius: 1.25rem; background: rgba(11, 27, 44, 0.88); border: 1px solid rgba(125, 211, 252, 0.12); box-shadow: 0 12px 40px rgba(0, 0, 0, 0.18); }
      .bubble--user { margin-left: auto; background: linear-gradient(135deg, rgba(8, 47, 73, 0.95), rgba(14, 116, 144, 0.7)); }
      header, footer { display: flex; justify-content: space-between; gap: 1rem; color: #7dd3fc; font-size: 0.72rem; }
      header { margin-bottom: 0.65rem; text-transform: uppercase; letter-spacing: 0.08em; }
      footer { margin-top: 0.9rem; color: #94a3b8; opacity: 0; transition: opacity 0.2s ease; }
      .bubble:hover footer { opacity: 1; }
      .content { color: #ecfeff; line-height: 1.6; }
      .attachments { display: flex; gap: 0.8rem; flex-wrap: wrap; margin-top: 0.9rem; }
      .attachment { width: min(16rem, 100%); margin: 0; display: grid; gap: 0.4rem; }
      .attachment img { width: 100%; border-radius: 1rem; }
      .attachment figcaption { color: #cbd5e1; font-size: 0.76rem; }
      .thinking, .tool { margin-top: 0.85rem; padding-top: 0.85rem; border-top: 1px solid rgba(125, 211, 252, 0.12); color: #cbd5e1; }
      .thinking { display: grid; gap: 0.5rem; }
      .thinking__badge { display: inline-flex; width: fit-content; gap: 0.45rem; align-items: center; border-radius: 999px; padding: 0.3rem 0.65rem; background: rgba(15, 23, 42, 0.7); color: #f8fafc; }
      .thinking__badge--active::before { content: ''; width: 0.55rem; height: 0.55rem; border-radius: 50%; background: #f59e0b; box-shadow: 0 0 0 0 rgba(245, 158, 11, 0.5); animation: pulse 1.4s infinite; }
      details summary { cursor: pointer; color: #e2e8f0; }
      .tool { display: grid; gap: 0.65rem; }
      .tool__block { display: grid; gap: 0.35rem; }
      .tool__title { color: #f8fafc; font-weight: 700; }
      .code-block { margin-top: 0.9rem; overflow: hidden; border-radius: 0.95rem; border: 1px solid rgba(148, 163, 184, 0.16); background: rgba(2, 6, 23, 0.64); }
      .code-block__bar { display: flex; justify-content: space-between; gap: 1rem; align-items: center; padding: 0.55rem 0.8rem; border-bottom: 1px solid rgba(148, 163, 184, 0.12); color: #cbd5e1; font-size: 0.72rem; text-transform: uppercase; letter-spacing: 0.08em; }
      .code-block__copy { border: 1px solid rgba(148, 163, 184, 0.2); border-radius: 999px; background: rgba(15, 23, 42, 0.8); color: #e2e8f0; padding: 0.3rem 0.7rem; cursor: pointer; text-transform: none; letter-spacing: normal; }
      pre { margin: 0; overflow: auto; padding: 0.85rem; white-space: pre-wrap; color: #f8fafc; }
      .scroll-pill { position: absolute; right: 0.75rem; bottom: 0.75rem; border: 1px solid rgba(251, 191, 36, 0.28); border-radius: 999px; background: rgba(124, 45, 18, 0.92); color: #fde68a; padding: 0.75rem 1rem; cursor: pointer; box-shadow: 0 10px 24px rgba(0, 0, 0, 0.25); }
      @keyframes pulse {
        0% { transform: scale(0.92); opacity: 0.75; }
        70% { transform: scale(1.08); opacity: 1; }
        100% { transform: scale(0.92); opacity: 0.75; }
      }
    `,
  ],
})
export class MessageListComponent implements OnChanges, AfterViewChecked {
  @Input() messages: ChatMessage[] = [];
  @Input() streaming = false;

  @ViewChild('viewport') private readonly viewport?: ElementRef<HTMLDivElement>;

  private pendingScroll = false;
  pinnedToBottom = true;

  constructor(private readonly sanitizer: DomSanitizer) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['messages'] || changes['streaming']) {
      this.pendingScroll = true;
    }
  }

  ngAfterViewChecked(): void {
    if (this.pendingScroll) {
      this.pendingScroll = false;
      if (this.pinnedToBottom) {
        queueMicrotask(() => this.scrollToBottom());
      }
    }
  }

  trackByMessage(_index: number, message: ChatMessage): string {
    return message.id;
  }

  attachmentSource(attachment: ImageAttachment): string {
    return `data:${attachment.mime_type};base64,${attachment.data}`;
  }

  thinkingDuration(message: ChatMessage): string {
    const duration = message.thinking?.duration_ms ?? 0;
    return duration > 0 ? `${(duration / 1000).toFixed(1)}s` : 'in progress';
  }

  renderBlocks(content: string): RenderBlock[] {
    const source = content ?? '';
    if (!source.trim()) {
      return [];
    }

    const blocks: RenderBlock[] = [];
    const pattern = /```([\w-]+)?\n([\s\S]*?)```/g;
    let lastIndex = 0;
    let match: RegExpExecArray | null;

    while ((match = pattern.exec(source)) !== null) {
      const before = source.slice(lastIndex, match.index).trim();
      if (before) {
        blocks.push({ kind: 'text', html: this.renderInlineMarkdown(before) });
      }
      blocks.push({
        kind: 'code',
        language: match[1] ?? '',
        code: match[2].replace(/\n+$/, ''),
      });
      lastIndex = pattern.lastIndex;
    }

    const after = source.slice(lastIndex).trim();
    if (after) {
      blocks.push({ kind: 'text', html: this.renderInlineMarkdown(after) });
    }

    if (blocks.length === 0) {
      blocks.push({ kind: 'text', html: this.renderInlineMarkdown(source) });
    }
    return blocks;
  }

  copyCode(code: string): void {
    if (!code) {
      return;
    }
    if (navigator.clipboard?.writeText) {
      void navigator.clipboard.writeText(code);
      return;
    }
    const textarea = document.createElement('textarea');
    textarea.value = code;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.focus();
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
  }

  onScroll(): void {
    const element = this.viewport?.nativeElement;
    if (!element) {
      return;
    }
    const distanceFromBottom = element.scrollHeight - element.scrollTop - element.clientHeight;
    this.pinnedToBottom = distanceFromBottom < 48;
  }

  scrollToBottom(): void {
    const element = this.viewport?.nativeElement;
    if (!element) {
      return;
    }
    element.scrollTop = element.scrollHeight;
    this.pinnedToBottom = true;
  }

  private renderInlineMarkdown(content: string): SafeHtml {
    const escaped = this.escapeHTML(content ?? '');
    const inline = escaped
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
