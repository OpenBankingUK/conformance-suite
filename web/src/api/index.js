import config from './config';
import discovery from './discovery';
import results from './results';
import testcases from './testcases';

export default {
  ...config,
  ...discovery,
  ...results,
  ...testcases,
};
