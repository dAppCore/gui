import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ModelEntry } from './chat.types';

@Component({
  selector: 'chat-model-selector',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <label class="selector">
      <span>Model</span>
      <select [ngModel]="value" (ngModelChange)="valueChange.emit($event)">
        <option *ngFor="let model of models" [ngValue]="model.name">
          {{ model.name }} · {{ model.architecture }} · {{ model.backend }}
        </option>
      </select>
    </label>
  `,
  styles: [
    `
      .selector { display: grid; gap: 0.35rem; color: #cbd5e1; font-size: 0.82rem; }
      select { min-width: 16rem; border-radius: 0.8rem; border: 1px solid rgba(124, 156, 191, 0.2); background: rgba(8, 21, 35, 0.8); color: #e2e8f0; padding: 0.72rem 0.9rem; }
    `,
  ],
})
export class ModelSelectorComponent {
  @Input() models: ModelEntry[] = [];
  @Input() value = '';
  @Output() valueChange = new EventEmitter<string>();
}
