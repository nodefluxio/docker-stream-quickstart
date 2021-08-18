const webpack = require("webpack");
const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const ModuleFederationPlugin = require("webpack/lib/container/ModuleFederationPlugin");

const imageInlineSizeLimit = parseInt(
  process.env.IMAGE_INLINE_SIZE_LIMIT || "10000",
  10
);

module.exports = {
  entry: path.resolve(__dirname, "src/index.js"),
  output: {
    path: path.resolve(__dirname, "build"),
    filename: "bundle.[contenthash].js"
  },
  devServer: {
    contentBase: path.resolve(__dirname, "build"),
    hot: true,
    port: 8082
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: ["babel-loader"],
        resolve: {
          extensions: [".js", ".jsx"]
        }
      },
      {
        test: [/\.bmp$/, /\.gif$/, /\.jpe?g$/, /\.png$/],
        loader: require.resolve("url-loader"),
        options: {
          limit: imageInlineSizeLimit,
          name: "static/media/[name].[hash:8].[ext]"
        }
      },
      {
        loader: require.resolve("file-loader"),
        exclude: [/\.(js|mjs|jsx|ts|tsx)$/, /\.html$/, /\.json$/],
        options: {
          name: "static/media/[name].[hash:8].[ext]"
        }
      }
    ]
  },
  plugins: [
    new CleanWebpackPlugin(),
    new webpack.HotModuleReplacementPlugin(),
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, "public/index.html")
    }),
    new MiniCssExtractPlugin({
      filename: "styles.[contentHash].css"
    }),
    new ModuleFederationPlugin({
      name: "SearchDukcapil",
      filename: "remoteEntry.js",
      exposes: {
        "./Button": "./src/components/atoms/LinkButton.js",
        "./Page": "./src/components/pages/SearchHome.js"
      },
      shared: [
        { react: { singleton: true, eager: true } },
        { "styled-components": { singleton: true } }
      ]
    })
  ],
  resolve: {
    extensions: [".js", ".jsx"]
  }
};
