import { Routes } from '@angular/router';
import { ApplicationFrameComponent } from '../frame/application-frame.component';
import { DashboardComponent } from './dashboard.component';
import { SettingsComponent } from './settings.component';

export const routes: Routes = [
  {
    path: '',
    component: ApplicationFrameComponent,
    children: [
      { path: '', component: DashboardComponent },
      { path: 'settings', component: SettingsComponent },
    ],
  },
];
