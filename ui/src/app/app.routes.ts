import { Routes } from '@angular/router';
import { ApplicationFrameComponent } from '../frame/application-frame.component';
import { DashboardComponent } from './dashboard.component';
import { ProviderHostComponent } from '../components/provider-host.component';
import { SettingsComponent } from './settings.component';

export const routes: Routes = [
  {
    path: '',
    component: ApplicationFrameComponent,
    children: [
      { path: '', component: DashboardComponent },
      { path: 'provider/:provider', component: ProviderHostComponent },
      { path: 'settings', component: SettingsComponent },
    ],
  },
];
