import { Component, OnInit } from '@angular/core';
import {Transaction} from "../../../services/shared/transaction";
import {TransactionService} from "../../../services/transaction.service";

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.css']
})
export class TransactionComponent implements OnInit {
  transactions! : Transaction[]
  constructor(private transactionService : TransactionService) { }

  ngOnInit(): void {
    this.transactions = this.transactionService.getTransactions()
  }

}
