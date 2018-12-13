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
        discoveryModel: {
          end: { column: 24, line: 2 },
          start: { column: 8, line: 2 },
        },
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
        'discoveryModel': {  // eslint-disable-line
          end: { column: 24, line: 2 },
          start: { column: 8, line: 2 },
        },
        'discoveryModel.name': {
          end: { column: 16, line: 3 },
          start: { column: 10, line: 3 },
        },
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
        'discoveryModel': {end: { column: 24, line: 2 }, start: { column: 8, line: 2 } }, // eslint-disable-line
        'discoveryModel.discoveryItems': { end: { column: 26, line: 3 }, start: { column: 10, line: 3 } },
        'discoveryModel.discoveryItems[0]': { end: { column: 11, line: 5 }, start: { column: 29, line: 3 } },
        'discoveryModel.discoveryItems[0].resourceBaseUri': { end: { column: 29, line: 4 }, start: { column: 12, line: 4 } },
      });
    });
  });
});
