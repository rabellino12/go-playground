import planck from "planck-js";
import Phaser from "phaser";
import Centrifuge from "centrifuge";

import { WSClient } from "../../services/centrifuge";
import { IMove, MovementIO } from "../../services/movementIO";

const assetsPrefix = "/game/assets";

interface IPlayer {
  id: string;
  position: {
    x: number;
    y: number;
  };
}
interface IEnemy extends IPlayer {
  sprite: Phaser.Physics.Arcade.Sprite;
}

interface IAction {
  velocityY?: number;
  velocityX?: number;
}

interface IMoves {
  up: IAction;
  left: IAction;
  right: IAction;
  stop: IAction;
}

interface IEnemyMove {
  enemy: IEnemy;
  action: IAction;
  dir: "left" | "right" | "up" | "stop";
}

export class MatchScene extends Phaser.Scene {
  public ground!: planck.Body;
  public player!: planck.Body;
  public userId!: string;
  public enemies: IEnemy[] = [];
  public cursors!: Phaser.Types.Input.Keyboard.CursorKeys;
  public wsClient?: WSClient;
  public worldScale!: number;
  public world!: planck.World;
  public tick: number = 0;
  private personalSub?: Centrifuge.Subscription;
  private movementService?: MovementIO;
  private moves!: IMoves;
  private enemyMoves: IEnemyMove[] = [];

  constructor() {
    super({
      physics: {
        arcade: {
          debug: false,
          gravity: { y: 1000 }
        },
        default: "arcade"
      }
    });
  }

  public init(data: any) {
    this.personalSub = data.personalSub;
    if (this.personalSub) {
      this.personalSub.on("publish", this.handlePersonalPublish);
    }
    this.wsClient = data.wsClient;
    this.userId = data.userId;
    this.moves = {
      left: {
        velocityX: -260,
        velocityY: undefined
      },
      right: {
        velocityX: 260,
        velocityY: undefined
      },
      stop: {
        velocityX: 0,
        velocityY: undefined
      },
      up: {
        velocityX: undefined,
        velocityY: -630
      }
    };
  }

  public preload() {
    this.load.image("sky", `${assetsPrefix}/sky.png`);
    this.load.image("ground", `${assetsPrefix}/platform.png`);
    this.load.image("star", `${assetsPrefix}/star.png`);
    this.load.image("bomb", `${assetsPrefix}/bomb.png`);
    this.load.spritesheet("dude", `${assetsPrefix}/dude.png`, {
      frameWidth: 32,
      frameHeight: 48
    });
  }

  public create() {
    // this.add.image(400, 300, "sky");
    // this.add.image(400, 300, "star");
    // this.platforms = this.physics.add.staticGroup();
    // this.platforms.create(undefined, undefined, "ground");

    // this.platforms.create(600, 400, "ground");
    // this.platforms.create(50, 250, "ground");
    // this.platforms.create(750, 220, "ground");
    this.worldScale = 30;
    let gravity = planck.Vec2(0, 3);
    this.world = planck.World(gravity);
    this.ground = this.createEnvironment();
    this.player = this.createPlayer(100, 450);
    this.setAnimations();
  }

  public update() {
    this.world.step(1 / 30, 8, 3);
    this.world.clearForces();

    this.cursors = this.input.keyboard.createCursorKeys();
    if (!this.cursors.left || !this.cursors.right || !this.cursors.up) {
      return ;
    }
    const player: Phaser.GameObjects.Sprite = this.player.m_userData as Phaser.GameObjects.Sprite;
    if (this.cursors.up.isDown && this.playerTouchingFloor()) {
			const f = this.player.getWorldVector(planck.Vec2(0.0, -0.80));
      const p = this.player.getWorldPoint(planck.Vec2(0.0, 0.08));
      this.player.applyLinearImpulse(f, p, true);
		}
		if (this.cursors && this.cursors.left.isDown) {
      const f = this.player.getWorldVector(planck.Vec2(-0.10, 0));
      const p = this.player.getWorldPoint(planck.Vec2(0, 0));
      this.player.applyLinearImpulse(f, p, true);
			player.anims.play('left', true);
			// this.movementService.move('left');
		} else if (this.cursors.right.isDown) {
      const f = this.player.getWorldVector(planck.Vec2(0.10, 0));
      const p = this.player.getWorldPoint(planck.Vec2(0, 0));
      this.player.applyLinearImpulse(f, p, true);
			player.anims.play('right', true);
			// this.movementService.move('right');
		} else {
			player.anims.play('stop', true);
			// this.movementService.stop();
		}

    for (let b = this.world.getBodyList(); b; b = b.getNext()){
 
        // get body position
        let bodyPosition = b.getPosition();

        // get body angle, in radians
        let bodyAngle = b.getAngle();

        // get body user data, the graphics object
        let userData: any = b.getUserData();

        if (userData) {
          // adjust graphic object position and rotation
          userData.x = bodyPosition.x * this.worldScale;
          userData.y = bodyPosition.y * this.worldScale;
          userData.rotation = bodyAngle;
        }
    }
  }

