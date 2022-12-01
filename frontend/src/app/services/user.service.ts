import { Injectable } from '@angular/core';
import {JwtHelperService} from "@auth0/angular-jwt";

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private jwtHelper: JwtHelperService
  constructor() {
    this.jwtHelper = new JwtHelperService()
  }


  getAddress() : string {
    let token = localStorage.getItem("token")

    if(token != null) {
      let obj = this.jwtHelper.decodeToken(token)
      return obj["pub-key"]
    }
    return ""
  }

  getUserName() : string {
    let token = localStorage.getItem("token")
    if(token != null) {
      let obj = this.jwtHelper.decodeToken(token)
      return obj["username"]
    }
    return ""
  }
}
