import { Injectable } from '@angular/core';
import {Transaction, TransactionList} from "./shared/transaction";
import axios from "axios";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class TransactionService {

  constructor() { }

  public getTransactions() : Promise<any> {
    return axios.get(
      environment.routes.authenticationService + "/get-transactions",
      {
        headers: {
          "Authorization": "Bearer " + localStorage.getItem("token")
        }
      }
    )
  }

  public async transactionInitialize(){

  }

  public async transactionFinalize(jwt: string, receiver : string, amount: number){

  }
}
