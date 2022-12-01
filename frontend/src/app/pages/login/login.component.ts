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

  constructor(private authService : AuthenticationService, private fidoService : FidoService, private router : Router, private userSerivce : UserService) {
  }

  ngOnInit(): void {
  }


  async login(username:string) {
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
        })
      })
  }
}
