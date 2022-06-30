import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor() { }


  getAddress() : string {
    return "0xD3a74341aDAc943De6600468393Bb6Ca4431A7Fd"
  }


}
