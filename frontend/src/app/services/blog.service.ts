import { Injectable } from '@angular/core';
import {Post} from "./shared/blog";

@Injectable({
  providedIn: 'root'
})
export class BlogService {

  constructor() {

  }

  public getPosts() : Post[] {
    return [
      {
        title : "Passwordless sign in with FIDO2",
        preview: "",
        url : "",
        content: ""
      },
      {
        title : "Secure key management with Intel Software Guard Extension",
        preview: "",
        url : "",
        content: ""
      },
      {
        title : "The current state of the project",
        preview: "",
        url : "",
        content: ""
      },
      {
        title : "How to enhance your business",
        preview: "",
        url : "",
        content: ""
      }
    ]
  }
}
