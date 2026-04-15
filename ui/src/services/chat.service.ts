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
  status?: 'pending' | 'success' | 'error';
  startedAt?: string;
  endedAt?: string;
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
  settings?: ChatSettings;
}

interface RemoteModelEntry {
  name: string;
  architecture: string;
  quant_bits: number;
  size_bytes: number;
  loaded: boolean;
  backend: string;
  supports_vision?: boolean;
}

interface RemoteChatSettings {
  temperature: number;
  top_p: number;
  top_k: number;
  max_tokens: number;
  context_window: number;
  system_prompt: string;
  default_model: string;
}

interface RemoteImageAttachment {
  id?: string;
  filename: string;
  mime_type: string;
  data: string;
  width: number;
  height: number;
}

interface RemoteThinkingState {
  active: boolean;
  content: string;
  started_at?: string;
  finished_at?: string;
}

interface RemoteToolCall {
  id: string;
  name: string;
  arguments?: Record<string, unknown>;
}

interface RemoteToolResult {
  tool_call_id: string;
  content: string;
}

interface RemoteToolInvocation {
  call: RemoteToolCall;
  result: RemoteToolResult;
  started_at?: string;
  ended_at?: string;
  error?: string;
}

interface RemoteChatMessage {
  id: string;
  role: string;
  content: string;
  created_at: string;
  attachments?: RemoteImageAttachment[];
  thinking?: RemoteThinkingState;
  tool_calls?: RemoteToolInvocation[];
  streaming?: boolean;
  finish_reason?: string;
}

interface RemoteConversation {
  id: string;
  title: string;
  model: string;
  created_at: string;
  updated_at: string;
  messages: RemoteChatMessage[];
  settings?: RemoteChatSettings;
}

interface RemoteChatSnapshot {
  settings: RemoteChatSettings;
  selected_model: string;
  conversations: Record<string, RemoteConversation>;
  queued_images: Record<string, RemoteImageAttachment[]>;
  thinking: Record<string, RemoteThinkingState>;
  streaming_message: Record<string, string>;
  models: RemoteModelEntry[];
}

const STORAGE_KEY = 'core.gui.chat.state';
export const SUPPORTED_CHAT_IMAGE_MIME_TYPES = [
  'image/png',
  'image/jpeg',
  'image/webp',
  'image/gif',
] as const;
export const SUPPORTED_CHAT_IMAGE_LABEL = 'PNG, JPEG, WebP, or GIF';
export const SUPPORTED_CHAT_IMAGE_ACCEPT = SUPPORTED_CHAT_IMAGE_MIME_TYPES.join(',');

interface ChatToolDefinition {
  name: string;
  description: string;
  parameters: Record<string, string>;
}

const REGISTERED_CHAT_TOOLS: ChatToolDefinition[] = [
  {
    name: 'gui.chat.settings.load',
    description: 'Load the persisted inference defaults for the chat shell.',
    parameters: {},
  },
  {
    name: 'gui.chat.models',
    description: 'List the locally available chat models and which model is selected.',
    parameters: {},
  },
  {
    name: 'gui.chat.conversations.search',
    description:
      'Search saved conversation history by title, content, tool calls, and attachments.',
    parameters: {
      q: 'Search string',
    },
  },
  {
    name: 'gui.route.store',
    description: 'Search the local CoreGUI store surface for matching chat data.',
    parameters: {
      q: 'Search string',
    },
  },
];

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

export function normaliseChatImageMimeType(mimeType: string): string {
  return mimeType.trim().toLowerCase();
}

function inferChatImageMimeType(fileName: string): string {
  const normalized = fileName.trim().toLowerCase();
  if (normalized.endsWith('.png')) {
    return 'image/png';
  }
  if (normalized.endsWith('.jpg') || normalized.endsWith('.jpeg')) {
    return 'image/jpeg';
  }
  if (normalized.endsWith('.webp')) {
    return 'image/webp';
  }
  if (normalized.endsWith('.gif')) {
    return 'image/gif';
  }
  return '';
}

function resolveChatImageMimeType(file: Pick<File, 'name' | 'type'>): string {
  const mimeType = normaliseChatImageMimeType(file.type ?? '');
  if (mimeType) {
    return mimeType;
  }
  return inferChatImageMimeType(file.name ?? '');
}

export function isSupportedChatImageMimeType(mimeType: string): boolean {
  return SUPPORTED_CHAT_IMAGE_MIME_TYPES.includes(
    normaliseChatImageMimeType(mimeType) as (typeof SUPPORTED_CHAT_IMAGE_MIME_TYPES)[number],
  );
}

export function isSupportedChatImageFile(file: Pick<File, 'name' | 'type'>): boolean {
  return isSupportedChatImageMimeType(resolveChatImageMimeType(file));
}

