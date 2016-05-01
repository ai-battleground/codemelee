const gulp         = require('gulp');
const babelify     = require('babelify');
const browserify   = require('browserify');
const browserSync  = require('browser-sync');
const concat       = require('gulp-concat');
const eslint       = require('gulp-eslint');
const newer        = require('gulp-newer');
const notify       = require('gulp-notify');
const reload       = browserSync.reload;
const source       = require('vinyl-source-stream');
const sourcemaps   = require('gulp-sourcemaps');
const gutil        = require('gulp-util')
const watchify     = require('watchify')

gulp.task('lint', function() {
  gulp.src('src/js/src/components/CodemeleeTetris.jsx')
    .pipe(eslint())
    .pipe(eslint.format())
    .pipe(eslint.failAfterError());
});

gulp.task('copy-css', function() {
  gulp.src('src/css/*.css')
    .pipe(gulp.dest('dist/css'));
});

gulp.task('copy-html', function() {
  gulp.src('index.html')
    .pipe(gulp.dest('dist'));
});

gulp.task('build', ['copy-css', 'copy-html'], function() {
  const b = browserify({
    entries: 'src/js/src/components/CodemeleeTetris.jsx',
    extensions: ['.jsx'],
    cache: {},
    packageCache: {},
    plugin: [watchify],
    debug: true
  });

  function bundle() {
    b.transform(babelify, {presets: ['es2015', 'react']})
      .bundle()
      .on('error',gutil.log)
      .pipe(source('app.js'))
      .pipe(gulp.dest('dist/js'));
  }

  b.on('update', bundle);
  bundle();
});

// BrowserSync
gulp.task('browsersync', function() {
  browserSync({
    server: {
      baseDir: 'dist/'
    }
  });
});

gulp.task('default', ['build', 'browsersync']);
