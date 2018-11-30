import {
  SET_WEBSOCKET_CONNECTION_STATE,
  SET_WEBSOCKET_LAST_UPDATE,
  SET_WEBSOCKET_CONNECTION,
} from './mutation-types';

const makeUrl = () => {
  const { protocol, host } = window.location;

  if (protocol === 'https:') {
    return `wss://${host}/api/ws`;
  }
  return `ws://${host}/api/ws`;
};

export default {
  connect({ commit }) {
    const url = makeUrl();
    const connection = new WebSocket(url);

    return new Promise((resolve, reject) => {
      commit(SET_WEBSOCKET_CONNECTION, connection);

      connection.onclose = (ev) => {
        commit(SET_WEBSOCKET_CONNECTION_STATE, 'DISCONNECTED');
        commit(SET_WEBSOCKET_LAST_UPDATE, ev);
        commit(SET_WEBSOCKET_CONNECTION, null);

        reject(ev);
      };
      connection.onerror = (ev) => {
        commit(SET_WEBSOCKET_CONNECTION_STATE, 'DISCONNECTED');
        commit(SET_WEBSOCKET_LAST_UPDATE, ev);
        commit(SET_WEBSOCKET_CONNECTION, null);

        reject(ev);
      };
      connection.onmessage = (ev) => {
        const payload = ev.data;
        commit(SET_WEBSOCKET_LAST_UPDATE, payload);
      };
      connection.onopen = (ev) => {
        commit(SET_WEBSOCKET_CONNECTION_STATE, 'CONNECTED');
        commit(SET_WEBSOCKET_LAST_UPDATE, ev);

        resolve({});
      };
    });
  },
  async disconnect({ state, commit }) {
    const { connection } = state;

    if (!connection) {
      const err = new Error(`WebSocket connection null, state=${JSON.stringify(state)}`);
      return Promise.reject(err);
    }

    return new Promise((resolve) => {
      connection.onclose = (ev) => {
        commit(SET_WEBSOCKET_CONNECTION_STATE, 'DISCONNECTED');
        commit(SET_WEBSOCKET_LAST_UPDATE, ev);
        commit(SET_WEBSOCKET_CONNECTION, null);

        resolve(ev);
      };

      connection.close();
    });
  },
};