export function conversationSearchText(conversation: Conversation): string {
  return [
    conversation.title,
    conversation.model,
    ...conversation.messages.flatMap((message) => [
      message.role,
      message.content,
      message.thinking?.content ?? '',
      ...(message.attachments ?? []).flatMap((attachment) => [
        attachment.filename,
        attachment.mimeType,
      ]),
      ...(message.toolCalls ?? []).flatMap((toolCall) => [
        toolCall.name,
        JSON.stringify(toolCall.arguments ?? {}),
        toolCall.result,
        toolCall.error ?? '',
      ]),
    ]),
  ]
    .join(' ')
    .toLowerCase();
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
  private readonly queuedAttachmentsByConversationState = signal<Record<string, ImageAttachment[]>>(
    {},
  );
  private readonly busyState = signal(false);
  private bridgeMode = false;
  private modelSwitchTimer: number | null = null;

  readonly conversations = this.conversationsState.asReadonly();
  readonly activeConversationId = this.activeConversationIdState.asReadonly();
  readonly selectedModel = this.selectedModelState.asReadonly();
  readonly models = this.modelsState.asReadonly();
  readonly modelSwitching = this.modelSwitchingState.asReadonly();
  readonly settings = this.settingsState.asReadonly();
  readonly draft = this.draftState.asReadonly();
  readonly queuedAttachments = computed(() => {
    const conversationID = this.activeConversationIdState();
    return [...(this.queuedAttachmentsByConversationState()[conversationID] ?? [])];
  });
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
    return this.draftState().trim().length > 0 || this.queuedAttachments().length > 0;
  });

  constructor() {
    void this.bootstrap();
  }

  setDraft(value: string): void {
    this.draftState.set(value);
  }

  async updateSettings(patch: Partial<ChatSettings>): Promise<void> {
    const nextSettings = { ...this.settingsState(), ...patch };
    this.settingsState.set(nextSettings);

    if (this.bridgeMode) {
      const remoteSettings = await this.safeBridgeCall<RemoteChatSettings>(
        'chat:settings-save',
        toRemoteSettings(nextSettings),
      );
      if (remoteSettings) {
        this.settingsState.set(fromRemoteSettings(remoteSettings));
      }
    } else {
      this.persist();
    }

    if (patch.defaultModel) {
      await this.selectModel(patch.defaultModel);
    }
  }

  async resetSettings(): Promise<void> {
    const settings = defaultSettings();
    if (this.bridgeMode) {
      const remoteSettings = await this.safeBridgeCall<RemoteChatSettings>('chat:settings-reset');
      this.settingsState.set(remoteSettings ? fromRemoteSettings(remoteSettings) : settings);
    } else {
      this.settingsState.set(settings);
      this.persist();
    }
    await this.selectModel(this.settingsState().defaultModel);
  }

  async selectModel(name: string): Promise<void> {
    if (!this.modelsState().some((model) => model.name === name)) {
      return;
    }

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

    if (this.bridgeMode) {
      const remoteModels = await this.safeBridgeCall<RemoteModelEntry[]>('chat:model-select', {
        model: name,
      });
      if (remoteModels) {
        this.modelsState.set(remoteModels.map(fromRemoteModel));
        this.selectedModelState.set(resolveSelectedModel(this.modelsState(), name));
      }
      return;
    }

    this.persist();
  }

  async createConversation(): Promise<void> {
    if (this.bridgeMode) {
      const remoteConversation =
        await this.safeBridgeCall<RemoteConversation>('chat:conversation-new');
      if (remoteConversation) {
        this.applyConversation(fromRemoteConversation(remoteConversation), true);
        return;
      }
    }

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
    this.applyConversation(conversation, true);
  }

  selectConversation(id: string): void {
    if (!this.conversationsState().some((conversation) => conversation.id === id)) {
      return;
    }
    this.activeConversationIdState.set(id);
    this.persist();
  }

  async renameConversation(id: string, title: string): Promise<void> {
    const cleanTitle = title.trim() || 'Untitled conversation';
    if (this.bridgeMode) {
      const remoteConversation = await this.safeBridgeCall<RemoteConversation>(
        'chat:conversation-rename',
        { id, title: cleanTitle },
      );
      if (remoteConversation) {
        this.applyConversation(fromRemoteConversation(remoteConversation));
        return;
      }
    }

    this.updateConversation(id, (conversation) => ({ ...conversation, title: cleanTitle }));
  }

  async deleteConversation(id: string): Promise<void> {
    if (this.bridgeMode) {
      await this.safeBridgeCall('chat:conversation-delete', { id });
    }

    this.conversationsState.update((items) =>
      items.filter((conversation) => conversation.id !== id),
    );
    this.queuedAttachmentsByConversationState.update((current) => {
      const next = { ...current };
      delete next[id];
      return next;
    });

    if (this.activeConversationIdState() === id) {
      this.activeConversationIdState.set(this.conversationsState()[0]?.id ?? '');
    }

    if (!this.activeConversation()) {
      await this.createConversation();
    }
    this.persist();
  }

  async clearActiveConversation(): Promise<void> {
    const conversation = this.activeConversation();
    if (!conversation) {
      return;
    }

    if (this.bridgeMode) {
      const remoteConversation = await this.safeBridgeCall<RemoteConversation>('chat:clear', {
        conversation_id: conversation.id,
      });
      if (remoteConversation) {
        this.applyConversation(fromRemoteConversation(remoteConversation));
      }
    } else {
      this.updateConversation(conversation.id, (current) => ({
        ...current,
        title: 'New conversation',
        messages: [],
        updatedAt: new Date().toISOString(),
      }));
    }

    this.updateQueuedAttachments(conversation.id, () => []);
  }

  async addAttachment(file: File): Promise<void> {
    const conversation = this.activeConversation();
    if (!conversation) {
      return;
    }

    const mimeType = resolveChatImageMimeType(file);
    if (!isSupportedChatImageMimeType(mimeType)) {
      throw new Error(`Supported image formats: ${SUPPORTED_CHAT_IMAGE_LABEL}.`);
    }

    const dataUrl = await readAsDataURL(file);
    const dimensions = await readImageSize(dataUrl);
    const attachment: ImageAttachment = {
      id: crypto.randomUUID(),
      filename: file.name,
      mimeType,
      data: dataUrl,
      width: dimensions.width,
      height: dimensions.height,
    };

    if (this.bridgeMode) {
      const remoteAttachments = await this.safeBridgeCall<RemoteImageAttachment[]>(
        'chat:attach-image',
        {
          conversation_id: conversation.id,
          attachment: toRemoteAttachment(attachment),
        },
      );
      if (remoteAttachments) {
        this.setQueuedAttachments(conversation.id, remoteAttachments.map(fromRemoteAttachment));
        return;
      }
    }

    this.updateQueuedAttachments(conversation.id, (items) => [...items, attachment]);
  }

  async addAttachments(files: Iterable<File>): Promise<void> {
    for (const file of files) {
      await this.addAttachment(file);
    }
  }

  async removeAttachment(id: string): Promise<void> {
    const conversation = this.activeConversation();
    if (!conversation) {
      return;
    }

    if (this.bridgeMode) {
      const remoteAttachments = await this.safeBridgeCall<RemoteImageAttachment[]>(
        'chat:detach-image',
        {
          conversation_id: conversation.id,
          attachment_id: id,
        },
      );
      if (remoteAttachments) {
        this.setQueuedAttachments(conversation.id, remoteAttachments.map(fromRemoteAttachment));
        return;
      }
    }

    this.updateQueuedAttachments(conversation.id, (items) =>
      items.filter((attachment) => attachment.id !== id),
    );
  }

  async sendMessage(): Promise<void> {
    const conversation = this.activeConversation();
    const content = this.draftState().trim();
    const attachments = this.queuedAttachments();
    if (!conversation || (!content && attachments.length === 0) || this.busyState()) {
      return;
    }

    const now = new Date().toISOString();
    const toolCalls = inferToolCalls(content);
    const userMessage: ChatMessage = {
      id: crypto.randomUUID(),
      role: 'user',
      content,
      createdAt: now,
      attachments,
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
      toolCalls,
    };

    const title =
      conversation.messages.length === 0
        ? deriveTitleFromMessage(content, attachments)
        : conversation.title;
    const updatedConversation: Conversation = {
      ...conversation,
      title,
      model: this.selectedModelState(),
      updatedAt: now,
      messages: [...conversation.messages, userMessage, assistantMessage],
    };

    this.replaceConversation(updatedConversation);
    this.draftState.set('');
    this.updateQueuedAttachments(conversation.id, () => []);
    this.busyState.set(true);

    try {
      if (this.bridgeMode) {
        await this.safeBridgeCall('chat:thinking-start', {
          conversation_id: updatedConversation.id,
        });
        await this.safeBridgeCall('chat:thinking-append', {
          conversation_id: updatedConversation.id,
          content: assistantMessage.thinking?.content ?? '',
        });
        await this.safeBridgeCall('chat:send', {
          conversation_id: updatedConversation.id,
          content,
        });
        await this.safeBridgeCall('chat:stream-start', {
          conversation_id: updatedConversation.id,
        });
      }

      const toolOutputs = await this.runToolCalls(
        updatedConversation.id,
        assistantMessage.id,
        content,
      );
      const response = await generateAssistantResponse(
        content,
        this.selectedModelState(),
        this.settingsState(),
        toolOutputs,
      );
      let previousFragment = '';
      await streamIntoMessage(response, async (fragment, done) => {
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
        if (!this.bridgeMode) {
          return;
        }
        const delta = fragment.slice(previousFragment.length);
        previousFragment = fragment;
        if (delta) {
          await this.safeBridgeCall('chat:stream-append', {
            conversation_id: updatedConversation.id,
            content: delta,
          });
        }
        if (done) {
          await this.safeBridgeCall('chat:thinking-end', {
            conversation_id: updatedConversation.id,
          });
          await this.safeBridgeCall('chat:stream-finish', {
            conversation_id: updatedConversation.id,
            finish_reason: 'stop',
          });
        }
      });
      if (this.bridgeMode) {
        await this.refreshConversationFromBridge(updatedConversation.id);
      }
    } finally {
      this.busyState.set(false);
    }
  }

  async exportConversation(conversation: Conversation): Promise<string> {
    if (this.bridgeMode) {
      const exported = await this.safeBridgeCall<string>('chat:conversation-export', {
        id: conversation.id,
      });
      if (typeof exported === 'string' && exported.trim()) {
        return exported;
      }
    }
    return buildConversationMarkdown(conversation);
  }

  private async bootstrap(): Promise<void> {
    if (await this.tryBridgeBootstrap()) {
      if (!this.activeConversation()) {
        await this.createConversation();
      }
      return;
    }

    this.hydrate();
    if (!this.activeConversation()) {
      await this.createConversation();
    }
  }

  private async tryBridgeBootstrap(): Promise<boolean> {
    if (!hasCoreGUIBridge()) {
      return false;
    }

    try {
      const snapshot = await this.callBridge<RemoteChatSnapshot>('chat:snapshot');
      this.bridgeMode = true;
      this.loadRemoteSnapshot(snapshot);
      return true;
    } catch {
      this.bridgeMode = false;
      return false;
    }
  }

  private async callBridge<T>(action: string, payload?: unknown): Promise<T> {
    const bridge = window.__coreGUIBridge;
    const invoke = bridge?.invoke ?? bridge?.call;
    if (typeof invoke !== 'function') {
      throw new Error(`CoreGUI bridge is unavailable for ${action}.`);
    }
    return (await Promise.resolve(invoke.call(bridge, action, payload))) as T;
  }

  private async safeBridgeCall<T>(action: string, payload?: unknown): Promise<T | null> {
    if (!this.bridgeMode) {
      return null;
    }
    try {
      return await this.callBridge<T>(action, payload);
    } catch (error) {
      console.error(`CoreGUI bridge action failed: ${action}`, error);
      return null;
    }
  }

  private loadRemoteSnapshot(snapshot: RemoteChatSnapshot): void {
    this.settingsState.set(fromRemoteSettings(snapshot.settings));
    this.modelsState.set((snapshot.models ?? []).map(fromRemoteModel));
    this.selectedModelState.set(resolveSelectedModel(this.modelsState(), snapshot.selected_model));
    this.conversationsState.set(
      sortConversations(Object.values(snapshot.conversations ?? {}).map(fromRemoteConversation)),
    );
    this.queuedAttachmentsByConversationState.set(
      Object.fromEntries(
        Object.entries(snapshot.queued_images ?? {}).map(([conversationID, attachments]) => [
          conversationID,
          attachments.map(fromRemoteAttachment),
        ]),
      ),
    );
    this.activeConversationIdState.set(this.resolveActiveConversationID(this.conversationsState()));
  }

  private resolveActiveConversationID(conversations: Conversation[]): string {
    const current = this.activeConversationIdState();
    if (current && conversations.some((conversation) => conversation.id === current)) {
      return current;
    }
    return conversations[0]?.id ?? '';
  }

  private applyConversation(conversation: Conversation, select = false): void {
    this.conversationsState.update((items) =>
      sortConversations([
        ...items.filter((item) => item.id !== conversation.id),
        this.normaliseConversation(conversation),
      ]),
    );
    if (select || !this.activeConversationIdState()) {
      this.activeConversationIdState.set(conversation.id);
    }
    this.persist();
  }

  private setQueuedAttachments(conversationID: string, attachments: ImageAttachment[]): void {
    this.queuedAttachmentsByConversationState.update((current) => {
      const next = { ...current };
      if (attachments.length === 0) {
        delete next[conversationID];
      } else {
        next[conversationID] = attachments;
      }
      return next;
    });
    this.persist();
  }

  private updateQueuedAttachments(
    conversationID: string,
    updater: (items: ImageAttachment[]) => ImageAttachment[],
  ): void {
    const current = [...(this.queuedAttachmentsByConversationState()[conversationID] ?? [])];
    this.setQueuedAttachments(conversationID, updater(current));
  }

  private async refreshConversationFromBridge(conversationID: string): Promise<void> {
    const remoteConversation = await this.safeBridgeCall<RemoteConversation>(
      'chat:conversation-get',
      {
        id: conversationID,
      },
    );
    if (remoteConversation) {
      this.applyConversation(fromRemoteConversation(remoteConversation));
    }
  }

  private hydrate(): void {
    if (this.bridgeMode) {
      return;
    }

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
        queuedAttachmentsByConversation?: Record<string, ImageAttachment[]>;
        queuedAttachments?: ImageAttachment[];
      };
      this.conversationsState.set(
        sortConversations(
          (parsed.conversations ?? []).map((conversation) =>
            this.normaliseConversation(conversation),
          ),
        ),
      );
      this.activeConversationIdState.set(parsed.activeConversationId ?? '');
      this.selectedModelState.set(parsed.selectedModel ?? 'lemer');
      this.modelsState.set(parsed.models ?? defaultModels());
      this.settingsState.set(parsed.settings ?? defaultSettings());

      const queuedAttachmentsByConversation = parsed.queuedAttachmentsByConversation ?? {};
      if (
        Object.keys(queuedAttachmentsByConversation).length === 0 &&
        parsed.activeConversationId &&
        Array.isArray(parsed.queuedAttachments)
      ) {
        queuedAttachmentsByConversation[parsed.activeConversationId] = parsed.queuedAttachments;
      }
      this.queuedAttachmentsByConversationState.set(queuedAttachmentsByConversation);
    } catch {
      this.storage.removeItem(STORAGE_KEY);
    }
  }

  private persist(): void {
    if (this.bridgeMode) {
      return;
    }

    this.storage.setItem(
      STORAGE_KEY,
      JSON.stringify({
        conversations: this.conversationsState(),
        activeConversationId: this.activeConversationIdState(),
        selectedModel: this.selectedModelState(),
        models: this.modelsState(),
        settings: this.settingsState(),
        queuedAttachmentsByConversation: this.queuedAttachmentsByConversationState(),
      }),
    );
  }

  private replaceConversation(conversation: Conversation): void {
    this.conversationsState.update((items) =>
      sortConversations(
        items.map((item) =>
          item.id === conversation.id ? this.normaliseConversation(conversation) : item,
        ),
      ),
    );
    this.persist();
  }

  private updateConversation(
    id: string,
    updater: (conversation: Conversation) => Conversation,
  ): void {
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

  private normaliseConversation(conversation: Conversation): Conversation {
    return {
      ...conversation,
      messages: conversation.messages.map((message) => ({
        ...message,
        attachments: message.attachments ?? [],
        toolCalls: (message.toolCalls ?? []).map((toolCall) => ({
          ...toolCall,
          status:
            toolCall.status ?? (toolCall.error ? 'error' : toolCall.result ? 'success' : 'pending'),
        })),
      })),
    };
  }

  private async runToolCalls(
    conversationId: string,
    messageId: string,
    prompt: string,
  ): Promise<string[]> {
    const conversation = this.conversationsState().find((item) => item.id === conversationId);
    const message = conversation?.messages.find((item) => item.id === messageId);
    const toolCalls = message?.toolCalls ?? [];
    if (toolCalls.length === 0) {
      return [];
    }

    const outputs: string[] = [];
    for (const toolCall of toolCalls) {
      await pause(140);
      try {
        const result = this.executeToolCall(toolCall, prompt);
        outputs.push(`${toolCall.name}: ${result}`);
        this.updateToolCall(conversationId, messageId, toolCall.id, {
          result,
          status: 'success',
          endedAt: new Date().toISOString(),
        });
        await this.safeBridgeCall('chat:tool-call', {
          conversation_id: conversationId,
          call: {
            id: toolCall.id,
            name: toolCall.name,
            arguments: toolCall.arguments,
          },
          result: {
            tool_call_id: toolCall.id,
            content: result,
          },
        });
      } catch (error) {
        const messageText = error instanceof Error ? error.message : 'Tool execution failed.';
        outputs.push(`${toolCall.name}: ${messageText}`);
        this.updateToolCall(conversationId, messageId, toolCall.id, {
          error: messageText,
          status: 'error',
          endedAt: new Date().toISOString(),
        });
        await this.safeBridgeCall('chat:tool-call', {
          conversation_id: conversationId,
          call: {
            id: toolCall.id,
            name: toolCall.name,
            arguments: toolCall.arguments,
          },
          result: {
            tool_call_id: toolCall.id,
            content: '',
          },
          error: messageText,
        });
      }
    }

    return outputs;
  }

  private updateToolCall(
    conversationId: string,
    messageId: string,
    toolId: string,
    patch: Partial<ToolInvocation>,
  ): void {
    this.updateMessage(conversationId, messageId, (message) => ({
      ...message,
      toolCalls: (message.toolCalls ?? []).map((toolCall) =>
        toolCall.id === toolId ? { ...toolCall, ...patch } : toolCall,
      ),
    }));
  }

  private executeToolCall(toolCall: ToolInvocation, prompt: string): string {
    if (/\b(tool error|fail tool|tool failure)\b/i.test(prompt)) {
      throw new Error('Simulated MCP tool failure for the RFC chat shell.');
    }

    switch (toolCall.name) {
      case 'gui.chat.settings.load':
        return describeSettings(this.settingsState());
      case 'gui.chat.models':
        return describeModels(this.modelsState(), this.selectedModelState());
      case 'gui.chat.conversations.search':
        return describeConversationMatches(
          this.conversationsState(),
          String(toolCall.arguments['q'] ?? ''),
        );
      case 'gui.route.store':
        return describeStoreMatches(
          this.conversationsState(),
          String(toolCall.arguments['q'] ?? ''),
        );
      default:
        throw new Error(`Tool not registered in the chat shell: ${toolCall.name}`);
    }
  }
}