  public handlePersonalPublish = ({ data }: any) => {
    if (data && data.event === "join" && this.wsClient) {
      this.movementService = new MovementIO({
        c: this.wsClient,
        matchId: data.game,
        userId: this.userId
      });
      this.movementService.movements$.subscribe(pub => {
        if (!this.movementService) {
          return;
        }
        const enemy = this.enemies.find(
          en => !!(pub.info && en.id === pub.info.user)
        );
        if (!enemy) {
          return;
        }
        const move: IMove = pub.data;
        this.enemyMoves.push({
          action: this.moves[move.action],
          dir: pub.data.action,
          enemy
        });
      });
    }
  };

  private playerTouchingFloor = (): boolean => {
    const contact = this.world.getContactList();
    const list = [contact?.getFixtureA(), contact?.getFixtureB()].filter((c) => !!c);
    if (list.length < 2) {
      return false;
    }
    const player = list.find(fixture => fixture === this.player.getFixtureList())
    const ground = list.find(fixture => fixture === this.ground.getFixtureList())
    return Boolean(player && ground);
  }

  private setAnimations = () => {
    this.anims.create({
      frameRate: 10,
      frames: this.anims.generateFrameNumbers("dude", { start: 0, end: 3 }),
      key: "left",
      repeat: -1
    });

    this.anims.create({
      frameRate: 20,
      frames: [{ key: "dude", frame: 4 }],
      key: "stop"
    });

    this.anims.create({
      frameRate: 10,
      frames: this.anims.generateFrameNumbers("dude", { start: 5, end: 8 }),
      key: "right",
      repeat: -1
    });
  };

  private createPlayer = (
    x: number,
    y: number
  ): planck.Body => {
    const player = this.add.sprite(x, y, "dude");
    const body = this.world.createDynamicBody({
      type: 'dynamic',
      // position: planck.Vec2(2, 5)
    });
    const fix = body.createFixture({
      density: 1,
      friction: 0.3,
      shape: planck.Box((48/this.worldScale)/2, (48/this.worldScale)/2)
    });
    body.setPosition(planck.Vec2(100 / this.worldScale, 350 / this.worldScale));
    body.setMassData({
      mass: 1,
      center: planck.Vec2(),
      I: 1
    });
    body.setUserData(player);
    return body;
  };
  
  private createEnvironment = (): planck.Body => {
    const plat = this.add.sprite(400 / this.worldScale, 568 / this.worldScale, 'ground').setScale(2);
    const ground = this.world.createBody({
      type: 'static',
      position: planck.Vec2(0, -10)
    });
    ground.createFixture({
      density: 1,
      friction: 1,
      shape: planck.Box((800 / this.worldScale) / 2, (64 / this.worldScale)/2)
    });
    ground.setPosition(planck.Vec2(400 / this.worldScale, 568 / this.worldScale));
    ground.setMassData({
      mass: 1,
      center: planck.Vec2(),
      I: 1
    });
    ground.setUserData(plat);
    const borderLeft = this.world.createBody({
      type: 'static',
      position: planck.Vec2(0, 0)
    });
    borderLeft.createFixture({
      density: 1,
      friction: 0,
      shape: planck.Edge(planck.Vec2(0, 0), planck.Vec2(0, 568 / this.worldScale))
    });
    borderLeft.setPosition(planck.Vec2(0, 0));
    borderLeft.setMassData({
      mass: 1,
      center: planck.Vec2(),
      I: 1
    });
    const borderRight = this.world.createBody({
      type: 'static',
      position: planck.Vec2(0, 0)
    });
    borderRight.createFixture({
      density: 1,
      friction: 0,
      shape: planck.Edge(planck.Vec2(800 / this.worldScale, 0), planck.Vec2(800 / this.worldScale, 568 / this.worldScale))
    });
    borderRight.setMassData({
      mass: 1,
      center: planck.Vec2(),
      I: 1
    });
    return ground;
  }

}
