import {
  SET_WEBSOCKET_CONNECTION_STATE,
  SET_WEBSOCKET_CHANNEL_CONNECTION_STATE,
  SET_WEBSOCKET_LAST_UPDATE,
  SET_WEBSOCKET_CHANNEL_FOUND,
} from './mutation-types';

export default {
  [SET_WEBSOCKET_CONNECTION_STATE](state, connectionState) {
    state.connectionState = connectionState;
  },
  [SET_WEBSOCKET_CHANNEL_CONNECTION_STATE](state, channelState) {
    state.channelState = channelState;
  },
  [SET_WEBSOCKET_LAST_UPDATE](state, lastUpdate) {
    state.lastUpdate = lastUpdate;
  },
  [SET_WEBSOCKET_CHANNEL_FOUND](state, channelDetails) {
    state.channelDetails = channelDetails;
  },
};
