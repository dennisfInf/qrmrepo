import { Component, OnInit } from '@angular/core';
import {AuthenticationService} from "../../services/authentication.service";
import {FidoService} from "../../services/fido.service";
import {Router} from "@angular/router";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  name! : string
  username! :string

  constructor(private authService : AuthenticationService, private fidoService : FidoService, private router : Router) { }

  ngOnInit(): void {
  }
  async login() {
    this.authService.loginInitialize(this.username)
      .then(res => {
        let jsonObj = JSON.parse(res)
        let userId = jsonObj.user.id as BufferSource
        let challenge = jsonObj.challenge
        this.fidoService.getCredential(challenge, userId).then(res => {
          this.authService.loginFinalize(this.username, res).then(res => {
            console.log(res)
            let token = res
            if (this.authService.login(token)) {
              this.router.navigate(["/dashboard"])
            }

          })
        })
      })
  }
}