function deriveTitle(content: string): string {
  const trimmed = content.trim();
  if (!trimmed) {
    return 'New conversation';
  }
  return trimmed.length > 50 ? `${trimmed.slice(0, 50).trim()}...` : trimmed;
}

function deriveTitleFromMessage(content: string, attachments: ImageAttachment[]): string {
  const title = deriveTitle(content);
  if (title !== 'New conversation') {
    return title;
  }
  if (attachments.length === 0) {
    return title;
  }
  return attachments[0]?.filename
    ? deriveTitle(`Image: ${attachments[0].filename}`)
    : 'Image conversation';
}

function buildConversationMarkdown(conversation: Conversation): string {
  const body = conversation.messages
    .map((message) => {
      const heading = message.role === 'user' ? '## User' : '## Assistant';
      const attachments = (message.attachments ?? [])
        .map((attachment) => `- ${attachment.filename} (${attachment.mimeType})`)
        .join('\n');
      const thinking = message.thinking?.content
        ? ['### Thinking', message.thinking.content].join('\n\n')
        : '';
      const toolCalls = (message.toolCalls ?? [])
        .map((toolCall) =>
          [
            `#### ${toolCall.name}`,
            '```json',
            JSON.stringify(toolCall.arguments, null, 2),
            '```',
            toolCall.result,
            toolCall.error ? `Error: ${toolCall.error}` : '',
          ]
            .filter(Boolean)
            .join('\n\n'),
        )
        .join('\n\n');
      return [
        heading,
        message.content,
        attachments ? `### Attachments\n${attachments}` : '',
        thinking,
        toolCalls ? `### Tool Calls\n\n${toolCalls}` : '',
      ]
        .filter(Boolean)
        .join('\n\n');
    })
    .join('\n\n---\n\n');

  return [
    `# ${conversation.title}`,
    `- Conversation ID: ${conversation.id}`,
    `- Model: ${conversation.model}`,
    `- Updated: ${conversation.updatedAt}`,
    '',
    body,
  ].join('\n');
}

