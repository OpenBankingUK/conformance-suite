import api from './apiUtil';

let setShowLoading;

describe('api.get', () => {
  const expectedOptions = {
    method: 'GET',
    headers: {
      Accept: 'application/json; charset=UTF-8',
      'Content-Type': 'application/json; charset=UTF-8',
    },
  };
  const url = '/api/test-cases';
  const data = { some: 'data' };

  beforeEach(() => {
    setShowLoading = jest.fn();
    fetch.resetMocks();
    fetch.mockResponseOnce(JSON.stringify(data), { status: 200 });
  });

  it('calls fetch once with expected url and options and returns result', async () => {
    try {
      const response = await api.get(url);
      expect(await response.json()).toEqual(data);

      // assert on the times called and arguments given to fetch
      expect(fetch.mock.calls.length).toEqual(1);
      expect(fetch.mock.calls[0][0]).toEqual(url);
      expect(fetch.mock.calls[0][1]).toEqual(expectedOptions);
    } catch (err) {
      // Should not get here.
      expect(err).toBeFalsy();
    }
  });

  it('calls setShowLoading function when setShowLoading provided', async () => {
    try {
      await api.get(url, setShowLoading);
    } catch (err) {
      expect(err).toBeFalsy();
    }
    expect(setShowLoading.mock.calls.length).toEqual(2);
    expect(setShowLoading.mock.calls[0][0]).toEqual(true);
    expect(setShowLoading.mock.calls[1][0]).toEqual(false);
  });
});

describe('api.post', () => {
  const expectedOptions = ({ body }) => ({
    method: 'POST',
    headers: {
      Accept: 'application/json; charset=UTF-8',
      'Content-Type': 'application/json; charset=UTF-8',
    },
    body,
  });
  const url = '/api/test-cases';
  const data = { some: 'data' };

  beforeEach(() => {
    setShowLoading = jest.fn();
    fetch.resetMocks();
    fetch.mockResponseOnce(JSON.stringify(data), { status: 200 });
  });

  it('calls fetch once with expected url and options and returns result', async () => {
    try {
      const response = await api.post(url, data);
      expect(await response.json()).toEqual(data);

      expect(fetch.mock.calls.length).toEqual(1);
      expect(fetch.mock.calls[0][0]).toEqual(url);
      expect(fetch.mock.calls[0][1]).toEqual(expectedOptions({ body: JSON.stringify(data) }));
    } catch (err) {
      expect(err).toBeFalsy();
    }
  });

  it('calls fetch with null body when called without data', async () => {
    try {
      const response = await api.post(url);
      expect(await response.json()).toEqual(data);

      expect(fetch.mock.calls.length).toEqual(1);
      expect(fetch.mock.calls[0][0]).toEqual(url);
      expect(fetch.mock.calls[0][1]).toEqual(expectedOptions({ body: null }));
    } catch (err) {
      expect(err).toBeFalsy();
    }
  });

  it('calls setShowLoading function when setShowLoading provided', async () => {
    try {
      await api.post(url, data, setShowLoading);
    } catch (err) {
      expect(err).toBeFalsy();
    }
    expect(setShowLoading.mock.calls.length).toEqual(2);
    expect(setShowLoading.mock.calls[0][0]).toEqual(true);
    expect(setShowLoading.mock.calls[1][0]).toEqual(false);
  });
});
