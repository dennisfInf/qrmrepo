import {Component, OnInit} from '@angular/core';
import {TransactionService} from "../../../services/transaction.service";
import {ContactResponse, ContactsService} from "../../../services/contacts.service";
import {ShopsService} from "../../../services/shops.service";
import {Transaction} from "../../../services/shared/transaction";
import {Shop} from "../../../services/shared/shop";
import {Contact} from "../../../services/shared/contact";
import {UserService} from "../../../services/user.service";
import {EtherscanService} from "../../../services/etherscan.service";
import {timestamp} from "rxjs";
import {Router} from "@angular/router";

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  transactions!: any[]
  contacts!: Contact[]
  address: string
  etherscanTransactions: any[] = []
  balance! : number
  contactName: any;
  contactAddress: any;
  constructor(private transactionService: TransactionService,
              private contactService: ContactsService,
              private shopService: ShopsService,
              private userService: UserService,
              private etherscanService: EtherscanService,
              private router : Router) {
    transactionService.getTransactions().then(res => {
      console.log(res)
      this.transactions = res.data.transactions
    })
    contactService.getContacts().then(res => {
      console.log("contacts")
      console.log(res)
      let response = res.data as ContactResponse
        this.contacts = response.contacts
    })
    this.address = this.userService.getAddress()
    this.etherscanService.getAddressBalance(this.address).then(res => {
      console.log(res.data.result)
      this.balance = this.gweiToEth(res.data.result)
    })
  }

  ngOnInit(): void {
  }

  date(timestamp : string) : string{
    let time : number = +timestamp
    return new Date(time*1000).toLocaleDateString("en-US")
  }

  gweiToEth(gwei : string) : number {
    console.log(gwei)
    return Number(gwei)  * 0.000000000000000001
  }

  route() {
    this.router.navigate(["/dashboard/user-payment"])
  }


  addContact() {
    if(this.contactName != "" && this.contactAddress != "") {
      this.contactService.createContact(this.contactName, this.contactAddress).then(res => {
        console.log(res)
        this.contactService.getContacts().then(res => {
          console.log("contacts")
          console.log(res)
          let response = res.data as ContactResponse
          this.contacts = response.contacts
        })
      })
    }
  }

  goToEtherscan(hash: string) {

    let address = this.etherscanService.getTransactionLink(hash)
    document.location.href = address;
  }

  getValueFromHex(hex : string) : number {
    let dec = parseInt(hex, 16)
    return this.gweiToEth(dec+ "")
  }
}
