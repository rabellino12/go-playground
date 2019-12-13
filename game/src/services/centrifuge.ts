import Centrifuge from 'centrifuge';
import {Observable} from 'rxjs';

interface ConnectionEvent {
    client: string
    transport: string
    latency: number
}
interface DisconnectionEvent {
    reason: string
    reconnect: boolean
}

export function connect(): Centrifuge {
    var centrifuge = new Centrifuge('ws://localhost:8081/connection/websocket');

    centrifuge.subscribe("news", function(message) {
        console.log(message);
    });

    centrifuge.connect();
    return centrifuge
}

export function onConnect$(client: Centrifuge): Observable<ConnectionEvent> {
    return new Observable<ConnectionEvent>(sub => {
        client.on('connect', (context: ConnectionEvent /* Couldn't find correct type for this context */) => {
            sub.next(context)
        })
    });
}

export function onDisconnect(client: Centrifuge): Observable<DisconnectionEvent> {
    return new Observable<DisconnectionEvent>(sub => {
        client.on('disconnect', function(context) {
            sub.next(context)
        });
    })
}