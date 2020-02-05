import Centrifuge from "centrifuge";
import { Observable } from "rxjs";
import { WSClient } from "./centrifuge";

export interface IPosition {
  x: number;
  y: number;
}
export interface IPlayer {
  id: string;
  position: IPosition;
}
export interface IMove {
  userId: IPlayer["id"];
  action: "left" | "right" | "stop";
  position: IPosition;
  matchId: string;
  timestamp: number;
  jumping?: boolean;
}
export interface IParams {
  c: WSClient;
  userId: string;
  matchId: string;
}

export class MovementIO {
  public movements$!: Observable<Centrifuge.PublicationContext>;
  public enemies$!: Observable<Array<string | undefined>>;
  public matchSubscription: Centrifuge.Subscription;
  private c: WSClient;
  private userId: string;
  private match: string;
  private lastMove?: IMove;

  constructor(params: IParams) {
    this.c = params.c;
    this.userId = params.userId;
    this.match = `$match:${params.matchId}`;
    this.matchSubscription = this.c.cent.subscribe(this.match);
    this.initializeEvents();
  }

  public move = (action: "left" | "right" | "stop", position: IPosition, jumping?: boolean) => {
    const message: IMove = {
      action,
      matchId: this.match,
      timestamp: this.getTime(),
      userId: this.userId,
      position,
      jumping
    };
    if (this.lastMove?.action !== action || this.lastMove?.jumping !== jumping) {
      this.matchSubscription.publish(message);
      this.lastMove = message;
    }
  };

  public stop = (position: IPosition, jumping?: boolean) => {
    const message: IMove = {
      action: "stop",
      matchId: this.match,
      timestamp: this.getTime(),
      userId: this.userId,
      position
    };
    if (this.lastMove?.action !== "stop" && this.lastMove?.jumping !== jumping) {
      this.matchSubscription.publish(message);
      this.lastMove = message;
    }
  };

  private getTime() {
    return new Date().getTime();
  }

  private initializeEvents = () => {
    this.movements$ = new Observable<Centrifuge.PublicationContext>(sub => {
      this.matchSubscription.on(
        "publish",
        (e: Centrifuge.PublicationContext) => {
          sub.next(e);
        }
      );
    });
    this.enemies$ = new Observable<Array<string | undefined>>(sub => {
      this.matchSubscription.on("subscribe", e => {
        this.matchSubscription.presence().then(res => {
          const enemies = Object.keys(res.presence).map(
            key => res.presence[key].user
          );
          sub.next(enemies);
        });
      });
    });
  };
}