function hasCoreGUIBridge(): boolean {
  return (
    typeof window.__coreGUIBridge?.invoke === 'function' ||
    typeof window.__coreGUIBridge?.call === 'function'
  );
}

function sortConversations(conversations: Conversation[]): Conversation[] {
  return [...conversations].sort((left, right) => right.updatedAt.localeCompare(left.updatedAt));
}

function resolveSelectedModel(models: ModelEntry[], preferred: string): string {
  if (preferred && models.some((model) => model.name === preferred)) {
    return preferred;
  }
  return models.find((model) => model.loaded)?.name ?? models[0]?.name ?? 'lemer';
}

function fromRemoteModel(model: RemoteModelEntry): ModelEntry {
  return {
    name: model.name,
    architecture: model.architecture,
    quantBits: model.quant_bits,
    sizeBytes: model.size_bytes,
    loaded: model.loaded,
    backend: model.backend,
    supportsVision: model.supports_vision,
  };
}

function fromRemoteSettings(settings: RemoteChatSettings): ChatSettings {
  return {
    temperature: settings.temperature,
    topP: settings.top_p,
    topK: settings.top_k,
    maxTokens: settings.max_tokens,
    contextWindow: settings.context_window,
    systemPrompt: settings.system_prompt,
    defaultModel: settings.default_model,
  };
}

