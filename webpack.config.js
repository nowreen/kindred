const path = require('path');
const webpack = require('webpack');

const config = {
  devtool: 'source-map',
  entry: path.join(__dirname, './client/index.jsx'),
	
  output: {
    path: path.join(__dirname, 'bundles'),
    filename: 'bundle.js',
  },
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
				
        query: {
          presets: ['env', 'react']
        }
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader']
      }
    ],
  },
  plugins: [
    new webpack.DefinePlugin({
      environment: JSON.stringify(process.env.ENV)
    })
  ]
};

module.exports = config;