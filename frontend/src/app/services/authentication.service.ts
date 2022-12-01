import { Injectable } from '@angular/core';
import axios from "axios";
import { environment } from "../../environments/environment";
import { JwtHelperService } from "@auth0/angular-jwt";
import { buffer } from 'rxjs';

function bufferEncode(value:ArrayBuffer) {
  return btoa(String.fromCharCode(...new Uint8Array(value))).replace(/\+/g, "-").replace(/\//g, "_").replace(/=/g, "")
 /* var u8 = new Uint8Array(value);
  var decoder = new TextDecoder('utf8');
  var b64encoded = btoa(decoder.decode(u8));
  return b64encoded*/
}

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {
  private jwtHelper: JwtHelperService
  constructor() {
    this.jwtHelper = new JwtHelperService()
  }

  async challenge(): Promise<string> {
    return "challenge"
  }

  async response(credential: Credential): Promise<any> {

  }

  async registerInitialize(username: string, name: string): Promise<any> {
    return axios.get(
      environment.routes.authenticationService + "/register-init",
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  async registerFinalize(username: string, token: PublicKeyCredential): Promise<Response> {
    const authAttRes = token.response as AuthenticatorAttestationResponse
    return fetch(environment.routes.authenticationService + "/register-finalize",
      {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'x-username': username
        },
        method: "POST",
        body: JSON.stringify({
          id: token.id,
          rawId: bufferEncode(token.rawId),
          type: token.type,
          response: {
            attestationObject: bufferEncode(authAttRes.attestationObject),
            clientDataJSON: bufferEncode(token.response.clientDataJSON),
          },
        }),
      })

   /* return axios.post(
      environment.routes.authenticationService + "/register/finalize",
      token,
      {
        headers: {
          "x-username": username
        }
      }
    )*/
  }

  async loginInitialize(username: string): Promise<any> {
    return axios.get(
      environment.routes.authenticationService + "/login-init",
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  async loginFinalize(username: string, token: PublicKeyCredential): Promise<Response> {
    console.log(token)
    const assertionResponse = token.response as AuthenticatorAssertionResponse
    return fetch(environment.routes.authenticationService + "/login-finalize",
      {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'x-username': username
        },
        method: "POST",
        body: JSON.stringify({
          id: token.id,
          rawId: bufferEncode(token.rawId),
          type: token.type,
          response: {
            authenticatorData: bufferEncode(assertionResponse.authenticatorData),
            clientDataJSON: bufferEncode(assertionResponse.clientDataJSON),
            signature: bufferEncode(assertionResponse.signature),
            userHandle: bufferEncode(assertionResponse.userHandle??new ArrayBuffer(0)),
          },
      }),
    })
    /*return axios.post(
      environment.routes.authenticationService + "/login/finalize",
      token,
      {
        headers: {
          "x-username": username
        }
      }
    )*/
  }

  public isAuthenticated(): boolean {
    const token = localStorage.getItem('token');
    if (token == null) {
      return false
    }
    return true;
  }

  public login(token: string): boolean {
    console.log("logging in")
    console.log(token)
    localStorage.setItem("token", token)
    return true
  }

  public getToken(): string {
    const token = localStorage.getItem('token');
    if (token == null) {
      return ""
    }
    return token
  }


  async transactionInitialize(username: string, amount: string, receiver: string): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/transaction-init",
      JSON.stringify(
        {
          username: username,
          amount: amount,
          receiver: receiver
        }),
      {
        headers: {
          "x-username": username,
        }
      }
    )
  }

  async transactionFinalize(username: string, token: any): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/transaction-finalize",
      token
      ,
      {
        headers: {
          "x-username": username,
        }
      }
    )
  }


  logout() {
    localStorage.removeItem("token")
  }
}

