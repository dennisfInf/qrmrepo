import {Component, OnInit} from '@angular/core';
import {FidoService} from "../../../../services/fido.service";
import {AuthenticationService} from "../../../../services/authentication.service";

@Component({
  selector: 'app-register-card',
  templateUrl: './register-card.component.html',
  styleUrls: ['./register-card.component.css']
})
export class RegisterCardComponent implements OnInit {

  username : string = ""
  name : string = ""
  userId : string = "Ich bin eine UserId"
  constructor(private fidoService : FidoService, private authService : AuthenticationService) {

  }

  ngOnInit(): void {
  }

  async register() {
      this.authService.registerInitialize(this.username).then(res => {
        this.fidoService.createCredential(res,this.username,this.userId,this.name).then(res => {
          console.log(res)
        })
      })
  }

}
