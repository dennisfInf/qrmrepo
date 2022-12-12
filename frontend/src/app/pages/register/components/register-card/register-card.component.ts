import {Component, OnInit} from '@angular/core';
import {FidoService} from "../../../../services/fido.service";
import {AuthenticationService} from "../../../../services/authentication.service";
import {timeout} from "rxjs";
import {Router} from "@angular/router";

@Component({
  selector: 'app-register-card',
  templateUrl: './register-card.component.html',
  styleUrls: ['./register-card.component.css']
})
export class RegisterCardComponent implements OnInit {
  showError: boolean = false
  error: any = ""
  userId: string = "Ich bin eine UserId"
  credential: PublicKeyCredential | null = null


  constructor(private fidoService: FidoService,
              private authService: AuthenticationService,
              private router: Router

  ) {

  }

  ngOnInit(): void {
  }

  async register(username:string) {
    if(this.validateEmail(username)) {
      this.authService.registerInitialize(username, username)
        .then(res => {
          let jsonObj = res.data
          this.fidoService.createCredential(jsonObj).then(res => {
            this.authService.registerFinalize(username, res as PublicKeyCredential).then(res => {

              return res.json()
            }).then(data => {
              this.authService.login(data.token)
              this.router.navigate(["/dashboard"])
            })
          })
        }, err => {
          this.error = "Email already registered"
          setTimeout(() => {
            this.error = ""
          }, 5000)
        })
    }else {
      this.error = "Please enter a valid Email"
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
