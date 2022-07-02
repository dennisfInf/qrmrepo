import { Component, OnInit } from '@angular/core';
import {TransactionService} from "../../../services/transaction.service";
import {Transaction} from "../../../services/shared/transaction";
import {UserService} from "../../../services/user.service";
import {EtherscanService} from "../../../services/etherscan.service";
import {Router} from "@angular/router";

@Component({
  selector: 'app-transactions',
  templateUrl: './transactions.component.html',
  styleUrls: ['./transactions.component.css']
})
export class TransactionsComponent implements OnInit {
  etherscanTransactions : any[] = []
  transactions!:Transaction[]
  constructor(private transactionService : TransactionService,
              private userService : UserService,
              private etherscanService : EtherscanService,
              private router : Router) {
    this.etherscanService.getTransactions(this.userService.getAddress(), 0).then(res => {
      this.etherscanTransactions = res.data.result
      console.log(this.etherscanTransactions)
    })
  }

  ngOnInit(): void {
    this.transactions = this.transactionService.getTransactions()
  }
  date(timestamp : string) : string{
    let time : number = +timestamp
    return new Date(time*1000).toLocaleDateString("en-US")
  }

  gweiToEth(gwei : string) : number {
    console.log(gwei)
    return Number(gwei)  * 0.000000000000000001
  }


}
