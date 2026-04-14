import { Injectable, computed, signal } from '@angular/core';

export interface ModelEntry {
  name: string;
  architecture: string;
  quantBits: number;
  sizeBytes: number;
  loaded: boolean;
  backend: string;
  supportsVision?: boolean;
}

export interface ChatSettings {
  temperature: number;
  topP: number;
  topK: number;
  maxTokens: number;
  contextWindow: number;
  systemPrompt: string;
  defaultModel: string;
}

export interface ImageAttachment {
  id: string;
  filename: string;
  mimeType: string;
  data: string;
  width: number;
  height: number;
}

export interface ThinkingState {
  active: boolean;
  content: string;
  startedAt?: string;
  finishedAt?: string;
}

export interface ToolInvocation {
  id: string;
  name: string;
  arguments: Record<string, unknown>;
  result: string;
  error?: string;
}

export interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  createdAt: string;
  attachments?: ImageAttachment[];
  thinking?: ThinkingState;
  toolCalls?: ToolInvocation[];
  streaming?: boolean;
}

export interface Conversation {
  id: string;
  title: string;
  model: string;
  createdAt: string;
  updatedAt: string;
  messages: ChatMessage[];
}

const STORAGE_KEY = 'core.gui.chat.state';

function defaultSettings(): ChatSettings {
  return {
    temperature: 1.0,
    topP: 0.95,
    topK: 64,
    maxTokens: 2048,
    contextWindow: 8192,
    systemPrompt: 'You are a helpful assistant.',
    defaultModel: 'lemer',
  };
}

function defaultModels(): ModelEntry[] {
  return [
    {
      name: 'lemer',
      architecture: 'gemma3',
      quantBits: 4,
      sizeBytes: 1_500_000_000,
      loaded: true,
      backend: 'metal',
      supportsVision: true,
    },
    {
      name: 'lemma',
      architecture: 'gemma3',
      quantBits: 8,
      sizeBytes: 3_200_000_000,
      loaded: false,
      backend: 'metal',
      supportsVision: true,
    },
    {
      name: 'lemmy',
      architecture: 'qwen3',
      quantBits: 4,
      sizeBytes: 1_100_000_000,
      loaded: false,
      backend: 'ollama',
      supportsVision: false,
    },
  ];
}

@Injectable({ providedIn: 'root' })
export class ChatService {
  private readonly storage = window.localStorage;
  private readonly conversationsState = signal<Conversation[]>([]);
  private readonly activeConversationIdState = signal('');
  private readonly selectedModelState = signal('lemer');
  private readonly modelsState = signal<ModelEntry[]>(defaultModels());
  private readonly modelSwitchingState = signal(false);
  private readonly settingsState = signal<ChatSettings>(defaultSettings());
  private readonly draftState = signal('');
  private readonly queuedAttachmentsState = signal<ImageAttachment[]>([]);
  private readonly busyState = signal(false);
  private modelSwitchTimer: number | null = null;

  readonly conversations = this.conversationsState.asReadonly();
  readonly activeConversationId = this.activeConversationIdState.asReadonly();
  readonly selectedModel = this.selectedModelState.asReadonly();
  readonly models = this.modelsState.asReadonly();
  readonly modelSwitching = this.modelSwitchingState.asReadonly();
  readonly settings = this.settingsState.asReadonly();
  readonly draft = this.draftState.asReadonly();
  readonly queuedAttachments = this.queuedAttachmentsState.asReadonly();
  readonly busy = this.busyState.asReadonly();
  readonly selectedModelEntry = computed<ModelEntry | null>(() => {
    const selected = this.selectedModelState();
    return this.modelsState().find((model) => model.name === selected) ?? null;
  });

  readonly activeConversation = computed<Conversation | null>(() => {
    const id = this.activeConversationIdState();
    return this.conversationsState().find((conversation) => conversation.id === id) ?? null;
  });

  readonly canSend = computed(() => {
    return this.draftState().trim().length > 0 || this.queuedAttachmentsState().length > 0;
  });

  constructor() {
    this.hydrate();
    if (!this.activeConversation()) {
      this.createConversation();
    }
  }

  setDraft(value: string): void {
    this.draftState.set(value);
  }

  updateSettings(patch: Partial<ChatSettings>): void {
    this.settingsState.update((current) => ({ ...current, ...patch }));
    if (patch.defaultModel) {
      this.selectModel(patch.defaultModel);
    }
    this.persist();
  }

  resetSettings(): void {
    const settings = defaultSettings();
    this.settingsState.set(settings);
    this.selectModel(settings.defaultModel);
    this.persist();
  }

