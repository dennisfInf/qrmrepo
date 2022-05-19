import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {HomeComponent} from "./pages/home/home.component";
import {AboutComponent} from "./pages/about/about.component";
import {DevelpersComponent} from "./pages/about/develpers/develpers.component";
import {UsersComponent} from "./pages/about/users/users.component";
import {ProjectComponent} from "./pages/about/project/project.component";
import {RegisterComponent} from "./pages/register/register.component";

const routes: Routes = [
  {
    path: "",
    component: HomeComponent
  },
  {
    path : "about",
    component: AboutComponent,
    children : [
      {
        path : "",
        pathMatch: "prefix",
        redirectTo : "project"
      },
      {
        path: "project",
        component: ProjectComponent
      },
      {
        path : "developers",
        component : DevelpersComponent
      },
      {
        path : "users",
        component: UsersComponent
      }
    ]
  },
  {
    path : "register",
    component: RegisterComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
