import {Injectable} from '@angular/core';
import axios from "axios";
import {environment} from "../../environments/environment";
import * as moment from "moment";
@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {

  constructor() {
  }

  async challenge(): Promise<string> {
    return "challenge"
  }

  async response(credential: Credential): Promise<any> {

  }

  async registerInitialize(username: string, name: string): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/register/initialize",
      JSON.stringify({username: username, name: name}),
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  async registerFinalize(username: string, token: any): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/register/finalize",
      token,
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  async loginInitialize(username: string): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/login/initialize",
      null,
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  async loginFinalize(username: string, token: any): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/login/finalize",
      token,
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  public isLoggedIn() {
    return moment().isBefore(this.getExpiration());
  }

  getUserId(): string {
    return ""
  }

  getRole(): string {
    return ""
  }

  getToken(): string {
    return ""
  }

  private setSession(authResult : any) {
    const expiresAt = moment().add(authResult.expiresIn,'second');

    localStorage.setItem('id_token', authResult.idToken);
    localStorage.setItem("expires_at", JSON.stringify(expiresAt.valueOf()) );
  }

  logout() {
    localStorage.removeItem("id_token");
    localStorage.removeItem("expires_at");
  }


  isLoggedOut() {
    return !this.isLoggedIn();
  }

  getExpiration() {
    const expiration = localStorage.getItem("expires_at");
    let expiresAt = ""
    if (expiration != null) {
       expiresAt = JSON.parse(expiration);

    }
    return moment(expiresAt);
  }

}