  selectModel(name: string): void {
    const current = this.selectedModelState();
    if (name !== current) {
      this.modelSwitchingState.set(true);
      if (this.modelSwitchTimer !== null) {
        window.clearTimeout(this.modelSwitchTimer);
      }
      this.modelSwitchTimer = window.setTimeout(() => {
        this.modelSwitchingState.set(false);
        this.modelSwitchTimer = null;
      }, 420);
    }
    this.selectedModelState.set(name);
    this.modelsState.update((models) =>
      models.map((model) => ({ ...model, loaded: model.name === name })),
    );
    this.persist();
  }

  createConversation(): void {
    const now = new Date().toISOString();
    const id = crypto.randomUUID();
    const conversation: Conversation = {
      id,
      title: 'New conversation',
      model: this.selectedModelState(),
      createdAt: now,
      updatedAt: now,
      messages: [],
    };
    this.conversationsState.update((items) => [conversation, ...items]);
    this.activeConversationIdState.set(id);
    this.persist();
  }

  selectConversation(id: string): void {
    this.activeConversationIdState.set(id);
  }

  renameConversation(id: string, title: string): void {
    const cleanTitle = title.trim() || 'Untitled conversation';
    this.updateConversation(id, (conversation) => ({ ...conversation, title: cleanTitle }));
  }

  deleteConversation(id: string): void {
    this.conversationsState.update((items) => items.filter((conversation) => conversation.id !== id));
    if (this.activeConversationIdState() === id) {
      const nextConversation = this.conversationsState()[0];
      if (nextConversation) {
        this.activeConversationIdState.set(nextConversation.id);
      } else {
        this.createConversation();
      }
    }
    this.persist();
  }

  clearActiveConversation(): void {
    const conversation = this.activeConversation();
    if (!conversation) {
      return;
    }
    this.updateConversation(conversation.id, (current) => ({
      ...current,
      title: 'New conversation',
      messages: [],
      updatedAt: new Date().toISOString(),
    }));
  }

  async addAttachment(file: File): Promise<void> {
    const dataUrl = await readAsDataURL(file);
    const dimensions = await readImageSize(dataUrl);
    const attachment: ImageAttachment = {
      id: crypto.randomUUID(),
      filename: file.name,
      mimeType: file.type || 'application/octet-stream',
      data: dataUrl,
      width: dimensions.width,
      height: dimensions.height,
    };
    this.queuedAttachmentsState.update((items) => [...items, attachment]);
  }

  async addAttachments(files: Iterable<File>): Promise<void> {
    for (const file of files) {
      await this.addAttachment(file);
    }
  }

  removeAttachment(id: string): void {
    this.queuedAttachmentsState.update((items) => items.filter((attachment) => attachment.id !== id));
  }

  async sendMessage(): Promise<void> {
    const conversation = this.activeConversation();
    const content = this.draftState().trim();
    if (!conversation || (!content && this.queuedAttachmentsState().length === 0) || this.busyState()) {
      return;
    }

    const now = new Date().toISOString();
    const userMessage: ChatMessage = {
      id: crypto.randomUUID(),
      role: 'user',
      content,
      createdAt: now,
      attachments: this.queuedAttachmentsState(),
    };

    const assistantMessage: ChatMessage = {
      id: crypto.randomUUID(),
      role: 'assistant',
      content: '',
      createdAt: new Date(Date.now() + 320).toISOString(),
      streaming: true,
      thinking: {
        active: true,
        content: `Inspecting the request through ${this.selectedModelState()}...`,
        startedAt: now,
      },
      toolCalls: inferToolCalls(content),
    };

    const title = conversation.messages.length === 0 ? deriveTitle(content) : conversation.title;
    const updatedConversation: Conversation = {
      ...conversation,
      title,
      model: this.selectedModelState(),
      updatedAt: now,
      messages: [...conversation.messages, userMessage, assistantMessage],
    };

    this.replaceConversation(updatedConversation);
    this.draftState.set('');
    this.queuedAttachmentsState.set([]);
    this.busyState.set(true);

    const response = buildResponse(content, this.selectedModelState(), this.settingsState());
    await streamIntoMessage(response, (fragment, done) => {
      this.updateMessage(updatedConversation.id, assistantMessage.id, (message) => ({
        ...message,
        content: fragment,
        streaming: !done,
        thinking: message.thinking
          ? {
              ...message.thinking,
              active: !done,
              finishedAt: done ? new Date().toISOString() : undefined,
            }
          : undefined,
      }));
      if (done) {
        this.busyState.set(false);
      }
    });
  }