function toRemoteSettings(settings: ChatSettings): RemoteChatSettings {
  return {
    temperature: settings.temperature,
    top_p: settings.topP,
    top_k: settings.topK,
    max_tokens: settings.maxTokens,
    context_window: settings.contextWindow,
    system_prompt: settings.systemPrompt,
    default_model: settings.defaultModel,
  };
}

function fromRemoteAttachment(attachment: RemoteImageAttachment): ImageAttachment {
  return {
    id: attachment.id ?? crypto.randomUUID(),
    filename: attachment.filename,
    mimeType: attachment.mime_type,
    data: attachment.data,
    width: attachment.width,
    height: attachment.height,
  };
}

function toRemoteAttachment(attachment: ImageAttachment): RemoteImageAttachment {
  return {
    id: attachment.id,
    filename: attachment.filename,
    mime_type: attachment.mimeType,
    data: attachment.data,
    width: attachment.width,
    height: attachment.height,
  };
}

function fromRemoteThinking(thinking?: RemoteThinkingState): ThinkingState | undefined {
  if (!thinking) {
    return undefined;
  }
  return {
    active: thinking.active,
    content: thinking.content,
    startedAt: thinking.started_at,
    finishedAt: thinking.finished_at,
  };
}

function fromRemoteToolInvocation(invocation: RemoteToolInvocation): ToolInvocation {
  return {
    id: invocation.call.id,
    name: invocation.call.name,
    arguments: invocation.call.arguments ?? {},
    result: invocation.result.content,
    error: invocation.error,
    status: invocation.error ? 'error' : invocation.result.content ? 'success' : 'pending',
    startedAt: invocation.started_at,
    endedAt: invocation.ended_at,
  };
}

