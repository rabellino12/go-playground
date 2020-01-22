import Centrifuge = require('centrifuge');
import { Observable } from 'rxjs';
import { WSClient } from './centrifuge';

interface IParams {
	c: WSClient;
	userId: string;
	matchId: string;
}

export interface IMove {
	action: 'up' | 'left' | 'right' | 'stop';
	matchId: string;
	timestamp: number;
	userId: string;
}
export class MovementIO {
	public movements$!: Observable<Centrifuge.PublicationContext>;
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

	public move = (action: 'up' | 'left' | 'right' | 'stop') => {
		const message: IMove = {
			action,
			matchId: this.match,
			timestamp: this.getTime(),
			userId: this.userId
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
			timestamp: this.getTime(),
			userId: this.userId
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
		this.movements$ = new Observable<Centrifuge.PublicationContext>(sub => {
			this.matchSubscription.on('publish', (e: Centrifuge.PublicationContext) => {
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