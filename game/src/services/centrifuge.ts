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
	private client: Centrifuge;

	constructor() {
		this.client = new Centrifuge('ws://localhost:8081/connection/websocket');
	}
	public connect(token: string): void {
		this.client.setToken(token);
		this.client.connect();
		this.onConnect$ = new Observable<IConnectionEvent>(sub => {
			this.client.on('connect', (
				context: IConnectionEvent /* Couldn't find correct type for this context */
			) => {
				sub.next(context);
			});
		});
		this.onDisconnect$ = new Observable<IDisconnectionEvent>(sub => {
			this.client.on('disconnect', (context: IDisconnectionEvent) => {
				sub.next(context);
			});
		});
	}
}
