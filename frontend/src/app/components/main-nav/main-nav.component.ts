import { Component, OnInit } from '@angular/core';
import {AuthenticationService} from "../../services/authentication.service";
import {UserService} from "../../services/user.service";
import {Router} from "@angular/router";


@Component({
  selector: 'app-main-nav',
  templateUrl: './main-nav.component.html',
  styleUrls: ['./main-nav.component.css']
})
export class MainNavComponent implements OnInit {
  profileMenu: string = 'hidden';
  mainMenu : string = 'hidden';
  constructor(public authService : AuthenticationService, public userService : UserService, private router : Router) { }

  ngOnInit(): void {
  }

  toggleProfileMenu(){
    if(this.profileMenu == 'hidden'){
      this.profileMenu = ''
    }else {
      this.profileMenu = 'hidden'
    }
  }

  closeMainMenu(){
    setTimeout(()=>{
      this.mainMenu = 'hidden'
    }, 100)
  }

  closeProfileMenu(){
    setTimeout(()=>{
      this.profileMenu = 'hidden'
    }, 100)
  }

  toggleMainMenu() {
    if(this.mainMenu == 'hidden'){
      this.mainMenu = ''
    }else {
      this.mainMenu = 'hidden'
    }
  }

  logout() {
    this.authService.logout()
    this.router.navigate(["/login"])
  }
}
