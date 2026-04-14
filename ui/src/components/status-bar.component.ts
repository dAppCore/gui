import { Component, Input, OnDestroy, OnInit, computed, inject, signal } from '@angular/core';
import { ChatService } from '../services/chat.service';

@Component({
  selector: 'status-bar',
  standalone: true,
  template: `
    <footer class="status-bar" [style.--sidebar-width]="sidebarWidth">
      <div class="status-left">
        <span class="status-item version">{{ version }}</span>
        <span class="status-item">
          <i class="fa-regular fa-comments"></i>
          {{ conversationCount() }} conversations
        </span>
        <span class="status-item">
          <i class="fa-regular fa-microchip-ai"></i>
          {{ activeModel() }}
        </span>
      </div>
      <div class="status-right">
        <span class="status-item connection" [class.connected]="!chat.busy()">
          <span class="status-dot"></span>
          {{ chat.busy() ? 'Streaming' : 'Ready' }}
        </span>
        <span class="status-item time">{{ time() }}</span>
      </div>
    </footer>
  `,
  styles: [
    `
      .status-bar {
        position: fixed;
        left: 0;
        width: 100%;
        bottom: 0;
        z-index: 40;
        height: 2.75rem;
        border-top: 1px solid rgba(255, 255, 255, 0.06);
        background: linear-gradient(180deg, rgba(6, 10, 18, 0.88), rgba(6, 10, 18, 0.96));
        backdrop-filter: blur(18px);
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding-inline: 1rem 1.25rem;
        box-shadow: 0 -12px 40px rgba(0, 0, 0, 0.2);
      }

      @media (min-width: 1024px) {
        .status-bar {
          left: var(--sidebar-width, 0);
          width: calc(100% - var(--sidebar-width, 0));
        }
      }

      .status-left,
      .status-right {
        display: flex;
        align-items: center;
        gap: 1rem;
      }

      .status-item {
        font-size: 0.875rem;
        color: rgb(168 179 207);
      }

      .status-item i {
        margin-right: 0.25rem;
      }

      .version {
        letter-spacing: 0.08em;
        text-transform: uppercase;
      }

      .status-dot {
        display: inline-block;
        width: 7px;
        height: 7px;
        border-radius: 50%;
        background: rgb(249 115 22);
        margin-right: 0.375rem;
      }

      .connection.connected .status-dot {
        background: rgb(20 184 166);
        box-shadow: 0 0 8px rgba(20, 184, 166, 0.4);
      }

      .time {
        font-family: 'JetBrains Mono', 'Fira Code', monospace;
        color: rgb(244 247 251);
      }
    `,
  ],
})
export class StatusBarComponent implements OnInit, OnDestroy {
  @Input() version = 'v0.1.0';
  @Input() sidebarWidth = '5rem';

  protected readonly chat = inject(ChatService);
  protected readonly time = signal('');
  protected readonly conversationCount = computed(() => this.chat.conversations().length);
  protected readonly activeModel = computed(() => this.chat.selectedModel());

  private intervalId: ReturnType<typeof setInterval> | undefined;

  ngOnInit(): void {
    this.updateTime();
    this.intervalId = setInterval(() => this.updateTime(), 1000);
  }

  ngOnDestroy(): void {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }

  private updateTime(): void {
    this.time.set(new Date().toLocaleTimeString());
  }
}
