import getters from './getters';

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

  describe('getters', () => {
    it('connectionState returns connection state of the WebSocket', () => {
      expect(getters.connectionState(state)).toEqual('DISCONNECTED');

      state.connectionState = 'CONNECTED';
      expect(getters.connectionState(state)).toEqual('CONNECTED');
    });

    it('channelState returns connection state of the Channel', () => {
      expect(getters.channelState(state)).toEqual('DISCONNECTED');

      state.channelState = 'CONNECTED';
      expect(getters.channelState(state)).toEqual('CONNECTED');
    });

    it('lastUpdate returns the last update was received on the Channel', () => {
      expect(getters.lastUpdate(state)).toEqual(null);

      state.lastUpdate = {};
      expect(getters.lastUpdate(state)).toEqual({});
    });

    it('channelDetails returns details of the channel', () => {
      expect(getters.channelDetails(state)).toEqual(null);

      const channelDetails = '<channel_details>';
      state.channelDetails = channelDetails;
      expect(getters.channelDetails(state)).toEqual(channelDetails);
    });

    it('tests returns a list of endpoints or null', () => {
      expect(getters.tests(state)).toEqual(null);

      const payload = {
        '/open-banking/v1.1/payments': {
          total_calls: 2,
          path: '/open-banking/v1.1/payments',
          failures: [],
          failed_calls: 0,
        },
        '/open-banking/v1.1/payment-submissions': {
          total_calls: 2,
          path: '/open-banking/v1.1/payment-submissions',
          failures: [],
          failed_calls: 0,
        },
      };
      state.lastUpdate = payload;
      expect(getters.tests(state)).toEqual(payload);
    });
  });
});