function fromRemoteMessage(message: RemoteChatMessage): ChatMessage {
  return {
    id: message.id,
    role: message.role === 'assistant' ? 'assistant' : 'user',
    content: message.content,
    createdAt: message.created_at,
    attachments: (message.attachments ?? []).map(fromRemoteAttachment),
    thinking: fromRemoteThinking(message.thinking),
    toolCalls: (message.tool_calls ?? []).map(fromRemoteToolInvocation),
    streaming: message.streaming,
  };
}

function fromRemoteConversation(conversation: RemoteConversation): Conversation {
  return {
    id: conversation.id,
    title: conversation.title,
    model: conversation.model,
    createdAt: conversation.created_at,
    updatedAt: conversation.updated_at,
    messages: conversation.messages.map(fromRemoteMessage),
    settings: conversation.settings ? fromRemoteSettings(conversation.settings) : undefined,
  };
}

function inferToolCalls(content: string): ToolInvocation[] {
  const normalized = content.toLowerCase();
  const query = extractSearchQuery(content);
  const toolCalls: ToolInvocation[] = [];
  const push = (name: string, args: Record<string, unknown>) => {
    if (toolCalls.some((toolCall) => toolCall.name === name)) {
      return;
    }
    toolCalls.push({
      id: crypto.randomUUID(),
      name,
      arguments: args,
      result: '',
      status: 'pending',
      startedAt: new Date().toISOString(),
    });
  };

  if (/\b(setting|temperature|context|prompt)\b/.test(normalized)) {
    push('gui.chat.settings.load', {});
  }
  if (/\bmodel|models\b/.test(normalized)) {
    push('gui.chat.models', {});
  }
  if (/\b(history|conversation|conversations|find|search)\b/.test(normalized)) {
    push('gui.chat.conversations.search', { q: query });
  }
  if (/\b(store|storage|cache|local data|tool)\b/.test(normalized)) {
    push('gui.route.store', { q: query });
  }

  return toolCalls.slice(0, 3);
}

