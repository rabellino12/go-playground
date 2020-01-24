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
  public platforms?: Phaser.Physics.Arcade.StaticGroup;
  public player?: Phaser.Physics.Arcade.Sprite;
  public userId!: string;
  public enemies: IEnemy[] = [];
  public cursors?: Phaser.Types.Input.Keyboard.CursorKeys;
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
    const plat = this.add.sprite(0, 0, 'ground');
    // this.platforms = this.physics.add.staticGroup();
    // this.platforms.create(undefined, undefined, "ground");

    // this.platforms.create(600, 400, "ground");
    // this.platforms.create(50, 250, "ground");
    // this.platforms.create(750, 220, "ground");
    // this.player = this.createPlayer(100, 450);
    // this.setAnimations();
    this.worldScale = 30;

    // world gravity, as a Vec2 object. It's just a x, y vector
    let gravity = planck.Vec2(0, 3);

    // this is how we create a Box2D world
    this.world = planck.World(gravity);
    const ground = this.world.createBody({
      type: 'kinematic',
      position: planck.Vec2(2, 5)
    });
    ground.createFixture({
      shape: planck.Edge(planck.Vec2(-40.0, 0.0), planck.Vec2(40.0, 0.0))
    });
    ground.setPosition(planck.Vec2(25, 18));

    // time to set mass information
    ground.setMassData({
      mass: 1,
      center: planck.Vec2(),

      // I have to say I do not know the meaning of this "I", but if you set it to zero, bodies won't rotate
      I: 1
    });

    // const userData = this.add.graphics();
    // userData.fillStyle(color.color, 1);
    // userData.fillRect(- 300 / 2, - 300 / 2, 300, 300);

    // a body can have anything in its user data, normally it's used to store its sprite
    ground.setUserData(plat);
  }

  public update() {
    this.world.step(1 / 30);
    this.world.clearForces();

    for (let b = this.world.getBodyList(); b; b = b.getNext()){
 
        // get body position
        let bodyPosition = b.getPosition();

        // get body angle, in radians
        let bodyAngle = b.getAngle();

        // get body user data, the graphics object
        let userData: any = b.getUserData();

        // adjust graphic object position and rotation
        userData.x = bodyPosition.x * this.worldScale;
        userData.y = bodyPosition.y * this.worldScale;
        userData.rotation = bodyAngle;
    }
  }

  public handlePersonalPublish = ({ data }: any) => {
    if (data && data.event === "join" && this.wsClient) {
      this.movementService = new MovementIO({
        c: this.wsClient,
        matchId: data.game,
        userId: this.userId
      });
      this.handleEnemiesEvent(
        data.players.map((p: any) => {
          const [x, y] = p.position;
          return {
            ...p,
            position: {
              x: Number(x),
              y: Number(y)
            }
          };
        })
      );
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

  private createBox = (
    posX: number,
    posY: number,
    width: number,
    height: number,
    isDynamic: boolean
  ) => {
    // this is how we create a generic Box2D body
    let box = this.world.createBody();
    if (isDynamic) {
      // Box2D bodies born as static bodies, but we can make them dynamic
      box.setDynamic();
    }

    // a body can have one or more fixtures. This is how we create a box fixture inside a body
    box.createFixture(
      planck.Box(width / 2 / this.worldScale, height / 2 / this.worldScale)
    );

    // now we place the body in the world
    box.setPosition(
      planck.Vec2(posX / this.worldScale, posY / this.worldScale)
    );

    // time to set mass information
    box.setMassData({
      mass: 1,
      center: planck.Vec2(),

      // I have to say I do not know the meaning of this "I", but if you set it to zero, bodies won't rotate
      I: 1
    });

    // now we create a graphics object representing the body
    var color = new Phaser.Display.Color();
    color.random();
    color.brighten(50).saturate(100);
    let userData = this.add.graphics();
    userData.fillStyle(color.color, 1);
    userData.fillRect(-width / 2, -height / 2, width, height);

    // a body can have anything in its user data, normally it's used to store its sprite
    box.setUserData(userData);
  };

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
  ): Phaser.Physics.Arcade.Sprite => {
    if (!this.platforms) {
      throw new Error("No platforms created");
    }
    const player = this.physics.add.sprite(x, y, "dude");

    player.setBounce(0.2);
    player.setCollideWorldBounds(true);
    this.physics.add.collider(player, this.platforms);
    return player;
  };

  private handleEnemiesEvent = (players: IPlayer[]) => {
    if (!this.player || !players) {
      return;
    }
    if (this.userId) {
      const meIndex = players.findIndex(p => p.id === this.userId);
      if (meIndex > -1) {
        this.player.setPosition(
          players[meIndex].position.x,
          players[meIndex].position.y
        );
        players.splice(meIndex, 1);
      }
      players.forEach(player => {
        const enemy = this.enemies.find(e => e.id === player.id);
        if (!enemy) {
          this.enemies.push({
            id: player.id,
            position: player.position,
            sprite: this.createPlayer(player.position.x, player.position.y)
          });
        }
      });
    }
  };
  private handleMove = (
    player: Phaser.Physics.Arcade.Sprite,
    dir: IMove["action"]
  ) => {
    if (!this.movementService || !this.moves[dir]) {
      return;
    }
    const { velocityX, velocityY } = this.moves[dir];
    if (velocityX !== undefined) {
      player.setVelocityX(velocityX);
    }
    if (velocityY !== undefined) {
      player.setVelocityY(velocityY);
    }
    if (dir === "stop") {
      player.anims.play(dir);
    } else if (dir !== "up") {
      player.anims.play(dir, true);
    }
  };

  private checkMoves = () => {
    const moves = this.enemyMoves.splice(0, this.enemyMoves.length);
    for (let i = 0; i < moves.length; i++) {
      const move = moves[i];
      this.handleMove(move.enemy.sprite, move.dir);
    }
  };
}
