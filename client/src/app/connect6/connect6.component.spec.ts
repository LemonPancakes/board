import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { Connect6Component } from './connect6.component';

describe('Connect6Component', () => {
  let component: Connect6Component;
  let fixture: ComponentFixture<Connect6Component>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ Connect6Component ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(Connect6Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
