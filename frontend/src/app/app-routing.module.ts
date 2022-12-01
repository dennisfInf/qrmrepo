import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {RegisterComponent} from "./pages/register/register.component";
import {DashboardModule} from "./dashboard/dashboard.module";
import {AuthGuardService} from "./services/auth-guard.service";
import {LoginComponent} from "./pages/login/login.component";


const routes: Routes = [
  {
    path: "dashboard",
    loadChildren: () => DashboardModule,
    canActivate: [AuthGuardService]
  },
  {
    path: "register",
    component: RegisterComponent
  },
  {
    path : "login",
    component: LoginComponent
  },
  {
    path: "", redirectTo: "login", pathMatch: "full"
  },


];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
