import { Injectable } from '@angular/core';
import {Shop, ShopList} from "./shared/shop";

@Injectable({
  providedIn: 'root'
})
export class ShopsService {

  constructor() { }

  public getShops() : Shop[] {
    return ShopList
  }
}
