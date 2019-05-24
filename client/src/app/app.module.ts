import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';

import { StorageServiceModule } from 'ngx-webstorage-service';

import { AppComponent } from './app.component';
import { SocketService } from "./socket.service";
import { Connect6Component } from './connect6/connect6.component';
import { HomeComponent } from './home/home.component';

const routes: Routes = [
  { path: 'connect6/:id', component: Connect6Component },
  { path: '', component: HomeComponent },
  { path: '**', redirectTo: '', pathMatch: 'full' },
]

@NgModule({
  declarations: [
    AppComponent,
    Connect6Component,
    HomeComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    StorageServiceModule,
    RouterModule.forRoot(routes),
    HttpClientModule
  ],
  providers: [SocketService],
  bootstrap: [AppComponent]
})
export class AppModule { }
