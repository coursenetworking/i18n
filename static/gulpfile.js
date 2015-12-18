var gulp = require('gulp'),
    uglify = require('gulp-uglify');

gulp.task('default', function(){
    gulp.src('src/app.js')
        .pipe(uglify())
        .pipe(gulp.dest('dist/'));
});
