import {Route, Routes} from "@angular/router";
import {DashboardLayoutComponent} from "./dashboard-layout/dashboard-layout.component";
import {DashboardComponent} from "./pages/dashboard/dashboard.component";
import {TransactionsComponent} from "./pages/transactions/transactions.component";
import {ContactsComponent} from "./pages/contacts/contacts.component";
import {ProfileComponent} from "./pages/profile/profile.component";

export const DashboardRoutes : Routes= [
  {
    path: "",
    component : DashboardLayoutComponent,
    children : [
      {path : "transactions" , component : TransactionsComponent},
      {path : "contacts" , component : ContactsComponent},
      {path : "profile" , component : ProfileComponent},
      {path : "" , component : DashboardComponent},

    ]
  }
]
