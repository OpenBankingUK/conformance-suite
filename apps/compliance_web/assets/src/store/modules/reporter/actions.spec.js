import { EventEmitter } from 'events';
import actions from './actions';
import {
  SET_WEBSOCKET_CONNECTION_STATE,
  SET_WEBSOCKET_CHANNEL_CONNECTION_STATE,
  SET_WEBSOCKET_LAST_UPDATE,
  SET_WEBSOCKET_CHANNEL_FOUND,
} from './mutation-types';
import { Socket } from '../../../../src/plugins/socket';

jest.mock('../../../../src/plugins/socket');

describe('Reporter', () => {
  describe('actions', () => {
    let commit;
    let rootGetters;

    beforeEach(() => {
      commit = jest.fn();
      rootGetters = jest.fn();
    });

    afterEach(() => {
      jest.resetAllMocks();
    });

    describe('connect', () => {
      it('should call commit with CONNECTED if Socket successful', () => {
        Socket.connect.mockReturnValue();
        actions.connect({ commit, rootGetters });
        expect(commit).toHaveBeenCalledTimes(1);
        expect(commit).toHaveBeenCalledWith(SET_WEBSOCKET_CONNECTION_STATE, 'CONNECTED');
      });

      it('should call commit with ERROR if Socket error', () => {
        Socket.connect.mockImplementation(() => {
          throw new Error('Some error');
        });
        actions.connect({ commit, rootGetters });
        expect(commit).toHaveBeenCalledTimes(1);
        expect(commit).toHaveBeenCalledWith(SET_WEBSOCKET_CONNECTION_STATE, 'ERROR');
      });
    });

    describe('disconnect', () => {
      it('should return if no validationRunId', async () => {
        await actions.disconnect({ commit });
        expect(commit).not.toHaveBeenCalled();
      });

      it('should return if no Socket connection closed', async () => {
        Socket.connClosed.mockResolvedValue();
        await actions.disconnect({ commit }, 'validationRunId');
        expect(commit).not.toHaveBeenCalled();
      });

      it('should leaveChannel if no errors', async () => {
        Socket.leaveChannel.mockResolvedValue();
        Socket.disconnect.mockResolvedValue();
        await actions.disconnect({ commit }, 'validationRunId');
        expect(commit).toHaveBeenNthCalledWith(1, SET_WEBSOCKET_CHANNEL_CONNECTION_STATE, 'DISCONNECTED');
        expect(commit).toHaveBeenNthCalledWith(2, SET_WEBSOCKET_CHANNEL_FOUND, null);
        expect(commit).toHaveBeenNthCalledWith(3, SET_WEBSOCKET_CONNECTION_STATE, 'DISCONNECTED');
      });

      it('should return if Socket.leaveChannel errors', async () => {
        Socket.leaveChannel.mockRejectedValue({ error: 'Some error' });
        await actions.disconnect({ commit }, 'validationRunId');
        expect(commit).toHaveBeenCalledWith(SET_WEBSOCKET_CHANNEL_CONNECTION_STATE, JSON.stringify({ error: 'Some error' }));
      });
    });

    describe('subscribeToChannel', () => {
      it('should return if no validationRunId', async () => {
        await actions.subscribeToChannel({ commit });
        expect(commit).not.toHaveBeenCalled();
      });

      it('should return if no result.channel', async () => {
        Socket.findChannel.mockResolvedValue({ payload: {} });
        await actions.subscribeToChannel({ commit }, 'validationRunId');
        expect(commit).not.toHaveBeenCalled();
      });

      it('should set SET_WEBSOCKET_CHANNEL_CONNECTION_STATE to CONNECTED', async () => {
        const channel = new EventEmitter();
        const payload = {};
        channel.state = 'STATE';
        channel.topic = 'TOPIC';
        Socket.findChannel.mockResolvedValue({ payload: {}, channel });
        await actions.subscribeToChannel({ commit }, 'validationRunId');

        channel.emit('started', { payload });
        channel.emit('updated', { payload });
        channel.emit('completed', { payload });
        channel.emit('stopped', { payload });

        expect(commit).toHaveBeenCalledWith(SET_WEBSOCKET_LAST_UPDATE, payload);
        expect(commit).toHaveBeenCalledWith(SET_WEBSOCKET_LAST_UPDATE, payload);
        expect(commit).toHaveBeenCalledWith(SET_WEBSOCKET_LAST_UPDATE, payload);
        expect(commit).toHaveBeenCalledWith(SET_WEBSOCKET_LAST_UPDATE, payload);
        expect(commit).toHaveBeenNthCalledWith(1, SET_WEBSOCKET_CHANNEL_CONNECTION_STATE, 'CONNECTED');
        expect(commit).toHaveBeenNthCalledWith(2, SET_WEBSOCKET_CHANNEL_FOUND, `state=${channel.state}_topic=${channel.topic}`);
      });

      it('should return an error if Socket.findChannel errors', async () => {
        Socket.findChannel.mockRejectedValue({ error: 'Some error' });
        await actions.subscribeToChannel({ commit }, 'validationRunId');
        expect(commit).toHaveBeenNthCalledWith(1, SET_WEBSOCKET_CHANNEL_CONNECTION_STATE, JSON.stringify({ error: 'Some error' }));
        expect(commit).toHaveBeenNthCalledWith(2, SET_WEBSOCKET_CHANNEL_FOUND, JSON.stringify({ error: 'Some error' }));
      });
    });
  });
});
