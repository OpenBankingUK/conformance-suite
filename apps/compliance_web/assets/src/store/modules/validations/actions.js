import validations from '../../../api/validations';
import {
  SET_VALIDATION_PAYLOAD,
  SET_VALIDATION_RUN_ID,
  SET_VALIDATION_STATUS,
} from './mutation-types';

export default {
  async validate({ commit }, validation) {
    commit(SET_VALIDATION_PAYLOAD, validation.payload);
    try {
      const started = await validations.start(validation);
      commit(SET_VALIDATION_RUN_ID, started.data.id);
      const info = await validations.track(started.data.id);
      const { status } = info.data;
      commit(SET_VALIDATION_STATUS, status);
    } catch (e) {
      commit(SET_VALIDATION_STATUS, 'FAILED');
    }
  },
};
