import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MainNavComponent } from './components/main-nav/main-nav.component';
import { CardComponent } from './components/card/card.component';
import { BannerComponent } from './components/banner/banner.component';
import { HomeComponent } from './pages/home/home.component';
import { AboutComponent } from './pages/about/about.component';
import { DevelpersComponent } from './pages/about/develpers/develpers.component';
import { UsersComponent } from './pages/about/users/users.component';
import { ProjectComponent } from './pages/about/project/project.component';
import { NavigationComponent } from './pages/about/components/navigation/navigation.component';
import { PageContentComponent } from './pages/about/components/page-content/page-content.component';
import { RegisterComponent } from './pages/register/register.component';
import { RegisterCardComponent } from './pages/register/components/register-card/register-card.component';
import {FormsModule} from "@angular/forms";
import { PaymentComponent } from './pages/payment/payment.component';


import { TransactionPreviewComponent } from './components/transaction/transaction-preview/transaction-preview.component';
import { LoginComponent } from './pages/login/login.component';
import { LandingComponent } from './pages/landing/landing.component';


@NgModule({
  declarations: [
    AppComponent,
    MainNavComponent,
    CardComponent,
    BannerComponent,
    HomeComponent,
    AboutComponent,
    DevelpersComponent,
    UsersComponent,
    ProjectComponent,
    NavigationComponent,
    PageContentComponent,
    RegisterComponent,
    RegisterCardComponent,
    PaymentComponent,
    TransactionPreviewComponent,
    LoginComponent,
    LandingComponent,

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
