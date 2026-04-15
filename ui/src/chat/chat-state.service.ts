import { Injectable, computed, effect, inject, signal } from '@angular/core';
import { WebSocketService } from '../services/websocket.service';
import { ChatMessage, ChatSettings, Conversation, ConversationSummary, ModelEntry } from './chat.types';

declare global {
  interface Window {
    __CORE_GUI_INVOKE__?: (route: string, payload?: unknown) => Promise<unknown> | unknown;
  }
}

@Injectable({ providedIn: 'root' })
export class ChatStateService {
  private readonly ws = inject(WebSocketService);

  readonly conversations = signal<ConversationSummary[]>([]);
  readonly activeConversation = signal<Conversation | null>(null);
  readonly models = signal<ModelEntry[]>([]);
  readonly settings = signal<ChatSettings>({
    temperature: 1,
    top_p: 0.95,
    top_k: 64,
    max_tokens: 2048,
    context_window: 8192,
    system_prompt: 'You are a helpful assistant.',
    default_model: 'local-default',
  });
  readonly draft = signal('');
  readonly historyQuery = signal('');
  readonly settingsOpen = signal(false);
  readonly sending = signal(false);
  readonly selectedModel = signal('');
  readonly filteredConversations = computed(() => {
    const needle = this.historyQuery().trim().toLowerCase();
    if (!needle) {
      return this.conversations();
    }
    return this.conversations().filter((item) => item.title.toLowerCase().includes(needle));
  });

  constructor() {
    effect(() => {
      const current = this.activeConversation();
      if (current) {
        this.selectedModel.set(current.model || this.settings().default_model);
      }
    });
  }

  async init(): Promise<void> {
    this.ws.connect();
    this.bindEvents();
    await this.loadBootstrap();
  }

  async refreshConversation(id: string): Promise<void> {
    const conversation = await this.invoke<Conversation>('gui.chat.conversations.get', { id });
    if (conversation) {
      this.activeConversation.set(conversation);
    }
  }

  async startConversation(): Promise<void> {
    const conversation = await this.invoke<Conversation>('gui.chat.conversations.new');
    if (conversation) {
      this.activeConversation.set(conversation);
      this.conversations.update((items) => [conversation, ...items.filter((item) => item.id !== conversation.id)]);
    }
  }

  async deleteConversation(id: string): Promise<void> {
    await this.invoke('gui.chat.conversations.delete', { id });
    this.conversations.update((items) => items.filter((item) => item.id !== id));
    if (this.activeConversation()?.id === id) {
      this.activeConversation.set(null);
    }
  }

  async renameConversation(id: string, title: string): Promise<void> {
    const updated = await this.invoke<Conversation>('gui.chat.conversations.rename', { id, title });
    if (updated) {
      this.mergeConversation(updated);
    }
  }

  async saveSettings(settings: ChatSettings): Promise<void> {
    const saved = await this.invoke<ChatSettings>('gui.chat.settings.save', settings);
    if (saved) {
      this.settings.set(saved);
      if (saved.default_model) {
        this.selectedModel.set(saved.default_model);
      }
    }
  }

  async sendMessage(): Promise<void> {
    const content = this.draft().trim();
    if (!content && !(this.activeConversation()?.messages?.length)) {
      return;
    }
    this.sending.set(true);
    try {
      const response = await this.invoke<Conversation>('gui.chat.send', {
        conversation_id: this.activeConversation()?.id,
        content,
      });
      if (response) {
        this.activeConversation.set(response);
        this.mergeConversation(response);
        this.draft.set('');
      }
    } finally {
      this.sending.set(false);
    }
  }

  private async loadBootstrap(): Promise<void> {
    const [models, settings, conversations] = await Promise.all([
      this.invoke<ModelEntry[]>('gui.chat.models'),
      this.invoke<ChatSettings>('gui.chat.settings.load'),
      this.invoke<ConversationSummary[]>('gui.chat.conversations.list'),
    ]);

    if (models?.length) {
      this.models.set(models);
      const current = models.find((item) => item.loaded) ?? models[0];
      if (current) {
        this.selectedModel.set(current.name);
      }
    } else {
      this.models.set([
        { name: 'lemer', architecture: 'gemma3', quant_bits: 4, size_bytes: 1_500_000_000, loaded: true, backend: 'metal' },
      ]);
      this.selectedModel.set('lemer');
    }

    if (settings) {
      this.settings.set(settings);
      if (settings.default_model) {
        this.selectedModel.set(settings.default_model);
      }
    }

    if (conversations?.length) {
      this.conversations.set(conversations);
      await this.refreshConversation(conversations[0].id);
    } else {
      await this.startConversation();
    }
  }

