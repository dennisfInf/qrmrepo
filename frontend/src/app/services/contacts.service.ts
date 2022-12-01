import { Injectable } from '@angular/core';
import {Contact, ContactList} from "./shared/contact";
import axios from "axios";
import {environment} from "../../environments/environment";
import {data} from "autoprefixer";

export class ContactResponse {
  contacts!: Contact[]
}


@Injectable({
  providedIn: 'root'
})
export class ContactsService {

  constructor() { }

  async getContacts() : Promise<any> {
    return axios.get(
      environment.routes.authenticationService + "/get-contacts",
      {
        headers: {
          "Authorization": "Bearer " + localStorage.getItem("token")
        }
      }
    )
  }

  async createContact(name: string, publicKey: string) : Promise<any> {
    return fetch(environment.routes.authenticationService + "/create-contact",
      {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          "Authorization": "Bearer " + localStorage.getItem("token"),

        },
        method: "POST",
        body: JSON.stringify({
          name: name,
          public_key: publicKey
        }),
      })
  }
}
