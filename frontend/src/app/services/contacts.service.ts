import { Injectable } from '@angular/core';
import {Contact, ContactList} from "./shared/contact";

@Injectable({
  providedIn: 'root'
})
export class ContactsService {

  constructor() { }

  public getContacts() : Contact[] {
    return ContactList
  }
}
