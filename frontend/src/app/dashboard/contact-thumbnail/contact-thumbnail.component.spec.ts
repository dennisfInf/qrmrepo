import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ContactThumbnailComponent } from './contact-thumbnail.component';

describe('ContactThumbnailComponent', () => {
  let component: ContactThumbnailComponent;
  let fixture: ComponentFixture<ContactThumbnailComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ContactThumbnailComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ContactThumbnailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
