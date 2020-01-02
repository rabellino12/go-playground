import Phaser from 'phaser';

import { StartScene } from './scenes/Start';
import { fetchJWT } from './services/auth';
import { WSClient } from './services/centrifuge';
import { generateUsername } from './services/username';

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
	private userId: string;
	private token?: string;
	// private lobby$: Observable<any>;

	constructor() {
		super(config);
		this.userId = generateUsername();
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
	}
	private listenEvents() {
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
		const personalSub = this.client.cent.subscribe(`lobby#${this.userId}`);
		personalSub.on('subscribe', (e) => {
			console.log('personalSub:Subscribe', e);
		});
		personalSub.on('unsubscribe', (e) => {
			console.log('personalSub:unsubscribe', e);
		});
		personalSub.on('join', (e) => {
			console.log('personalSub:Join', e);
		});
		personalSub.on('error', (e) => {
			console.log('personalSub:Error', e);
		});
		personalSub.on('publish', (e) => {
			console.log('personalSub:publish', e);
		});
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
		return fetchJWT(this.userId);
	}
}

new PhaserApp();
