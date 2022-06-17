import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShopThumbnailComponent } from './shop-thumbnail.component';

describe('ShopThumbnailComponent', () => {
  let component: ShopThumbnailComponent;
  let fixture: ComponentFixture<ShopThumbnailComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ShopThumbnailComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ShopThumbnailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
