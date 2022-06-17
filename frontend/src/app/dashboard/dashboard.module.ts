import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DashboardComponent } from './pages/dashboard/dashboard.component';
import { TransactionsComponent } from './pages/transactions/transactions.component';
import { TransactionComponent } from './pages/transaction/transaction.component';
import { ContactsComponent } from './pages/contacts/contacts.component';
import { ProfileComponent } from './pages/profile/profile.component';
import { DashboardLayoutComponent } from './dashboard-layout/dashboard-layout.component';
import {DashboardRoutes} from "./dashboard.routes";
import {RouterModule} from "@angular/router";
import {ContactThumbnailComponent} from "./contact-thumbnail/contact-thumbnail.component";
import {ShopThumbnailComponent} from "./shop-thumbnail/shop-thumbnail.component";


@NgModule({
  declarations: [
    DashboardComponent,
    TransactionsComponent,
    TransactionComponent,
    ContactsComponent,
    ProfileComponent,
    DashboardLayoutComponent,
    ContactThumbnailComponent,
    ShopThumbnailComponent
  ],
  imports: [
    CommonModule,
    RouterModule.forChild(DashboardRoutes)
  ],
  exports : [
    RouterModule
  ]
})
export class DashboardModule { }
