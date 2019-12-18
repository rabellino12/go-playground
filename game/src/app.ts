import Phaser from 'phaser';
import { Subscription, Observable } from 'rxjs';

import { StartScene } from './scenes/Start';
import { fetchJWT } from './services/auth';
import { WSClient } from './services/centrifuge';

const config = {
	height: 600,
	type: Phaser.AUTO,
	width: 800
};

interface Lobby {
	status: string;
	id?: string;
}

class PhaserApp extends Phaser.Game {
	private client: WSClient;
	private token?: string;
	// private lobby$: Observable<any>;

	constructor() {
		super(config);
		this.client = new WSClient();
		this.getToken()
			.then(token => {
				this.token = token;
				this.client.connect(this.token);
				this.scene.add('Start', StartScene);
				this.scene.start('Start', {
					wsClient: this.client
				});
				this.listenEvents();
			})
			.catch(err => {
				console.log(err);
			});
		const sub = this.client.cent.subscribe('$lobby');
		sub.on('subscribe', console.log);
		sub.on('unsubscribe', console.log);
		sub.on('join', console.log);
		sub.on('error', console.log);

		// this.lobby$ = this.client.subscribe('lobby');
		// this.lobby$.subscribe(e => {
		// 	console.log(e);
		// });
	}
	private listenEvents() {
		if (this.client.onConnect$) {
			this.client.onConnect$.subscribe(context => {
				console.log('WS Connected');
			});
		}
		if (this.client.onDisconnect$) {
			this.client.onDisconnect$.subscribe(context => {
				// this.lobby$.unsubscribe();
				console.log(context);
			});
		}
	}

	private getToken(): Promise<string> {
		return fetchJWT('user1');
	}
}

new PhaserApp();
