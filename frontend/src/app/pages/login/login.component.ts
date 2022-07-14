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

  constructor(private authService : AuthenticationService, private fidoService : FidoService, private router : Router) { }

  ngOnInit(): void {
  }
  async login(username:string) {
    this.authService.loginInitialize(username)
      .then(res => {
        console.log("moin")
        console.log("login init")
        console.log(res.data)
        this.fidoService.getCredential(res.data).then(res => {
          this.authService.loginFinalize(username, res as PublicKeyCredential).then(res => {

            if (this.authService.login(username)) {
              this.router.navigate(["/dashboard"])
            }

          })
        })
      })
  }
}
