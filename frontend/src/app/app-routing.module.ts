import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {RegisterComponent} from "./pages/register/register.component";
import {DashboardModule} from "./dashboard/dashboard.module";
import {AuthGuardService} from "./services/auth-guard.service";
import {LoginComponent} from "./pages/login/login.component";
import {MainLayoutComponent} from "./layouts/main-layout/main-layout.component";
import {HomeComponent} from "./pages/home/home.component";


const routes: Routes = [
  {
    path: "dashboard",
    component: MainLayoutComponent,
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
  {
    path: "landing", component: HomeComponent
  }


];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
