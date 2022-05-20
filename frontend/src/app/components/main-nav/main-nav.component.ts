import { Component, OnInit } from '@angular/core';


@Component({
  selector: 'app-main-nav',
  templateUrl: './main-nav.component.html',
  styleUrls: ['./main-nav.component.css']
})
export class MainNavComponent implements OnInit {
  profileMenu: string = 'hidden';
  mainMenu : string = 'hidden';
  constructor() { }

  ngOnInit(): void {
  }

  toggleProfileMenu(){
    if(this.profileMenu == 'hidden'){
      this.profileMenu = ''
    }else {
      this.profileMenu = 'hidden'
    }
  }

  closeMainMenu(){
    setTimeout(()=>{
      this.mainMenu = 'hidden'
    }, 100)
  }

  closeProfileMenu(){
    setTimeout(()=>{
      this.profileMenu = 'hidden'
    }, 100)
  }

  toggleMainMenu() {
    if(this.mainMenu == 'hidden'){
      this.mainMenu = ''
    }else {
      this.mainMenu = 'hidden'
    }
  }
}
