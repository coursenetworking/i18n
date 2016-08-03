var gulp = require('gulp'),
    uglify = require('gulp-uglify'),
    rev = require('gulp-rev'),
    revCollector = require('gulp-rev-collector');

gulp.task('rev', ['js', 'css', 'bower_components'], function() {
    return gulp.src(['./rev/**/*.json', './src/index.html'])
        .pipe(revCollector({
            replaceReved: true,
            dirReplacements: {
                'css/': 'static/css/',
                'js/': 'static/js/',
                'bower_components/': 'static/bower_components/'
                'img/': 'static/img/'
            }
        }))
        .pipe(gulp.dest('./dist'))
});

gulp.task('revjs', ['tpl', 'rev'], function() {
    return gulp.src(['./rev/tpl/*.json', './dist/js/*.js'])
        .pipe(revCollector({
            replaceReved: true,
            dirReplacements: {
                'tpl/': 'static/tpl/'
            }
        }))
        .pipe(gulp.dest('./dist/js'))
});

gulp.task('js', function(){
    return gulp.src(['src/js/*.js', '!js/config.js'])
        .pipe(rev())
        .pipe(gulp.dest('dist/js'))
        .pipe(rev.manifest())
        .pipe(gulp.dest('rev/js'));
});

gulp.task('css', function(){
    return gulp.src(['src/css/*.css'])
        .pipe(rev())
        .pipe(gulp.dest('./dist/css'))
        .pipe(rev.manifest())
        .pipe(gulp.dest('rev/css'));
});

gulp.task('tpl', function(){
    return gulp.src(['src/tpl/*.html'])
        .pipe(rev())
        .pipe(gulp.dest('./dist/tpl'))
        .pipe(rev.manifest())
        .pipe(gulp.dest('rev/tpl'));
});

gulp.task('bower_components', function(){
    return gulp.src(['src/bower_components/**/*'])
        .pipe(rev())
        .pipe(gulp.dest('./dist/bower_components'))
        .pipe(rev.manifest())
        .pipe(gulp.dest('rev/bower_components'));
});

gulp.task('image', function(){
    return gulp.src(['src/img/*']).pipe(gulp.dest('./dist/img'));
});

gulp.task('default', ['image', 'js', 'css', 'tpl', 'bower_components', 'rev', 'revjs']);
