import { Component, OnInit } from '@angular/core';
import {BlogService} from "../../services/blog.service";
import {Post} from "../../services/shared/blog";

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  posts : Post[];
  constructor(private blogService : BlogService) {
    this.posts = blogService.getPosts()
  }

  ngOnInit(): void {
  }

}
