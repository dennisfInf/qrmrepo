import { Component, OnInit } from '@angular/core';
import { AuthenticationService } from "../../services/authentication.service";
import { FidoService } from "../../services/fido.service";
import { Router } from "@angular/router";
import { UserService } from "../../services/user.service";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  error = ""
  constructor(private authService: AuthenticationService, private fidoService: FidoService, private router: Router, private userSerivce: UserService) {
    if (this.authService.isAuthenticated()) {
      this.router.navigate(["/dashboard"])
    }
  }

  ngOnInit(): void {
  }


  async login(username: string, password: string) {
    if (this.validateEmail(username)) {
      if (username == localStorage.getItem('email') && password == localStorage.getItem('password')) {
        this.router.navigate(["/dashboard"])
      } else {
        this.error = "wrong email or password"
        setTimeout(() => {
          this.error = ""
        }, 5000)
      }
      localStorage.setItem('password', password);
    } else {
      this.error = "Enter a valid email address"
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
}
