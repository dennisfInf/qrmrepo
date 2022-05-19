import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DevelpersComponent } from './develpers.component';

describe('DevelpersComponent', () => {
  let component: DevelpersComponent;
  let fixture: ComponentFixture<DevelpersComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DevelpersComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(DevelpersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
