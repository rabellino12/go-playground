import Phaser from 'phaser';

import { StartScene } from './scenes/Start';
import { fetchJWT } from './services/auth';
import { WSClient } from './services/centrifuge';

const config = {
	height: 600,
	type: Phaser.AUTO,
	width: 800
};
class PhaserApp extends Phaser.Game {
	private client: WSClient;
	private token?: string;

	constructor() {
		super(config);
		this.client = new WSClient();
		this.getToken().then(token => {
			this.token = token;
			this.client.connect(this.token);
			this.scene.add('Start', StartScene);
			this.scene.start('Start', {
				wsClient: this.client
			});
			this.listenEvents();
		})
		.catch(err => {
			console.log(err)
		});
	}
	private listenEvents() {
		if (this.client.onConnect$) {
			this.client.onConnect$.subscribe(context => {
				console.log(context);
			});
		}
		if (this.client.onDisconnect$) {
			this.client.onDisconnect$.subscribe(context => {
				console.log(context);
			});
		}
	}

	private getToken(): Promise<string> {
		return fetchJWT('user1');
	}
}

new PhaserApp();
