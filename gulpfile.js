/* ----------------------------------------------------------------------------
 * Imports
 * ------------------------------------------------------------------------- */

var gulp       = require('gulp');
var $          = require('gulp-load-plugins')();
var path       = require('path');
var merge      = require('merge-stream');
var addsrc     = require('gulp-add-src');
var args       = require('yargs').argv;
var autoprefix = require('autoprefixer-core');
var clean      = require('del');
var collect    = require('gulp-rev-collector');
var concat     = require('gulp-concat');
var ignore     = require('gulp-ignore');
var mincss     = require('gulp-minify-css');
var minhtml    = require('gulp-htmlmin');
var modernizr  = require('gulp-modernizr');
var mqpacker   = require('css-mqpacker');
var notifier   = require('node-notifier');
var gulpif     = require('gulp-if');
var pixrem     = require('pixrem');
var plumber    = require('gulp-plumber');
var postcss    = require('gulp-postcss');
var reload     = require('gulp-livereload');
var rev        = require('gulp-rev');
var sass       = require('gulp-sass');
var sourcemaps = require('gulp-sourcemaps');
var sync       = require('gulp-sync')(gulp).sync;
var child      = require('child_process');
var uglify     = require('gulp-uglify');
var util       = require('gulp-util');
var vinyl      = require('vinyl-paths');

/* ----------------------------------------------------------------------------
 * Locals
 * ------------------------------------------------------------------------- */

/* Application server */
var server = null;

var DIST = 'dist';

var dist = function(subpath) {
  return !subpath ? DIST : path.join(DIST, subpath);
};


/* ----------------------------------------------------------------------------
 * Overrides
 * ------------------------------------------------------------------------- */

/*
 * Override gulp.src() for nicer error handling.
 */
var src = gulp.src;
gulp.src = function() {
  return src.apply(gulp, arguments)
    .pipe(plumber(function(error) {
      util.log(util.colors.red(
        'Error (' + error.plugin + '): ' + error.message
      ));
      notifier.notify({
        title: 'Error (' + error.plugin + ')',
        message: error.message.split('\n')[0]
      });
      this.emit('end');
    })
  );
};



/* ----------------------------------------------------------------------------
 * Assets pipeline
 * ------------------------------------------------------------------------- */

/*
 * Build stylesheets from SASS source.
 */
gulp.task('assets:stylesheets', function() {
  return gulp.src('app/styles/*.scss')
    .pipe(gulpif(args.sourcemaps, sourcemaps.init()))
    .pipe(sass({
      includePaths: [
        /* Your SASS dependencies via bower_components */
      ]}))
    .pipe(gulpif(args.production,
      postcss([
        autoprefix(),
        mqpacker,
        pixrem('10px')
      ])))
    .pipe(gulpif(args.sourcemaps, sourcemaps.write()))
    .pipe(gulpif(args.production, mincss()))
    .pipe(gulp.dest('public/styles/'))
    .pipe(reload());
});

// Copy all files at the root level (app)
gulp.task('copy', function() {

  // Add components to .tmp dir so they can get concatenated
  // when we vulcanize
  var app = gulp.src(['app/bower_components/**/*'])
    .pipe(gulp.dest(dist('bower_components')));

  return merge(app)
    .pipe($.size({
      title: 'copy'
    }));
});


// Transpile all JS to ES5.
gulp.task('js', function () {
 return gulp.src(['app/**/*.{js,html}', '!app/bower_components/**/*'])
   .pipe($.if('*.html', $.crisper({scriptInHead:false}))) // Extract JS from .html files
   .pipe($.sourcemaps.init())
   .pipe($.if('*.js', $.babel({
     presets: ['es2015']
   })))
   .pipe($.sourcemaps.write())
   .pipe(gulp.dest('.tmp/'))
   .pipe(gulp.dest('dist/'));
});


/*
 * Build javascripts from Bower components and source.
 */
gulp.task('assets:javascripts', function() {
  return gulp.src([
    /* Your JS dependencies via bower_components */
    /* Your JS libraries */
  ]).pipe(gulpif(args.sourcemaps, sourcemaps.init()))
    .pipe(concat('application.js'))
    .pipe(gulpif(args.sourcemaps, sourcemaps.write()))
    .pipe(gulpif(args.production, uglify()))
    .pipe(gulp.dest('public/javascripts/'))
    .pipe(reload());
});

/*
 * Create a customized modernizr build.
 */
gulp.task('assets:modernizr', function() {
  return gulp.src([
    'public/styles/style.css',
    'public/javascripts/application.js'
  ]).pipe(
      modernizr({
        options: [
          'addTest',                   /* Add custom tests */
          'fnBind',                    /* Use function.bind */
          'html5printshiv',            /* HTML5 support for IE */
          'setClasses',                /* Add CSS classes to root tag */
          'testProp'                   /* Test for properties */
        ]
      }))
    .pipe(addsrc.append('bower_components/respond/dest/respond.src.js'))
    .pipe(concat('modernizr.js'))
    .pipe(gulpif(args.production, uglify()))
    .pipe(gulp.dest('public/javascripts'));
});

