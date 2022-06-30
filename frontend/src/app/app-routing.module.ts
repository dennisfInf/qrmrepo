import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {HomeComponent} from "./pages/home/home.component";
import {AboutComponent} from "./pages/about/about.component";
import {DevelpersComponent} from "./pages/about/develpers/develpers.component";
import {UsersComponent} from "./pages/about/users/users.component";
import {ProjectComponent} from "./pages/about/project/project.component";
import {RegisterComponent} from "./pages/register/register.component";
import {PaymentComponent} from "./pages/payment/payment.component";
import {DashboardModule} from "./dashboard/dashboard.module";
import {AuthGuardService} from "./services/auth-guard.service";


const routes: Routes = [
  {
    path: "payment",
    component: PaymentComponent
  },
  {
    path: "dashboard",
    loadChildren: () => DashboardModule,
    // canActivate : [AuthGuardService]
  },
  {
    path: "register",
    component: RegisterComponent
  },
  {
    path: "about",
    component: AboutComponent,
    children: [
      {
        path: "",
        pathMatch: "prefix",
        redirectTo: "project"
      },
      {
        path: "project",
        component: ProjectComponent
      },
      {
        path: "developers",
        component: DevelpersComponent
      },
      {
        path: "users",
        component: UsersComponent
      }
    ]
  },
  {
    path: "",
    component: HomeComponent
  },


];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