  exportConversation(conversation: Conversation): string {
    return conversation.messages
      .map((message) => {
        const heading = message.role === 'user' ? '## User' : '## Assistant';
        const attachments = (message.attachments ?? [])
          .map((attachment) => `- ${attachment.filename} (${attachment.mimeType})`)
          .join('\n');
        return [heading, message.content, attachments].filter(Boolean).join('\n\n');
      })
      .join('\n\n---\n\n');
  }

  private hydrate(): void {
    try {
      const raw = this.storage.getItem(STORAGE_KEY);
      if (!raw) {
        return;
      }
      const parsed = JSON.parse(raw) as {
        conversations?: Conversation[];
        activeConversationId?: string;
        selectedModel?: string;
        models?: ModelEntry[];
        settings?: ChatSettings;
      };
      this.conversationsState.set(parsed.conversations ?? []);
      this.activeConversationIdState.set(parsed.activeConversationId ?? '');
      this.selectedModelState.set(parsed.selectedModel ?? 'lemer');
      this.modelsState.set(parsed.models ?? defaultModels());
      this.settingsState.set(parsed.settings ?? defaultSettings());
    } catch {
      this.storage.removeItem(STORAGE_KEY);
    }
  }

  private persist(): void {
    this.storage.setItem(
      STORAGE_KEY,
      JSON.stringify({
        conversations: this.conversationsState(),
        activeConversationId: this.activeConversationIdState(),
        selectedModel: this.selectedModelState(),
        models: this.modelsState(),
        settings: this.settingsState(),
      }),
    );
  }

  private replaceConversation(conversation: Conversation): void {
    this.conversationsState.update((items) =>
      items
        .map((item) => (item.id === conversation.id ? conversation : item))
        .sort((left, right) => right.updatedAt.localeCompare(left.updatedAt)),
    );
    this.persist();
  }

  private updateConversation(id: string, updater: (conversation: Conversation) => Conversation): void {
    const current = this.conversationsState().find((conversation) => conversation.id === id);
    if (!current) {
      return;
    }
    const updated = updater(current);
    this.replaceConversation(updated);
  }

  private updateMessage(
    conversationId: string,
    messageId: string,
    updater: (message: ChatMessage) => ChatMessage,
  ): void {
    this.updateConversation(conversationId, (conversation) => ({
      ...conversation,
      updatedAt: new Date().toISOString(),
      messages: conversation.messages.map((message) =>
        message.id === messageId ? updater(message) : message,
      ),
    }));
  }
}

function deriveTitle(content: string): string {
  const trimmed = content.trim();
  if (!trimmed) {
    return 'New conversation';
  }
  return trimmed.length > 52 ? `${trimmed.slice(0, 52).trim()}...` : trimmed;
}

function inferToolCalls(content: string): ToolInvocation[] {
  const normalized = content.toLowerCase();
  if (!normalized.includes('tool')) {
    return [];
  }
  return [
    {
      id: crypto.randomUUID(),
      name: 'gui.store.search',
      arguments: { q: content.slice(0, 40) },
      result: 'Local tool manifest execution is simulated in the chat shell.',
    },
  ];
}

function buildResponse(content: string, model: string, settings: ChatSettings): string {
  const prompt = content.trim() || 'your multimodal prompt';
  return [
    `Using ${model} with temperature ${settings.temperature.toFixed(1)}.`,
    '',
    `I stored this exchange locally and can keep the conversation context across sessions.`,
    '',
    `Prompt summary: ${prompt}`,
    '',
    '```ts',
    `const response = await window.core.ml.generate(${JSON.stringify(prompt)});`,
    '```',
    '',
    'The actual inference bridge is not present in this workspace, so this shell demonstrates the RFC flow with local state, progressive rendering, tool-call blocks, and multimodal attachments.',
  ].join('\n');
}

async function streamIntoMessage(
  content: string,
  onUpdate: (fragment: string, done: boolean) => void,
): Promise<void> {
  const chunks = content.match(/.{1,18}/g) ?? [content];
  let assembled = '';
  for (const chunk of chunks) {
    assembled += chunk;
    onUpdate(assembled, false);
    await new Promise((resolve) => window.setTimeout(resolve, 45));
  }
  onUpdate(assembled, true);
}

function readAsDataURL(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onerror = () => reject(reader.error);
    reader.onload = () => resolve(String(reader.result));
    reader.readAsDataURL(file);
  });
}

function readImageSize(dataUrl: string): Promise<{ width: number; height: number }> {
  return new Promise((resolve) => {
    const img = new Image();
    img.onload = () => resolve({ width: img.width, height: img.height });
    img.onerror = () => resolve({ width: 0, height: 0 });
    img.src = dataUrl;
  });
}
