import jsonLocation from './jsonLocation';


describe('parse', () => {
  describe('empty JSON object', () => {
    it('returns empty array', () => {
      expect(jsonLocation.parse('{}')).toEqual({});
    });
  });

  describe('JSON object with single property', () => {
    it('returns object with property key and location', () => {
      expect(jsonLocation.parse(`{
        "discoveryModel": "example"
      }`)).toEqual({
        'discoveryModel': { column: 8, line: 2 }, // eslint-disable-line
      });
    });
  });

  describe('JSON object with nested object property', () => {
    it('returns object with property key and location', () => {
      expect(jsonLocation.parse(`{
        "discoveryModel": {
          "name": "example"
        }
      }`)).toEqual({
        'discoveryModel': { column: 8, line: 2 }, // eslint-disable-line
        'discoveryModel.name': { column: 10, line: 3 },
      });
    });
  });

  describe('JSON object with nested object property with nested object array', () => {
    it('returns object with property key and location', () => {
      expect(jsonLocation.parse(`{
        "discoveryModel": {
          "discoveryItems": [{
            "resourceBaseUri": "example",
          }]
        }
      }`)).toEqual({
        'discoveryModel': { column: 8, line: 2 }, // eslint-disable-line
        'discoveryModel.discoveryItems': { column: 10, line: 3 },
        'discoveryModel.discoveryItems[0]': { column: 29, line: 3 },
        'discoveryModel.discoveryItems[0].resourceBaseUri': { column: 12, line: 4 },
      });
    });
  });
});
