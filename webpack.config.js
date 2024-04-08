const path = require('path');  

module.exports = {
  entry: './frontend/js/index.js',
  output: {
    path: path.resolve(__dirname, 'static', ".dist"),
    filename: 'main.js',
  }, 
};