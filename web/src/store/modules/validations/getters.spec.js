import getters from './getters';

const processingStatus = 'PROCESSING';
const authorisationServerId = 'auth-server-id';

describe('Validations', () => {
  const SELECTED_ASPSP = 'ASPSP_1_ID';
  const EXPECTED_AUTHORISATION_SERVERS_LIST = [
    {
      text: 'AAA Example Bank',
      value: {
        name: 'AAA Example Bank',
        logoUri: '',
        id: 'aaaj4NmBD8lQxmLh2O',
      },
      disabled: false,
    },
    {
      text: 'BBB Example Bank',
      value: {
        name: 'BBB Example Bank',
        logoUri: '',
        id: 'bbbX7tUB4fPIYB0k1m',
      },
      disabled: false,
    },
    {
      text: 'CCC Example Bank',
      value: {
        name: 'CCC Example Bank',
        LogoUri: '',
        id: 'cccbN8iAsMh74sOXhk',
      },
      disabled: false,
    },
  ];

  describe('getters', () => {
    it('validationStatus returns the last validation status', () => {
      const state = { payload: [], status: processingStatus };
      expect(getters.validationStatus(state)).toEqual(processingStatus);
    });

    it('authorisationServers returns the last the fetched list of ASPSP', () => {
      const stateOld = {
        payload: [],
        status: processingStatus,
      };
      expect(getters.authorisationServers(stateOld)).toEqual(undefined);

      const stateNew = {
        payload: [],
        status: processingStatus,
        authorisationServers: EXPECTED_AUTHORISATION_SERVERS_LIST,
      };
      expect(getters.authorisationServers(stateNew)).toEqual(EXPECTED_AUTHORISATION_SERVERS_LIST);
    });

    it('selectedAspsp returns the last selected ASPSP', () => {
      const stateOld = {
        payload: [],
        status: processingStatus,
      };
      expect(getters.selectedAspsp(stateOld)).toEqual(undefined);

      const stateNew = {
        payload: [],
        status: processingStatus,
        selectedAspsp: SELECTED_ASPSP,
      };
      expect(getters.selectedAspsp(stateNew)).toEqual(SELECTED_ASPSP);
    });

    it('selectedAspsp returns null when no aspsp has been selected', () => {
      const state = { payload: [], status: null, selectedAspsp: null };
      expect(getters.selectedAspsp(state)).toEqual(null);
    });

    it('selectedAspsp returns the id when as aspsp has been selected', () => {
      const state = { payload: [], status: null, selectedAspsp: authorisationServerId };
      expect(getters.selectedAspsp(state)).toEqual(authorisationServerId);
    });

    it('validationRunId returns the validationRunId', () => {
      expect(getters.validationRunId({})).toEqual(undefined);

      const state = { validationRunId: '<validationRunId>' };
      expect(getters.validationRunId(state)).toEqual('<validationRunId>');
    });
  });
});
