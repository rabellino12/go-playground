import Centrifuge from 'centrifuge';
// import Phaser from 'phaser';

// import { WSClient } from '../../services/centrifuge';
// import { IMove, MovementIO } from '../../services/movementIO';

// const assetsPrefix = '/game/assets';

// interface IPlayer {
// 	id: string;
// 	position: {
// 		x: number;
// 		y: number;
// 	};
// }
// interface IEnemy extends IPlayer {
// 	sprite: Phaser.Physics.Arcade.Sprite;
// }

// interface IAction {
// 	velocityY?: number;
// 	velocityX?: number;
// }

// interface IMoves {
// 	up: IAction;
// 	left: IAction;
// 	right: IAction;
// 	stop: IAction;
// }

// interface IEnemyMove {
// 	enemy: IEnemy;
// 	action: IAction;
// 	dir: 'left' | 'right' | 'up' | 'stop';
// }

// export class StartScene extends Phaser.Scene {
// 	public platforms?: Phaser.Physics.Arcade.StaticGroup;
// 	public player?: Phaser.Physics.Arcade.Sprite;
// 	public userId!: string;
// 	public enemies: IEnemy[] = [];
// 	public cursors?: Phaser.Types.Input.Keyboard.CursorKeys;
// 	public wsClient?: WSClient;
// 	private personalSub?: Centrifuge.Subscription;
// 	private movementService?: MovementIO;
// 	private moves!: IMoves;
// 	private enemyMoves: IEnemyMove[] = [];
// 	constructor() {
// 		super({
// 			physics: {
// 				arcade: {
// 					debug: true,
// 					gravity: { y: 1000 }
// 				},
// 				default: 'arcade'
// 			}
// 		});
// 	}

// 	public init(data: any) {
// 		this.personalSub = data.personalSub;
// 		if (this.personalSub) {
// 			this.personalSub.on('publish', this.handlePersonalPublish);
// 		}
// 		this.wsClient = data.wsClient;
// 		this.userId = data.userId;
// 		this.moves = {
// 			left: {
// 				velocityX: -260,
// 				velocityY: undefined
// 			},
// 			right: {
// 				velocityX: 260,
// 				velocityY: undefined
// 			},
// 			stop: {
// 				velocityX: 0,
// 				velocityY: undefined
// 			},
// 			up: {
// 				velocityX: undefined,
// 				velocityY: -630
// 			}
// 		};
// 	}

// 	public preload() {
// 		this.load.image('sky', `${assetsPrefix}/sky.png`);
// 		this.load.image('ground', `${assetsPrefix}/platform.png`);
// 		this.load.image('star', `${assetsPrefix}/star.png`);
// 		this.load.image('bomb', `${assetsPrefix}/bomb.png`);
// 		this.load.spritesheet('dude', `${assetsPrefix}/dude.png`, { frameWidth: 32, frameHeight: 48 });
// 	}

// 	public create() {
// 		this.add.image(400, 300, 'sky');
// 		this.add.image(400, 300, 'star');
// 		this.platforms = this.physics.add.staticGroup();
// 		this.platforms
// 			.create(400, 568, 'ground')
// 			.setScale(2)
// 			.refreshBody();

// 		this.platforms.create(600, 400, 'ground');
// 		this.platforms.create(50, 250, 'ground');
// 		this.platforms.create(750, 220, 'ground');
// 		this.player = this.createPlayer(100, 450);
// 		this.setAnimations();
// 	}
// 	public update() {
// 		if (!this.movementService) {
// 			return;
// 		}
// 		this.cursors = this.input.keyboard.createCursorKeys();
// 		if (
// 			!this.movementService ||
// 			!this.cursors ||
// 			!this.player ||
// 			!this.cursors.left ||
// 			!this.cursors.right ||
// 			!this.cursors.up
// 		) {
// 			return;
// 		}
// 		if (this.cursors && this.cursors.left.isDown) {
// 			this.player.setVelocityX(-260);
// 			this.player.anims.play('left', true);
// 			this.movementService.move('left');
// 		} else if (this.cursors.right.isDown) {
// 			this.player.setVelocityX(260);
// 			this.player.anims.play('right', true);
// 			this.movementService.move('right');
// 		} else {
// 			this.player.setVelocityX(0);
// 			this.player.anims.play('stop');
// 			this.movementService.stop();
// 		}

