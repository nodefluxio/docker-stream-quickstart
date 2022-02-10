const webpack = require("webpack");
const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const ModuleFederationPlugin = require("webpack/lib/container/ModuleFederationPlugin");
const Dotenv = require("dotenv-webpack");
const CopyPlugin = require("copy-webpack-plugin");

const imageInlineSizeLimit = parseInt(
  process.env.IMAGE_INLINE_SIZE_LIMIT || "10000",
  10
);

module.exports = {
  entry: path.resolve(__dirname, "src/index.js"),
  output: {
    path: path.resolve(__dirname, "build"),
    filename: "bundle.[contenthash].js",
    publicPath: "auto"
  },
  devServer: {
    contentBase: path.resolve(__dirname, "build"),
    hot: true,
    port: 8082,
    historyApiFallback: true,
    headers: {
      // Enable wide open CORS
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, PATCH, OPTIONS",
      "Access-Control-Allow-Headers":
        "X-Requested-With, content-type, Authorization"
    }
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
        test: /\.css$/,
        use: ["style-loader", "css-loader"]
      },
      {
        test: [/\.bmp$/, /\.gif$/, /\.jpe?g$/, /\.png$/],
        loader: require.resolve("url-loader"),
        options: {
          limit: imageInlineSizeLimit,
          name: "[name].[hash:8].[ext]"
        }
      },
      {
        loader: require.resolve("file-loader"),
        exclude: [/\.(js|mjs|jsx|ts|tsx)$/, /\.html$/, /\.json$/, /\.css$/],
        options: {
          name: "[name].[hash:8].[ext]"
        }
      }
    ]
  },
  plugins: [
    new Dotenv({
      path: path.resolve(__dirname, ".env"),
      ignoreStub: true
    }),
    new CleanWebpackPlugin(),
    new webpack.HotModuleReplacementPlugin(),
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, "public/index.html")
    }),
    new webpack.ProvidePlugin({
      process: "process/browser",
      Buffer: ["buffer", "Buffer"]
    }),
    new CopyPlugin({
      patterns: [{ from: "public/env-config.json", to: "env-config.json" }]
    }),
    new MiniCssExtractPlugin({
      filename: "styles.[contentHash].css"
    }),
    new ModuleFederationPlugin({
      name: "Search",
      filename: "remoteEntry.js",
      exposes: {
        "./Button": "./src/components/molecules/Menu/index.js",
        "./Person": "./src/components/pages/SearchFace.js",
        "./Vehicle": "./src/components/pages/SearchPlate.js"
      },
      shared: [
        { react: { singleton: true, eager: true, requiredVersion: "^17.0.2" } },
        {
          "styled-components": {
            singleton: true,
            eager: true,
            requiredVersion: "^5.3.0"
          }
        },
        { "react-router-dom": { singleton: true } },
        { "prop-types": { singleton: true } },
        { "core-js": { singleton: true, eager: true } }
      ]
    })
  ],
  resolve: {
    extensions: [".js", ".jsx"]
  }
};