async function generateAssistantResponse(
  content: string,
  model: string,
  settings: ChatSettings,
  toolOutputs: string[] = [],
): Promise<string> {
  const prompt = content.trim() || 'your multimodal prompt';
  const promptWithTools = buildToolAwarePrompt(prompt, toolOutputs);
  try {
    if (typeof window.core?.ml?.generate === 'function') {
      const generated = await window.core.ml.generate(promptWithTools);
      if (generated.trim()) {
        return generated;
      }
    }
  } catch {
    // Fall back to the local RFC demo response below.
  }
  return buildFallbackResponse(prompt, model, settings, toolOutputs);
}

function buildFallbackResponse(
  prompt: string,
  model: string,
  settings: ChatSettings,
  toolOutputs: string[],
): string {
  const toolSection =
    toolOutputs.length > 0
      ? ['', 'Tool results:', ...toolOutputs.map((output) => `- ${output}`), '']
      : [];
  return [
    `Using ${model} with temperature ${settings.temperature.toFixed(1)}.`,
    '',
    `I stored this exchange locally and can keep the conversation context across sessions.`,
    '',
    `Prompt summary: ${prompt}`,
    ...toolSection,
    '',
    '```ts',
    `const response = await window.core.ml.generate(${JSON.stringify(prompt)});`,
    '```',
    '',
    'The live CoreGUI inference bridge was unavailable, so this shell fell back to the local RFC demo response with progressive rendering, tool-call blocks, and multimodal attachments.',
  ].join('\n');
}

async function streamIntoMessage(
  content: string,
  onUpdate: (fragment: string, done: boolean) => void | Promise<void>,
): Promise<void> {
  const chunks = content.match(/.{1,18}/g) ?? [content];
  let assembled = '';
  for (const chunk of chunks) {
    assembled += chunk;
    await onUpdate(assembled, false);
    await new Promise((resolve) => window.setTimeout(resolve, 45));
  }
  await onUpdate(assembled, true);
}

