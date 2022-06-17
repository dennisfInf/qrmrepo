import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-transaction-preview',
  templateUrl: './transaction-preview.component.html',
  styleUrls: ['./transaction-preview.component.css']
})
export class TransactionPreviewComponent implements OnInit {

  @Input("from")
  from!:string

  @Input('to')
  to!:string

  @Input('name')
  name! : string

  @Input('amount')
  amount! : string

  @Input('date')
  date! : string
  constructor() { }

  ngOnInit(): void {
  }

}
