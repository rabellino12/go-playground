import Phaser from 'phaser';
import { Observable } from 'rxjs';

import bomb from '../assets/bomb.png';
import dude from '../assets/dude.png';
import platform from '../assets/platform.png';
import sky from '../assets/sky.png';
import star from '../assets/star.png';
import { IConnectionEvent, WSClient } from '../services/centrifuge';

export class StartScene extends Phaser.Scene {
	public platforms?: Phaser.Physics.Arcade.StaticGroup;
	public player?: Phaser.Physics.Arcade.Sprite;
	public cursors?: Phaser.Types.Input.Keyboard.CursorKeys;

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
		this.player = this.physics.add.sprite(100, 450, 'dude');

		this.player.setBounce(0.2);
		this.player.setCollideWorldBounds(true);

		this.anims.create({
			frameRate: 10,
			frames: this.anims.generateFrameNumbers('dude', { start: 0, end: 3 }),
			key: 'left',
			repeat: -1
		});

		this.anims.create({
			key: 'turn',
			frames: [{ key: 'dude', frame: 4 }],
			frameRate: 20
		});

		this.anims.create({
			frameRate: 10,
			frames: this.anims.generateFrameNumbers('dude', { start: 5, end: 8 }),
			key: 'right',
			repeat: -1
		});
		this.physics.add.collider(this.player, this.platforms);
	}
	public update() {
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
		} else if (this.cursors.right.isDown) {
			this.player.setVelocityX(260);

			this.player.anims.play('right', true);
		} else {
			this.player.setVelocityX(0);

			this.player.anims.play('turn');
		}

		if (this.cursors.up.isDown && this.player.body.touching.down) {
			this.player.setVelocityY(-630);
		}
	}
}
