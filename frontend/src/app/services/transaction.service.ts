import { Injectable } from '@angular/core';
import {Transaction, TransactionList} from "./shared/transaction";

@Injectable({
  providedIn: 'root'
})
export class TransactionService {

  constructor() { }

  public getTransactions() : Transaction[] {
    return TransactionList
  }

  public async transactionInitialize(){

  }

  public async transactionFinalize(jwt: string, receiver : string, amount: number){

  }
}
