import Centrifuge from 'centrifuge';
import { Observable } from 'rxjs';

export interface IConnectionEvent {
	client: string;
	transport: string;
	latency: number;
}
export interface IDisconnectionEvent {
	reason: string;
	reconnect: boolean;
}

export class WSClient {
	public onConnect$?: Observable<IConnectionEvent>;
	public onDisconnect$?: Observable<IDisconnectionEvent>;
	public cent: Centrifuge;

	constructor() {
		this.cent = new Centrifuge('ws://localhost:8081/connection/websocket', {
			onPrivateSubscribe: this.onPrivateSubscribe,
			subscribeEndpoint: 'http://localhost:8080/auth/centrifuge'
		});
	}

	public connect(token: string): void {
		this.cent.setToken(token);
		this.cent.connect();
		this.onConnect$ = new Observable<IConnectionEvent>(sub => {
			this.cent.on('connect', (
				context: IConnectionEvent /* Couldn't find correct type for this context */
			) => {
				sub.next(context);
			});
		});
		this.onDisconnect$ = new Observable<IDisconnectionEvent>(sub => {
			this.cent.on('disconnect', (context: IDisconnectionEvent) => {
				sub.next(context);
			});
		});
	}

	public subscribe(channel: string): Observable<any> {
		return new Observable<any>(sub => {
			this.cent.subscribe(channel, e => {
				sub.next(e);
			});
		});
	}

	private onPrivateSubscribe(
		{ data }: Centrifuge.SubscribePrivateContext,
		cb: (res: Centrifuge.SubscribePrivateResponse) => void
	): void {
		fetch('http://localhost:8080/auth/centrifuge', {
			body: JSON.stringify(data),
			method: 'POST'
		})
			.then(res => {
				if (!res.ok) {
					throw new Error('Authorization request failed');
				}
				return res.json();
			})
			.then(data => {
				return {
					data,
					status: 200
				};
			})
			.then(cb)
			.catch(console.log);
	}
}
