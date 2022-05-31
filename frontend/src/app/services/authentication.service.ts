import {Injectable} from '@angular/core';
import axios from "axios";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {

  constructor() {
  }

  async challenge(): Promise<string> {
    return "challenge"
  }

  async response(credential : Credential) : Promise<any> {

  }

  async registerInitialize(username : string): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/register/initialize",
      null,
      {
        headers : {
          "x-username" : username
        }
      }
    )
  }

  async registerFinalize(username: string,token : any) : Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/register/finalize",
      token,
      {
        headers : {
          "x-username" : username
        }
      }
    )
  }

  async loginInitialize(username : string) : Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/login/initialize",
      null,
      {
        headers : {
          "x-username" : username
        }
      }
    )
  }

  async loginFinalize(username : string, token : any) : Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/login/finalize",
      token,
      {
        headers : {
          "x-username" : username
        }
      }
    )
  }

}
