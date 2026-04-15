import { CommonModule } from '@angular/common';
import { Component, CUSTOM_ELEMENTS_SCHEMA, OnInit, inject } from '@angular/core';
import { ConversationSidebarComponent } from './conversation-sidebar.component';
import { InputAreaComponent } from './input-area.component';
import { MessageListComponent } from './message-list.component';
import { ModelSelectorComponent } from './model-selector.component';
import { SettingsPanelComponent } from './settings-panel.component';
import { ChatStateService } from './chat-state.service';

@Component({
  selector: 'core-chat-panel',
  standalone: true,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  imports: [
    CommonModule,
    ConversationSidebarComponent,
    InputAreaComponent,
    MessageListComponent,
    ModelSelectorComponent,
    SettingsPanelComponent,
  ],
  template: `
    <div class="workspace">
      <chat-conversation-sidebar
        [conversations]="state.filteredConversations()"
        [activeId]="state.activeConversation()?.id || ''"
        [query]="state.historyQuery()"
        (queryChange)="state.historyQuery.set($event)"
        (create)="state.startConversation()"
        (select)="state.refreshConversation($event)"
      />

      <main class="chat-shell">
        <header class="chat-shell__header">
          <div>
            <p class="eyebrow">CoreGUI Chat</p>
            <h1>{{ state.activeConversation()?.title || 'Local chat' }}</h1>
          </div>
          <div class="chat-shell__controls">
            <chat-model-selector
              [models]="state.models()"
              [value]="state.selectedModel()"
              (valueChange)="state.selectedModel.set($event)"
            />
            <button type="button" class="settings" (click)="state.settingsOpen.set(!state.settingsOpen())">Settings</button>
          </div>
        </header>

        <chat-settings-panel
          [open]="state.settingsOpen()"
          [models]="state.models()"
          [settings]="state.settings()"
          (saved)="state.saveSettings($event)"
          (closed)="state.settingsOpen.set(false)"
        />

        <section class="chat-shell__thread">
          <chat-message-list [messages]="state.activeConversation()?.messages || []" />
        </section>

        <chat-input-area
          [value]="state.draft()"
          [disabled]="state.sending()"
          (valueChange)="state.draft.set($event)"
          (submit)="state.sendMessage()"
        />
      </main>
    </div>
  `,
  styles: [
    `
      :host { display: block; min-height: 100vh; color: #f8fafc; background:
        radial-gradient(circle at top left, rgba(245, 158, 11, 0.18), transparent 30%),
        radial-gradient(circle at right, rgba(14, 165, 233, 0.16), transparent 24%),
        linear-gradient(160deg, #020617 0%, #081121 46%, #111827 100%);
        font-family: 'Iowan Old Style', 'Palatino Linotype', 'Book Antiqua', serif; }
      .workspace { min-height: 100vh; display: grid; grid-template-columns: 20rem 1fr; }
      .chat-shell { display: grid; grid-template-rows: auto auto 1fr auto; gap: 1rem; padding: 1.5rem; }
      .chat-shell__header { display: flex; justify-content: space-between; gap: 1rem; align-items: end; }
      .chat-shell__controls { display: flex; flex-wrap: wrap; gap: 0.75rem; align-items: center; }
      .chat-shell__thread { overflow: auto; padding: 1rem 0.2rem 1rem 0; }
      .eyebrow { margin: 0; color: #f59e0b; text-transform: uppercase; letter-spacing: 0.18em; font-size: 0.72rem; }
      h1 { margin: 0.2rem 0 0; font-size: clamp(2rem, 3vw, 3rem); line-height: 1; }
      .settings { border: 1px solid rgba(251, 191, 36, 0.22); border-radius: 999px; background: rgba(124, 45, 18, 0.25); color: #fde68a; padding: 0.85rem 1.2rem; cursor: pointer; }
      @media (max-width: 960px) {
        .workspace { grid-template-columns: 1fr; }
      }
    `,
  ],
})
export class ChatPanelComponent implements OnInit {
  readonly state = inject(ChatStateService);

  async ngOnInit(): Promise<void> {
    await this.state.init();
  }
}
