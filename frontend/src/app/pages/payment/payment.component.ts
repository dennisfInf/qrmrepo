import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {FidoService} from "../../services/fido.service";

@Component({
  selector: 'app-payment',
  templateUrl: './payment.component.html',
  styleUrls: ['./payment.component.css']
})
export class PaymentComponent implements OnInit {
  username: string = "username"
  name: string = "username"
  showError: boolean = false
  error: any = ""
  userId: string = "Ich bin eine UserId"
  credential: ArrayBuffer = Uint8Array.from([
    86,
    -45,
    110,
    -26,
    -41,
    71,
    119,
    -13,
    44,
    2,
    52,
    118,
    -118,
    66,
    -95,
    112,
    96,
    25,
    51,
    -122,
    -125,
    -1,
    -71,
    -103,
    -67,
    -6,
    4,
    88,
    -7,
    21,
    -120,
    -92
  ])
  constructor(private route : ActivatedRoute, private fidoService : FidoService) { }
  amount! :string | null;
  receiver! : string | null;
  ngOnInit(): void {
    this.route.queryParamMap.subscribe(params => {
      this.amount = params.get("amount")
      this.receiver = params.get("receiver")
      window.opener.postMessage("ADRESSE DES KUNDEN", "*")
    })
  }

  async login() {
    if (this.credential != null)
      /*this.fidoService.getCredential("test", this.credential).then(res => {
        this.error = res
      })*/

    this.error = "credential is null"
  }

}
