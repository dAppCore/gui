import { CommonModule } from '@angular/common';
import { Component, computed, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ChatService } from '../services/chat.service';

@Component({
  selector: 'settings-view',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <main class="settings-shell">
      <section class="settings-hero">
        <div>
          <p class="eyebrow">Settings</p>
          <h1>Inference defaults</h1>
          <p>
            These values seed new conversations and mirror the RFC settings panel:
            temperature, context limits, system prompt, and default model.
          </p>
        </div>

        <div class="settings-summary">
          <article>
            <span>Default model</span>
            <strong>{{ settings().defaultModel }}</strong>
          </article>
          <article>
            <span>Context window</span>
            <strong>{{ settings().contextWindow }}</strong>
          </article>
          <article>
            <span>Conversations</span>
            <strong>{{ conversationCount() }}</strong>
          </article>
        </div>
      </section>

      <section class="settings-grid">
        <article class="settings-card">
          <label>
            <span>Temperature</span>
            <input
              type="range"
              min="0"
              max="2"
              step="0.1"
              [ngModel]="settings().temperature"
              (ngModelChange)="chat.updateSettings({ temperature: +$event })"
            />
            <strong>{{ settings().temperature.toFixed(1) }}</strong>
          </label>

          <label>
            <span>Top-P</span>
            <input
              type="range"
              min="0"
              max="1"
              step="0.05"
              [ngModel]="settings().topP"
              (ngModelChange)="chat.updateSettings({ topP: +$event })"
            />
            <strong>{{ settings().topP.toFixed(2) }}</strong>
          </label>

          <label>
            <span>Top-K</span>
            <input
              type="number"
              min="0"
              max="200"
              [ngModel]="settings().topK"
              (ngModelChange)="chat.updateSettings({ topK: +$event })"
            />
          </label>

          <label>
            <span>Max tokens</span>
            <input
              type="number"
              min="64"
              max="32768"
              [ngModel]="settings().maxTokens"
              (ngModelChange)="chat.updateSettings({ maxTokens: +$event })"
            />
          </label>

          <label>
            <span>Context window</span>
            <select
              [ngModel]="settings().contextWindow"
              (ngModelChange)="chat.updateSettings({ contextWindow: +$event })"
            >
              @for (size of [2048, 4096, 8192, 16384, 32768]; track size) {
                <option [value]="size">{{ size }}</option>
              }
            </select>
          </label>
        </article>

        <article class="settings-card">
          <label class="block-field">
            <span>System prompt</span>
            <textarea
              rows="8"
              [ngModel]="settings().systemPrompt"
              (ngModelChange)="chat.updateSettings({ systemPrompt: $event })"
            ></textarea>
          </label>

          <label>
            <span>Default model</span>
            <select
              [ngModel]="settings().defaultModel"
              (ngModelChange)="chat.updateSettings({ defaultModel: $event })"
            >
              @for (model of chat.models(); track model.name) {
                <option [value]="model.name">{{ model.name }} · {{ model.architecture }}</option>
              }
            </select>
          </label>

          <div class="settings-actions">
            <button type="button" class="ghost-button" (click)="chat.resetSettings()">
              Reset defaults
            </button>
          </div>
        </article>
      </section>
    </main>
  `,
  styles: [
    `
      .settings-shell {
        min-height: calc(100vh - 2.75rem);
        padding: 2rem;
        background:
          radial-gradient(circle at top left, rgba(234, 179, 8, 0.16), transparent 24%),
          linear-gradient(180deg, #07101d 0%, #111827 100%);
        color: #f8fafc;
      }

      .settings-hero,
      .settings-summary,
      .settings-grid,
      .settings-actions {
        display: grid;
        gap: 1rem;
      }

      .settings-hero {
        grid-template-columns: minmax(0, 1.4fr) minmax(18rem, 0.8fr);
        align-items: start;
        margin-bottom: 1.5rem;
      }

      .eyebrow {
        margin: 0;
        text-transform: uppercase;
        letter-spacing: 0.16em;
        font-size: 0.72rem;
        color: rgba(248, 250, 252, 0.55);
      }

      h1 {
        margin: 0.35rem 0 0.75rem;
        font-family: 'Iowan Old Style', 'Palatino Linotype', serif;
        font-size: clamp(2rem, 4vw, 3rem);
      }

      .settings-summary {
        grid-template-columns: repeat(3, 1fr);
      }

      .settings-summary article,
      .settings-card {
        border-radius: 1.4rem;
        padding: 1.2rem;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.08);
      }

      .settings-summary span,
      label span {
        display: block;
        margin-bottom: 0.35rem;
        color: rgba(226, 232, 240, 0.74);
        font-size: 0.85rem;
      }

      .settings-grid {
        grid-template-columns: repeat(2, minmax(0, 1fr));
      }

      .settings-card {
        display: grid;
        gap: 1rem;
      }

      label,
      .block-field {
        display: grid;
        gap: 0.45rem;
      }

      input,
      select,
      textarea,
      .ghost-button {
        border-radius: 1rem;
        border: 0;
        font: inherit;
      }

      input,
      select,
      textarea {
        padding: 0.85rem 1rem;
        background: rgba(15, 23, 42, 0.86);
        color: #f8fafc;
      }

      textarea {
        resize: vertical;
      }

      .ghost-button {
        cursor: pointer;
        padding: 0.9rem 1.1rem;
        background: rgba(255, 255, 255, 0.08);
        color: inherit;
      }

      @media (max-width: 960px) {
        .settings-hero,
        .settings-grid,
        .settings-summary {
          grid-template-columns: 1fr;
        }
      }
    `,
  ],
})
export class SettingsComponent {
  protected readonly chat = inject(ChatService);
  protected readonly settings = this.chat.settings;
  protected readonly conversationCount = computed(() => this.chat.conversations().length);
}
