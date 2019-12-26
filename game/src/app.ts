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

const userId = 'user1';

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
		const sub = this.client.cent.subscribe('$lobby:index');
		sub.on('subscribe', (e) => {
			console.log('Subscribe', e);
		});
		sub.on('unsubscribe', (e) => {
			console.log('unsubscribe', e);
		});
		sub.on('join', (e) => {
			console.log('Join', e);
		});
		sub.on('error', (e) => {
			console.log('Error', e);
		});
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
		return fetchJWT(userId);
	}
}

new PhaserApp();
