import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

import { environment } from '../../environments/environment';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  public constructor(private router: Router, private http: HttpClient) {
  }

  public ngOnInit() {
  }

  public startNewGame() {
    this.http.post(environment.BACKEND + "/connect6/", {}).subscribe(
      (gameId: number) => {
        console.log(gameId);
        this.router.navigate(['/connect6/', gameId]);
      }
    )
  }
}
