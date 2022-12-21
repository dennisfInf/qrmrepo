import { Component, OnInit } from '@angular/core';
import { FidoService } from "../../../../services/fido.service";
import { AuthenticationService } from "../../../../services/authentication.service";
import { timeout } from "rxjs";
import { Router } from "@angular/router";

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

  async register(username: string, password: string) {
    if (this.validateEmail(username)) {
      if (this.checkPassword(password)) {
        localStorage.setItem('email', username);
        localStorage.setItem('password', password);
        this.router.navigate(["/dashboard"])
      } else {
        this.error = "Password must contain at least one numeric digit, one uppercase, one lowercase letter and 6 characters long"
        setTimeout(() => {
          this.error = ""
        }, 5000)
      }
    } else {
      this.error = "Please enter a valid Email"
      setTimeout(() => {
        this.error = ""
      }, 5000)
    }

  }


  validateEmail(email: string) {
    return String(email)
      .toLowerCase()
      .match(
        /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
      );
  };

  checkPassword(password: string) {
    var passw = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,20}$/;
    if (password.match(passw)) {
      return true;
    }
    else {
      return false;
    }
  }



}