  private bindEvents(): void {
    this.ws.on('chat.conversation', (payload) => {
      const data = payload as { conversation?: Conversation; conversation_id?: string };
      if (data.conversation) {
        this.mergeConversation(data.conversation);
      }
      if (data.conversation_id && this.activeConversation()?.id === data.conversation_id) {
        this.activeConversation.set(null);
      }
    });
    this.ws.on('chat.message', (payload) => {
      const data = payload as { conversation_id: string; message: ChatMessage };
      if (!data?.message) {
        return;
      }
      if (this.activeConversation()?.id === data.conversation_id) {
        this.activeConversation.update((conversation) =>
          conversation
            ? {
                ...conversation,
                messages: [...conversation.messages, data.message],
                updated_at: data.message.created_at,
              }
            : conversation,
        );
      }
    });
    this.ws.on('chat.token', (payload) => {
      const data = payload as { conversation_id: string; message_id: string; content: string };
      if (this.activeConversation()?.id !== data.conversation_id) {
        return;
      }
      this.activeConversation.update((conversation) => {
        if (!conversation) {
          return conversation;
        }
        const messages = conversation.messages.map((message) =>
          message.id === data.message_id ? { ...message, content: `${message.content ?? ''}${data.content ?? ''}` } : message,
        );
        return { ...conversation, messages };
      });
    });
  }

  private mergeConversation(conversation: Conversation): void {
    const summary: ConversationSummary = {
      id: conversation.id,
      title: conversation.title,
      model: conversation.model,
      created_at: conversation.created_at,
      updated_at: conversation.updated_at,
      message_count: conversation.messages?.length ?? conversation.message_count ?? 0,
    };
    this.conversations.update((items) => [summary, ...items.filter((item) => item.id !== summary.id)]);
    if (this.activeConversation()?.id === conversation.id || !this.activeConversation()) {
      this.activeConversation.set(conversation);
    }
  }

  private async invoke<T>(route: string, payload?: unknown): Promise<T> {
    if (typeof window.__CORE_GUI_INVOKE__ === 'function') {
      return (await window.__CORE_GUI_INVOKE__(route, payload)) as T;
    }
    return this.mockInvoke<T>(route, payload);
  }

  private async mockInvoke<T>(route: string, payload?: unknown): Promise<T> {
    if (route === 'gui.chat.models') {
      return [
        { name: 'lemer', architecture: 'gemma3', quant_bits: 4, size_bytes: 1500000000, loaded: true, backend: 'metal' },
        { name: 'lemma', architecture: 'qwen3', quant_bits: 8, size_bytes: 3200000000, loaded: false, backend: 'ollama' },
      ] as T;
    }
    if (route === 'gui.chat.settings.load' || route === 'gui.chat.settings.save') {
      return (payload ?? this.settings()) as T;
    }
    if (route === 'gui.chat.conversations.list') {
      return [] as T;
    }
    if (route === 'gui.chat.conversations.new') {
      const now = new Date().toISOString();
      return {
        id: `conv-${Date.now().toString(36)}`,
        title: 'New Chat',
        model: this.selectedModel() || 'lemer',
        created_at: now,
        updated_at: now,
        message_count: 0,
        messages: [],
      } as T;
    }
    if (route === 'gui.chat.send') {
      const conversation = this.activeConversation() ?? ((await this.mockInvoke('gui.chat.conversations.new')) as Conversation);
      const content = (payload as { content?: string })?.content ?? '';
      const now = new Date().toISOString();
      const response: Conversation = {
        ...conversation,
        updated_at: now,
        title: conversation.title === 'New Chat' ? content.slice(0, 48) || 'New Chat' : conversation.title,
        messages: [
          ...conversation.messages,
          { id: `user-${Date.now()}`, role: 'user', content, created_at: now, model: this.selectedModel() },
          {
            id: `assistant-${Date.now() + 1}`,
            role: 'assistant',
            content: `Local mock response for: ${content}`,
            created_at: now,
            model: this.selectedModel(),
          },
        ],
      };
      return response as T;
    }
    return {} as T;
  }
}
