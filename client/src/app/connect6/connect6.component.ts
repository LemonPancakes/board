import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { SocketService } from "../socket.service";

@Component({
  selector: 'app-game',
  templateUrl: './connect6.component.html',
  styleUrls: ['./connect6.component.css']
})
export class Connect6Component implements OnInit {
  public board: Array<Array<number>>;
  public player: number;
  public currentPlayer: number;
  public firstMove: boolean;
  public finished: boolean = false;

  public error: string;
  public info: string;

  public constructor(private socket: SocketService, private route: ActivatedRoute) {
    this.board = [];
    for (var i = 0; i < 19; i++) {
      this.board[i] = [];
      for (var j = 0; j < 19; j++) {
        this.board[i][j] = 0;
      }
    }
  }

  public ngOnInit() {
    this.socket.getEventListener().subscribe(event => {
      if (event.type == "message") {
        this.parseMessage(event.data);
      }
      if (event.type == "close") {
        console.log("connection closed");
        this.error = "connection closed";
      }
      if (event.type == "open") {
        console.log("connection opened");
        this.route.params.subscribe((params) => {
          console.log(params);
          this.socket.send(params[0]);
        });
      }
    });
  }

  public ngOnDestroy() {
    this.socket.close();
  }

  public parseMessage(message: any) {
    let type = message.type;
    let data = message.content;
    console.log(type, data);
    switch (type) {
      case 'GameState':
        var board: number[];
        [this.player, this.currentPlayer, this.firstMove, ...board] = data.split(',');
        this.finished = false;

        for (let i = 0; i < 19; i++) {
          for (let j = 0; j < 19; j++) {
            this.board[i][j] = board[i * 19 + j]
          }
        }

        break;
      case 'Move':
        var [i, j, p] = data.split(',');
        this.board[i][j] = p;
        break;
      case 'CurrentPlayer':
        this.currentPlayer = data;
        break;
      case 'Finished':
        this.finished = true;
        break;
      case 'NewGame':
        var board: number[];
        [this.currentPlayer, this.firstMove, ...board] = data.split(',');
        this.finished = false;

        for (let i = 0; i < 19; i++) {
          for (let j = 0; j < 19; j++) {
            this.board[i][j] = 0
          }
        }
        break;
      case "Resign":
        this.finished = true;
        this.info = "player " + (3 - data) + " won!"
        break;
      default:
        console.log("unknown");
        break;
    }
  }

  public makeMove(i: number, j: number) {
    if (this.finished) {
      this.setError("Game is finished");
      return;
    }

    if (this.player != this.currentPlayer) {
      this.setError("you're player: " + this.player + ", not the current player: " + this.currentPlayer)
      return;
    }

    this.board[i][j] = this.player;
    this.socket.send(i + ',' + j);
  }

  public dismissAlert() {
    this.error = "";
  }

  public setError(err: string) {
    this.error = err;
    setTimeout(() => {
      this.error = "";
    }, 2000);
  }

  public newGame() {
    this.info = "";
    if (!this.finished) {
      this.setError("Game is not finished yet");
      return;
    }

    this.socket.send("NewGame");
  }

  public resign() {
    this.socket.send("Resign");
  }

  public playerStatus(): string {
    if (this.player == 1 || this.player == 2) {
      return `Player ${this.player}`
    } else {
      return "Spectator"
    }
  }
}
