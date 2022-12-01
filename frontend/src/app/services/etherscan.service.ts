import { Injectable } from '@angular/core';
import axios, {AxiosResponse} from "axios";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class EtherscanService {
  apiKey : string = "4PPXSAC2FQTKNEDKXI23V5MQ87PR4RB8IV"
  constructor() { }

  public async getAddressBalance(address : string) :  Promise<AxiosResponse<any, any>>{
    let url = "https://api-goerli.etherscan.io/api\n" +
      "?module=account" +
      "&action=balance" +
      "&address=" + address +
      "&tag=latest" +
      "&apikey="+this.apiKey
    return axios.get(url)
  }

  public async getTransactions(address : string, offset:number){
    let url = "https://api-goerli.etherscan.io/api" +
      "?module=account" +
      "&action=txlist" +
      "&address=" + address +
      "&startblock=0" +
      "&endblock=99999999" +
      "&page=1" +
      "&offset="+ offset +
      "&sort=asc" +
      "&apikey=" + this.apiKey

    return axios.get(url)
  }

  public async getTransactionReceiptStatus(txHash : string){
    let url = "https://api-goerli.etherscan.io/api" +
      "?module=transaction" +
      "&action=gettxreceiptstatus" +
      "&txhash=" + txHash+
      "&apikey=" +this.apiKey
    return axios.get(url)
  }

  public getTransactionLink(txHash : string) {
    return "https://goerli.etherscan.io/tx/" + txHash
  }
}
