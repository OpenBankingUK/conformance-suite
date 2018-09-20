import mutations from './mutations';
import {
  SET_WEBSOCKET_CHANNEL_FOUND,
  SET_WEBSOCKET_CHANNEL_CONNECTION_STATE,
  SET_WEBSOCKET_CONNECTION_STATE,
  SET_WEBSOCKET_LAST_UPDATE,
} from './mutation-types';

describe('Reporter', () => {
  let state;

  beforeEach(() => {
    state = {
      connectionState: 'DISCONNECTED',
      channelState: 'DISCONNECTED',
      lastUpdate: null,
      channelDetails: null,
    };
  });

  describe('mutations', () => {
    it('SET_WEBSOCKET_CONNECTION_STATE commits WebSocket connection state to the state', () => {
      const connectionState = 'CONNECTED';

      expect(state.connectionState).not.toEqual(connectionState);
      mutations[SET_WEBSOCKET_CONNECTION_STATE](state, connectionState);
      expect(state.connectionState).toEqual(connectionState);
    });

    it('SET_WEBSOCKET_CHANNEL_CONNECTION_STATE commits the channel connection state to the state', () => {
      const channelState = 'CONNECTED';

      expect(state.channelState).not.toEqual(channelState);
      mutations[SET_WEBSOCKET_CHANNEL_CONNECTION_STATE](state, channelState);
      expect(state.channelState).toEqual(channelState);
    });

    it('SET_WEBSOCKET_LAST_UPDATE commits last update to the state', () => {
      const lastUpdate = Date.now().toString();

      expect(state.lastUpdate).not.toEqual(lastUpdate);
      mutations[SET_WEBSOCKET_LAST_UPDATE](state, lastUpdate);
      expect(state.lastUpdate).toEqual(lastUpdate);
    });

    it('SET_WEBSOCKET_CHANNEL_FOUND commits details of the channel found to the state', () => {
      const channelDetails = '<channel_details>';

      expect(state.channelDetails).not.toEqual(channelDetails);
      mutations[SET_WEBSOCKET_CHANNEL_FOUND](state, channelDetails);
      expect(state.channelDetails).toEqual(channelDetails);
    });
  });
});
