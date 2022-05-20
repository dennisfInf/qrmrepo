import {Injectable} from '@angular/core';

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

}
