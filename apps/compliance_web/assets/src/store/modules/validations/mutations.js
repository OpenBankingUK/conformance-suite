import {
  SET_VALIDATION_PAYLOAD,
  SET_VALIDATION_RUN_ID,
  SET_VALIDATION_STATUS,
} from './mutation-types';

export default {
  [SET_VALIDATION_PAYLOAD](state, payload) {
    state.payload = payload;
  },
  [SET_VALIDATION_STATUS](state, status) {
    state.status = status;
  },
  [SET_VALIDATION_RUN_ID](state, validationRunId) {
    state.validationRunId = validationRunId;
  },
};
