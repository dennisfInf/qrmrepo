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

  username: string = "username"
  name: string = "username"
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

  async register() {
    this.authService.registerInitialize(this.username, this.name)
      .then(res => {
        let jsonObj = JSON.parse(res)
        this.userId = jsonObj.user.id
        let challenge = jsonObj.challenge
        this.fidoService.createCredential(challenge, this.username, this.userId, this.name).then(res => {
          this.authService.registerFinalize(this.username, res).then(res => {
            console.log(res)
          })
        })
      })
  }

  displayError(message: string): void {

  }




}
