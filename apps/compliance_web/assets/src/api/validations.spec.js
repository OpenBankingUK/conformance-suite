import Vue from 'vue';
import validations from './validations';

describe('Validations api', () => {
  describe('start', () => {
    it('should throw an error if status != 202', async () => {
      Vue.axios.post.mockResolvedValue({ status: 400 });
      try {
        await validations.start('some validation');
      } catch (e) {
        expect(e).toEqual(new Error('Expected 202 Accepted Status.'));
      }
    });

    it('should return data if no error', async () => {
      Vue.axios.post.mockResolvedValue({ status: 202, data: {} });
      const start = await validations.start('some validation');
      expect(start).toEqual({});
    });
  });

  describe('track', () => {
    it('should throw an error if status != 200', async () => {
      Vue.axios.get.mockResolvedValue({ status: 400 });
      try {
        await validations.track('some id');
      } catch (e) {
        expect(e).toEqual(new Error('Expected 200 Ok Status.'));
      }
    });

    it('should return data if no error', async () => {
      Vue.axios.get.mockResolvedValue({ status: 200, data: {} });
      const track = await validations.track('some id');
      expect(track).toEqual({});
    });
  });
});
