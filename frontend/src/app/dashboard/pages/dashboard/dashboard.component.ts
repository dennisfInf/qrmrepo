import {Component, OnInit} from '@angular/core';
import {TransactionService} from "../../../services/transaction.service";
import {ContactResponse, ContactsService} from "../../../services/contacts.service";
import {ShopsService} from "../../../services/shops.service";

import {Shop} from "../../../services/shared/shop";
import {Contact} from "../../../services/shared/contact";
import {UserService} from "../../../services/user.service";
import {EtherscanService} from "../../../services/etherscan.service";
import {timestamp} from "rxjs";
import {Router} from "@angular/router";

class Transaction {
  block_hash!: string
  block_number!: string
  block_timestamp!: string
  from_address!: string
  gas!: string
  gas_price!: string
  hash!: string
  input!: string
  nonce!: string
  receipt_contract_address! : string
  receipt_cumulative_gas_used! : string
  receipt_gas_used! : string
  receipt_root!: string
  receipt_status!: string
  to_address!:string
  transaction_index!: string
  transfer_index! : number[]
  value!: string
}

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  transactions!: Transaction[]
  contacts!: Contact[]
  address: string
  etherscanTransactions: any[] = []
  balance! : number
  contactName: any;
  contactAddress: any;
  pending : boolean = false
  constructor(private transactionService: TransactionService,
              private contactService: ContactsService,
              private shopService: ShopsService,
              private userService: UserService,
              private etherscanService: EtherscanService,
              private router : Router,
             ) {
    contactService.getContacts().then(res => {
      let response = res.data as ContactResponse
      this.contacts = response.contacts
    })
    transactionService.getTransactions().then(res => {
      this.transactions = res.data.transactions as Transaction[]
    })

    this.address = this.userService.getAddress()
    this.etherscanService.getAddressBalance(this.address).then(res => {
      this.balance = this.gweiToEth(res.data.result)
    })

    this.getBalanceTimer()
    this.getTransactionTimer()
  }

  ngOnInit(): void {
  }

  date(timestamp : string) : string{
    let time : number = +timestamp
    return new Date(time*1000).toLocaleDateString("en-US")
  }

  gweiToEth(gwei : string) : number {
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

  getBalanceTimer() {
    setInterval(() => {
      this.etherscanService.getAddressBalance(this.address).then(res => {
        if(this.balance != this.gweiToEth(res.data.result) ) {
          this.balance = this.gweiToEth(res.data.result)
        }
      })
    }, 10000)

  }

  private getTransactionTimer() {
    setInterval( () => {
      this.transactionService.getTransactions().then(res => {
        if(res.data.transactions.length > this.transactions.length) {
          this.transactions = res.data.transactions
        }
      })
    }, 10000)
  }

  shortAddress(address : string) : string {
    let firstPart =  address.slice(0,8)
    let lastPart = address.slice(address.length -4 , address.length)
    return "(" + firstPart + "..." + lastPart + ")"
  }

  getContactName(address : string) : string {
    let res = this.contacts.filter(item => item.public_key.toLowerCase() == address.toLowerCase())
    if(res.length > 0) {
      return res[0].name
    }
    return "Unknown"
  }

  getTimeString(date : string) : string {
    return new Date(date).toISOString().
    replace(/T/, ' ').      // replace T with a space
      replace(/\..+/, '')
  }
}
