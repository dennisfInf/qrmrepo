import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-contact-thumbnail',
  templateUrl: './contact-thumbnail.component.html',
  styleUrls: ['./contact-thumbnail.component.css']
})
export class ContactThumbnailComponent implements OnInit {
  @Input('imgUrl')
  imgUrl!: string

  @Input('name')
  name!:string
  constructor() { }

  ngOnInit(): void {
  }

}
