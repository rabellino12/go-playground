import {createStore} from 'redux';

import {reducers} from './reducers';

export interface IGameReduxAction {
	type: string;
	payload: string;
}

export function setupStore() {
	const store = createStore(reducers);
	return store;
}
