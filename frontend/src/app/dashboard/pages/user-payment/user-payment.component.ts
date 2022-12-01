import { Component, OnInit } from '@angular/core';
import {AuthenticationService} from "../../../services/authentication.service";
import {FidoService} from "../../../services/fido.service";
import {Contact} from "../../../services/shared/contact";
import {ContactResponse, ContactsService} from "../../../services/contacts.service";

@Component({
  selector: 'app-user-payment',
  templateUrl: './user-payment.component.html',
  styleUrls: ['./user-payment.component.css']
})
export class UserPaymentComponent implements OnInit {
  receiver!: string;
  amount! : string;
  username : string
  contacts!: Contact[]
  constructor(private authService : AuthenticationService, private fidoService : FidoService, private contactService: ContactsService) {
    this.username = authService.getToken()
  }

  ngOnInit(): void {
     this.contactService.getContacts().then(res => {
       let result = res.data as ContactResponse
       this.contacts = result.contacts
     })
  }

  async makePayment() {

    this.authService.transactionInitialize(this.username, this.amount, this.receiver)
      .then(res => {
        let jsonObj = JSON.parse(res)
        let userId = jsonObj.user.id
        let challenge = jsonObj.challenge
        /* this.fidoService.getCredential(challenge, userId).then(res => {
          this.authService.transactionFinalize(this.username, res).then(res => {
            console.log(res)
          })
        })*/
      })
  }


  selectContact(address : string) {
    this.receiver = address
  }
}
