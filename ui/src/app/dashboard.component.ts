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
  SUPPORTED_CHAT_IMAGE_ACCEPT,
  SUPPORTED_CHAT_IMAGE_LABEL,
  ToolInvocation,
  conversationSearchText,
  isSupportedChatImageFile,
  isSupportedChatImageMimeType,
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
          <div>
            <span>Active model</span>
            <strong>{{ selectedModelLabel() }}</strong>
          </div>
          <div class="model-status">
            <span class="status-pill" [class.loading]="chat.modelSwitching()">
              <span class="status-dot"></span>
              {{ chat.modelSwitching() ? 'Loading' : 'Loaded' }}
            </span>
            <span class="status-pill subtle">
              {{ selectedModelSupportsVision() ? 'Vision' : 'Text only' }}
            </span>
          </div>
        </div>

        <div class="history-list">
          @for (group of filteredGroups(); track group.label) {
            <section class="history-group">
              <h2>{{ group.label }}</h2>
              @for (conversation of group.conversations; track conversation.id) {
                <article
                  class="history-row"
                  [class.active]="activeConversation()?.id === conversation.id"
                >
                  @if (editingConversationId() === conversation.id) {
                    <div class="history-edit">
                      <input
                        #renameInput
                        type="text"
                        class="history-title-input"
                        [ngModel]="renameDraft()"
                        (ngModelChange)="renameDraft.set($event)"
                        (keydown)="onRenameKeydown($event, conversation.id)"
                        (click)="$event.stopPropagation()"
                        (blur)="commitRename(conversation.id)"
                      />
                      <div class="history-actions">
                        <button
                          type="button"
                          class="row-icon save"
                          (mousedown)="$event.preventDefault()"
                          (click)="commitRename(conversation.id)"
                        >
                          Save
                        </button>
                        <button
                          type="button"
                          class="row-icon"
                          (mousedown)="$event.preventDefault()"
                          (click)="cancelRename()"
                        >
                          Cancel
                        </button>
                      </div>
                    </div>
                  } @else {
                    <button
                      type="button"
                      class="history-main"
                      (click)="chat.selectConversation(conversation.id)"
                    >
                      <span class="history-title">{{ conversation.title }}</span>
                      <span class="history-meta">
                        {{ conversation.model }} · {{ conversation.updatedAt | date: 'shortTime' }}
                      </span>
                    </button>
                    <div class="history-actions">
                      <button
                        type="button"
                        class="row-icon"
                        (click)="startRename(conversation, $event)"
                      >
                        Rename
                      </button>
                      <button
                        type="button"
                        class="row-icon danger"
                        (click)="deleteConversation(conversation, $event)"
                      >
                        Delete
                      </button>
                    </div>
                  }
                </article>
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
              class="ghost-button"
              [disabled]="!activeConversation()"
              (click)="clearActiveConversation()"
            >
              Clear
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
                  <time
                    [attr.datetime]="message.createdAt"
                    [title]="message.createdAt | date: 'medium'"
                  >
                    {{ message.createdAt | date: 'shortTime' }}
                  </time>
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
                      <button
                        type="button"
                        class="collapse-toggle"
                        (click)="toggleThinking(message.id)"
                      >
                        <span class="thinking-label" [class.active]="message.thinking?.active">
                          @if (message.thinking?.active) {
                            <span class="thinking-pulse" aria-hidden="true"></span>
                          }
                          Thinking{{ message.thinking?.active ? '...' : '' }}
                        </span>
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
                        <pre><code [innerHTML]="renderCodeBlock(segment)"></code></pre>
                      </div>
                    }
                  }

                  @if (message.toolCalls?.length) {
                    @for (tool of message.toolCalls || []; track tool.id) {
                      <section class="tool-panel">
                        <button type="button" class="collapse-toggle" (click)="toggleTool(tool.id)">
                          <span>{{ tool.name }}</span>
                          <small
                            class="tool-status"
                            [class.pending]="toolStatus(tool) === 'pending'"
                            [class.error]="toolStatus(tool) === 'error'"
                          >
                            @if (toolStatus(tool) === 'pending') {
                              <span class="tool-spinner" aria-hidden="true"></span>
                            }
                            {{ toolStatusLabel(tool) }}
                          </small>
                        </button>
                        @if (toolExpanded(tool.id)) {
                          <div class="tool-body">
                            <pre>{{ tool.arguments | json }}</pre>
                            @if (tool.result) {
                              <p>{{ tool.result }}</p>
                            } @else if (toolStatus(tool) === 'pending') {
                              <p>Waiting for the local MCP tool result...</p>
                            }
                            @if (tool.error) {
                              <strong class="error-text">{{ tool.error }}</strong>
                            }
                          </div>
                        }
                      </section>
                    }
                  }

                  @if (message.streaming) {
                    <span class="stream-cursor" aria-hidden="true"></span>
                  }
                </div>
              </article>
            }
          }
        </div>

        @if (!autoScroll()) {
          <button type="button" class="scroll-pill" (click)="jumpToBottom()">
            Scroll to bottom
          </button>
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

          @if (composerNotice(); as notice) {
            <p class="composer-notice">{{ notice }}</p>
          }

          <div
            class="composer-shell"
            [class.drag-active]="dragActive()"
            [class.attach-disabled]="!selectedModelSupportsVision()"
            (dragenter)="onComposerDragEnter($event)"
            (dragover)="onComposerDragOver($event)"
            (dragleave)="onComposerDragLeave($event)"
            (drop)="onComposerDrop($event)"
          >
            <textarea
              #composer
              rows="1"
              [ngModel]="chat.draft()"
              (ngModelChange)="chat.setDraft($event)"
              (input)="resizeComposer(composer)"
              (keydown)="onComposerKeydown($event)"
              (paste)="onComposerPaste($event)"
              placeholder="Ask locally, attach images, or describe the task."
            ></textarea>

            <div class="composer-actions">
              <input
                #filePicker
                type="file"
                [attr.accept]="attachmentAccept"
                hidden
                [disabled]="!selectedModelSupportsVision()"
                (change)="onFilePicked($event)"
              />
              <button
                type="button"
                class="icon-button"
                [disabled]="!selectedModelSupportsVision()"
                [title]="attachmentButtonTitle()"
                (click)="openFilePicker(filePicker)"
              >
                <i class="fa-regular fa-image"></i>
              </button>
              <button
                type="button"
                class="icon-button"
                [disabled]="chat.busy()"
                (click)="showPasteHint()"
                [title]="attachmentPasteTitle()"
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

      .row-icon.danger {
        color: #fca5a5;
      }

      .model-chip {
        margin: 1.25rem 0;
        border-radius: 1rem;
        padding: 1rem;
        background: linear-gradient(135deg, rgba(249, 115, 22, 0.28), rgba(59, 130, 246, 0.16));
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 1rem;
      }

      .model-status,
      .history-actions,
      .history-edit {
        display: flex;
        align-items: center;
        gap: 0.5rem;
      }

      .status-pill {
        display: inline-flex;
        align-items: center;
        gap: 0.45rem;
        padding: 0.38rem 0.7rem;
        border-radius: 999px;
        background: rgba(15, 23, 42, 0.42);
        font-size: 0.78rem;
      }

      .status-pill.subtle {
        background: rgba(255, 255, 255, 0.08);
      }

      .status-dot {
        width: 0.58rem;
        height: 0.58rem;
        border-radius: 999px;
        background: #86efac;
        box-shadow: 0 0 0.8rem rgba(134, 239, 172, 0.45);
      }

      .status-pill.loading .status-dot {
        background: #facc15;
        box-shadow: 0 0 0.8rem rgba(250, 204, 21, 0.5);
        animation: pulse-dot 0.9s ease-in-out infinite;
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
        display: grid;
        grid-template-columns: minmax(0, 1fr) auto;
        align-items: center;
        gap: 0.75rem;
      }

      .history-row.active {
        background: rgba(251, 191, 36, 0.14);
        box-shadow: inset 0 0 0 1px rgba(251, 191, 36, 0.18);
      }

      .history-main {
        border: 0;
        background: transparent;
        color: inherit;
        text-align: left;
        display: grid;
        gap: 0.3rem;
        min-width: 0;
      }

      .history-main,
      .history-title-input {
        width: 100%;
      }

      .history-edit {
        min-width: 0;
        grid-column: 1 / -1;
      }

      .history-title-input {
        border: 0;
        border-radius: 0.8rem;
        padding: 0.65rem 0.85rem;
        background: rgba(15, 23, 42, 0.86);
        color: inherit;
        font: inherit;
      }

      .history-title {
        font-weight: 700;
        min-width: 0;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .row-icon {
        border: 0;
        border-radius: 999px;
        padding: 0.45rem 0.7rem;
        background: rgba(255, 255, 255, 0.08);
        color: rgba(226, 232, 240, 0.9);
        cursor: pointer;
        font-size: 0.78rem;
      }

      .row-icon.save {
        color: #fde68a;
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

      .message-meta time {
        opacity: 0;
        transition: opacity 140ms ease;
      }

      .message:hover .message-meta time,
      .message:focus-within .message-meta time {
        opacity: 1;
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

      .markdown ul,
      .markdown ol {
        margin: 0 0 0.8rem 1.25rem;
        padding: 0;
      }

      .markdown li + li {
        margin-top: 0.35rem;
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

      .stream-cursor {
        display: inline-block;
        width: 0.7rem;
        height: 1.15rem;
        border-radius: 0.15rem;
        background: linear-gradient(180deg, #f59e0b, #fb7185);
        animation: blink-cursor 0.9s steps(1, end) infinite;
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
        transition:
          border-color 140ms ease,
          background-color 140ms ease,
          box-shadow 140ms ease;
        border: 1px dashed transparent;
      }

      .composer-shell.drag-active {
        background: rgba(14, 116, 144, 0.12);
        border-color: rgba(125, 211, 252, 0.65);
        box-shadow: 0 0 0 1px rgba(125, 211, 252, 0.15);
      }

      .composer-shell.attach-disabled {
        border-color: rgba(148, 163, 184, 0.2);
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
      .ghost-button:disabled,
      .icon-button:disabled {
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

      .tool-status,
      .thinking-label {
        display: inline-flex;
        align-items: center;
        gap: 0.45rem;
      }

      .tool-status.pending,
      .thinking-label.active {
        color: #fde68a;
      }

      .tool-status.error {
        color: #fda4af;
      }

      .tool-spinner,
      .thinking-pulse {
        width: 0.72rem;
        height: 0.72rem;
        border-radius: 999px;
        flex: 0 0 auto;
      }

      .tool-spinner {
        border: 2px solid rgba(253, 224, 71, 0.26);
        border-top-color: #facc15;
        animation: spin 0.8s linear infinite;
      }

      .thinking-pulse {
        background: #f59e0b;
        box-shadow: 0 0 0.8rem rgba(245, 158, 11, 0.45);
        animation: pulse-dot 0.9s ease-in-out infinite;
      }

      :host ::ng-deep .code-block code .token-keyword,
      :host ::ng-deep .code-block code .token-boolean {
        color: #93c5fd;
      }

      :host ::ng-deep .code-block code .token-string {
        color: #86efac;
      }

      :host ::ng-deep .code-block code .token-comment {
        color: #94a3b8;
      }

      :host ::ng-deep .code-block code .token-number {
        color: #fca5a5;
      }

      .composer-notice {
        margin: 0 0 0.85rem;
        color: #fcd34d;
        font-size: 0.84rem;
      }

      @keyframes blink-cursor {
        0%,
        49% {
          opacity: 1;
        }
        50%,
        100% {
          opacity: 0;
        }
      }

      @keyframes pulse-dot {
        0%,
        100% {
          transform: scale(1);
          opacity: 1;
        }
        50% {
          transform: scale(0.82);
          opacity: 0.7;
        }
      }

      @keyframes spin {
        from {
          transform: rotate(0deg);
        }
        to {
          transform: rotate(360deg);
        }
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
  readonly dragActive = signal(false);
  readonly composerNotice = signal('');
  readonly editingConversationId = signal('');
  readonly renameDraft = signal('');
  private readonly expandedThinkingIds = signal<Set<string>>(new Set());
  private readonly expandedToolIds = signal<Set<string>>(new Set());
  private composerNoticeTimer: number | null = null;

  protected readonly searchQuery = this.uiState.searchQuery;
  protected readonly activeConversation = this.chat.activeConversation;
  protected readonly attachmentAccept = SUPPORTED_CHAT_IMAGE_ACCEPT;
  protected readonly selectedModelSupportsVision = computed(
    () => this.chat.selectedModelEntry()?.supportsVision !== false,
  );

  protected readonly filteredGroups = computed<ConversationGroup[]>(() => {
    const query = this.searchQuery().toLowerCase();
    const grouped = new Map<string, Conversation[]>();

    for (const conversation of this.chat.conversations()) {
      if (query) {
        if (!conversationSearchText(conversation).includes(query)) {
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
    void this.chat.createConversation();
    this.uiState.clearSearchQuery();
    this.cancelRename();
  }

  protected async sendMessage(textarea: HTMLTextAreaElement): Promise<void> {
    await this.chat.sendMessage();
    this.resizeComposer(textarea, true);
    this.jumpToBottom();
  }

  protected async deleteActiveConversation(): Promise<void> {
    const active = this.chat.activeConversation();
    if (!active) {
      return;
    }
    if (!this.confirmDeleteConversation(active.title)) {
      return;
    }
    await this.chat.deleteConversation(active.id);
    this.cancelRename();
  }

  protected async clearActiveConversation(): Promise<void> {
    const active = this.chat.activeConversation();
    if (!active) {
      return;
    }
    if (!window.confirm(`Clear every message from "${active.title}"?`)) {
      return;
    }
    await this.chat.clearActiveConversation();
    this.cancelRename();
  }

  protected async exportActiveConversation(): Promise<void> {
    const active = this.chat.activeConversation();
    if (!active) {
      return;
    }
    const blob = new Blob([await this.chat.exportConversation(active)], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement('a');
    anchor.href = url;
    anchor.download = `${active.title.replace(/\s+/g, '-').toLowerCase() || 'conversation'}.md`;
    anchor.click();
    URL.revokeObjectURL(url);
  }

  protected async deleteConversation(conversation: Conversation, event: Event): Promise<void> {
    event.stopPropagation();
    if (!this.confirmDeleteConversation(conversation.title)) {
      return;
    }
    await this.chat.deleteConversation(conversation.id);
    this.cancelRename();
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
    await this.attachFiles(Array.from(input.files ?? []));
    input.value = '';
  }

  protected showPasteHint(): void {
    if (!this.selectedModelSupportsVision()) {
      this.showComposerNotice('The selected model does not accept image input.');
      return;
    }
    this.showComposerNotice('Paste an image directly into the composer with Cmd/Ctrl+V.');
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

  protected attachmentButtonTitle(): string {
    return this.selectedModelSupportsVision()
      ? 'Attach PNG, JPEG, WebP, or GIF'
      : 'This model does not support vision input';
  }

  protected attachmentPasteTitle(): string {
    return this.selectedModelSupportsVision()
      ? 'Paste an image from the clipboard'
      : 'Switch to a vision-capable model to paste images';
  }

  protected openFilePicker(input: HTMLInputElement): void {
    if (!this.selectedModelSupportsVision()) {
      this.showComposerNotice('Switch to Lemer or Lemma to attach images.');
      return;
    }
    input.click();
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
    return renderMarkdownContent(content);
  }

  protected renderCodeBlock(segment: MessageSegment): string {
    return renderCodeContent(segment.content, segment.language);
  }

  protected async copyText(value: string): Promise<void> {
    await navigator.clipboard.writeText(value);
  }

  protected startRename(conversation: Conversation, event: Event): void {
    event.stopPropagation();
    this.editingConversationId.set(conversation.id);
    this.renameDraft.set(conversation.title);
    this.chat.selectConversation(conversation.id);
  }

  protected commitRename(id: string): void {
    if (this.editingConversationId() !== id) {
      return;
    }
    void this.chat.renameConversation(id, this.renameDraft());
    this.cancelRename();
  }

  protected cancelRename(): void {
    this.editingConversationId.set('');
    this.renameDraft.set('');
  }

  protected onRenameKeydown(event: KeyboardEvent, id: string): void {
    if (event.key === 'Enter') {
      event.preventDefault();
      this.commitRename(id);
      return;
    }
    if (event.key === 'Escape') {
      event.preventDefault();
      this.cancelRename();
    }
  }

  protected onComposerDragEnter(event: DragEvent): void {
    if (!hasImageFiles(event.dataTransfer)) {
      return;
    }
    event.preventDefault();
    this.dragActive.set(true);
  }

  protected onComposerDragOver(event: DragEvent): void {
    if (!hasImageFiles(event.dataTransfer)) {
      return;
    }
    event.preventDefault();
    this.dragActive.set(true);
  }

  protected onComposerDragLeave(event: DragEvent): void {
    const nextTarget = event.relatedTarget as Node | null;
    if (nextTarget && (event.currentTarget as HTMLElement | null)?.contains(nextTarget)) {
      return;
    }
    this.dragActive.set(false);
  }

  protected async onComposerDrop(event: DragEvent): Promise<void> {
    event.preventDefault();
    this.dragActive.set(false);
    await this.attachFiles(Array.from(event.dataTransfer?.files ?? []));
  }

  protected async onComposerPaste(event: ClipboardEvent): Promise<void> {
    const files = Array.from(event.clipboardData?.files ?? []).filter((file) =>
      file.type.startsWith('image/'),
    );
    if (files.length === 0) {
      return;
    }
    event.preventDefault();
    await this.attachFiles(files);
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
    return `${(Math.max(end - start, 0) / 1000).toFixed(1)}s`;
  }

  protected toggleTool(id: string): void {
    this.expandedToolIds.update((current) => toggleSet(current, id));
  }

  protected toolExpanded(id: string): boolean {
    return this.expandedToolIds().has(id);
  }

  protected toolStatus(tool: ToolInvocation): 'pending' | 'success' | 'error' {
    if (tool.status) {
      return tool.status;
    }
    if (tool.error) {
      return 'error';
    }
    if (tool.result) {
      return 'success';
    }
    return 'pending';
  }

  protected toolStatusLabel(tool: ToolInvocation): string {
    const status = this.toolStatus(tool);
    if (status === 'pending') {
      return 'Running';
    }
    const duration = toolRuntime(tool);
    if (status === 'error') {
      return duration ? `Error · ${duration}` : 'Error';
    }
    return duration ? `Done · ${duration}` : 'Done';
  }

  private async attachFiles(files: File[]): Promise<void> {
    const imageFiles = files.filter((file) => file.type.startsWith('image/'));
    if (imageFiles.length === 0) {
      if (files.length > 0) {
        this.showComposerNotice('Only image attachments are supported in the chat composer.');
      }
      return;
    }

    const supportedFiles = imageFiles.filter((file) => isSupportedChatImageFile(file));
    const skippedFiles = imageFiles.length - supportedFiles.length;
    if (supportedFiles.length === 0) {
      this.showComposerNotice(`Supported image formats: ${SUPPORTED_CHAT_IMAGE_LABEL}.`);
      return;
    }
    if (!this.selectedModelSupportsVision()) {
      this.showComposerNotice('The selected model does not support image input.');
      return;
    }

    try {
      await this.chat.addAttachments(supportedFiles);
      if (skippedFiles > 0) {
        this.showComposerNotice(
          `Skipped ${skippedFiles} unsupported attachment${skippedFiles === 1 ? '' : 's'}. Supported formats: ${SUPPORTED_CHAT_IMAGE_LABEL}.`,
        );
      }
    } catch (error) {
      const message =
        error instanceof Error
          ? error.message
          : `Supported image formats: ${SUPPORTED_CHAT_IMAGE_LABEL}.`;
      this.showComposerNotice(message);
    }
  }

  private showComposerNotice(message: string): void {
    this.composerNotice.set(message);
    if (this.composerNoticeTimer !== null) {
      window.clearTimeout(this.composerNoticeTimer);
    }
    this.composerNoticeTimer = window.setTimeout(() => {
      this.composerNotice.set('');
      this.composerNoticeTimer = null;
    }, 2600);
  }

  private confirmDeleteConversation(title: string): boolean {
    return window.confirm(`Delete "${title}"?`);
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

function hasImageFiles(dataTransfer: DataTransfer | null): boolean {
  if (!dataTransfer) {
    return false;
  }
  return Array.from(dataTransfer.items).some((item) => isSupportedChatImageMimeType(item.type));
}

function renderMarkdownContent(content: string): string {
  const formatted = formatInlineMarkdown(escapeHtml(content));
  const lines = formatted.split('\n');
  const blocks: string[] = [];
  let paragraph: string[] = [];
  let listType: 'ul' | 'ol' | null = null;
  let listItems: string[] = [];

  const flushParagraph = (): void => {
    if (paragraph.length === 0) {
      return;
    }
    blocks.push(`<p>${paragraph.join('<br />')}</p>`);
    paragraph = [];
  };

  const flushList = (): void => {
    if (!listType || listItems.length === 0) {
      listType = null;
      listItems = [];
      return;
    }
    blocks.push(`<${listType}>${listItems.join('')}</${listType}>`);
    listType = null;
    listItems = [];
  };

  for (const rawLine of lines) {
    const line = rawLine.trim();
    if (!line) {
      flushParagraph();
      flushList();
      continue;
    }

    const unordered = line.match(/^[-*]\s+(.*)$/);
    if (unordered) {
      flushParagraph();
      if (listType !== 'ul') {
        flushList();
        listType = 'ul';
      }
      listItems.push(`<li>${unordered[1]}</li>`);
      continue;
    }

    const ordered = line.match(/^\d+\.\s+(.*)$/);
    if (ordered) {
      flushParagraph();
      if (listType !== 'ol') {
        flushList();
        listType = 'ol';
      }
      listItems.push(`<li>${ordered[1]}</li>`);
      continue;
    }

    flushList();
    paragraph.push(line);
  }

  flushParagraph();
  flushList();

  return blocks.join('');
}

function formatInlineMarkdown(value: string): string {
  return value
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noreferrer">$1</a>')
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/\*([^*]+)\*/g, '<em>$1</em>')
    .replace(/`([^`]+)`/g, '<code>$1</code>');
}

function renderCodeContent(content: string, language?: string): string {
  const normalizedLanguage = (language ?? 'text').toLowerCase();
  const escaped = escapeHtml(content);

  if (normalizedLanguage === 'json') {
    return highlightJsonCode(escaped);
  }

  const keywords = languageKeywords(normalizedLanguage);
  if (keywords.length === 0) {
    return escaped;
  }

  return highlightSourceCode(
    escaped,
    keywords,
    normalizedLanguage === 'sh' || normalizedLanguage === 'bash',
  );
}

function languageKeywords(language: string): string[] {
  switch (language) {
    case 'ts':
    case 'tsx':
    case 'typescript':
    case 'js':
    case 'jsx':
    case 'javascript':
      return [
        'async',
        'await',
        'break',
        'case',
        'catch',
        'class',
        'const',
        'continue',
        'default',
        'else',
        'export',
        'extends',
        'false',
        'finally',
        'for',
        'function',
        'if',
        'import',
        'interface',
        'let',
        'new',
        'null',
        'return',
        'static',
        'switch',
        'throw',
        'true',
        'try',
        'type',
        'undefined',
      ];
    case 'go':
      return [
        'break',
        'case',
        'chan',
        'const',
        'continue',
        'default',
        'defer',
        'else',
        'fallthrough',
        'false',
        'for',
        'func',
        'go',
        'if',
        'import',
        'interface',
        'map',
        'package',
        'range',
        'return',
        'select',
        'struct',
        'switch',
        'true',
        'type',
        'var',
      ];
    case 'sh':
    case 'bash':
      return [
        'case',
        'do',
        'done',
        'echo',
        'elif',
        'else',
        'esac',
        'export',
        'fi',
        'for',
        'function',
        'if',
        'in',
        'local',
        'return',
        'then',
        'while',
      ];
    default:
      return [];
  }
}

function highlightSourceCode(
  escapedCode: string,
  keywords: string[],
  hashComments: boolean,
): string {
  const stashedTokens: string[] = [];
  const stash = (source: string, pattern: RegExp, className: string): string =>
    source.replace(pattern, (match) => {
      const token = `@@${stashedTokens.length}@@`;
      stashedTokens.push(`<span class="${className}">${match}</span>`);
      return token;
    });

  let highlighted = escapedCode;
  highlighted = stash(
    highlighted,
    /(`[^`]*`|&quot;(?:\\.|[\s\S])*?&quot;|&#39;(?:\\.|[\s\S])*?&#39;)/g,
    'token-string',
  );
  highlighted = stash(
    highlighted,
    hashComments ? /(\/\/.*$|#.*$)/gm : /\/\/.*$/gm,
    'token-comment',
  );
  highlighted = highlighted.replace(
    new RegExp(`\\b(${keywords.join('|')})\\b`, 'g'),
    '<span class="token-keyword">$1</span>',
  );
  highlighted = highlighted.replace(/\b(\d+(?:\.\d+)?)\b/g, '<span class="token-number">$1</span>');

  return highlighted.replace(/@@(\d+)@@/g, (_, index) => stashedTokens[Number(index)] ?? '');
}

function highlightJsonCode(escapedCode: string): string {
  return escapedCode
    .replace(/(&quot;.*?&quot;)(?=\s*:)/g, '<span class="token-keyword">$1</span>')
    .replace(/(:\s*)(&quot;.*?&quot;)/g, '$1<span class="token-string">$2</span>')
    .replace(/\b(true|false|null)\b/g, '<span class="token-boolean">$1</span>')
    .replace(/\b(\d+(?:\.\d+)?)\b/g, '<span class="token-number">$1</span>');
}

function toolRuntime(tool: ToolInvocation): string {
  if (!tool.startedAt || !tool.endedAt) {
    return '';
  }

  const duration = Math.max(
    new Date(tool.endedAt).getTime() - new Date(tool.startedAt).getTime(),
    0,
  );
  return `${(duration / 1000).toFixed(1)}s`;
}
