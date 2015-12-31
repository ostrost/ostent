// required Promises https://github.com/postcss/postcss-nested/issues/30
require('es6-promise').polyfill();

var ExtractTextPlugin = require('extract-text-webpack-plugin'),
    gulp              = require('gulp'),
    jade              = require('gulp-jade'),
    rename            = require('gulp-rename'),
    gutil             = require('gulp-util'),
    _                 = require('lodash'),
    path              = require('path'),
    webpack           = require('webpack'),
    argv              = require('yargs').argv;
require('./gulpfile.watch.js');

var bower_components = path.join(__dirname, './bower_components'),
    node_modules     = path.join(__dirname, './node_modules');
var wpconf = {
  resolve: {
    root: [bower_components, node_modules],
    //? extensions: ['', '.js', '.jsx', '.css', '.scss'],
    alias: {
      jquery:     'jquery/dist/jquery'  +'.js',
     'react-dom': 'react/lib/ReactDOM'  +'.js',
     'react-prm': 'react/lib/ReactComponentWithPureRenderMixin'+'.js'
    }
  },
  module: {
    loaders: [
      {test: /\.css$/,  loader: ExtractTextPlugin.extract('style-loader', 'css-loader')},
      {test: /\.scss$/, loader: ExtractTextPlugin.extract('style-loader', 'css-loader!sass-loader')},
      {
        test: /\.jsx?$/,
        loader: 'babel',
        exclude: /(node_modules|bower_components)/,
        query: {
          presets: ['react', 'es2015'],
          plugins: ['transform-react-jsx',
                    'babel-plugin-transform-react-constant-elements',
                    'babel-plugin-transform-react-inline-elements'],
          cacheDirectory: './share/cache'
        }
      }
    ],
    postLoaders: [{loader: 'transform?envify'}]
  },
  sassLoader: {includePaths: [
    bower_components+'/foundation-sites/scss/',
    bower_components+'/foundation-apps/scss/'
  ]},
  plugins: [
    // new webpack.ResolverPlugin(new webpack.ResolverPlugin.DirectoryDescriptionFilePlugin(
    //     'bower.json', ['main'])), // this will resolve with bower.json:"main"
    new ExtractTextPlugin('index.css', {allChunks: true})
  ]
};

gulp.task('jade', function() {
  gulp.src(argv.input)
    .pipe(jade({
      pretty: true,
      locals: {}
    }))
    .pipe(rename(argv.output))
    .pipe(gulp.dest('.'));
});

gulp.task('webpack', [], function(callback) {
  var wparg = wpconf;
  wparg = _.merge(wparg, {entry: argv.input});
  wparg = _.merge(wparg, {output: {}});
  wparg.output.filename = path.basename(argv.output);
  wparg.output.path     = path.join(__dirname, path.dirname(argv.output));
  if (argv.output.match(/\.min\.js($|\?)/i) !== null) {
    wparg.plugins.push(new webpack.optimize.UglifyJsPlugin({mangle: true}));
    process.env.NODE_ENV = 'production';
  }
  webpack(wparg).run(function(err, stats) {
    if (err != null) {
      throw new gutil.PluginError('webpack', err);
    }
    var statsString = stats.toString({chunks: false});
    if (stats.hasErrors() == true) {
      throw new gutil.PluginError('webpack', statsString);
    }
    gutil.log('[webpack]', statsString);
    callback();
  });
});

// Local variables:
// js-indent-level: 2
// js2-basic-offset: 2
// End:
