import { Component, OnInit } from '@angular/core';
import { AuthenticationService } from "../../../services/authentication.service";
import { FidoService } from "../../../services/fido.service";
import { Contact } from "../../../services/shared/contact";
import { ContactResponse, ContactsService } from "../../../services/contacts.service";
import {Router} from "@angular/router";
import {timeout} from "rxjs";
import {EtherscanService} from "../../../services/etherscan.service";


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
  gasPrice!: string
  constructor(private authService: AuthenticationService,
              private fidoService: FidoService,
              private contactService: ContactsService,
              private router: Router,
              private etherScanService : EtherscanService) {
    this.token = authService.getToken()
  }

  ngOnInit(): void {
    this.contactService.getContacts().then(res => {
      let result = res.data as ContactResponse
      this.contacts = result.contacts

    })

    this.etherScanService.getGasPrice().then(res => {
      console.log(res)
      let price = +res.data.result.SafeGasPrice
      let total = price * 21000
      this.gasPrice =  this.gweiToEth(total+"") + ""
      console.log(this.gasPrice)
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

  shortAddress(address : string) : string {
    let firstPart =  address.slice(0,8)
    let lastPart = address.slice(address.length -4 , address.length)
    return "(" + firstPart + "..." + lastPart + ")"
  }

  getTotal(amount: string, fees: string) : string {
    let a = +amount
    let b = +fees
    return a+ b + ""
  }

  gweiToEth(gwei : string) : number {
    return Number(gwei)  * 0.000000000000000001
  }

}
