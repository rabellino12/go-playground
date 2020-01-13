import Centrifuge from 'centrifuge';
import Phaser from 'phaser';

import { MovementIO } from '../services/movementIO';
import { WSClient } from '../services/centrifuge';

import bomb from '../assets/bomb.png';
import dude from '../assets/dude.png';
import platform from '../assets/platform.png';
import sky from '../assets/sky.png';
import star from '../assets/star.png';

interface IEnemy {
	id: string;
	sprite: Phaser.Physics.Arcade.Sprite;
}

export class StartScene extends Phaser.Scene {
	public platforms?: Phaser.Physics.Arcade.StaticGroup;
	public player?: Phaser.Physics.Arcade.Sprite;
	public userId!: string;
	public enemies: IEnemy[] = [];
	public cursors?: Phaser.Types.Input.Keyboard.CursorKeys;
	public wsClient?: WSClient;
	private personalSub?: Centrifuge.Subscription;
	private movementService?: MovementIO;
	constructor() {
		super({
			physics: {
				arcade: {
					debug: false,
					gravity: { y: 1000 }
				},
				default: 'arcade'
			}
		});
	}

	public init(data: any) {
		this.personalSub = data.personalSub;
		if (this.personalSub) {
			this.personalSub.on('publish', this.handlePersonalPublish);
		}
		this.wsClient = data.wsClient;
		this.userId = data.userId;
	}

	public preload() {
		this.load.image('sky', sky);
		this.load.image('ground', platform);
		this.load.image('star', star);
		this.load.image('bomb', bomb);
		this.load.spritesheet('dude', dude, { frameWidth: 32, frameHeight: 48 });
	}

	public create() {
		this.add.image(400, 300, 'sky');
		this.add.image(400, 300, 'star');
		this.platforms = this.physics.add.staticGroup();
		this.platforms
			.create(400, 568, 'ground')
			.setScale(2)
			.refreshBody();

		this.platforms.create(600, 400, 'ground');
		this.platforms.create(50, 250, 'ground');
		this.platforms.create(750, 220, 'ground');
		this.player = this.createPlayer(100, 450);
		this.setAnimations();
	}
	public update() {
		if (!this.movementService || !this.wsClient) {
			return;
		}
		this.cursors = this.input.keyboard.createCursorKeys();
		if (
			!this.cursors ||
			!this.player ||
			!this.cursors.left ||
			!this.cursors.right ||
			!this.cursors.up
		) {
			return;
		}
		if (this.cursors && this.cursors.left.isDown) {
			this.player.setVelocityX(-260);
			this.player.anims.play('left', true);
			this.movementService.move('left');
		} else if (this.cursors.right.isDown) {
			this.player.setVelocityX(260);
			this.player.anims.play('right', true);
			this.movementService.move('right');
		} else {
			this.player.setVelocityX(0);
			this.player.anims.play('turn');
			this.movementService.stop();
		}

		if (this.cursors.up.isDown && this.player.body.touching.down) {
			this.player.setVelocityY(-630);
			this.movementService.move('jump');
		}
	}

	public handlePersonalPublish = ({ data }: any) => {
		if (data && data.event === 'join' && this.wsClient) {
			this.movementService = new MovementIO({
				c: this.wsClient,
				matchId: data.game,
				userId: this.userId
			});
			this.movementService.matchSubscription.on('publish', (e: any) => {
				console.log(e);
			});
			this.movementService.enemies$.subscribe(e => {
				e.forEach(this.handleEnemiesEvent);
			});
		}
	}

	private setAnimations = () => {
		this.anims.create({
			frameRate: 10,
			frames: this.anims.generateFrameNumbers('dude', { start: 0, end: 3 }),
			key: 'left',
			repeat: -1
		});

		this.anims.create({
			frameRate: 20,
			frames: [{ key: 'dude', frame: 4 }],
			key: 'turn'
		});

		this.anims.create({
			frameRate: 10,
			frames: this.anims.generateFrameNumbers('dude', { start: 5, end: 8 }),
			key: 'right',
			repeat: -1
		});
	}

	private createPlayer = (x: number, y: number): Phaser.Physics.Arcade.Sprite => {
		if (!this.platforms) {
			throw new Error('No platforms created');
		}
		const player = this.physics.add.sprite(x, y, 'dude');

		player.setBounce(0.2);
		player.setCollideWorldBounds(true);
		this.physics.add.collider(player, this.platforms);
		return player;
	}

	private handleEnemiesEvent = (data: string | undefined) => {
		if (!this.player || !data) {
			return;
		}
		if (this.userId && this.userId !== data) {
			const enemy = this.enemies.find(e => e.id === data);
			if (!enemy) {
				this.enemies.push({
					id: data,
					sprite: this.createPlayer(
						this.player.x + 10 * (1 + this.enemies.length),
						450
					)
				});
			}
		}
	}
}
