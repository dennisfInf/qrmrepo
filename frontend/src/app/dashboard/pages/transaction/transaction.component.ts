import { Component, OnInit } from '@angular/core';
import {Transaction} from "../../../services/shared/transaction";
import {TransactionService} from "../../../services/transaction.service";
import {ActivatedRoute} from "@angular/router";
import {EtherscanService} from "../../../services/etherscan.service";

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.css']
})
export class TransactionComponent implements OnInit {
  transactionHash : string | null = ""
  transactions : any[] = []
  constructor(private transactionService : TransactionService,
              private router : ActivatedRoute,
              private etherscanService : EtherscanService
  ) { }

  ngOnInit(): void {
    this.transactions = this.transactionService.getTransactions()
    this.router.paramMap.subscribe( paramMap => {
      this.transactionHash = paramMap.get('hash');
      if(this.transactionHash != null) {
        this.etherscanService
      }
    })
  }

}
