import { Component, OnInit } from '@angular/core';
import {AuthenticationService} from "../../services/authentication.service";
import {FidoService} from "../../services/fido.service";
import {Router} from "@angular/router";
import {UserService} from "../../services/user.service";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  error = ""
  constructor(private authService : AuthenticationService, private fidoService : FidoService, private router : Router, private userSerivce : UserService) {
    if(this.authService.isAuthenticated()){
      this.router.navigate(["/dashboard"])
    }
  }

  ngOnInit(): void {
  }


  async login(username:string) {
    if(this.validateEmail(username)) {
      this.authService.loginInitialize(username)
        .then(res => {
          this.fidoService.getCredential(res.data).then(res => {
            this.authService.loginFinalize(username, res as PublicKeyCredential).then(res => {
              return res.json()

            }).then(data => {
              console.log(data)
              if (this.authService.login(data.token)) {
                this.router.navigate(["/dashboard"])
              }
            })
          }, err => {
            this.error = "User not found"
            setTimeout(() => {
              this.error = ""
            }, 5000)
          })
        })
    }else {
      this.error = "Enter a valid email address"
      setTimeout(() => {
        this.error = ""
      }, 5000)
    }

  }

  validateEmail(email : string)  {
    return String(email)
      .toLowerCase()
      .match(
        /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
      );
  };
}
