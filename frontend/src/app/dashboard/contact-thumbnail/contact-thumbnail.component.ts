import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-contact-thumbnail',
  templateUrl: './contact-thumbnail.component.html',
  styleUrls: ['./contact-thumbnail.component.css']
})
export class ContactThumbnailComponent implements OnInit {
  @Input('address')
  address!: string

  @Input('name')
  name!:string
  constructor() { }

  ngOnInit(): void {
  }

  shortAddress(address : string) : string {
    let firstPart =  address.slice(0,8)
    let lastPart = address.slice(address.length -4 , address.length)
    return "(" + firstPart + "..." + lastPart + ")"
  }

}
