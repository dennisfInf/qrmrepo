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
        this.fidoService.getCredential(res.data as PublicKeyCredentialRequestOptions).then(res => {
          this.authService.loginFinalize(this.username, res as PublicKeyCredential).then(res => {

            if (this.authService.login(this.name)) {
              this.router.navigate(["/dashboard"])
            }

          })
        })
      })
  }
}
