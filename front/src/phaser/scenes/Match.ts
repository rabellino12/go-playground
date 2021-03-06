import planck from "planck-js";
import Phaser from "phaser";
import Centrifuge from "centrifuge";

import { WSClient } from "../../services/centrifuge";
import {
  IMove,
  MovementIO,
  IPlayer,
  IMoveInput
} from "../../services/movementIO";

const assetsPrefix = "/game/assets";

export class MatchScene extends Phaser.Scene {
  public ground!: planck.Body;
  public player!: planck.Body;
  public platforms: planck.Body[] = [];
  public userId!: string;
  public cursors!: Phaser.Types.Input.Keyboard.CursorKeys;
  public wsClient?: WSClient;
  public worldScale!: number;
  public world!: planck.World;
  private moves!: IMove[];
  private personalSub?: Centrifuge.Subscription;
  private movementService?: MovementIO;
  private players!: IPlayer[];

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
    this.worldScale = 30;
    let gravity = planck.Vec2(0, 5);
    this.world = planck.World(gravity);
    this.createEnvironment();
    this.setAnimations();
  }

  public update() {
    this.world.step(1 / 30, 8, 3);
    this.world.clearForces();

    this.cursors = this.input.keyboard.createCursorKeys();
    for (let b = this.world.getBodyList(); b; b = b.getNext()) {
      let bodyPosition = b.getPosition();
      let bodyAngle = b.getAngle();
      let userData: any = b.getUserData();

      if (userData) {
        userData.x = bodyPosition.x * this.worldScale;
        userData.y = bodyPosition.y * this.worldScale;
        userData.rotation = bodyAngle;
      }
    }
    if (
      !this.cursors.left ||
      !this.cursors.right ||
      !this.cursors.up ||
      !this.movementService
    ) {
      return;
    }
    this.players.forEach(p => {
      if (p.id === this.userId) {
        if (!p.body) {
          return;
        }
        const player: Phaser.GameObjects.Sprite = p.body.m_userData as Phaser.GameObjects.Sprite;
        const vel = p.body.getLinearVelocity();
        const { x, y } = p.body.getPosition();
        var force = 0;
        const move: IMoveInput = {
          position: {
            x,
            y
          },
          action: "stop",
          jumping: false
        };
        if (this.cursors.up?.isDown && this.playerTouchingFloor(p.body)) {
          const f = p.body.getWorldVector(planck.Vec2(0.0, -1));
          const point = p.body.getWorldPoint(planck.Vec2(0.0, 0.1));
          p.body.applyLinearImpulse(f, point, true);
          move.jumping = true;
        }
        if (this.cursors && this.cursors.left?.isDown) {
          if (vel.x > -5) {
            force = -50;
          }
          player.anims.play("left", true);
          move.action = "left";
        } else if (this.cursors && this.cursors.right?.isDown) {
          if (vel.x < 5) {
            force = 50;
          }
          player.anims.play("right", true);
          move.action = "right";
        } else {
          if (vel.x) {
            force = vel.x * -10;
          }
          player.anims.play("stop", true);
          move.action = "stop";
        }
        this.movementService?.move(move);
        p.body.applyForce(
          planck.Vec2(force, 0),
          p.body.getWorldCenter(),
          true
        );
      }
    });
  }

  public handlePersonalPublish = ({ data }: any) => {
    console.log(data);
    if (data && data.event === "join" && this.wsClient) {
      this.movementService = new MovementIO({
        c: this.wsClient,
        matchId: data.game,
        userId: this.userId
      });
      this.players = data.players as IPlayer[];
      this.players.forEach((p, i) => {
        this.players[i].body = this.createPlayer(p.position.x, p.position.y);
      });
      this.movementService.matchSubscription.on("publish", e => {
        console.log(e);
      });
      this.movementService.snapshot$.subscribe(pub => {
        console.log(pub);
        if (!this.movementService) {
          return;
        }
        const enemy = this.players.find(
          en => !!(pub.info && en.id === pub.info.user)
        );
        if (!enemy) {
          return;
        }
        const move: IMove = pub.data;
        this.moves.push(move);
      });
    }
  };

  private playerTouchingFloor = (playerBody: planck.Body): boolean => {
    const contact = this.world.getContactList();
    const list = [contact?.getFixtureA(), contact?.getFixtureB()].filter(
      c => !!c
    );
    if (list.length < 2) {
      return false;
    }
    const player = list.find(
      fixture => fixture === playerBody.getFixtureList()
    );
    const platform = list.find(
      fixture =>
        fixture &&
        this.platforms.map(body => body.getFixtureList()).includes(fixture)
    );
    return Boolean(player && platform);
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

  private createPlayer = (x: number, y: number): planck.Body => {
    const player = this.add.sprite(x, y, "dude");
    const body = this.world.createDynamicBody({
      type: "dynamic",
      fixedRotation: true
    });
    const fix = body.createFixture({
      density: 1,
      friction: 0,
      shape: planck.Box(28 / this.worldScale / 2, 48 / this.worldScale / 2)
    });
    body.setPosition(planck.Vec2(x / this.worldScale, y / this.worldScale));
    body.setMassData({
      mass: 1,
      center: planck.Vec2(),
      I: 1
    });
    body.setUserData(player);
    return body;
  };

  private createGround = (
    xPx: number,
    yPx: number,
    widthPx: number,
    heightPx: number,
    scale: number = 1
  ): planck.Body => {
    const plat = this.add
      .sprite(xPx / this.worldScale, yPx / this.worldScale, "ground")
      .setScale(scale);
    const ground = this.world.createBody({
      type: "static",
      position: planck.Vec2(xPx / this.worldScale, yPx / this.worldScale)
    });
    ground.createFixture({
      density: 1,
      friction: 0,
      shape: planck.Box(
        widthPx / this.worldScale / 2,
        heightPx / this.worldScale / 2
      )
    });
    ground.setMassData({
      mass: 1,
      center: planck.Vec2(),
      I: 1
    });
    ground.setUserData(plat);
    return ground;
  };

  private createEnvironment = () => {
    const ground = this.createGround(400, 578, 800, 64, 2);
    const plat1 = this.createGround(600, 400, 400, 32);
    this.platforms = [ground, plat1];
    const borderLeft = this.world.createBody({
      type: "static",
      position: planck.Vec2(0, 0)
    });
    borderLeft.createFixture({
      density: 1,
      friction: 0,
      shape: planck.Edge(
        planck.Vec2(0, 0),
        planck.Vec2(0, 568 / this.worldScale)
      )
    });
    borderLeft.setPosition(planck.Vec2(0, 0));
    borderLeft.setMassData({
      mass: 1,
      center: planck.Vec2(),
      I: 1
    });
    const borderRight = this.world.createBody({
      type: "static",
      position: planck.Vec2(0, 0)
    });
    borderRight.createFixture({
      density: 1,
      friction: 0,
      shape: planck.Edge(
        planck.Vec2(800 / this.worldScale, 0),
        planck.Vec2(800 / this.worldScale, 568 / this.worldScale)
      )
    });
    borderRight.setMassData({
      mass: 1,
      center: planck.Vec2(),
      I: 1
    });
  };
}
