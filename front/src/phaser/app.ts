import Phaser from 'phaser';
import Centrifuge from 'centrifuge';

import { StartScene } from './scenes/Start';
import { MatchScene } from './scenes/Match';
import { fetchJWT } from '../services/auth';
import { WSClient } from '../services/centrifuge';
import { generateUsername } from '../services/username';

interface Lobby {
	status: string;
	id?: string;
}

export class PhaserApp extends Phaser.Game {
	private client: WSClient;
	private userId: string;
	private token?: string;
	private sub: Centrifuge.Subscription | undefined;
	private personalSub: Centrifuge.Subscription | undefined;

	constructor(parent: HTMLElement) {
		super({
			height: 600,
			type: Phaser.AUTO,
			width: 800,
			parent
		});
		this.userId = generateUsername();
		this.client = new WSClient();
		this.getToken()
			.then(token => {
				this.token = token;
				this.client.connect(this.token);
				this.listenEvents();
				this.scene.add('Start', StartScene);
				this.scene.add('Match', MatchScene);
				this.scene.start('Match', {
					personalSub: this.personalSub,
					userId: this.userId,
					wsClient: this.client
				});
			})
			.catch(err => {
				console.log(err);
			});
	}
	private listenEvents() {
		this.lobbySubscription();
		this.personalSubscription();
	}

	private lobbySubscription() {
		this.sub = this.client.cent.subscribe('$lobby:index');
		this.sub.on('subscribe', (e) => {
			console.log('Subscribe', e);
		});
		this.sub.on('unsubscribe', (e) => {
			console.log('unsubscribe', e);
		});
		this.sub.on('join', (e) => {
			console.log('Join', e);
		});
		this.sub.on('error', (e) => {
			console.log('Error', e);
		});
	}

	private personalSubscription() {
		this.personalSub = this.client.cent.subscribe(`lobby#${this.userId}`);
		this.personalSub.on('subscribe', (e) => {
			console.log('personalSub:Subscribe', e);
		});
		this.personalSub.on('unsubscribe', (e) => {
			console.log('personalSub:unsubscribe', e);
		});
		this.personalSub.on('join', (e) => {
			console.log('personalSub:Join', e);
		});
		this.personalSub.on('error', (e) => {
			console.log('personalSub:Error', e);
		});
		this.personalSub.on('publish', (e) => {
			console.log('personalSub:publish', e);
		});
	}

	private getToken(): Promise<string> {
		return fetchJWT(this.userId);
	}
}