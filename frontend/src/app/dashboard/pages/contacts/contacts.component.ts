import { Component, OnInit } from '@angular/core';
import {ContactsService} from "../../../services/contacts.service";
import {Contact} from "../../../services/shared/contact";

@Component({
  selector: 'app-contacts',
  templateUrl: './contacts.component.html',
  styleUrls: ['./contacts.component.css']
})
export class ContactsComponent implements OnInit {

  contacts! : Contact[]
  constructor(private contactService : ContactsService) { }

  ngOnInit(): void {
    this.contacts = this.contactService.getContacts()
  }

}
