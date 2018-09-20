import actions from './actions';
import {
  SET_VALIDATION_PAYLOAD,
  SET_VALIDATION_RUN_ID,
  SET_VALIDATION_STATUS,
} from './mutation-types';
import validations from '../../../api/validations';

jest.mock('../../../api/validations');

describe('Validations', () => {
  describe('actions', () => {
    let commit;

    beforeEach(() => {
      commit = jest.fn();
    });

    afterEach(() => {
      jest.resetAllMocks();
    });

    describe('validate', () => {
      const validation = {
        payload: {},
      };

      it('should fail if there is an error in api/validations', async () => {
        validations.start.mockRejectedValue();
        await actions.validate({ commit }, validation);
        expect(commit).toHaveBeenNthCalledWith(1, SET_VALIDATION_PAYLOAD, validation.payload);
        expect(commit).toHaveBeenNthCalledWith(2, SET_VALIDATION_STATUS, 'FAILED');
      });

      it('should call SET_VALIDATION_RUN_ID and SET_VALIDATION_STATUS', async () => {
        const startValue = { data: { id: 1 } };
        const trackValue = { data: { status: 'Some status' } };

        validations.start.mockResolvedValue(startValue);
        validations.track.mockResolvedValue(trackValue);
        await actions.validate({ commit }, validation);
        expect(commit).toHaveBeenNthCalledWith(1, SET_VALIDATION_PAYLOAD, validation.payload);
        expect(commit).toHaveBeenNthCalledWith(2, SET_VALIDATION_RUN_ID, startValue.data.id);
        expect(commit).toHaveBeenNthCalledWith(3, SET_VALIDATION_STATUS, trackValue.data.status);
      });
    });
  });
});
