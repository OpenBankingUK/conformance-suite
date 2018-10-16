import {
  SET_WEBSOCKET_CONNECTION_STATE,
  SET_WEBSOCKET_LAST_UPDATE,
  SET_WEBSOCKET_CONNECTION,
} from './mutation-types';

export default {
  [SET_WEBSOCKET_CONNECTION_STATE](state, connectionState) {
    state.connectionState = connectionState;
  },
  [SET_WEBSOCKET_LAST_UPDATE](state, lastUpdate) {
    state.lastUpdate = lastUpdate;
  },
  [SET_WEBSOCKET_CONNECTION](state, connection) {
    state.connection = connection;
  },
};
