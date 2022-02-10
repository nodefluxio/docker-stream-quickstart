const webpack = require("webpack");
const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");
const CopyPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const Dotenv = require("dotenv-webpack");
const ModuleFederationPlugin = require("webpack/lib/container/ModuleFederationPlugin");

const imageInlineSizeLimit = parseInt(
  process.env.IMAGE_INLINE_SIZE_LIMIT || "10000",
  10
);

module.exports = webpackEnv => {
  const isEnvProduction = !webpackEnv.WEBPACK_SERVE;

  return {
    entry: path.resolve(__dirname, "src/index.js"),
    output: {
      path: path.resolve(__dirname, "build"),
      filename: "bundle.[contenthash].js",
      publicPath: "/"
    },
    devServer: {
      hot: true,
      port: 8083,
      historyApiFallback: true
    },
    module: {
      rules: [
        {
          // "oneOf" will traverse all following loaders until one will
          // match the requirements. When no loader matches it will fall
          // back to the "file" loader at the end of the loader list.
          oneOf: [
            // TODO: Merge this config once `image/avif` is in the mime-db
            // https://github.com/jshttp/mime-db
            {
              test: [/\.avif$/],
              loader: require.resolve("url-loader"),
              options: {
                limit: imageInlineSizeLimit,
                mimetype: "image/avif",
                name: "[name].[hash:8].[ext]"
              }
            },
            // "url" loader works like "file" loader except that it embeds assets
            // smaller than specified limit in bytes as data URLs to avoid requests.
            // A missing `test` is equivalent to a match.
            {
              test: [/\.bmp$/, /\.gif$/, /\.jpe?g$/, /\.png$/],
              loader: require.resolve("url-loader"),
              options: {
                limit: imageInlineSizeLimit,
                name: "[name].[hash:8].[ext]"
              }
            },
            // Process application JS with Babel.
            // The preset includes JSX, Flow, TypeScript, and some ESnext features.
            {
              test: /\.(js|mjs|jsx|ts|tsx)$/,
              include: path.resolve(__dirname, "src"),
              loader: require.resolve("babel-loader"),
              options: {
                plugins: [
                  [
                    require.resolve("babel-plugin-named-asset-import"),
                    {
                      loaderMap: {
                        svg: {
                          ReactComponent:
                            "@svgr/webpack?-svgo,+titleProp,+ref![path]"
                        }
                      }
                    }
                  ]
                ].filter(Boolean),
                // This is a feature of `babel-loader` for webpack (not Babel itself).
                // It enables caching results in ./node_modules/.cache/babel-loader/
                // directory for faster rebuilds.
                cacheDirectory: true,
                // See #6846 for context on why cacheCompression is disabled
                cacheCompression: false,
                compact: isEnvProduction
              }
            },
            // Process any JS outside of the app with Babel.
            // Unlike the application JS, we only compile the standard ES features.
            {
              test: /\.(js|mjs)$/,
              exclude: /@babel(?:\/|\\{1,2})runtime/,
              loader: require.resolve("babel-loader"),
              options: {
                babelrc: false,
                configFile: false,
                compact: false,
                cacheDirectory: true,
                // See #6846 for context on why cacheCompression is disabled
                cacheCompression: false
              }
            },
            // "file" loader makes sure those assets get served by WebpackDevServer.
            // When you `import` an asset, you get its (virtual) filename.
            // In production, they would get copied to the `build` folder.
            // This loader doesn't use a "test" so it will catch all modules
            // that fall through the other loaders.
            {
              loader: require.resolve("file-loader"),
              // Exclude `js` files to keep "css" loader working as it injects
              // its runtime that would otherwise be processed through "file" loader.
              // Also exclude `html` and `json` extensions so they get processed
              // by webpacks internal loaders.
              exclude: [/\.(js|mjs|jsx|ts|tsx)$/, /\.html$/, /\.json$/],
              options: {
                name: "static/media/[name].[hash:8].[ext]"
              }
            }
            // ** STOP ** Are you adding a new loader?
            // Make sure to add the new loader(s) before the "file" loader.
          ]
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
      new CopyPlugin({
        patterns: [
          {
            from: path.resolve(__dirname, "public"),
            to: path.resolve(__dirname, "build"),
            globOptions: {
              ignore: ["**/index.html"]
            }
          }
        ]
      }),
      new webpack.ProvidePlugin({
        process: "process/browser",
        Buffer: ["buffer", "Buffer"]
      }),
      new MiniCssExtractPlugin({
        filename: "styles.[contentHash].css"
      }),
      new ModuleFederationPlugin({
        name: "VanillaDashboard",
        filename: "remoteEntry.js",
        shared: {
          react: {
            import: "react", // the "react" package will be used a provided and fallback module
            shareKey: "react", // under this name the shared module will be placed in the share scope
            shareScope: "default", // share scope with this name will be used
            singleton: true, // only a single version of the shared module is allowed
            eager: true,
            requiredVersion: "^17.0.2"
          },
          "styled-components": {
            singleton: true,
            eager: true
          },
          "core-js": { singleton: true, eager: true }
        }
      })
    ],
    resolve: {
      extensions: [".js", ".jsx"],
      fallback: {
        stream: require.resolve("stream-browserify"),
        util: require.resolve("util/"),
        buffer: require.resolve("buffer/")
      }
    }
  };
};
