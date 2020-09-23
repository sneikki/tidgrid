var gulp = require('gulp')
var sass = require('gulp-dart-sass')
var rename = require('gulp-rename')
var cleancss = require('gulp-clean-css')
var sourcemaps = require('gulp-sourcemaps')
var autoprefixer = require('gulp-autoprefixer')
var del = require('del')

gulp.task('transpile', function transpile() {
    return gulp
        .src('scss/tidgrid.scss')
        .pipe(sourcemaps.init())
        .pipe(autoprefixer({ cascade: false }))
        .pipe(sass({ includePaths: ['node_modules'] }).on('error', sass.logError))
        .pipe(rename('tidgrid.css'))
        .pipe(sourcemaps.write('.'))
        .pipe(gulp.dest('./dist'))
})

gulp.task('minify', function() {
    return gulp
        .src('./dist/tidgrid.css')
        .pipe(sourcemaps.init())
        .pipe(cleancss({ level: 2 }))
        .pipe(rename('tidgrid.min.css'))
        .pipe(sourcemaps.write('.'))
        .pipe(gulp.dest('./dist'))
})

gulp.task('build', gulp.series('transpile', 'minify'))

gulp.task('transpile:watch', function() {
    gulp.watch('./scss/**/*.scss', gulp.task('transpile'))
})

gulp.task('build:watch', function() {
    gulp.watch('./scss/**/*.scss', gulp.task('build'))
})

gulp.task('clean', async function() {
    return await del('dist')
})
