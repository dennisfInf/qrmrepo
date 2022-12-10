import { Component, OnInit } from '@angular/core';
import { AuthenticationService } from "../../../services/authentication.service";
import { FidoService } from "../../../services/fido.service";
import { Contact } from "../../../services/shared/contact";
import { ContactResponse, ContactsService } from "../../../services/contacts.service";
import {Router} from "@angular/router";
import {timeout} from "rxjs";

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
  error! : string
  success: boolean = false;
  constructor(private authService: AuthenticationService, private fidoService: FidoService, private contactService: ContactsService, private router: Router) {
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
            console.log(res)
            if(res.data.message) {
              // TODO: Error handling
              console.log(res.data)
            }else {
              this.success = true

              setTimeout( () =>{
                this.router.navigate(["/dashboard"])
              },1000)
            }
          }, err => {
            console.log(err)
            this.error = err.response.data.message
          })
        })
      })
  }


  selectContact(address: string) {
    this.receiver = address
  }
}
