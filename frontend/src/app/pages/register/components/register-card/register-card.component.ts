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
    this.authService.registerInitialize(username, username)
      .then(res => {
        let jsonObj = res.data
        this.fidoService.createCredential(jsonObj).then(res => {
          this.authService.registerFinalize(username, res).then(res => {
            this.router.navigate(["/login"])
          })
        })
      })
  }

  displayError(message: string): void {

  }




}
