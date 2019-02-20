const esprima = require('esprima');

const addProperty = (paths, property, parents) => {
  const { value } = property.key;
  const { loc } = property.key;
  parents = parents.concat([value]); // eslint-disable-line
  const key = parents.join('.');
  paths[key] = loc; // eslint-disable-line

  if (property.value.properties) {
    property.value.properties.forEach((p) => {
      addProperty(paths, p, parents);
    });
  } else if (property.value.elements) {
    property.value.elements.forEach((e, i) => {
      const elementKey = `${key}[${i}]`;
      paths[elementKey] = e.loc; // eslint-disable-line
      e.properties.forEach((p) => {
        addProperty(paths, p, [elementKey]);
      });
    });
  }
};

export default {
  // Uses esprima to parse JSON string and
  // returns mapping of JSON path to start and end location column and line:
  // e.g. {
  //   "discoveryModel.discoveryItem[0].openidConfigurationUri":
  //      {"start":{"line":2,"column":8},"end":{"line":2,"column":24}}
  // }
  parse(string) {
    const tree = esprima.parseScript(`(${string})`, { loc: true });
    const { properties } = tree.body[0].expression;
    if (properties.length === 0) {
      return {};
    }
    const paths = {};
    const property = properties[0];
    addProperty(paths, property, []);
    return paths;
  },
};
