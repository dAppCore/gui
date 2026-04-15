import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'chat-input-area',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="composer">
      <textarea
        [ngModel]="value"
        (ngModelChange)="valueChange.emit($event)"
        (keydown.enter)="submitOnEnter($event)"
        placeholder="Ask the local model something useful"
      ></textarea>
      <div class="composer__meta">
        <span>{{ value.length }} chars</span>
        <button type="button" [disabled]="disabled" (click)="submit.emit()">Send</button>
      </div>
    </div>
  `,
  styles: [
    `
      .composer { display: grid; gap: 0.75rem; padding: 1rem; border-radius: 1.35rem; background: rgba(8, 21, 35, 0.86); border: 1px solid rgba(125, 211, 252, 0.12); }
      textarea { width: 100%; min-height: 7rem; resize: vertical; border: 0; background: transparent; color: #f8fafc; font: inherit; outline: 0; }
      .composer__meta { display: flex; justify-content: space-between; align-items: center; color: #94a3b8; font-size: 0.82rem; }
      button { border: 0; border-radius: 999px; padding: 0.8rem 1.2rem; background: linear-gradient(135deg, #f59e0b, #fb7185); color: #111827; font-weight: 800; cursor: pointer; }
      button:disabled { opacity: 0.4; cursor: not-allowed; }
    `,
  ],
})
export class InputAreaComponent {
  @Input() value = '';
  @Input() disabled = false;
  @Output() valueChange = new EventEmitter<string>();
  @Output() submit = new EventEmitter<void>();

  submitOnEnter(event: Event): void {
    const keyboard = event as KeyboardEvent;
    if (!keyboard.shiftKey) {
      keyboard.preventDefault();
      this.submit.emit();
    }
  }
}
