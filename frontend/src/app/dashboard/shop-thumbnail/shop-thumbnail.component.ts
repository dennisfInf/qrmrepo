import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-shop-thumbnail',
  templateUrl: './shop-thumbnail.component.html',
  styleUrls: ['./shop-thumbnail.component.css']
})
export class ShopThumbnailComponent implements OnInit {
  @Input("name")
  name!: string
  constructor() { }

  ngOnInit(): void {
  }

}
