import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MainNavComponent } from './components/main-nav/main-nav.component';
import { CardComponent } from './components/card/card.component';
import { BannerComponent } from './components/banner/banner.component';
import { RegisterComponent } from './pages/register/register.component';
import { RegisterCardComponent } from './pages/register/components/register-card/register-card.component';
import {FormsModule} from "@angular/forms";


import { TransactionPreviewComponent } from './components/transaction/transaction-preview/transaction-preview.component';
import { LoginComponent } from './pages/login/login.component';
import { MainLayoutComponent } from './layouts/main-layout/main-layout.component';
import { HomeComponent } from './pages/home/home.component';


@NgModule({
  declarations: [
    AppComponent,
    MainNavComponent,
    CardComponent,
    BannerComponent,
    RegisterComponent,
    RegisterCardComponent,
    TransactionPreviewComponent,
    LoginComponent,
    MainLayoutComponent,
    HomeComponent,

  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule
  ],
  providers: [],
  exports: [

  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
