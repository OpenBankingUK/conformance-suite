// https://stackoverflow.com/a/42191018/241993
if (process.env.NODE_ENV === 'test') {
  /* eslint-disable */
  // Implement require.context for tests
  if (typeof require.context === 'undefined') {
    const fs = require('fs');
    const path = require('path');

    require.context = (base = '.', scanSubDirectories = false, regularExpression = /\.js$/) => {
      const files = {};

      function readDirectory(directory) {
        fs.readdirSync(directory).forEach((file) => {
          const fullPath = path.resolve(directory, file);

          if (fs.statSync(fullPath).isDirectory()) {
            if (scanSubDirectories) readDirectory(fullPath);

            return;
          }

          if (!regularExpression.test(fullPath)) return;

          files[fullPath] = true;
        });
      }

      readDirectory(path.resolve(__dirname, base));

      function Module(file) {
        return require(file);
      }

      Module.keys = () => Object.keys(files);

      return Module;
    };
  }
  /* eslint-enable */
}

// Use webpack require.context to import templates and images.
// See: https://webpack.js.org/guides/dependency-management/#require-context
//      https://vuejs.org/v2/guide/components-registration.html#Automatic-Global-Registration-of-Base-Components
const requireTemplates = require.context('../../../../../pkg/discovery/templates/', false, /.+\.json$/);
const discoveryTemplates = requireTemplates.keys().map(file => requireTemplates(file));

const requireImages = require.context('./images/', false, /.+\.png$/);
const discoveryImages = {};
requireImages.keys().forEach(file => discoveryImages[file] = requireImages(file)); // eslint-disable-line

export default () => ({ discoveryTemplates, discoveryImages });
