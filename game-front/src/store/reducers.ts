import { combineReducers } from 'redux';

import { SET_USER_NAME } from './actions/base';
import { IGameReduxAction } from './store';

interface IInitialState {
	userName: string;
}

const initialState: IInitialState = {
	userName: ''
};

function baseReducer(state: IInitialState, action: IGameReduxAction) {
	if (typeof state === 'undefined') {
		return initialState;
	}

	switch (action.type) {
		case SET_USER_NAME:
			return {
				...state,
				userName: action.payload
			};
		default:
			break;
	}

	// Por ahora, no maneja ninguna acción
	// y solo devuelve el estado que recibimos.
	return state;
}

export const reducers = combineReducers({
	base: baseReducer
});