/*
 * Minify views.
 */
gulp.task('assets:views', args.production ? [
  'assets:revisions:clean',
  'assets:revisions'
] : [], function() {
  return gulp.src([
    'manifest.json',
    'views/**/*.tmpl'
  ]).pipe(gulpif(args.production, collect()))
    .pipe(
      minhtml({
        collapseBooleanAttributes: true,
        collapseWhitespace: true,
        removeComments: true,
        removeScriptTypeAttributes: true,
        removeStyleLinkTypeAttributes: true,
        minifyCSS: true,
        minifyJS: true
      }))
    .pipe(gulp.dest('.views'));
});

/*
 * Clean outdated revisions.
 */
gulp.task('assets:revisions:clean', function() {
  return gulp.src(['public/**/*.{css,js}'])
    .pipe(ignore.include(/-[a-f0-9]{8}\.(css|js)$/))
    .pipe(vinyl(clean));
});

/*
 * Revision assets after build.
 */
gulp.task('assets:revisions', [
  'assets:revisions:clean'
], function() {
  return gulp.src(['public/**/*.{css,js}'])
    .pipe(ignore.exclude(/-[a-f0-9]{8}\.(css|js)$/))
    .pipe(rev())
    .pipe(gulp.dest('public'))
    .pipe(rev.manifest('manifest.json'))
    .pipe(gulp.dest('.'));
})

/*
 * Build assets.
 */
gulp.task('assets:build', [
  'copy',
  'js',
  'assets:stylesheets',
  'assets:javascripts',
  'assets:modernizr',
  'assets:views'
]);

/*
 * Watch assets for changes and rebuild on the fly.
 */
gulp.task('assets:watch', function() {

  /* Rebuild stylesheets on-the-fly */
  gulp.watch([
    'app/styles/**/*.scss'
  ], ['assets:stylesheets']);

  /* Rebuild javascripts on-the-fly */
  gulp.watch([
    'app/**/*.js',
    'bower.json'
  ], ['assets:javascripts']);

  /* Minify views on-the-fly */
  gulp.watch([
    'views/**/*.tmpl'
  ], ['assets:views']);
});

/* ----------------------------------------------------------------------------
 * Application server
 * ------------------------------------------------------------------------- */

/*
 * Build application server.
 */
gulp.task('server:build', function() {
  var build = child.spawnSync('go', ['install', './backend/server']);
  if (build.stderr.length) {
    var lines = build.stderr.toString()
      .split('\n').filter(function(line) {
        return line.length
      });
    for (var l in lines)
      util.log(util.colors.red(
        'Error (go install): ' + lines[l]
      ));
    notifier.notify({
      title: 'Error (go install)',
      message: lines
    });
  }
  return build;
});

/*
 * Restart application server.
 */
gulp.task('server:spawn', function() {
  if (server)
    server.kill();

  /* Spawn application server */
  server = child.spawn('server');

  /* Trigger reload upon server start */
  server.stdout.once('data', function() {
    reload.reload('/');
  });

  /* Pretty print server log output */
  server.stdout.on('data', function(data) {
    var lines = data.toString().split('\n')
    for (var l in lines)
      if (lines[l].length)
        util.log(lines[l]);
  });

  /* Print errors to stdout */
  server.stderr.on('data', function(data) {
    process.stdout.write(data.toString());
  });
});

/*
 * Watch source for changes and restart application server.
 */
gulp.task('server:watch', function() {

  /* Restart application server */
  gulp.watch([
    '.views/**/*.tmpl',
    'locales/*.json'
  ], ['server:spawn']);

  gulp.watch([
    'app/**/*.html'
  ], ['js', 'refresh'])

  /* Rebuild and restart application server */
  gulp.watch([
    '*/**/*.go',
  ], sync([
    'server:build',
    'server:spawn'
  ], 'server'));
});

gulp.task('refresh', ['js'], function() {
  reload.reload();
  console.log('livereload is triggered');
})

/* ----------------------------------------------------------------------------
 * Interface
 * ------------------------------------------------------------------------- */

/*
 * Build assets and application server.
 */
gulp.task('build', [
  'assets:build',
  'server:build'
]);

/*
 * Start asset and server watchdogs and initialize livereload.
 */
gulp.task('watch', [
  'assets:build',
  'server:build'
], function() {
  reload.listen();
  return gulp.start([
    'assets:watch',
    'server:watch',
    'server:spawn'
  ]);
});

/*
 * Build assets by default.
 */
gulp.task('default', ['build']);
