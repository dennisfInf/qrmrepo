import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DashboardComponent } from './pages/dashboard/dashboard.component';

import { DashboardLayoutComponent } from './dashboard-layout/dashboard-layout.component';
import {DashboardRoutes} from "./dashboard.routes";
import {RouterModule} from "@angular/router";
import {ContactThumbnailComponent} from "./contact-thumbnail/contact-thumbnail.component";
import {ShopThumbnailComponent} from "./shop-thumbnail/shop-thumbnail.component";
import { UserPaymentComponent } from './pages/user-payment/user-payment.component';
import {FormsModule} from "@angular/forms";


@NgModule({
  declarations: [
    DashboardComponent,
    DashboardLayoutComponent,
    ContactThumbnailComponent,
    ShopThumbnailComponent,
    UserPaymentComponent
  ],
  imports: [
    CommonModule,
    RouterModule.forChild(DashboardRoutes),
    FormsModule
  ],
  exports : [
    RouterModule
  ]
})
export class DashboardModule { }
