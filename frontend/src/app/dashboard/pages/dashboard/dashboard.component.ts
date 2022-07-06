import {Component, OnInit} from '@angular/core';
import {TransactionService} from "../../../services/transaction.service";
import {ContactsService} from "../../../services/contacts.service";
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

  transactions: Transaction[]
  shops: Shop[]
  contacts: Contact[]
  address: string
  etherscanTransactions: any[] = []

  constructor(private transactionService: TransactionService,
              private contactService: ContactsService,
              private shopService: ShopsService,
              private userService: UserService,
              private etherscanService: EtherscanService,
              private router : Router) {
    this.transactions = transactionService.getTransactions()
    this.contacts = contactService.getContacts()
    this.shops = shopService.getShops()
    this.address = this.userService.getAddress()
    this.etherscanService.getTransactions(this.userService.getAddress(), 0).then(res => {
      this.etherscanTransactions = res.data.result.splice(0, 5)
      console.log(this.etherscanTransactions)
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
}