function pause(ms: number): Promise<void> {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

function extractSearchQuery(content: string): string {
  const quoted = content.match(/"([^"]+)"/);
  if (quoted?.[1]) {
    return quoted[1];
  }

  const afterFor = content.match(/\b(?:for|about|named)\s+(.+)/i);
  if (afterFor?.[1]) {
    return afterFor[1].trim().slice(0, 48);
  }

  return content.trim().slice(0, 48);
}

function describeSettings(settings: ChatSettings): string {
  return [
    `temperature=${settings.temperature.toFixed(1)}`,
    `topP=${settings.topP.toFixed(2)}`,
    `topK=${settings.topK}`,
    `maxTokens=${settings.maxTokens}`,
    `contextWindow=${settings.contextWindow}`,
    `defaultModel=${settings.defaultModel}`,
  ].join(', ');
}

function describeModels(models: ModelEntry[], selectedModel: string): string {
  const summary = models
    .map((model) => {
      const marker = model.name === selectedModel ? '[selected]' : '[available]';
      const vision = model.supportsVision === false ? 'text-only' : 'vision';
      return `${marker} ${model.name} (${model.architecture}, ${vision}, ${model.backend})`;
    })
    .join('; ');
  return summary || 'No local models are currently listed.';
}

function describeConversationMatches(conversations: Conversation[], query: string): string {
  const needle = query.trim().toLowerCase();
  const matches = conversations.filter((conversation) => {
    if (!needle) {
      return true;
    }
    return conversationSearchText(conversation).includes(needle);
  });

  if (matches.length === 0) {
    return `No conversations matched "${query.trim() || 'the current query'}".`;
  }

  return `Found ${matches.length} conversation(s): ${matches
    .slice(0, 3)
    .map((conversation) => `"${conversation.title}"`)
    .join(', ')}.`;
}

function messageStoreSearchText(conversation: Conversation, message: ChatMessage): string {
  return [
    conversation.title,
    message.role,
    message.content,
    message.thinking?.content ?? '',
    ...(message.attachments ?? []).flatMap((attachment) => [
      attachment.filename,
      attachment.mimeType,
    ]),
    ...(message.toolCalls ?? []).flatMap((toolCall) => [
      toolCall.name,
      JSON.stringify(toolCall.arguments ?? {}),
      toolCall.result,
      toolCall.error ?? '',
    ]),
  ]
    .join(' ')
    .toLowerCase();
}

function messageStoreSnippet(message: ChatMessage): string {
  if (message.content.trim().length > 0) {
    return message.content.trim().slice(0, 80);
  }

  if ((message.thinking?.content ?? '').trim().length > 0) {
    return message.thinking?.content.trim().slice(0, 80) ?? '';
  }

  for (const toolCall of message.toolCalls ?? []) {
    const result = toolCall.result.trim();
    if (result.length > 0) {
      return result.slice(0, 80);
    }
    const error = (toolCall.error ?? '').trim();
    if (error.length > 0) {
      return error.slice(0, 80);
    }
  }

  const attachmentNames = (message.attachments ?? []).map((attachment) => attachment.filename);
  if (attachmentNames.length > 0) {
    return `[attachments] ${attachmentNames.join(', ')}`.slice(0, 80);
  }

  return '[message]';
}

function describeStoreMatches(conversations: Conversation[], query: string): string {
  const needle = query.trim().toLowerCase();
  const snippets = conversations.flatMap((conversation) =>
    conversation.messages.flatMap((message) => {
      const attachments = (message.attachments ?? []).map((attachment) => attachment.filename);
      const haystack = messageStoreSearchText(conversation, message);
      if (needle && !haystack.includes(needle)) {
        return [];
      }
      return [
        `${conversation.title}: ${messageStoreSnippet(message)}${
          attachments.length > 0 ? ` (attachments: ${attachments.join(', ')})` : ''
        }`,
      ];
    }),
  );

  if (snippets.length === 0) {
    return `The local store search found no matches for "${query.trim() || 'the current query'}".`;
  }

  return `Store hits (${snippets.length}): ${snippets.slice(0, 3).join(' | ')}.`;
}

function buildToolAwarePrompt(prompt: string, toolOutputs: string[]): string {
  const sections = ['System tool manifest:', buildToolManifest(), '', 'User prompt:', prompt];

  if (toolOutputs.length > 0) {
    sections.push('', 'Tool context:', ...toolOutputs.map((output) => `- ${output}`));
  }

  return sections.join('\n');
}

function buildToolManifest(): string {
  return REGISTERED_CHAT_TOOLS.map((tool) => {
    const params = Object.entries(tool.parameters)
      .map(([name, description]) => `${name}: ${description}`)
      .join(', ');
    return params
      ? `- ${tool.name}: ${tool.description} Parameters: ${params}.`
      : `- ${tool.name}: ${tool.description}`;
  }).join('\n');
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
