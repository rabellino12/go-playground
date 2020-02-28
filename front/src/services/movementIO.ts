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

export interface IMoveInput {
  action: "left" | "right" | "stop";
  position: IPosition;
  jumping?: boolean;
}

export interface IMove extends IMoveInput {
  userId: IPlayer["id"];
  timestamp: number;
}
export interface IParams {
  c: WSClient;
  userId: string;
  matchId: string;
}

export class MovementIO {
  public snapshot$!: Observable<Centrifuge.PublicationContext>;
  public enemies$!: Observable<Array<string | undefined>>;
  public matchSubscription: Centrifuge.Subscription;
  public snapshotSubscription: Centrifuge.Subscription;
  private c: WSClient;
  private userId: string;
  private matchChannel: string;
  private snapshotChannel: string;
  private lastMove?: IMove;

  constructor(params: IParams) {
    this.c = params.c;
    this.userId = params.userId;
    this.matchChannel = `$match:${params.matchId}`;
    this.snapshotChannel = `$snapshot:${params.matchId}`;
    this.matchSubscription = this.c.cent.subscribe(this.matchChannel);
    this.snapshotSubscription = this.c.cent.subscribe(this.snapshotChannel);
    this.initializeEvents();
  }

  public move = (move: IMoveInput) => {
    const message: IMove = {
      ...move,
      timestamp: this.getTime(),
      userId: this.userId
    };
    if (this.lastMove?.action !== move.action || this.lastMove?.jumping !== move.jumping) {
      return this.matchSubscription.publish(message)
        .then((res) => {
          this.lastMove = message;
        })
        .catch(console.log);
    }
  };

  private getTime() {
    return new Date().getTime();
  }

  private initializeEvents = () => {
    this.snapshotSubscription.on('error', console.log);
    this.snapshotSubscription.on('unsubscribe', console.log);
    this.snapshotSubscription.on('subscribe', console.log);
    this.snapshot$ = new Observable<Centrifuge.PublicationContext>(sub => {
      this.snapshotSubscription.on(
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
