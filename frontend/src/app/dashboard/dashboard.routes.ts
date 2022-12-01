import {Route, Routes} from "@angular/router";
import {DashboardLayoutComponent} from "./dashboard-layout/dashboard-layout.component";
import {DashboardComponent} from "./pages/dashboard/dashboard.component";
import {UserPaymentComponent} from "./pages/user-payment/user-payment.component";

export const DashboardRoutes : Routes= [
  {
    path: "",
    component : DashboardLayoutComponent,
    children : [
      {path : "user-payment" , component : UserPaymentComponent},
      {path : "" , component : DashboardComponent},

    ]
  }
]