// 		if (this.cursors.up.isDown && this.player.body.touching.down) {
// 			this.player.setVelocityY(-630);
// 			this.movementService.move('up');
// 		}
// 		this.checkMoves();
// 	}

// 	public handlePersonalPublish = ({ data }: any) => {
// 		if (data && data.event === 'join' && this.wsClient) {
// 			this.movementService = new MovementIO({
// 				c: this.wsClient,
// 				matchId: data.game,
// 				userId: this.userId
// 			});
// 			this.handleEnemiesEvent(
// 				data.players.map((p: any) => {
// 					const [x, y] = p.position;
// 					return {
// 						...p,
// 						position: {
// 							x: Number(x),
// 							y: Number(y)
// 						}
// 					};
// 				})
// 			);
// 			this.movementService.movements$.subscribe(pub => {
// 				if (!this.movementService) {
// 					return;
// 				}
// 				const enemy = this.enemies.find(
// 					en => !!(pub.info && en.id === pub.info.user)
// 				);
// 				if (!enemy) {
// 					return;
// 				}
// 				const move: IMove = pub.data;
// 				this.enemyMoves.push({
// 					action: this.moves[move.action],
// 					dir: pub.data.action,
// 					enemy
// 				});
// 			});
// 		}
// 	};

// 	private setAnimations = () => {
// 		this.anims.create({
// 			frameRate: 10,
// 			frames: this.anims.generateFrameNumbers('dude', { start: 0, end: 3 }),
// 			key: 'left',
// 			repeat: -1
// 		});

// 		this.anims.create({
// 			frameRate: 20,
// 			frames: [{ key: 'dude', frame: 4 }],
// 			key: 'stop'
// 		});

// 		this.anims.create({
// 			frameRate: 10,
// 			frames: this.anims.generateFrameNumbers('dude', { start: 5, end: 8 }),
// 			key: 'right',
// 			repeat: -1
// 		});
// 	};

// 	private createPlayer = (
// 		x: number,
// 		y: number
// 	): Phaser.Physics.Arcade.Sprite => {
// 		if (!this.platforms) {
// 			throw new Error('No platforms created');
// 		}
// 		const player = this.physics.add.sprite(x, y, 'dude');

// 		player.setBounce(0.2);
// 		player.setCollideWorldBounds(true);
// 		this.physics.add.collider(player, this.platforms);
// 		return player;
// 	};

// 	private handleEnemiesEvent = (players: IPlayer[]) => {
// 		if (!this.player || !players) {
// 			return;
// 		}
// 		if (this.userId) {
// 			const meIndex = players.findIndex((p) => p.id === this.userId);
// 			if (meIndex > -1) {
// 				this.player.setPosition(players[meIndex].position.x, players[meIndex].position.y);
// 				players.splice(meIndex, 1);
// 			}
// 			players.forEach(player => {
// 				const enemy = this.enemies.find(e => e.id === player.id);
// 				if (!enemy) {
// 					this.enemies.push({
// 						id: player.id,
// 						position: player.position,
// 						sprite: this.createPlayer(
// 							player.position.x,
// 							player.position.y
// 						)
// 					});
// 				}
// 			});
// 		}
// 	};
// 	private handleMove = (
// 		player: Phaser.Physics.Arcade.Sprite,
// 		dir: IMove['action']
// 	) => {
// 		if (!this.movementService || !this.moves[dir]) {
// 			return;
// 		}
// 		const { velocityX, velocityY } = this.moves[dir];
// 		if (velocityX !== undefined) {
// 			player.setVelocityX(velocityX);
// 		}
// 		if (velocityY !== undefined) {
// 			player.setVelocityY(velocityY);
// 		}
// 		if (dir === 'stop') {
// 			player.anims.play(dir);
// 		} else if (dir !== 'up') {
// 			player.anims.play(dir, true);
// 		}
// 	};

// 	private checkMoves = () => {
// 		const moves = this.enemyMoves.splice(0, this.enemyMoves.length);
// 		for (let i = 0; i < moves.length; i++) {
// 			const move = moves[i];
// 			this.handleMove(move.enemy.sprite, move.dir);
// 		}
// 	}
// }
