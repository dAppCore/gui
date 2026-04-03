// SPDX-Licence-Identifier: EUPL-1.2

import { DoBootstrap, Injector, NgModule, provideBrowserGlobalErrorListeners } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { createCustomElement } from '@angular/elements';
import { RouterModule } from '@angular/router';

import { App } from './app';
import { routes } from './app.routes';

@NgModule({
  imports: [BrowserModule, App, RouterModule.forRoot(routes)],
  providers: [provideBrowserGlobalErrorListeners()],
})
export class AppModule implements DoBootstrap {
  constructor(private injector: Injector) {
    const el = createCustomElement(App, { injector });
    customElements.define('core-display', el);
  }

  ngDoBootstrap() {}
}
