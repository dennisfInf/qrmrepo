import { Component, OnInit } from '@angular/core';
import {TransactionService} from "../../../services/transaction.service";
import {ContactsService} from "../../../services/contacts.service";
import {ShopsService} from "../../../services/shops.service";
import {Transaction} from "../../../services/shared/transaction";
import {Shop} from "../../../services/shared/shop";
import {Contact} from "../../../services/shared/contact";

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  transactions: Transaction[]
  shops: Shop[]
  contacts : Contact[]

  constructor(private transactionService : TransactionService,
              private contactService : ContactsService,
              private shopService : ShopsService)
  {
   this.transactions = transactionService.getTransactions()
   this.contacts = contactService.getContacts()
   this.shops = shopService.getShops()
  }

  ngOnInit(): void {
  }

}
