import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ImageAttachment } from './chat.types';

@Component({
  selector: 'chat-input-area',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="composer" (dragover)="onDragOver($event)" (drop)="onDrop($event)">
      <input #filePicker type="file" accept="image/*" multiple hidden (change)="onFileSelection($event)" />
      <div class="composer__attachments" *ngIf="attachments.length">
        <figure *ngFor="let attachment of attachments; let index = index" class="attachment">
          <img [src]="attachmentSource(attachment)" [alt]="attachment.filename" />
          <figcaption>{{ attachment.filename }}</figcaption>
          <button type="button" (click)="removeAttachment.emit(index)">Remove</button>
        </figure>
      </div>
      <textarea
        #textarea
        [ngModel]="value"
        (ngModelChange)="onValueChange($event)"
        (keydown.enter)="submitOnEnter($event)"
        (paste)="onPaste($event)"
        placeholder="Ask the local model something useful"
      ></textarea>
      <div class="composer__meta">
        <span>{{ value.length }} chars · {{ attachments.length }} image(s)</span>
        <div class="composer__actions">
          <button type="button" class="ghost" (click)="filePicker.click()">Attach</button>
          <button type="button" [disabled]="disabled" (click)="submit.emit()">Send</button>
        </div>
      </div>
    </div>
  `,
  styles: [
    `
      .composer { display: grid; gap: 0.75rem; padding: 1rem; border-radius: 1.35rem; background: rgba(8, 21, 35, 0.86); border: 1px solid rgba(125, 211, 252, 0.12); }
      .composer__attachments { display: flex; gap: 0.75rem; overflow: auto; }
      .attachment { min-width: 8rem; margin: 0; display: grid; gap: 0.45rem; padding: 0.6rem; border-radius: 1rem; background: rgba(2, 6, 23, 0.76); }
      .attachment img { width: 100%; aspect-ratio: 1.3; object-fit: cover; border-radius: 0.8rem; }
      .attachment figcaption { color: #cbd5e1; font-size: 0.78rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
      textarea { width: 100%; min-height: 3.5rem; max-height: 10.5rem; resize: none; border: 0; background: transparent; color: #f8fafc; font: inherit; outline: 0; line-height: 1.6; }
      .composer__meta { display: flex; justify-content: space-between; align-items: center; gap: 1rem; color: #94a3b8; font-size: 0.82rem; }
      .composer__actions { display: flex; gap: 0.6rem; }
      button { border: 0; border-radius: 999px; padding: 0.8rem 1.2rem; background: linear-gradient(135deg, #f59e0b, #fb7185); color: #111827; font-weight: 800; cursor: pointer; }
      .ghost { background: rgba(15, 23, 42, 0.8); color: #e2e8f0; border: 1px solid rgba(148, 163, 184, 0.22); }
      button:disabled { opacity: 0.4; cursor: not-allowed; }
    `,
  ],
})
export class InputAreaComponent implements AfterViewInit {
  @ViewChild('textarea') private readonly textarea?: ElementRef<HTMLTextAreaElement>;

  @Input() value = '';
  @Input() disabled = false;
  @Input() attachments: ImageAttachment[] = [];
  @Output() valueChange = new EventEmitter<string>();
  @Output() attachFiles = new EventEmitter<FileList | File[]>();
  @Output() removeAttachment = new EventEmitter<number>();
  @Output() submit = new EventEmitter<void>();

  ngAfterViewInit(): void {
    this.resizeTextarea();
  }

  onValueChange(value: string): void {
    this.valueChange.emit(value);
    queueMicrotask(() => this.resizeTextarea());
  }

  submitOnEnter(event: Event): void {
    const keyboard = event as KeyboardEvent;
    if (!keyboard.shiftKey) {
      keyboard.preventDefault();
      this.submit.emit();
    }
  }

  onDragOver(event: DragEvent): void {
    event.preventDefault();
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    if (event.dataTransfer?.files?.length) {
      this.attachFiles.emit(event.dataTransfer.files);
    }
  }

  onPaste(event: ClipboardEvent): void {
    const files = Array.from(event.clipboardData?.files ?? []);
    if (files.length > 0) {
      this.attachFiles.emit(files);
    }
  }

  onFileSelection(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files?.length) {
      this.attachFiles.emit(input.files);
      input.value = '';
    }
  }

  attachmentSource(attachment: ImageAttachment): string {
    return `data:${attachment.mime_type};base64,${attachment.data}`;
  }

  private resizeTextarea(): void {
    const element = this.textarea?.nativeElement;
    if (!element) {
      return;
    }
    element.style.height = 'auto';
    element.style.height = `${Math.min(element.scrollHeight, 168)}px`;
  }
}
