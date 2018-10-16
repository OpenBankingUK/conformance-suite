import mutations from './mutations';
import {
  SET_VALIDATION_PAYLOAD,
  SET_VALIDATION_RUN_ID,
  SET_VALIDATION_STATUS,
} from './mutation-types';

const failedStatus = 'FAILED';
const payload = [
  {
    name: 'Sam Morse',
    sort_code: '111111',
    account_number: '12345678',
    amount: '10.00',
    type: 'payments',
  },
];
const validationRunId = 'validation-run-id';

describe('Validations', () => {
  describe('mutations', () => {
    it('SET_VALIDATION_PAYLOAD commits validation payload to the state', () => {
      const state = { payload: [], status: null };
      mutations[SET_VALIDATION_PAYLOAD](state, payload);
      expect(state.payload).toEqual(payload);
    });

    it('SET_VALIDATION_RUN_ID commits the validation run id to the state', () => {
      const state = { validationRunId: null };
      mutations[SET_VALIDATION_RUN_ID](state, validationRunId);
      expect(state.validationRunId).toEqual(validationRunId);
    });

    it('SET_VALIDATION_STATUS commits validation status to the state', () => {
      const state = { payload: [], status: null };
      mutations[SET_VALIDATION_STATUS](state, failedStatus);
      expect(state.status).toEqual(failedStatus);
    });
  });
});
