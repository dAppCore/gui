export interface ChatMessage {
  id: string;
  role: string;
  content: string;
  created_at: string;
  model?: string;
  finish_reason?: string;
  thinking?: {
    active: boolean;
    content: string;
    duration_ms?: number;
  };
  tool_calls?: Array<{ id: string; name: string; arguments: Record<string, unknown> }>;
  tool_results?: Array<{ tool_call_id: string; content: string }>;
  attachments?: Array<{ filename: string; mime_type: string; data: string; width: number; height: number }>;
}

export interface ConversationSummary {
  id: string;
  title: string;
  model: string;
  created_at: string;
  updated_at: string;
  message_count: number;
}

export interface Conversation extends ConversationSummary {
  messages: ChatMessage[];
}

export interface ModelEntry {
  name: string;
  architecture: string;
  quant_bits: number;
  size_bytes: number;
  loaded: boolean;
  backend: string;
}

export interface ChatSettings {
  temperature: number;
  top_p: number;
  top_k: number;
  max_tokens: number;
  context_window: number;
  system_prompt: string;
  default_model: string;
}
