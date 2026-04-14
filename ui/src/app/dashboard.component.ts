import { CommonModule } from '@angular/common';
import {
  AfterViewChecked,
  Component,
  ElementRef,
  ViewChild,
  computed,
  inject,
  signal,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  ChatMessage,
  ChatService,
  Conversation,
  ImageAttachment,
  ToolInvocation,
} from '../services/chat.service';
import { UiStateService } from '../services/ui-state.service';

interface MessageSegment {
  kind: 'markdown' | 'code';
  content: string;
  language?: string;
}

interface ConversationGroup {
  label: string;
  conversations: Conversation[];
}

@Component({
  selector: 'dashboard-view',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <main class="chat-shell">
      <aside class="history-rail">
        <div class="rail-head">
          <div>
            <p class="eyebrow">Chat history</p>
            <h1>Core Chat</h1>
          </div>
          <button type="button" class="ghost-button" (click)="createConversation()">
            New chat
          </button>
        </div>

        <div class="model-chip">
          <span>Active model</span>
          <strong>{{ selectedModelLabel() }}</strong>
        </div>

        <div class="history-list">
          @for (group of filteredGroups(); track group.label) {
            <section class="history-group">
              <h2>{{ group.label }}</h2>
              @for (conversation of group.conversations; track conversation.id) {
                <button
                  type="button"
                  class="history-row"
                  [class.active]="activeConversation()?.id === conversation.id"
                  (click)="chat.selectConversation(conversation.id)"
                >
                  <span class="history-title">{{ conversation.title }}</span>
                  <span class="history-meta">
                    {{ conversation.model }} · {{ conversation.updatedAt | date: 'shortTime' }}
                  </span>
                </button>
              }
            </section>
          }
        </div>
      </aside>

      <section class="thread-shell">
        <header class="thread-head">
          <div>
            <p class="eyebrow">Conversation</p>
            <h2>{{ activeConversation()?.title || 'New conversation' }}</h2>
          </div>

          <div class="thread-actions">
            <label class="model-select">
              <span>Model</span>
              <select [ngModel]="chat.selectedModel()" (ngModelChange)="chat.selectModel($event)">
                @for (model of chat.models(); track model.name) {
                  <option [value]="model.name">{{ model.name }} · {{ model.architecture }}</option>
                }
              </select>
            </label>

            <button
              type="button"
              class="ghost-button"
              [disabled]="!activeConversation()"
              (click)="exportActiveConversation()"
            >
              Export
            </button>
            <button
              type="button"
              class="ghost-button danger"
              [disabled]="!activeConversation()"
              (click)="deleteActiveConversation()"
            >
              Delete
            </button>
          </div>
        </header>

        <div class="thread" #thread (scroll)="onThreadScroll()">
          @if (activeConversation(); as conversation) {
            @if (conversation.messages.length === 0) {
              <section class="empty-thread">
                <p class="eyebrow">Local-first</p>
                <h3>Start a conversation</h3>
                <p>
                  This shell implements the RFC chat surface with local history, progressive
                  rendering, model selection, settings, tool-call blocks, and multimodal inputs.
                </p>
              </section>
            }

            @for (message of conversation.messages; track message.id) {
              <article class="message" [class.user]="message.role === 'user'">
                <div class="message-meta">
                  <span>{{ message.role === 'user' ? 'You' : 'Assistant' }}</span>
                  <time>{{ message.createdAt | date: 'shortTime' }}</time>
                </div>

                <div class="bubble">
                  @if (message.attachments?.length) {
                    <div class="attachment-grid">
                      @for (attachment of message.attachments || []; track attachment.id) {
                        <button
                          type="button"
                          class="attachment-preview"
                          (click)="lightboxImage.set(attachment)"
                        >
                          <img [src]="attachment.data" [alt]="attachment.filename" />
                          <span>{{ attachment.filename }}</span>
                        </button>
                      }
                    </div>
                  }

                  @if (message.thinking?.content) {
                    <section class="thinking-panel">
                      <button type="button" class="collapse-toggle" (click)="toggleThinking(message.id)">
                        <span>Thinking{{ message.thinking?.active ? '...' : '' }}</span>
                        <small>{{ thinkingDuration(message) }}</small>
                      </button>
                      @if (thinkingExpanded(message.id)) {
                        <pre>{{ message.thinking?.content }}</pre>
                      }
                    </section>
                  }

                  @for (segment of segmentsFor(message); track $index) {
                    @if (segment.kind === 'markdown') {
                      <div class="markdown" [innerHTML]="renderMarkdown(segment.content)"></div>
                    } @else {
                      <div class="code-block">
                        <div class="code-head">
                          <span>{{ segment.language || 'text' }}</span>
                          <button type="button" (click)="copyText(segment.content)">Copy</button>
                        </div>
                        <pre><code>{{ segment.content }}</code></pre>
                      </div>
                    }
                  }

                  @if (message.toolCalls?.length) {
                    @for (tool of message.toolCalls || []; track tool.id) {
                      <section class="tool-panel">
                        <button
                          type="button"
                          class="collapse-toggle"
                          (click)="toggleTool(tool.id)"
                        >
                          <span>{{ tool.name }}</span>
                          <small>Tool call</small>
                        </button>
                        @if (toolExpanded(tool.id)) {
                          <div class="tool-body">
                            <pre>{{ tool.arguments | json }}</pre>
                            <p>{{ tool.result }}</p>
                            @if (tool.error) {
                              <strong class="error-text">{{ tool.error }}</strong>
                            }
                          </div>
                        }
                      </section>
                    }
                  }
                </div>
              </article>
            }
          }
        </div>

        @if (!autoScroll()) {
          <button type="button" class="scroll-pill" (click)="jumpToBottom()">Scroll to bottom</button>
        }

        <footer class="composer">
          @if (chat.queuedAttachments().length) {
            <div class="queued-attachments">
              @for (attachment of chat.queuedAttachments(); track attachment.id) {
                <article class="queued-card">
                  <img [src]="attachment.data" [alt]="attachment.filename" />
                  <div>
                    <strong>{{ attachment.filename }}</strong>
                    <span>{{ attachment.width }} × {{ attachment.height }}</span>
                  </div>
                  <button type="button" (click)="chat.removeAttachment(attachment.id)">
                    Remove
                  </button>
                </article>
              }
            </div>
          }

          <div class="composer-shell">
            <textarea
              #composer
              rows="1"
              [ngModel]="chat.draft()"
              (ngModelChange)="chat.setDraft($event)"
              (input)="resizeComposer(composer)"
              (keydown)="onComposerKeydown($event)"
              placeholder="Ask locally, attach images, or describe the task."
            ></textarea>

            <div class="composer-actions">
              <label class="icon-button" title="Attach image">
                <input type="file" accept="image/*" hidden (change)="onFilePicked($event)" />
                <i class="fa-regular fa-image"></i>
              </label>
              <button
                type="button"
                class="icon-button"
                [disabled]="chat.busy()"
                (click)="pasteImagePlaceholder()"
                title="Paste hint"
              >
                <i class="fa-regular fa-clipboard"></i>
              </button>
              <button
                type="button"
                class="send-button"
                [disabled]="!chat.canSend() || chat.busy()"
                (click)="sendMessage(composer)"
              >
                {{ chat.busy() ? 'Streaming' : 'Send' }}
              </button>
            </div>
          </div>

          <div class="composer-meta">
            <span>{{ chat.draft().length }} characters</span>
            <span>Enter to send · Shift+Enter for newline</span>
          </div>
        </footer>
      </section>
    </main>

    @if (lightboxImage(); as image) {
      <div class="lightbox" (click)="lightboxImage.set(null)">
        <figure class="lightbox-card" (click)="$event.stopPropagation()">
          <img [src]="image.data" [alt]="image.filename" />
          <figcaption>
            <strong>{{ image.filename }}</strong>
            <span>{{ image.width }} × {{ image.height }}</span>
          </figcaption>
        </figure>
      </div>
    }
  `,
  styles: [
    `
      .chat-shell {
        min-height: calc(100vh - 2.75rem);
        display: grid;
        grid-template-columns: minmax(18rem, 22rem) minmax(0, 1fr);
        background:
          radial-gradient(circle at top left, rgba(245, 158, 11, 0.16), transparent 24%),
          radial-gradient(circle at top right, rgba(56, 189, 248, 0.12), transparent 28%),
          linear-gradient(180deg, #08111c 0%, #0f172a 54%, #111827 100%);
        color: #f8fafc;
      }

      @media (max-width: 1100px) {
        .chat-shell {
          grid-template-columns: 1fr;
        }
      }

      .history-rail {
        border-right: 1px solid rgba(255, 255, 255, 0.08);
        padding: 1.5rem 1rem 1.1rem;
        background: rgba(7, 12, 20, 0.82);
        backdrop-filter: blur(18px);
      }

      .rail-head,
      .thread-head,
      .composer-meta,
      .composer-actions,
      .thread-actions,
      .history-row,
      .attachment-preview,
      .queued-card,
      .tool-body,
      .code-head,
      .message-meta,
      .collapse-toggle,
      .lightbox figcaption {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 0.75rem;
      }

      .eyebrow {
        margin: 0;
        text-transform: uppercase;
        letter-spacing: 0.16em;
        font-size: 0.72rem;
        color: rgba(248, 250, 252, 0.55);
      }

      h1,
      h2,
      h3 {
        margin: 0.3rem 0 0;
        font-family: 'Iowan Old Style', 'Palatino Linotype', serif;
      }

      .ghost-button,
      .send-button,
      .icon-button,
      .history-row,
      .attachment-preview,
      .collapse-toggle {
        border: 0;
        cursor: pointer;
      }

      .ghost-button,
      .icon-button,
      .collapse-toggle {
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.08);
        color: inherit;
        padding: 0.7rem 0.95rem;
      }

      .ghost-button.danger {
        color: #fca5a5;
      }

      .model-chip {
        margin: 1.25rem 0;
        border-radius: 1rem;
        padding: 1rem;
        background: linear-gradient(135deg, rgba(249, 115, 22, 0.28), rgba(59, 130, 246, 0.16));
      }

      .model-chip span,
      .history-meta,
      .composer-meta,
      .queued-card span,
      .tool-body p,
      .message-meta time,
      .lightbox span,
      .thinking-panel small {
        color: rgba(226, 232, 240, 0.7);
        font-size: 0.84rem;
      }

      .history-list {
        display: grid;
        gap: 1rem;
        max-height: calc(100vh - 15rem);
        overflow: auto;
        padding-right: 0.3rem;
      }

      .history-group {
        display: grid;
        gap: 0.5rem;
      }

      .history-group h2 {
        margin: 0;
        font-size: 0.8rem;
        text-transform: uppercase;
        letter-spacing: 0.12em;
        color: rgba(255, 255, 255, 0.5);
      }

      .history-row {
        width: 100%;
        padding: 0.9rem 1rem;
        border-radius: 1rem;
        text-align: left;
        background: rgba(255, 255, 255, 0.04);
        flex-direction: column;
        align-items: flex-start;
      }

      .history-row.active {
        background: rgba(251, 191, 36, 0.14);
        box-shadow: inset 0 0 0 1px rgba(251, 191, 36, 0.18);
      }

      .history-title {
        font-weight: 700;
      }

      .thread-shell {
        display: grid;
        grid-template-rows: auto minmax(0, 1fr) auto;
        min-height: calc(100vh - 2.75rem);
        position: relative;
      }

      .thread-head,
      .composer {
        padding: 1.4rem 1.6rem;
        border-bottom: 1px solid rgba(255, 255, 255, 0.06);
        background: rgba(10, 15, 26, 0.58);
        backdrop-filter: blur(18px);
      }

      .composer {
        border-bottom: 0;
        border-top: 1px solid rgba(255, 255, 255, 0.06);
      }

      .thread {
        overflow: auto;
        padding: 1.5rem;
        display: grid;
        gap: 1rem;
      }

      .empty-thread,
      .bubble,
      .tool-panel,
      .thinking-panel {
        border-radius: 1.2rem;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.06);
      }

      .empty-thread {
        padding: 2rem;
        max-width: 48rem;
      }

      .message {
        display: grid;
        gap: 0.4rem;
        max-width: min(48rem, 100%);
      }

      .message.user {
        margin-left: auto;
      }

      .bubble {
        padding: 1rem;
        display: grid;
        gap: 0.9rem;
      }

      .message.user .bubble {
        background: linear-gradient(135deg, rgba(251, 146, 60, 0.18), rgba(245, 158, 11, 0.08));
      }

      .markdown {
        color: #e2e8f0;
        line-height: 1.7;
      }

      .markdown p {
        margin: 0 0 0.8rem;
      }

      .markdown a {
        color: #7dd3fc;
      }

      .code-block,
      .tool-panel,
      .thinking-panel {
        overflow: hidden;
      }

      .code-block pre,
      .thinking-panel pre,
      .tool-body pre {
        margin: 0;
        overflow: auto;
        padding: 1rem;
        background: rgba(2, 6, 23, 0.86);
        color: #dbeafe;
      }

      .code-head {
        padding: 0.75rem 1rem;
        background: rgba(15, 23, 42, 0.86);
        text-transform: uppercase;
        letter-spacing: 0.08em;
        font-size: 0.75rem;
      }

      .code-head button,
      .queued-card button {
        border: 0;
        background: transparent;
        color: #7dd3fc;
        cursor: pointer;
      }

      .attachment-grid,
      .queued-attachments {
        display: flex;
        gap: 0.8rem;
        flex-wrap: wrap;
      }

      .attachment-preview,
      .queued-card {
        border-radius: 1rem;
        padding: 0.6rem;
        background: rgba(255, 255, 255, 0.05);
        min-width: 10rem;
      }

      .attachment-preview img,
      .queued-card img,
      .lightbox img {
        width: 100%;
        border-radius: 0.9rem;
        object-fit: cover;
      }

      .attachment-preview {
        flex-direction: column;
        align-items: flex-start;
      }

      .composer-shell {
        display: grid;
        grid-template-columns: minmax(0, 1fr) auto;
        gap: 1rem;
        padding: 1rem;
        border-radius: 1.25rem;
        background: rgba(255, 255, 255, 0.04);
      }

      textarea,
      select {
        width: 100%;
        border: 0;
        background: transparent;
        color: inherit;
        resize: none;
        outline: none;
        font: inherit;
      }

      textarea {
        min-height: 3rem;
        max-height: 10.5rem;
      }

      .send-button {
        border-radius: 999px;
        padding: 0.95rem 1.15rem;
        background: linear-gradient(135deg, #f59e0b, #fb7185);
        color: #111827;
        font-weight: 800;
      }

      .send-button:disabled,
      .ghost-button:disabled {
        opacity: 0.45;
        cursor: default;
      }

      .model-select {
        min-width: 14rem;
        display: grid;
        gap: 0.3rem;
      }

      .model-select span {
        font-size: 0.78rem;
        color: rgba(226, 232, 240, 0.65);
      }

      .scroll-pill {
        position: absolute;
        right: 1.5rem;
        bottom: 10rem;
        border: 0;
        border-radius: 999px;
        padding: 0.75rem 1rem;
        background: rgba(15, 23, 42, 0.92);
        color: #f8fafc;
        cursor: pointer;
      }

      .lightbox {
        position: fixed;
        inset: 0;
        background: rgba(2, 6, 23, 0.82);
        backdrop-filter: blur(12px);
        display: grid;
        place-items: center;
        padding: 2rem;
        z-index: 90;
      }

      .lightbox-card {
        max-width: min(62rem, 100%);
        margin: 0;
        display: grid;
        gap: 0.75rem;
      }

      .error-text {
        color: #fda4af;
      }
    `,
  ],
})
export class DashboardComponent implements AfterViewChecked {
  private readonly uiState = inject(UiStateService);
  protected readonly chat = inject(ChatService);

  @ViewChild('thread') private thread?: ElementRef<HTMLDivElement>;

  readonly lightboxImage = signal<ImageAttachment | null>(null);
  readonly autoScroll = signal(true);
  private readonly expandedThinkingIds = signal<Set<string>>(new Set());
  private readonly expandedToolIds = signal<Set<string>>(new Set());

  protected readonly searchQuery = this.uiState.searchQuery;
  protected readonly activeConversation = this.chat.activeConversation;

  protected readonly filteredGroups = computed<ConversationGroup[]>(() => {
    const query = this.searchQuery().toLowerCase();
    const grouped = new Map<string, Conversation[]>();

    for (const conversation of this.chat.conversations()) {
      if (query) {
        const haystack = `${conversation.title} ${conversation.messages.map((message) => message.content).join(' ')}`.toLowerCase();
        if (!haystack.includes(query)) {
          continue;
        }
      }
      const bucket = bucketLabel(conversation.updatedAt);
      grouped.set(bucket, [...(grouped.get(bucket) ?? []), conversation]);
    }

    return Array.from(grouped.entries()).map(([label, conversations]) => ({
      label,
      conversations,
    }));
  });

  ngAfterViewChecked(): void {
    if (this.autoScroll()) {
      this.jumpToBottom();
    }
  }

  protected createConversation(): void {
    this.chat.createConversation();
    this.uiState.clearSearchQuery();
  }

  protected async sendMessage(textarea: HTMLTextAreaElement): Promise<void> {
    await this.chat.sendMessage();
    this.resizeComposer(textarea, true);
    this.jumpToBottom();
  }

  protected deleteActiveConversation(): void {
    const active = this.chat.activeConversation();
    if (!active) {
      return;
    }
    this.chat.deleteConversation(active.id);
  }

  protected exportActiveConversation(): void {
    const active = this.chat.activeConversation();
    if (!active) {
      return;
    }
    const blob = new Blob([this.chat.exportConversation(active)], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement('a');
    anchor.href = url;
    anchor.download = `${active.title.replace(/\s+/g, '-').toLowerCase() || 'conversation'}.md`;
    anchor.click();
    URL.revokeObjectURL(url);
  }

  protected resizeComposer(textarea: HTMLTextAreaElement, reset = false): void {
    if (reset) {
      textarea.style.height = 'auto';
      return;
    }
    textarea.style.height = 'auto';
    textarea.style.height = `${Math.min(textarea.scrollHeight, 168)}px`;
  }

  protected onComposerKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      const textarea = event.target as HTMLTextAreaElement;
      void this.sendMessage(textarea);
    }
  }

  protected async onFilePicked(event: Event): Promise<void> {
    const input = event.target as HTMLInputElement;
    const files = Array.from(input.files ?? []);
    for (const file of files) {
      await this.chat.addAttachment(file);
    }
    input.value = '';
  }

  protected pasteImagePlaceholder(): void {
    this.chat.setDraft(`${this.chat.draft()}\n[Paste image from clipboard here when running inside the desktop shell.]`.trim());
  }

  protected onThreadScroll(): void {
    const element = this.thread?.nativeElement;
    if (!element) {
      return;
    }
    const distance = element.scrollHeight - element.scrollTop - element.clientHeight;
    this.autoScroll.set(distance < 48);
  }

  protected jumpToBottom(): void {
    const element = this.thread?.nativeElement;
    if (!element) {
      return;
    }
    element.scrollTop = element.scrollHeight;
    this.autoScroll.set(true);
  }

  protected selectedModelLabel(): string {
    const selected = this.chat.models().find((model) => model.name === this.chat.selectedModel());
    return selected ? `${selected.name} · ${selected.architecture}` : this.chat.selectedModel();
  }

  protected segmentsFor(message: ChatMessage): MessageSegment[] {
    const source = message.content || (message.streaming ? '...' : '');
    const parts: MessageSegment[] = [];
    const pattern = /```([\w-]*)\n([\s\S]*?)```/g;
    let lastIndex = 0;
    let match: RegExpExecArray | null;

    while ((match = pattern.exec(source)) !== null) {
      if (match.index > lastIndex) {
        parts.push({ kind: 'markdown', content: source.slice(lastIndex, match.index) });
      }
      parts.push({
        kind: 'code',
        language: match[1] || 'text',
        content: match[2].trim(),
      });
      lastIndex = match.index + match[0].length;
    }

    if (lastIndex < source.length) {
      parts.push({ kind: 'markdown', content: source.slice(lastIndex) });
    }

    return parts.length > 0 ? parts : [{ kind: 'markdown', content: source }];
  }

  protected renderMarkdown(content: string): string {
    const escaped = escapeHtml(content);
    return escaped
      .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noreferrer">$1</a>')
      .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
      .replace(/\*([^*]+)\*/g, '<em>$1</em>')
      .replace(/`([^`]+)`/g, '<code>$1</code>')
      .split(/\n{2,}/)
      .map((block) => `<p>${block.replace(/\n/g, '<br />')}</p>`)
      .join('');
  }

  protected async copyText(value: string): Promise<void> {
    await navigator.clipboard.writeText(value);
  }

  protected toggleThinking(id: string): void {
    this.expandedThinkingIds.update((current) => toggleSet(current, id));
  }

  protected thinkingExpanded(id: string): boolean {
    return this.expandedThinkingIds().has(id);
  }

  protected thinkingDuration(message: ChatMessage): string {
    const thinking = message.thinking;
    if (!thinking?.startedAt) {
      return '';
    }
    const start = new Date(thinking.startedAt).getTime();
    const end = thinking.finishedAt ? new Date(thinking.finishedAt).getTime() : Date.now();
    return `${Math.max(end - start, 0) / 1000}s`;
  }

  protected toggleTool(id: string): void {
    this.expandedToolIds.update((current) => toggleSet(current, id));
  }

  protected toolExpanded(id: string): boolean {
    return this.expandedToolIds().has(id);
  }
}

function bucketLabel(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const startOfDate = new Date(date.getFullYear(), date.getMonth(), date.getDate());
  const diff = (startOfToday.getTime() - startOfDate.getTime()) / 86_400_000;
  if (diff <= 0) {
    return 'Today';
  }
  if (diff <= 1) {
    return 'Yesterday';
  }
  if (diff < 7) {
    return 'Previous 7 Days';
  }
  return 'Older';
}

function escapeHtml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');
}

function toggleSet(current: Set<string>, value: string): Set<string> {
  const next = new Set(current);
  if (next.has(value)) {
    next.delete(value);
  } else {
    next.add(value);
  }
  return next;
}
