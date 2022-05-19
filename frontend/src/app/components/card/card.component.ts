import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-card',
  templateUrl: './card.component.html',
  styleUrls: ['./card.component.css']
})
export class CardComponent implements OnInit {
  @Input() title: string
  @Input() content: string
  @Input() target : string
    @Input() justify : string
  constructor() {
    this.title = 'Headline'
    this.content = 'Bacon ipsum dolor amet tail filet mignon buffalo cupim burgdoggen, bresaola spare ribs boudin turkey flank capicola. Shank rump shankle burgdoggen. T-bone frankfurter picanha ball tip tongue boudin sirloin meatball leberkas kielbasa jerky porchetta filet mignon salami. Meatloaf boudin tongue shoulder burgdoggen porchetta. Corned beef jowl kielbasa, bacon sausage porchetta flank pancetta ham hock boudin andouille jerky meatball tenderloin shoulder.'

    this.target = 'fasdf'
    this.justify = 'justify-end'
  }

  ngOnInit(): void {
  }

}
