var child             = require('child_process'),
    gulp              = require('gulp');
//  reload            = require('gulp-livereload');

var makevars = {}, server;

function dup(std, proc) {
  // return proc[std].on('data', function(d) { process[std].write(d.toString()); });
  return proc[std].on('data', function(d) {
    return process[std].write(d.toString());
    var lines = d.toString().split('\n');
    for (var i in lines) {
      util.log(lines[i]); // util is gulp-util
    }
  });
}

gulp.task('make dev', function() {
  var make = child.spawn('make', '-C . dev'.split(' '));
  dup('stdout', make);
  dup('stderr', make);
});

gulp.task('make print', function(cb) {
  var make = child.spawn('make', 'print-package print-devpackagefiles'.split(' '));
  make.stdout.on('end', cb);
  make.stdout.on('data', function(d) {
    var words = d.toString().split('=');
    makevars[words[0]] = words[1].replace(/\n$/, '').split(' ');
  });
  dup('stderr', make);
});

gulp.task('server build', function(cb) {
  var build = child.spawn('go', 'install -race'.split(' '));
  build.stdout.on('end', cb);
  dup('stdout', build);
  dup('stderr', build);
});

gulp.task('server run', ['server build'], function() {
  var run = function() {
    // console.log('makevars (run)', makevars);
    server = child.spawn(makevars.package[0].replace(/.*\//, ''),
                         '-b 127.0.0.1'.split(' '));
    dup('stdout', server);
    dup('stderr', server);
    // server.stdout.once('data', function() { reload.reload('/'); });
  };
  if (server != null) {
    server.kill();
    setTimeout(run, 1000);
  } else {
    run();
  }
});

gulp.task('server watch', function() {
  gulp.watch(makevars.devpackagefiles.map(function(x) { return __dirname + '/' + x; }),
             ['server run']);
});

gulp.task('watch', ['make print'], function() {
  // reload.listen();
  gulp.start(['server watch', 'server run']);
  gulp.watch([
    __dirname+'/share/ace.templates/*',
    __dirname+'/share/js/*',
    __dirname+'/share/style/*'
  ], ['make dev']);
});

// Local variables:
// js-indent-level: 2
// js2-basic-offset: 2
// End:
