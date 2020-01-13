import Centrifuge = require('centrifuge');
import { Observable } from 'rxjs';
import { WSClient } from './centrifuge';

interface IParams {
	c: WSClient;
	userId: string;
	matchId: string;
}

interface IMove {
	action: string;
	matchId: string;
	timestamp: number;
}
export class MovementIO {
	public movements$!: Observable<IMove>;
	public enemies$!: Observable<Array<string | undefined>>;
	public matchSubscription: Centrifuge.Subscription;
	private c: WSClient;
	private userId: string;
	private match: string;
	private lastMove?: string;

	constructor(params: IParams) {
		this.c = params.c;
		this.userId = params.userId;
		this.match = `$match:${params.matchId}`;
		this.matchSubscription = this.c.cent.subscribe(this.match);
		this.initializeEvents();
	}

	public move = (action: string) => {
		const message: IMove = {
			action,
			matchId: this.match,
			timestamp: this.getTime()
		};
		if (this.lastMove !== action) {
			this.matchSubscription.publish(message);
			this.lastMove = action;
		}
	};

	public stop = () => {
		const message: IMove = {
			action: 'stop',
			matchId: this.match,
			timestamp: this.getTime()
		};
		if (this.lastMove !== 'stop') {
			this.matchSubscription.publish(message);
			this.lastMove = 'stop';
		}
	};

	private getTime() {
		return new Date().getTime();
	}

	private initializeEvents = () => {
		console.log(this.matchSubscription);
		this.movements$ = new Observable<IMove>(sub => {
			this.matchSubscription.on('publish', (e: IMove) => {
				sub.next(e);
			});
		});
		this.enemies$ = new Observable<Array<string | undefined>>(sub => {
			this.matchSubscription.on('subscribe', e => {
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
