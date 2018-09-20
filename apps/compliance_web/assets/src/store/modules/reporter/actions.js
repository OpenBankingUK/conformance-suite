import { Socket } from '../../../plugins/socket';
import {
  SET_WEBSOCKET_CONNECTION_STATE,
  SET_WEBSOCKET_CHANNEL_CONNECTION_STATE,
  SET_WEBSOCKET_LAST_UPDATE,
  SET_WEBSOCKET_CHANNEL_FOUND,
} from './mutation-types';

const CHANNEL_PREFIX = 'report';

export default {
  connect({ commit, rootGetters }) {
    try {
      Socket.connect(rootGetters['user/getAccessToken'], false);
      commit(SET_WEBSOCKET_CONNECTION_STATE, 'CONNECTED');
    } catch (err) {
      commit(SET_WEBSOCKET_CONNECTION_STATE, 'ERROR');
    }
  },
  async disconnect({ commit }, validationRunId) {
    if (!validationRunId) return;
    if (Socket.connClosed()) return;

    try {
      await Socket.leaveChannel(validationRunId, CHANNEL_PREFIX);
    } catch (err) {
      commit(SET_WEBSOCKET_CHANNEL_CONNECTION_STATE, JSON.stringify(err));
      return;
    }

    Socket.disconnect();
    commit(SET_WEBSOCKET_CHANNEL_CONNECTION_STATE, 'DISCONNECTED');
    commit(SET_WEBSOCKET_CHANNEL_FOUND, null);
    commit(SET_WEBSOCKET_CONNECTION_STATE, 'DISCONNECTED');
  },
  async subscribeToChannel({ commit }, validationRunId) {
    if (!validationRunId) return;

    try {
      const result = await Socket.findChannel(validationRunId, CHANNEL_PREFIX);
      if (!result || !result.channel) return;

      const { channel } = result;

      channel.on('started', ({ payload }) => commit(SET_WEBSOCKET_LAST_UPDATE, payload));
      channel.on('updated', ({ payload }) => commit(SET_WEBSOCKET_LAST_UPDATE, payload));
      channel.on('completed', ({ payload }) => commit(SET_WEBSOCKET_LAST_UPDATE, payload));
      channel.on('stopped', ({ payload }) => commit(SET_WEBSOCKET_LAST_UPDATE, payload));

      commit(SET_WEBSOCKET_CHANNEL_CONNECTION_STATE, 'CONNECTED');
      commit(SET_WEBSOCKET_CHANNEL_FOUND, `state=${channel.state}_topic=${channel.topic}`);
    } catch (err) {
      commit(SET_WEBSOCKET_CHANNEL_CONNECTION_STATE, JSON.stringify(err));
      commit(SET_WEBSOCKET_CHANNEL_FOUND, JSON.stringify(err));
    }
  },
};
