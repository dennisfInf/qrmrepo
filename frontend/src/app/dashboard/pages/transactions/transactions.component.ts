import { Component, OnInit } from '@angular/core';
import {TransactionService} from "../../../services/transaction.service";
import {Transaction} from "../../../services/shared/transaction";

@Component({
  selector: 'app-transactions',
  templateUrl: './transactions.component.html',
  styleUrls: ['./transactions.component.css']
})
export class TransactionsComponent implements OnInit {

  transactions!:Transaction[]
  constructor(private transactionService : TransactionService) { }

  ngOnInit(): void {
    this.transactions = this.transactionService.getTransactions()
  }

}
