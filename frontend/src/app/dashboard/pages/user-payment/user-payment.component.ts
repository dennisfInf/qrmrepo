import { Component, OnInit } from '@angular/core';
import { AuthenticationService } from "../../../services/authentication.service";
import { FidoService } from "../../../services/fido.service";
import { Contact } from "../../../services/shared/contact";
import { ContactResponse, ContactsService } from "../../../services/contacts.service";

@Component({
  selector: 'app-user-payment',
  templateUrl: './user-payment.component.html',
  styleUrls: ['./user-payment.component.css']
})
export class UserPaymentComponent implements OnInit {
  receiver!: string;
  amount!: string;
  token: string
  contacts!: Contact[]
  constructor(private authService: AuthenticationService, private fidoService: FidoService, private contactService: ContactsService) {
    this.token = authService.getToken()
  }

  ngOnInit(): void {
    this.contactService.getContacts().then(res => {
      let result = res.data as ContactResponse
      this.contacts = result.contacts
    })
  }

  async makePayment() {

    this.authService.transactionInitialize(this.token)
      .then(res => {
        this.fidoService.getCredential(res.data).then(res => {
          this.authService.transactionFinalize(this.receiver, this.token, this.amount, res as PublicKeyCredential).then(res => {
            console.log(res.data.transaction_hash)
          })
        })
      })
  }


  selectContact(address: string) {
    this.receiver = address
  }
}
