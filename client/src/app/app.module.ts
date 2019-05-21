import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';

import { StorageServiceModule } from 'ngx-webstorage-service';

import { AppComponent } from './app.component';
import { SocketService } from "./socket.service";

@NgModule({
    declarations: [
        AppComponent
    ],
    imports: [
        BrowserModule,
        FormsModule,
        StorageServiceModule,
    ],
    providers: [SocketService],
    bootstrap: [AppComponent]
})
export class AppModule { }
